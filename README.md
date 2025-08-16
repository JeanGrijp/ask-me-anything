# Ask Me Anything - API

Uma aplicação de perguntas e respostas construída em Go com PostgreSQL, WebSockets e SQLC.

## 🚀 Tecnologias

- **Backend**: Go 1.25+
- **Database**: PostgreSQL
- **Router**: Chi
- **WebSocket**: Gorilla WebSocket
- **Database Query**: SQLC
- **Migrations**: golang-migrate
- **Logs**: zap (uber-go)
- **Containerização**: Docker & Docker Compose

## 🔄 Mudanças Recentes

### Schema Alignment (Migration 002)

**Principais mudanças implementadas:**

1. **Simplificação do Schema**: Removidas tabelas complexas (`categories`, `votes`) para focar nas entidades principais
2. **Alinhamento com Models**: Schema do banco agora reflete exatamente os structs Go
3. **Queries Atualizadas**: Todas as queries SQLC foram atualizadas para usar os novos campos
4. **Migração Criada**: Nova migração `002_align_with_models` pronta para aplicação

**Estrutura Final:**
- Users com roles simplificados (admin, user, guest)
- Questions com like_count integrado
- Answers vinculadas a questions
- Rooms com ownership por user
- Magic links para autenticação sem senha

**Próximos Passos:**
```bash
# Aplicar a migração
make migrate-up

# Gerar código SQLC atualizado
make sqlc-generate

# Testar a aplicação
make run
```

## 📝 TODO

- [x] ~~Implementar autenticação JWT~~
- [x] ~~Adicionar middleware de CORS~~
- [x] ~~Implementar rate limiting~~
- [x] ~~Implementar autenticação Magic Link~~
- [x] ~~Separar Routes de Handlers~~
- [x] ~~Alinhar Schema com Models Go~~a WebSocket
- **Database Query**: SQLC
- **Migrations**: golang-migrate
- **Logs**: slog (built-in)
- **Containerização**: Docker & Docker Compose

## 📁 Estrutura do Projeto

```
ask-me-anything/
├── cmd/                    # Aplicação principal
│   └── main.go
├── internal/               # Código interno da aplicação
│   ├── config/            # Configurações
│   ├── database/          # Conexão com o banco
│   ├── handlers/          # Handlers HTTP e WebSocket
│   ├── models/            # Modelos de dados
│   └── services/          # Lógica de negócio
├── migrations/            # Migrações do banco de dados
├── queries/               # Queries SQL para SQLC
├── docker-compose.yml     # Configuração Docker
├── Dockerfile            # Build da aplicação
├── Makefile              # Comandos de automação
└── sqlc.yaml             # Configuração do SQLC
```

## 🛠️ Configuração

### Pré-requisitos

- Go 1.24+
- Docker & Docker Compose
- Make

### Instalação das Ferramentas

```bash
make install-tools
```

### Variáveis de Ambiente

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Edite o arquivo `.env` conforme necessário.

## 🏃‍♂️ Executando o Projeto

### Desenvolvimento Local

1. **Subir o banco de dados:**
```bash
make docker-up
```

2. **Executar migrações:**
```bash
make migrate-up
```

3. **Gerar código SQLC:**
```bash
make sqlc-generate
```

4. **Executar a aplicação:**
```bash
make run
```

### Comando Único para Desenvolvimento

```bash
make dev
```

## 🐳 Docker

### Executar com Docker Compose

```bash
docker-compose up -d
```

### Build da aplicação

```bash
make build
```

## 📊 API Endpoints

### Health Check
- `GET /health` - Status da aplicação

### Usuários
- `GET /api/v1/users` - Listar usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users/{id}` - Obter usuário
- `PUT /api/v1/users/{id}` - Atualizar usuário
- `DELETE /api/v1/users/{id}` - Deletar usuário

### Perguntas
- `GET /api/v1/questions` - Listar perguntas
- `POST /api/v1/questions` - Criar pergunta
- `GET /api/v1/questions/{id}` - Obter pergunta
- `PUT /api/v1/questions/{id}` - Atualizar pergunta
- `DELETE /api/v1/questions/{id}` - Deletar pergunta

### Respostas
- `GET /api/v1/answers` - Listar respostas
- `POST /api/v1/answers` - Criar resposta
- `GET /api/v1/answers/{id}` - Obter resposta
- `PUT /api/v1/answers/{id}` - Atualizar resposta
- `DELETE /api/v1/answers/{id}` - Deletar resposta

### WebSocket
- `WS /ws` - Conexão WebSocket para atualizações em tempo real

## 🗄️ Banco de Dados

### Entidades Principais (Models)

- **Users**: Usuários da plataforma com roles (admin, user, guest)
- **Questions**: Perguntas com sistema de likes
- **Answers**: Respostas às perguntas
- **Rooms**: Salas de perguntas gerenciadas por proprietários
- **Magic Links**: Sistema de autenticação sem senha

### Schema Atualizado

O banco foi alinhado com os models Go definidos em `internal/models/`:

**User Model:**
```go
type User struct {
    ID        int64    `json:"id" db:"id"`
    Email     string   `json:"email" db:"email"`
    Name      string   `json:"name" db:"name"`
    Role      UserRole `json:"role" db:"role"`
    CreatedAt string   `json:"created_at" db:"created_at"`
}
```

**Question Model:**
```go
type Question struct {
    ID        int64  `json:"id" db:"id"`
    Content   string `json:"content" db:"content"`
    UserID    int64  `json:"user_id" db:"user_id"`
    LikeCount int64  `json:"like_count" db:"like_count"`
}
```

**Answer Model:**
```go
type Answer struct {
    ID       int64    `json:"id" db:"id"`
    Question Question `json:"question" db:"question"`
    Answer   string   `json:"answer" db:"answer"`
    UserID   int64    `json:"user_id" db:"user_id"`
}
```

**Room Model:**
```go
type Room struct {
    ID      int64  `json:"id" db:"id"`
    Name    string `json:"name" db:"name"`
    OwnerID int64  `json:"owner_id" db:"owner_id"`
}
```

### Migrações

Criar nova migração:
```bash
make migrate-create name=nome_da_migracao
```

Aplicar migrações:
```bash
make migrate-up
```

Reverter migrações:
```bash
make migrate-down
```

## 🧪 Testes

```bash
make test
```

## 📝 Comandos Makefile

- `make build` - Compilar a aplicação
- `make run` - Executar a aplicação
- `make test` - Executar testes
- `make clean` - Limpar binários
- `make docker-up` - Subir containers
- `make docker-down` - Parar containers
- `make migrate-up` - Aplicar migrações
- `make migrate-down` - Reverter migrações
- `make sqlc-generate` - Gerar código SQLC
- `make dev` - Ambiente de desenvolvimento completo
- `make tidy` - Organizar dependências Go

## 📋 TODO

- [ ] Implementar autenticação JWT
- [ ] Adicionar middleware de CORS
- [ ] Implementar rate limiting
- [ ] Adicionar cache Redis
- [ ] Implementar notificações WebSocket
- [ ] Adicionar testes unitários
- [ ] Documentação Swagger/OpenAPI
- [ ] Métricas e monitoramento

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
