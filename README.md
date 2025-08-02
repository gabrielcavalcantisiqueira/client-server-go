# Projeto Cotação Dólar - Go

Este projeto consiste em dois sistemas em Go que se comunicam via HTTP para consultar a cotação do dólar:

- `server.go`: servidor HTTP que consulta uma API externa, persiste a cotação no SQLite e retorna o valor ao cliente.
- `client.go`: cliente HTTP que solicita a cotação ao servidor, recebe o valor e salva em arquivo `cotacao.txt`.

---

## Funcionalidades

- O servidor consome a API pública de câmbio:  
  `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- O servidor usa `context` para controlar timeout na chamada da API (200ms) e na gravação no banco SQLite (10ms).
- O servidor persiste cada cotação recebida no banco SQLite local (`cotacoes.db`).
- O servidor expõe o endpoint `/cotacao` na porta `8080`.
- O cliente realiza requisição HTTP ao servidor com timeout de 300ms.
- O cliente recebe o valor do câmbio (`bid`) e salva no arquivo `cotacao.txt` com o formato:  
  `Dólar: {valor}`.
- Logs registram erros de timeout em todas as etapas.

## Como rodar localmente

### Server

Para facilitar o ambiente, o servidor pode ser executado via Docker Compose dentro da pasta `server`.

Passos:

1. Navegue até a pasta `server`:

```
cd server
```

2. Execute o comando para buildar a imagem e subir o container em background:

```
docker-compose up -d --build
```

### Client

```
go run ./client/.
```




