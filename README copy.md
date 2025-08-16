# API de Salas e Mensagens - Go + React Server

Esta é uma API RESTful construída em Go para gerenciamento de salas de chat em tempo real com funcionalidades de mensagens, reações e WebSockets.

## 📋 Índice

- [Características](#-características)
- [Tecnologias](#-tecnologias)
- [Configuração](#-configuração)
- [Execução](#-execução)
- [Autenticação](#-autenticação)
- [Endpoints da API](#-endpoints-da-api)
- [WebSocket](#-websocket)
- [Modelos de Dados](#-modelos-de-dados)
- [Exemplos de Uso](#-exemplos-de-uso)
- [Documentação Adicional](#-documentação-adicional)

## ✨ Características

- ✅ Criação e gerenciamento de salas de chat
- ✅ Sistema de mensagens em tempo real
- ✅ Sistema de reações para mensagens
- ✅ **Autenticação simples de host** (sem JWT)
- ✅ **Controle de permissões** para marcar mensagens como respondidas
- ✅ Marcar mensagens como respondidas
- ✅ WebSocket para atualizações em tempo real
- ✅ Banco de dados PostgreSQL
- ✅ CORS configurado para desenvolvimento
- ✅ Logging estruturado

## 🛠 Tecnologias

- **Go 1.21+**
- **Chi Router** - Roteamento HTTP
- **PostgreSQL** - Banco de dados
- **pgx/v5** - Driver PostgreSQL
- **Gorilla WebSocket** - WebSockets
- **UUID** - Identificadores únicos
- **godotenv** - Variáveis de ambiente
- **Zap** - Logger estruturado de alta performance
- **Lumberjack** - Rotação automática de logs

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto (baseado no `.env.example`):

```env
# Database Configuration
WSRS_DATABASE_USER=seu_usuario
WSRS_DATABASE_PASSWORD=sua_senha
WSRS_DATABASE_HOST=localhost
WSRS_DATABASE_PORT=5432
WSRS_DATABASE_NAME=nome_do_banco

# Logger Configuration (optional)
LOG_LEVEL=info  # debug, info, warn, error
```

### Banco de Dados

O projeto usa migrações SQL localizadas em `internal/store/pgstore/migrations/`:

1. `001_create_rooms_table.sql` - Cria a tabela de salas
2. `002_create_messages_table.sql` - Cria a tabela de mensagens

## 🚀 Execução

```bash
# Instalar dependências
go mod download

# Executar o servidor
go run cmd/wsrs/main.go
```

O servidor será iniciado na porta `:8080`.

## 🔐 Autenticação

A API implementa um **sistema de autenticação simples** para identificar hosts das salas:

### Como Funciona

1. **Criar Sala** → Gera token de host automaticamente
2. **Token no Header** → Use `X-Host-Token` para se identificar como host
3. **Permissões** → Apenas hosts podem marcar mensagens como respondidas
4. **Expiração** → Token válido por 24 horas

### Exemplo de Uso

```bash
# 1. Criar sala (recebe token de host)
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"theme": "Discussão sobre Go"}'

# Resposta:
# {
#   "id": "550e8400-e29b-41d4-a716-446655440000",
#   "host_token": "123e4567-e89b-12d3-a456-426614174001"
# }

# 2. Usar token para ações de host
curl -X PATCH http://localhost:8080/api/rooms/{room_id}/messages/{message_id}/answer \
  -H "X-Host-Token: 123e4567-e89b-12d3-a456-426614174001"
```

### Permissões

- ✅ **Qualquer pessoa**: Criar salas, enviar mensagens, reagir
- ✅ **Apenas host**: Marcar mensagens como respondidas

📖 **Documentação completa**: [docs/authentication.md](docs/authentication.md)

## 📊 Sistema de Logging

A aplicação utiliza um sistema de logging estruturado baseado no **Zap** para análise completa de logs.

### Características do Logger

- **Logging Estruturado**: Todos os logs são estruturados em JSON para facilitar análises
- **Níveis de Log**: Debug, Info, Warn, Error, Fatal
- **Contexto de Requisição**: Cada log contém informações da requisição HTTP
- **Rotação Automática**: Logs são automaticamente rotacionados e comprimidos
- **Saída Dupla**: Console (desenvolvimento) + Arquivo (produção)

### Localização dos Logs

Os logs são salvos em `internal/logger/logs/api.log` com rotação automática:
- Tamanho máximo por arquivo: 10MB
- Backup de 5 arquivos antigos
- Logs mantidos por 30 dias
- Compressão automática dos arquivos antigos

### Configuração do Nível de Log

Configure o nível através da variável de ambiente `LOG_LEVEL`:

```bash
export LOG_LEVEL=debug    # Para desenvolvimento
export LOG_LEVEL=info     # Para produção (padrão)
export LOG_LEVEL=warn     # Apenas warnings e erros
export LOG_LEVEL=error    # Apenas erros
```

### Campos de Contexto Incluídos

Cada log inclui automaticamente:
- `request_id` - ID único da requisição
- `client_ip` - IP do cliente
- `user_agent` - User agent do cliente
- `method` - Método HTTP (GET, POST, etc.)
- `path` - Caminho da requisição
- `query` - Query parameters
- `referer` - Referer da requisição
- `host` - Host da requisição
- `latency` - Tempo de resposta
- `status_code` - Código de status HTTP
- `user_id` - ID do usuário autenticado (quando aplicável)

### Exemplo de Log

```json
{
  "timestamp": "2024-08-15T10:30:45.123Z",
  "level": "INFO",
  "message": "message created successfully",
  "caller": "api/api.go:245",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_ip": "192.168.1.100",
  "method": "POST",
  "path": "/api/rooms/abc123/messages",
  "room_id": "abc123",
  "message_id": "def456"
}
```

## 📡 Endpoints da API

### Base URL

```text
http://localhost:8080/api
```

### 🏠 Salas (Rooms)

#### Criar uma nova sala

```http
POST /api/rooms
Content-Type: application/json

{
    "theme": "Tópico da sala"
}
```

**Resposta:**

```json
{
    "id": "uuid-da-sala"
}
```

#### Listar todas as salas

```http
GET /api/rooms
```

**Resposta:**

```json
[
    {
        "id": "uuid-da-sala",
        "theme": "Tópico da sala"
    }
]
```

#### Obter uma sala específica

```http
GET /api/rooms/{room_id}
```

**Resposta:**

```json
{
    "id": "uuid-da-sala",
    "theme": "Tópico da sala"
}
```

### 💬 Mensagens (Messages)

#### Criar uma nova mensagem

```http
POST /api/rooms/{room_id}/messages
Content-Type: application/json

{
    "message": "Conteúdo da mensagem"
}
```

**Resposta:**

```json
{
    "id": "uuid-da-mensagem"
}
```

#### Listar mensagens de uma sala

```http
GET /api/rooms/{room_id}/messages
```

**Resposta:**

```json
[
    {
        "id": "uuid-da-mensagem",
        "room_id": "uuid-da-sala",
        "message": "Conteúdo da mensagem",
        "reaction_count": 5,
        "answered": false
    }
]
```

#### Obter uma mensagem específica

```http
GET /api/rooms/{room_id}/messages/{message_id}
```

**Resposta:**

```json
{
    "id": "uuid-da-mensagem",
    "room_id": "uuid-da-sala",
    "message": "Conteúdo da mensagem",
    "reaction_count": 5,
    "answered": false
}
```

#### Reagir a uma mensagem (adicionar reação)

```http
PATCH /api/rooms/{room_id}/messages/{message_id}/react
```

**Resposta:**

```json
{
    "count": 6
}
```

#### Remover reação de uma mensagem

```http
DELETE /api/rooms/{room_id}/messages/{message_id}/react
```

**Resposta:**

```json
{
    "count": 5
}
```

#### Marcar mensagem como respondida

```http
PATCH /api/rooms/{room_id}/messages/{message_id}/answer
```

**Resposta:** Status 200 OK (sem corpo)

## 🔌 WebSocket

### Conectar ao WebSocket

```text
ws://localhost:8080/subscribe/{room_id}
```

### Eventos em Tempo Real

O WebSocket envia eventos JSON com a seguinte estrutura:

```json
{
    "kind": "tipo_do_evento",
    "value": { /* dados do evento */ }
}
```

#### Tipos de Eventos

1. **message_created** - Nova mensagem criada

```json
{
    "kind": "message_created",
    "value": {
        "id": "uuid-da-mensagem",
        "message": "Conteúdo da mensagem"
    }
}
```

2. **message_reaction_increased** - Reação adicionada

```json
{
    "kind": "message_reaction_increased",
    "value": {
        "id": "uuid-da-mensagem",
        "count": 6
    }
}
```

3. **message_reaction_decreased** - Reação removida

```json
{
    "kind": "message_reaction_decreased",
    "value": {
        "id": "uuid-da-mensagem",
        "count": 5
    }
}
```

4. **message_answered** - Mensagem marcada como respondida

```json
{
    "kind": "message_answered",
    "value": {
        "id": "uuid-da-mensagem"
    }
}
```

## 📊 Modelos de Dados

### Room (Sala)

```go
type Room struct {
    ID    uuid.UUID `json:"id"`
    Theme string    `json:"theme"`
}
```

### Message (Mensagem)

```go
type Message struct {
    ID            uuid.UUID `json:"id"`
    RoomID        uuid.UUID `json:"room_id"`
    Message       string    `json:"message"`
    ReactionCount int64     `json:"reaction_count"`
    Answered      bool      `json:"answered"`
}
```

## 🔍 Exemplos de Uso

### Exemplo com cURL

#### Criar uma sala

```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"theme": "Discussão sobre Go"}'
```

#### Criar uma mensagem

```bash
curl -X POST http://localhost:8080/api/rooms/{room_id}/messages \
  -H "Content-Type: application/json" \
  -d '{"message": "Olá, pessoal!"}'
```

#### Reagir a uma mensagem

```bash
curl -X PATCH http://localhost:8080/api/rooms/{room_id}/messages/{message_id}/react
```

### Exemplo com JavaScript (WebSocket)

```javascript
const roomId = 'uuid-da-sala';
const ws = new WebSocket(`ws://localhost:8080/subscribe/${roomId}`);

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    switch(data.kind) {
        case 'message_created':
            console.log('Nova mensagem:', data.value);
            break;
        case 'message_reaction_increased':
            console.log('Reação adicionada:', data.value);
            break;
        case 'message_reaction_decreased':
            console.log('Reação removida:', data.value);
            break;
        case 'message_answered':
            console.log('Mensagem respondida:', data.value);
            break;
    }
};

ws.onopen = function() {
    console.log('Conectado ao WebSocket');
};

ws.onclose = function() {
    console.log('Desconectado do WebSocket');
};
```

## 📝 Códigos de Status HTTP

- `200 OK` - Requisição bem-sucedida
- `400 Bad Request` - Dados inválidos (JSON malformado, UUID inválido, etc.)
- `404 Not Found` - Recurso não encontrado (sala ou mensagem inexistente)
- `500 Internal Server Error` - Erro interno do servidor

## 🔧 Desenvolvimento

### Estrutura do Projeto

```text
├── cmd/
│   └── wsrs/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── api/
│   │   ├── api.go              # Handlers da API
│   │   └── utils.go            # Utilitários
│   └── store/
│       └── pgstore/
│           ├── db.go           # Conexão com banco
│           ├── models.go       # Modelos de dados
│           ├── queries.sql.go  # Queries SQL geradas
│           └── migrations/     # Migrações do banco
├── docs/
│   └── Rooms.postman_collection.json  # Coleção Postman
└── README.md                   # Esta documentação
```

### Collection do Postman

Uma coleção completa do Postman está disponível em `docs/Rooms.postman_collection.json` com todos os endpoints configurados para teste.

## � Documentação Adicional

- **[Especificação OpenAPI/Swagger](docs/api.yaml)** - Documentação completa da API em formato OpenAPI 3.0
- **[Exemplos Práticos](docs/examples.md)** - Cenários reais de uso da API com scripts e códigos
- **[Coleção Postman](docs/Rooms.postman_collection.json)** - Collection para importar no Postman

## �📄 Licença

Este projeto foi desenvolvido como parte da Semana Tech da Rocketseat.
