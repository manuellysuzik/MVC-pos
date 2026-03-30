# E-commerce API

API REST em Go que gerencia Clientes, Produtos e Pedidos.
Construída com **GIN** (HTTP), **GORM** (ORM) e **SQLite** (banco embarcado),
seguindo **Clean Architecture** e princípios **SOLID**.

## Pré-requisitos

- Go 1.21+
- GCC (necessário para o driver SQLite com CGO)
  - macOS: `xcode-select --install`
  - Ubuntu: `sudo apt install gcc`

## Configuração

```bash
cp .env.example .env
# Edite .env e defina API_KEY
```

Conteúdo do `.env`:
```
API_KEY=minha-chave-secreta
DB_PATH=./ecommerce.db
PORT=8080
```

## Como rodar

```bash
go run main.go
```

O servidor sobe em `http://localhost:8080`.
O banco `ecommerce.db` é criado automaticamente na primeira execução.

## Como testar

```bash
go test ./... -v
```

## Autenticação

Todas as rotas exigem o header:
```
X-API-Key: <valor configurado em API_KEY>
```

## Endpoints

### Clientes

| Método | Rota | Descrição |
|---|---|---|
| POST | `/api/v1/clientes` | Criar cliente |
| GET | `/api/v1/clientes` | Listar todos |
| GET | `/api/v1/clientes/count` | Contar total |
| GET | `/api/v1/clientes/search?nome=` | Buscar por nome |
| GET | `/api/v1/clientes/:id` | Buscar por ID |
| PUT | `/api/v1/clientes/:id` | Atualizar |
| DELETE | `/api/v1/clientes/:id` | Deletar |

### Produtos

| Método | Rota | Descrição |
|---|---|---|
| POST | `/api/v1/produtos` | Criar produto |
| GET | `/api/v1/produtos` | Listar todos |
| GET | `/api/v1/produtos/count` | Contar total |
| GET | `/api/v1/produtos/search?nome=` | Buscar por nome |
| GET | `/api/v1/produtos/:id` | Buscar por ID |
| PUT | `/api/v1/produtos/:id` | Atualizar |
| DELETE | `/api/v1/produtos/:id` | Deletar |

### Pedidos

| Método | Rota | Descrição |
|---|---|---|
| POST | `/api/v1/pedidos` | Criar pedido (com itens) |
| GET | `/api/v1/pedidos` | Listar todos |
| GET | `/api/v1/pedidos/count` | Contar total |
| GET | `/api/v1/pedidos/search?status=` | Buscar por status |
| GET | `/api/v1/pedidos/cliente/:id` | Pedidos de um cliente |
| GET | `/api/v1/pedidos/:id` | Buscar por ID |
| PUT | `/api/v1/pedidos/:id` | Atualizar |
| DELETE | `/api/v1/pedidos/:id` | Deletar |

## Exemplos curl

```bash
# Criar cliente
curl -s -X POST http://localhost:8080/api/v1/clientes \
  -H "X-API-Key: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{"nome":"João Silva","email":"joao@email.com","telefone":"11999999999"}' | jq .

# Criar produto
curl -s -X POST http://localhost:8080/api/v1/produtos \
  -H "X-API-Key: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{"nome":"Notebook","descricao":"Notebook i7","preco":3500.00,"estoque":10}' | jq .

# Criar pedido
curl -s -X POST http://localhost:8080/api/v1/pedidos \
  -H "X-API-Key: minha-chave-secreta" \
  -H "Content-Type: application/json" \
  -d '{"cliente_id":1,"itens":[{"produto_id":1,"quantidade":2,"preco_unitario":3500.00}]}' | jq .

# Listar clientes
curl -s http://localhost:8080/api/v1/clientes \
  -H "X-API-Key: minha-chave-secreta" | jq .

# Buscar por nome
curl -s "http://localhost:8080/api/v1/clientes/search?nome=João" \
  -H "X-API-Key: minha-chave-secreta" | jq .

# Contar produtos
curl -s http://localhost:8080/api/v1/produtos/count \
  -H "X-API-Key: minha-chave-secreta" | jq .

# Pedidos de um cliente
curl -s http://localhost:8080/api/v1/pedidos/cliente/1 \
  -H "X-API-Key: minha-chave-secreta" | jq .

# Sem API Key — deve retornar 401
curl -s http://localhost:8080/api/v1/clientes | jq .
```

## Estrutura de Pastas

```
.
├── main.go                        # Entry point — DI manual, wiring de todas as camadas
├── Makefile                       # Atalhos: make run, make test, make build, make tidy
├── .env                           # Configurações locais (não commitado)
├── .env.example                   # Template de configuração
├── internal/
│   ├── domain/
│   │   └── entity/                # Structs de domínio (Cliente, Produto, Pedido, ItemPedido)
│   ├── repository/
│   │   ├── cliente/               # package clienterepo — implementação GORM
│   │   ├── produto/               # package produtorepo — implementação GORM
│   │   └── pedido/                # package pedidorepo — implementação GORM
│   ├── service/
│   │   ├── cliente/               # package clientesvc — regras de negócio + interface do repo
│   │   ├── produto/               # package produtosvc — regras de negócio + interface do repo
│   │   └── pedido/                # package pedidosvc — regras de negócio + interface do repo
│   └── handler/
│       ├── router.go              # SetupRouter — registra todas as rotas com middleware
│       ├── cliente/               # package clienthdl — endpoints HTTP de clientes
│       ├── produto/               # package produtohdl — endpoints HTTP de produtos
│       └── pedido/                # package pedidohdl — endpoints HTTP de pedidos
├── pkg/
│   └── middleware/                # Middleware X-API-Key
└── docs/                          # Diagramas C4 (Mermaid + PlantUML) e specs
```

> **Padrão de interfaces:** cada pacote de service declara a interface do repository que consome (`interfaces.go`), seguindo o princípio Go de "aceite interfaces, retorne structs". Os handlers declaram sua própria interface de service inline. Não há pacote centralizado de interfaces.

> **Construtores:** todos os pacotes expõem `New(...)` como construtor. No `main.go`, os aliases de import resolvem a ambiguidade (`clienterepo.New(db)`, `clientesvc.New(&repo)`, `clienthdl.New(svc)`).

## Arquitetura

Ver [docs/architecture.md](docs/architecture.md) para diagramas C4.

**Princípios SOLID aplicados:**
- **S** — cada arquivo tem uma única responsabilidade
- **O** — nova fonte de dados = nova impl da interface, sem alterar service
- **L** — qualquer impl de `ClienteRepository` é intercambiável
- **I** — interfaces pequenas por entidade
- **D** — handlers/services dependem de interfaces, não de implementações concretas
