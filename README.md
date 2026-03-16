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
| `AUCTION_INTERVAL` | Duração do leilão em segundos | `20s` |
| `BATCH_INSERT_INTERVAL` | Intervalo de escrita em lote | `7m` |
| `MAX_BATCH_SIZE` | Limite de itens por batch | `10` |

---

## 🛠️ Como Executar

### 1. Requisitos

- Docker
- Docker Compose

### 2. Clonar o repositório

```bash
git clone -b sua-branch https://github.com/seu-repo/auction-api.git
cd auction-api
```

### 3. Subir a aplicação com Docker

O projeto utiliza **Docker Compose** para subir a API e o MongoDB automaticamente.

```bash
docker compose up --build
```

Após subir os containers:

- API disponível em: http://localhost:8080  
- MongoDB rodando na porta: 27017

### 4. Parar os containers

```bash
docker compose down
```

---

## Estrutura de Containers

O `docker-compose` sobe dois serviços:

- **app** → API escrita em Go  
- **mongodb** → Banco de dados MongoDB  

Os dados do MongoDB são persistidos em um volume Docker chamado **mongo-data**.

---

## Variáveis de Ambiente

As variáveis de ambiente são carregadas a partir do arquivo:

```
cmd/auction/.env
```

---

## Executar Testes

Os testes estão localizados em:

```
internal/infra/database/auction/create_auction_test.go
```

Para executar todos os testes do módulo:

```bash
go test ./internal/infra/database/auction/...
```
