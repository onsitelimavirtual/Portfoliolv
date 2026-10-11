# NatorVM 10.3 Legacy — Windows 7 32 bits

## O que esta edição entrega
Código-fonte Go ajustado para detecção de `qemu-system-i386`, com limites de RAM (até 3072 MB por VM) e CPU (até 4 vCPUs), interface local HTML, criação de discos QCOW2 e comandos de inicialização e parada de emergência. Os recursos dependem do QEMU executável no host.

## Limitações decisivas
- **Não inclui binário homologado para Windows 7.** O compilador desta sessão é Go 1.23.2, incompatível com o Windows 7. Para criar o executável legado utilize **Go 1.20.x** para Windows 32 bits e `COMPILAR-WIN7-32.bat`, depois execute os testes em Windows 7.
- Inclui apenas um executável **de referência para Windows 32 bits mais recente**, compilado com Go 1.23.2, em `bin/REFERENCIA-Windows-x86-Go123.exe`. **NÃO o utilize como executável para Windows 7.**
- QEMU e `qemu-img` **não estão incluídos**. Obtenha uma compilação Win32 que suporte Windows 7 a partir de uma fonte confiável e licenciada, e coloque os arquivos executáveis necessários em `qemu/`, incluindo DLLs exigidas pela compilação. Não misture executáveis e DLLs de releases diferentes.
- No Windows 7 32 bits, limitações de memória e desempenho podem inviabilizar convidados modernos. Recomenda-se convidado Linux leve e host com RAM suficiente.
- O navegador do Windows 7 pode não suportar recursos modernos do painel; teste com navegador compatível instalado no host.
- O botão **Parada de emergência** encerra à força, com risco de corrupção. Faça backups antes de testar.
- O painel local é acessível apenas em 127.0.0.1. Nunca exponha serviços Windows 7 diretamente à Internet.

## Sequência de uso
1. Execute `DIAGNOSTICO-WIN7-32.bat` no computador.
2. Instale ou prepare **Go 1.20.x 32-bit** e compile `COMPILAR-WIN7-32.bat` (ou compile em ambiente de build compatível).
3. Providencie os binários QEMU Win32 compatíveis e coloque em `qemu/`.
4. Execute `INICIAR-WIN7-32.bat` e abra o URL temporário mostrado no console, no mesmo PC.
5. Crie VM com ISO legítima de sistema operacional leve, 512–1024 MB RAM e 1 CPU para primeiro teste.
6. Verifique criação do disco, inicialização real, entrada no BIOS/ISO e desligamento; execute backup do disco antes de trabalho importante.

## Situação de testes
Compilação cruzada de referência e testes de sanidade podem ser realizados em Linux, mas **não substituem teste de execução no Windows 7 32 bits com QEMU real**.

Copyright © 2026 Nailson Rodrigues de Lima. Consulte Termos de Uso e Política de Privacidade incluídos.
