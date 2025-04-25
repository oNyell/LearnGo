# Pingador de Sites

Um aplicativo web simples para monitorar o status de sites em tempo real usando WebSocket.

## 🚀 Funcionalidades

- Interface web para adicionar URLs para monitoramento
- Verificação automática do status dos sites a cada 10 segundos
- Atualização em tempo real do status usando WebSocket
- Interface simples e intuitiva

## 📋 Pré-requisitos

- Go 1.24.2 ou superior
- Dependência: github.com/gorilla/websocket v1.5.3

## 🔧 Instalação

1. Clone o repositório:
```bash
git clone https://github.com/oNyell/LearnGo.git
cd LearnGo
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute o servidor:
```bash
go run main.go
```

4. Acesse a aplicação em seu navegador:
```
http://localhost:8080
```

## 🛠️ Como usar

1. Acesse a página inicial
2. Digite a URL do site que deseja monitorar no campo de texto
3. Clique em "Adicionar"
4. O status do site será exibido na lista abaixo
5. O status será atualizado automaticamente a cada 10 segundos

## 📝 Status possíveis

- OK: Site está respondendo corretamente
- FORA DO AR: Site não está respondendo ou retornou erro
- Aguardando verificação...: Status inicial após adicionar um novo site

## 🤝 Contribuindo

Contribuições são sempre bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes. 