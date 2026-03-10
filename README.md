# Lab Auction API

API desenvolvida em **Golang** para gerenciamento de leilões, utilizando **Clean Architecture** e **MongoDB**. O sistema gerencia o ciclo de vida de leilões, desde a criação até o fechamento automático baseado em timers.

---

## 🏗️ Arquitetura

O projeto segue os princípios da **Clean Architecture**, garantindo desacoplamento e facilidade de testes:
- **Entities:** Regras de negócio centrais.
- **Use Cases:** Orquestração da lógica da aplicação.
- **Controllers:** Adaptadores de entrada (HTTP).
- **Infra:** Implementações de banco de dados (MongoDB) e configurações externas.

---

## 🕒 Fechamento Automático (Timers)

Diferente de sistemas que dependem de CRON jobs, esta API utiliza **timers em memória** para cada leilão:
- Ao criar um leilão, um timer é iniciado.
- A duração é definida pela variável `AUCTION_END_TIME`.
- Quando o tempo expira, o status é alterado para `Completed` de forma assíncrona.

> **Localização da lógica:** `cmd/internal/infra/database/auction/create_auction.go`

---

## 🚀 Endpoints Principais

### Leilões
- **POST `/auction`**: Cria um novo leilão.
- **GET `/auction/{id}`**: Detalhes de um leilão específico.
- **GET `/auctions`**: Lista todos os leilões cadastrados.
- **GET `/auction/winner/:auctionId`**: Retorna o maior lance e o vencedor do leilão.

---

## ⚙️ Variáveis de Ambiente (.env)

| Variável | Descrição | Valor Padrão |
| :--- | :--- | :--- |
| `MONGO_URL` | String de conexão MongoDB | `mongodb://localhost:27017` |
| `MONGO_DB` | Nome do banco de dados | `auctions` |
| `AUCTION_END_TIME` | Duração do leilão em segundos | `60` |
| `BATCH_INSERT_INTERVAL` | Intervalo de escrita em lote | `7m` |
| `MAX_BATCH_SIZE` | Limite de itens por batch | `10` |

---

## 🛠️ Como Executar

### 1. Requisitos
- Go 1.25+
- MongoDB rodando na porta 27017

### 2. Instalação e Execução
```bash
# Clone o repositório
git clone -b sua-branch [https://github.com/seu-repo/auction-api.git](https://github.com/seu-repo/auction-api.git)
cd auction-api

# Instale as dependências
go mod tidy

# Execute a aplicação
go run ./cmd/auction/main.go


Testes localizados em cmd/internal/infra/database/auction/create_auction_test.go
:

go test ./cmd/internal/infra/database/auction/...