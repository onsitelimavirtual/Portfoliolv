# POLÍTICA DE PRIVACIDADE — NatorVM 10.0 Native
Versão: 10.0 | Atualização: 09/10/2026
Responsável pelo conteúdo original: Nailson Rodrigues de Lima.

## 1. Escopo
Esta política descreve o funcionamento do gerenciador NatorVM 10.0 Native entregue neste pacote, em modo local. Ela não abrange diretamente sistemas convidados nem eventuais serviços e ferramentas de terceiros que o usuário instale.

## 2. Dados tratados
O programa mantém localmente um catálogo de VMs (nomes, identificadores, memória, CPUs, capacidade de disco, caminhos das ISOs e dos discos), arquivos de disco virtual e informações necessárias à operação. O painel utiliza cookie temporário de sessão para acesso local e token antifalsificação (CSRF). O terminal exibe uma URL temporária de acesso; não a compartilhe.

## 3. Finalidades e armazenamento
Os dados são utilizados para listar, configurar e executar as VMs. Em configuração padrão, ficam no diretório de configuração do usuário, na pasta NatorVM10, ou no caminho definido pela variável NATORVM_HOME. Discos virtuais podem conter dados pessoais inseridos pelo próprio usuário nos sistemas convidados.

## 4. Envio, terceiros e acesso
O código entregue não implementa telemetria, conta online, analytics, sincronização em nuvem ou transmissão intencional de dados a um servidor do desenvolvedor. A conexão do painel fica restrita ao computador local (127.0.0.1). Sistemas convidados e QEMU, quando instalados, podem usar rede própria conforme sua configuração: avalie separadamente esse tráfego.

## 5. Retenção e exclusão
Os arquivos locais permanecem enquanto o usuário os conservar. A exclusão de uma VM pelo aplicativo pode remover seu disco; portanto, faça backup antes. A desinstalação do executável não garante a exclusão automática dos dados no diretório NatorVM10. Remoção segura exige verificar cópias e backups.

## 6. Segurança
A aplicação utiliza sessão temporária, cookie HttpOnly e SameSite=Strict, proteção CSRF e acesso vinculado a loopback. Essas medidas reduzem riscos, mas não equivalem a auditoria completa de segurança. Proteja o computador, os discos virtuais e as credenciais de acesso.

## 7. Direitos e LGPD
Quando forem tratados dados pessoais, devem ser observadas as obrigações da Lei nº 13.709/2018 (LGPD), conforme o contexto e o agente de tratamento responsável. O operador local deve gerenciar acesso, conservação, compartilhamento e eventuais solicitações dos titulares relativos aos dados sob seu controle. Antes de comercializar ou operar institucionalmente, identifique os agentes de tratamento e os canais de contato aplicáveis.

## 8. Revisões
Esta política poderá ser atualizada conforme o software evoluir. Este texto é um modelo inicial sujeito a adequação jurídica, inclusive definição de contato e informações do controlador em eventual distribuição pública.
