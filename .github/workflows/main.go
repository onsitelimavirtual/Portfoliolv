package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type VM struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	RAM    int    `json:"ram"`
	CPU    int    `json:"cpu"`
	DiskGB int    `json:"diskGB"`
	ISO    string `json:"iso"`
	Disk   string `json:"disk"`
}
type App struct {
	mu        sync.Mutex
	VMs       []VM `json:"vms"`
	processes map[string]*exec.Cmd
	dir       string
	token     string
	csrf      string
	qemu      string
	img       string
	port      string
}

func random() string {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func dataDir() string {
	if d := os.Getenv("NATORVM_HOME"); d != "" {
		return d
	}
	base, e := os.UserConfigDir()
	if e != nil {
		return ".natorvm"
	}
	return filepath.Join(base, "NatorVM10")
}
func (a *App) save() error {
	b, e := json.MarshalIndent(a.VMs, "", "  ")
	if e != nil {
		return e
	}
	f := filepath.Join(a.dir, "machines.json")
	t := f + ".tmp"
	if e = os.WriteFile(t, b, 0600); e != nil {
		return e
	}
	return os.Rename(t, f)
}
func (a *App) auth(w http.ResponseWriter, r *http.Request) bool {
	c, e := r.Cookie("nator_session")
	if e != nil || subtle.ConstantTimeCompare([]byte(c.Value), []byte(a.token)) != 1 {
		http.Error(w, "Acesso negado", 403)
		return false
	}
	if r.Method == "POST" {
		if subtle.ConstantTimeCompare([]byte(r.FormValue("csrf")), []byte(a.csrf)) != 1 {
			http.Error(w, "Token CSRF inválido", 403)
			return false
		}
	}
	return true
}
func (a *App) find(id string) (int, error) {
	for i, v := range a.VMs {
		if v.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("Máquina não localizada")
}
func (a *App) running(id string) bool { return a.processes[id] != nil }
func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	if !a.auth(w, r) {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	type Entry struct {
		VM      VM
		Running bool
	}
	var vs []Entry
	for _, v := range a.VMs {
		vs = append(vs, Entry{v, a.running(v.ID)})
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
	t := template.Must(template.New("p").Parse(page))
	_ = t.Execute(w, map[string]any{"VMs": vs, "CSRF": a.csrf, "Qemu": a.qemu != "", "Img": a.img != "", "Host": runtime.GOOS, "Data": a.dir})
}
func (a *App) message(w http.ResponseWriter, msg string) {
	http.Redirect(w, &http.Request{}, "/", http.StatusSeeOther)
	log.Println(msg)
}
func (a *App) action(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método inválido", 405)
		return
	}
	if !a.auth(w, r) {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	act := r.FormValue("action")
	id := r.FormValue("id")
	var err error
	switch act {
	case "create":
		name := strings.TrimSpace(r.FormValue("name"))
		ram, _ := strconv.Atoi(r.FormValue("ram"))
		cpu, _ := strconv.Atoi(r.FormValue("cpu"))
		gb, _ := strconv.Atoi(r.FormValue("diskgb"))
		iso := strings.TrimSpace(r.FormValue("iso"))
		if len(name) < 1 || len(name) > 80 || ram < 256 || ram > 3072 || cpu < 1 || cpu > 4 || gb < 1 || gb > 2048 {
			err = errors.New("Nome, RAM, CPU ou disco fora dos limites")
			break
		}
		if iso != "" {
			inf, e := os.Stat(iso)
			if e != nil || inf.IsDir() {
				err = errors.New("ISO não encontrada")
				break
			}
		}
		if a.img == "" {
			err = errors.New("qemu-img indisponível: configure o QEMU antes de criar discos")
			break
		}
		id = random()[:16]
		disk := filepath.Join(a.dir, "disks", id+".qcow2")
		cmd := exec.Command(a.img, "create", "-f", "qcow2", disk, fmt.Sprintf("%dG", gb))
		if output, e := cmd.CombinedOutput(); e != nil {
			err = fmt.Errorf("Falha qemu-img: %s %v", string(output), e)
			break
		}
		a.VMs = append(a.VMs, VM{id, name, ram, cpu, gb, iso, disk})
		err = a.save()
		if err != nil {
			a.VMs = a.VMs[:len(a.VMs)-1]
			_ = os.Remove(disk)
		}
	case "start":
		i, e := a.find(id)
		if e != nil {
			err = e
			break
		}
		if a.running(id) {
			err = errors.New("VM já executando")
			break
		}
		if a.qemu == "" {
			err = errors.New("QEMU indisponível")
			break
		}
		v := a.VMs[i]
		args := []string{"-name", v.Name, "-m", strconv.Itoa(v.RAM), "-smp", strconv.Itoa(v.CPU), "-drive", "file=" + v.Disk + ",format=qcow2,if=ide", "-nic", "user,model=e1000", "-display", "default"}
		if v.ISO != "" {
			args = append(args, "-cdrom", v.ISO, "-boot", "order=d")
		}
		if runtime.GOOS == "linux" {
			if _, e := os.Stat("/dev/kvm"); e == nil {
				args = append(args, "-enable-kvm")
			}
		}
		cmd := exec.Command(a.qemu, args...)
		if e := cmd.Start(); e != nil {
			err = e
			break
		}
		a.processes[id] = cmd
		go func() {
			e := cmd.Wait()
			a.mu.Lock()
			delete(a.processes, id)
			a.mu.Unlock()
			log.Printf("VM %s encerrada: %v", id, e)
		}()
	case "stop":
		if !a.running(id) {
			err = errors.New("VM não está ligada")
			break
		}
		err = a.processes[id].Process.Kill() // emergency only; no graceful shutdown promised
	case "delete":
		i, e := a.find(id)
		if e != nil {
			err = e
			break
		}
		if a.running(id) {
			err = errors.New("Desligue a VM antes de excluir")
			break
		}
		v := a.VMs[i]
		a.VMs = append(a.VMs[:i], a.VMs[i+1:]...)
		if e = a.save(); e != nil {
			err = e
			break
		}
		err = os.Remove(v.Disk)
	default:
		err = errors.New("Ação desconhecida")
	}
	if err != nil {
		log.Println("Erro:", err)
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/", 303)
}
func locateBinary(name string) string {
	candidates := []string{name}
	if runtime.GOOS == "windows" {
		candidates = []string{name + ".exe"}
	}
	for _, n := range candidates {
		if p, e := exec.LookPath(n); e == nil {
			return p
		}
	}
	exe, _ := os.Executable()
	dirs := []string{filepath.Dir(exe), filepath.Join(filepath.Dir(exe), "qemu"), filepath.Join(filepath.Dir(exe), "bin")}
	if runtime.GOOS == "windows" {
		for _, env := range []string{"ProgramFiles", "ProgramFiles(x86)"} {
			if b := os.Getenv(env); b != "" {
				dirs = append(dirs, filepath.Join(b, "qemu"), filepath.Join(b, "QEMU"))
			}
		}
	}
	for _, d := range dirs {
		for _, n := range candidates {
			p := filepath.Join(d, n)
			if st, e := os.Stat(p); e == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func main() {
	dir := dataDir()
	if e := os.MkdirAll(filepath.Join(dir, "disks"), 0700); e != nil {
		log.Fatal(e)
	}
	a := &App{dir: dir, token: random(), csrf: random(), processes: map[string]*exec.Cmd{}}
	if b, e := os.ReadFile(filepath.Join(dir, "machines.json")); e == nil {
		if e = json.Unmarshal(b, &a.VMs); e != nil {
			log.Fatal("Catálogo inválido: ", e)
		}
	}
	a.qemu = locateBinary("qemu-system-i386")
	if a.qemu == "" {
		a.qemu = locateBinary("qemu-system-x86_64")
	}
	a.img = locateBinary("qemu-img")
	port := "0"
	if len(os.Args) > 1 && os.Args[1] == "--diagnostico" {
		fmt.Printf("NatorVM 10.3 Legacy Windows 7 (32 bits) | %s/%s\nQEMU: %v\nqemu-img: %v\nPasta: %s\n", runtime.GOOS, runtime.GOARCH, a.qemu, a.img, dir)
		return
	}
	ln, e := net.Listen("tcp", "127.0.0.1:"+port)
	if e != nil {
		log.Fatal(e)
	}
	url := "http://" + ln.Addr().String()
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.dashboard)
	mux.HandleFunc("/action", a.action)
	mux.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login/"+a.token {
			http.NotFound(w, r)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "nator_session", Value: a.token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode})
		http.Redirect(w, r, "/", 303)
	})
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 30 * time.Second}
	fmt.Println("NatorVM 10.3 Legacy Windows 7 (32 bits) — execução local")
	fmt.Println("Abra no NAVEGADOR DESTE COMPUTADOR:", url+"/login/"+a.token)
	fmt.Println("Mantenha esta janela aberta. Não divulgue o link com token.")
	log.Fatal(srv.Serve(ln))
}

const page = `<!doctype html><html lang="pt-BR"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>NatorVM 10.3 Legacy Windows 7 (32 bits)</title><style>body{margin:0;background:#0b1220;color:#eaf1ff;font:16px system-ui,Arial}header{padding:22px;background:#17243c;border-bottom:1px solid #34445a}main{max-width:1040px;margin:auto;padding:20px}h1{margin:0}small,p{color:#c4d2e6}section{background:#17243c;border:1px solid #34445a;border-radius:12px;padding:18px;margin:16px 0}input{padding:12px;background:#0d1729;color:white;border:1px solid #607491;border-radius:7px;width:100%;box-sizing:border-box}label{display:block;margin:9px 0}button{background:#31a6dd;color:#07131f;border:0;padding:11px 15px;border-radius:8px;font-weight:bold;margin:4px;cursor:pointer}.danger{background:#f49d92}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(170px,1fr));gap:8px}.vm{border-top:1px solid #40516b;padding:14px 0}code{word-break:break-all}strong{color:#fff}</style></head><body><header><h1>◈ NatorVM 10.3 Legacy Windows 7 (32 bits)</h1><small>Gerenciador local • Sem Python ou Node.js • Copyright © 2026 Nailson Rodrigues de Lima</small></header><main><section><h2>Saúde do ambiente</h2><p>Host: <strong>{{.Host}}</strong> · QEMU: <strong>{{if .Qemu}}Detectado{{else}}Ausente{{end}}</strong> · qemu-img: <strong>{{if .Img}}Detectado{{else}}Ausente{{end}}</strong></p><p>Dados: <code>{{.Data}}</code></p><small>Compatibilidade Win7 requer binário produzido com Go 1.20 e QEMU 32-bit legado; limite sugerido até 3072 MB RAM por VM. A execução requer QEMU compatível.</small></section><section><h2>Criar máquina virtual</h2><form method="post" action="/action"><input type="hidden" name="csrf" value="{{.CSRF}}"><input type="hidden" name="action" value="create"><label>Nome<input name="name" maxlength="80" required placeholder="Minha máquina virtual"></label><div class="grid"><label>RAM (MB)<input type="number" name="ram" min="256" max="3072" value="1024"></label><label>CPUs<input type="number" name="cpu" min="1" max="4" value="1"></label><label>Disco (GB)<input type="number" name="diskgb" min="1" max="2048" value="20"></label></div><label>Caminho absoluto da ISO (opcional)<input name="iso" placeholder="/caminho/sistema.iso"></label><button>Criar VM e disco QCOW2</button></form></section><section><h2>Máquinas virtuais</h2>{{range .VMs}}<div class="vm"><h3>{{.VM.Name}} · {{if .Running}}🟢 ligada{{else}}⚪ desligada{{end}}</h3><p>{{.VM.RAM}} MB · {{.VM.CPU}} vCPU · {{.VM.DiskGB}} GB</p><div class="grid"><form method="post" action="/action"><input type="hidden" name="csrf" value="{{$.CSRF}}"><input type="hidden" name="id" value="{{.VM.ID}}"><button name="action" value="start">Iniciar</button></form><form method="post" action="/action" onsubmit="return confirm('FORÇAR parada? Pode corromper dados.')"><input type="hidden" name="csrf" value="{{$.CSRF}}"><input type="hidden" name="id" value="{{.VM.ID}}"><button name="action" value="stop" class="danger">Parada de emergência</button></form><form method="post" action="/action" onsubmit="return confirm('EXCLUIR definitivamente a VM e o disco?')"><input type="hidden" name="csrf" value="{{$.CSRF}}"><input type="hidden" name="id" value="{{.VM.ID}}"><button name="action" value="delete" class="danger">Excluir VM</button></form></div></div>{{else}}<p>Nenhuma VM cadastrada.</p>{{end}}</section><section><h2>Limites da edição nativa</h2><p>Não inclui QEMU, noVNC ou sistemas operacionais convidados. Snapshots, backup avançado e acesso remoto permanecem disponíveis na base legada 9.1 (requer Python) ou precisam ser portados para o núcleo nativo. Não exponha o painel à Internet. Windows 7 está sem suporte de segurança; mantenha o host isolado.</p></section><section id="termos"><h2>Termos de Uso</h2><p>O NatorVM 10.0 é um gerenciador local em edição de avaliação. O usuário é responsável por possuir licença das ISOs e sistemas convidados, realizar backups e validar compatibilidade. A exclusão de VMs e a parada forçada podem causar perda de dados. O motor QEMU não está incluído neste executável. Não utilize em sistemas críticos sem homologação. Consulte <code>TERMOS_DE_USO.md</code> no pacote para as condições completas.</p></section><section id="privacidade"><h2>Política de Privacidade</h2><p>O NatorVM mantém o catálogo e os discos das VMs localmente na pasta de dados do usuário, ou em <code>NATORVM_HOME</code>. A versão nativa não inclui telemetria nem envio intencional dos dados a servidores do desenvolvedor. O painel usa cookie temporário de sessão e fica restrito a 127.0.0.1. O tráfego dos sistemas convidados deve ser avaliado separadamente. Consulte <code>POLITICA_DE_PRIVACIDADE.md</code> no pacote.</p></section><footer><p><a href="#termos" style="color:#75cffd">Termos de Uso</a> · <a href="#privacidade" style="color:#75cffd">Política de Privacidade</a></p><small>Copyright © 2026 Nailson Rodrigues de Lima. Direitos do código original reservados, respeitadas as licenças de terceiros.</small></footer></main></body></html>`
