# Ask Me Anything - API

Uma aplicação moderna de perguntas e respostas em tempo real construída em Go com PostgreSQL, WebSockets e autenticação por sessão.

## 🚀 Tecnologias

- **Backend**: Go 1.25+
- **Database**: PostgreSQL
- **Router**: Chi
- **WebSocket**: Gorilla WebSocket  
- **Database Query**: SQLC
- **Migrations**: golang-migrate
- **Logs**: slog (built-in)
- **Session Management**: Cookie-based com auto-renewal
- **Containerização**: Docker & Docker Compose

## ✨ Funcionalidades Implementadas

### 🔐 Sistema de Autenticação
- **Sessões por Cookie**: Autenticação automática via cookies HttpOnly
- **Auto-renovação**: Sessões renovadas automaticamente por 24 horas
- **Rastreamento de Usuários**: Sistema completo de sessões persistentes no banco

### 🏠 Gerenciamento de Salas
- **Criação de Salas**: Usuários podem criar salas de perguntas
- **Host Sessions**: Sistema de tokens de host para criadores de salas
- **Deleção Segura**: Apenas criadores podem deletar suas salas
- **Verificação de Ownership**: Validação de propriedade das salas

### 💬 Sistema de Mensagens
- **Envio de Perguntas**: Usuários podem enviar perguntas nas salas
- **Reações**: Sistema de "likes" nas mensagens
- **Rastreamento de Reações**: Usuários sabem quais mensagens já reagiram
- **Respostas do Host**: Hosts podem marcar mensagens como respondidas

### � WebSocket em Tempo Real
- **Mensagens em Tempo Real**: Novas mensagens aparecem instantaneamente
- **Contadores de Reações**: Atualização de likes em tempo real
- **Notificações de Sala Deletada**: Usuários são notificados quando sala é removida
- **Sincronização**: Todos os clientes conectados recebem atualizações

### 🎯 Recursos Avançados
- **Estado de Reações**: Frontend sabe quais mensagens o usuário reagiu
- **Deleção em Cascata**: Remover sala deleta mensagens e reações automaticamente
- **CORS Configurado**: Suporte para cookies entre domínios
- **Logging Estruturado**: Logs limpos e informativos

## 📁 Estrutura do Projeto

```
ask-me-anything/
├── cmd/                    # Aplicação principal
│   ├── wsrs/              # WebSocket Rooms Server
│   └── tools/             # Ferramentas auxiliares
├── internal/               # Código interno da aplicação
│   ├── api/               # Handlers HTTP e WebSocket
│   ├── auth/              # Sistema de autenticação e sessões
│   ├── logger/            # Sistema de logs estruturados
│   ├── middleware/        # Middlewares da aplicação
│   ├── responses/         # Helpers para respostas HTTP
│   ├── store/pgstore/     # Queries e models do PostgreSQL
│   └── validators/        # Validadores de entrada
├── migrations/            # Migrações do banco de dados
├── docs/                  # Documentação da API
├── docker-compose.yml     # Configuração Docker
├── Dockerfile            # Build da aplicação
├── Makefile              # Comandos de automação
└── sqlc.yaml             # Configuração do SQLC
```

## �️ Schema do Banco de Dados

### Tabelas Principais

- **`rooms`**: Salas de perguntas com tema
- **`messages`**: Mensagens/perguntas enviadas nas salas
- **`user_sessions`**: Sessões de usuários com cookies
- **`user_reactions`**: Reações dos usuários nas mensagens
- **`room_creators`**: Relacionamento entre usuários e salas criadas

### Relacionamentos

```sql
rooms (1) ←→ (N) messages
rooms (1) ←→ (1) room_creators ←→ (1) user_sessions
messages (1) ←→ (N) user_reactions ←→ (1) user_sessions
```

## 📊 API Endpoints

### 🏠 Salas (Rooms)
- `GET /api/rooms/` - Listar todas as salas
- `POST /api/rooms/` - Criar nova sala (retorna host_token)
- `GET /api/rooms/{room_id}/` - Obter detalhes da sala
- `DELETE /api/rooms/{room_id}/` - Deletar sala (apenas criador)
- `GET /api/rooms/{room_id}/host-status` - Verificar se é host da sala

### 💬 Mensagens
- `GET /api/rooms/{room_id}/messages/` - Listar mensagens da sala
- `POST /api/rooms/{room_id}/messages/` - Enviar nova mensagem
- `GET /api/rooms/{room_id}/messages/{message_id}/` - Obter mensagem específica
- `PATCH /api/rooms/{room_id}/messages/{message_id}/answer` - Marcar como respondida (host)
- `PATCH /api/rooms/{room_id}/messages/{message_id}/react` - Reagir à mensagem
- `DELETE /api/rooms/{room_id}/messages/{message_id}/react` - Remover reação

### 👤 Usuário
- `GET /api/user/rooms` - Listar salas criadas pelo usuário
- `DELETE /api/user/logout` - Fazer logout (invalidar sessão)

### 🔄 WebSocket
- `WS /subscribe/{room_id}` - Conexão WebSocket para atualizações em tempo real

### 🩺 Sistema
- `GET /health` - Health check da aplicação
- `GET /status` - Status detalhado do sistema

## 📡 WebSocket Events

### Eventos Enviados pelo Servidor

```json
// Nova mensagem criada
{
  "kind": "message_created",
  "room_id": "uuid",
  "value": {
    "id": "uuid", 
    "message": "texto da mensagem"
  }
}

// Reação adicionada
{
  "kind": "message_reaction_increased",
  "room_id": "uuid", 
  "value": {
    "id": "message_uuid",
    "count": 5
  }
}

// Reação removida
{
  "kind": "message_reaction_decreased", 
  "room_id": "uuid",
  "value": {
    "id": "message_uuid", 
    "count": 4
  }
}

// Mensagem marcada como respondida
{
  "kind": "message_answered",
  "room_id": "uuid",
  "value": {
    "id": "message_uuid"
  }
}

// Sala foi deletada
{
  "kind": "room_deleted",
  "room_id": "uuid", 
  "value": {
    "id": "room_uuid",
    "reason": "deleted_by_creator"
  }
}
```

## 🔐 Autenticação

### Sistema de Sessões

A autenticação é baseada em cookies HttpOnly que são automaticamente gerenciados:

```javascript
// Todas as requisições incluem cookies automaticamente
fetch('/api/rooms/', {
  credentials: 'include'  // Importante para cookies
})
```

### Headers Especiais

```javascript
// Para operações de host (verificar se é criador da sala)
fetch('/api/rooms/uuid/host-status', {
  headers: {
    'X-Host-Token': 'token-retornado-na-criacao-da-sala'
  }
})
```

## 🛠️ Configuração e Execução

### Pré-requisitos

- Docker & Docker Compose
- Make

### Execução com Docker (Recomendado)

**Iniciar tudo de uma vez:**
```bash
docker-compose up -d
```

**Comandos do Makefile:**
```bash
# Recarregar aplicação preservando dados
make docker-reload

# Parar containers
make docker-down

# Ver logs da aplicação
make docker-logs

# Limpar tudo e recomeçar
make docker-clean
```

### Desenvolvimento Local (Opcional)

Se quiser rodar localmente sem Docker:

```bash
# 1. Subir apenas o banco PostgreSQL
make docker-db

# 2. Aplicar migrações
make migrate-up

# 3. Executar aplicação localmente
make run
```

### Comandos Úteis

```bash
# Regenerar código SQLC após mudanças nas queries
make sqlc-generate

# Criar nova migração
make migrate-create name=nome_da_migracao

# Ver status das migrações
make migrate-status
```

## 🔗 URLs da Aplicação

Após iniciar com `docker-compose up -d`:

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **WebSocket**: ws://localhost:8080/subscribe/{room_id}

## 📋 Exemplos de Uso

### Criar Sala e Enviar Mensagem

```bash
# 1. Criar sala
curl -X POST "http://localhost:8080/api/rooms/" \
  -H "Content-Type: application/json" \
  -d '{"theme":"Minha sala de perguntas"}' \
  --cookie-jar cookies.txt

# Resposta: {"id":"uuid","host_token":"token"}

# 2. Enviar mensagem
curl -X POST "http://localhost:8080/api/rooms/{room_id}/messages/" \
  -H "Content-Type: application/json" \
  -d '{"message":"Minha pergunta"}' \
  --cookie cookies.txt

# 3. Reagir à mensagem  
curl -X PATCH "http://localhost:8080/api/rooms/{room_id}/messages/{msg_id}/react" \
  --cookie cookies.txt

# 4. Verificar se é host
curl -X GET "http://localhost:8080/api/rooms/{room_id}/host-status" \
  -H "X-Host-Token: {host_token}"
```

## 🧪 Testes

```bash
# Executar testes
make test

# Testes com cobertura
make test-coverage
```

## 📝 TODO Implementado

- ✅ **Sistema de Sessões**: Autenticação por cookies HttpOnly
- ✅ **WebSocket Real-time**: Mensagens e reações em tempo real  
- ✅ **Gerenciamento de Salas**: Criação, deleção e ownership
- ✅ **Sistema de Reações**: Like/unlike com rastreamento por usuário
- ✅ **CORS Configurado**: Suporte para frontend separado
- ✅ **Logs Estruturados**: Sistema de logging limpo e informativo
- ✅ **Deleção em Cascata**: Remoção segura de salas com dados relacionados
- ✅ **Docker Completo**: Ambiente totalmente containerizado

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)  
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
