# Libras Ao Vivo 5.1 — Correção da captação

## Mudanças
- Botão **Testar microfone**, com solicitação de permissão e liberação do dispositivo após o teste.
- Retomada controlada do reconhecimento, com limite de quatro tentativas e mensagens de erro orientativas.
- Feedback do início de áudio e de fala, e exibição de resultados parciais.
- Mantidos os recursos da versão 5.0 (biblioteca de vídeos próprios, histórico, exportação).

## Instalação e uso
1. Extraia o ZIP e publique o `index.html` em hospedagem HTTPS, como GitHub Pages. Abrir com `file://` não é recomendado para captação.
2. Abra a página no celular em navegador compatível e atualizado.
3. Toque em **Testar microfone**, autorize o acesso e confirme que o teste foi bem-sucedido.
4. Toque em **Iniciar microfone** e fale. Use **Parar** para encerrar.
5. Caso o recurso de reconhecimento não esteja disponível, ainda é possível digitar a frase manualmente. O teste de microfone **não garante** reconhecimento de fala, que depende de outro recurso do navegador e pode exigir conexão à internet.

## Limitações
- Web Speech API não é compatível com todos os navegadores e pode usar serviço remoto.
- Biblioteca audiovisual de sinais não é fornecida; registre apenas vídeos licenciados e linguisticamente validados.
- O aplicativo não faz tradução automática contextual de português para Libras.
- Teste de fala efetivo em um aparelho Android físico permanece necessário.

Copyright © 2026 Nailson Rodrigues de Lima.
