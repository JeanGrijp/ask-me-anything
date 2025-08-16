# Sistema de Autenticação Simples - Guia de Uso

## 📋 Como Funciona

O sistema de autenticação é **simples e sem complexidade**:

1. **Criação de Sala** → Gera token de host
2. **Token no Header** → `X-Host-Token` para identificar host
3. **Controle de Permissões** → Apenas host pode marcar mensagens como respondidas
4. **Expiração** → Token válido por 24 horas

## 🚀 Fluxo de Uso

### 1. Criar uma Sala (Gera Token de Host)

```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"theme": "Discussão sobre Go"}'
```

**Resposta:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "host_token": "123e4567-e89b-12d3-a456-426614174001"
}
```

### 2. Verificar Status de Host

```bash
curl -X GET http://localhost:8080/api/rooms/550e8400-e29b-41d4-a716-446655440000/host-status \
  -H "X-Host-Token: 123e4567-e89b-12d3-a456-426614174001"
```

**Resposta:**
```json
{
  "is_host": true,
  "room_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 3. Marcar Mensagem como Respondida (Apenas Host)

```bash
curl -X PATCH http://localhost:8080/api/rooms/550e8400-e29b-41d4-a716-446655440000/messages/abc-123/answer \
  -H "X-Host-Token: 123e4567-e89b-12d3-a456-426614174001"
```

**Resposta de Sucesso:**
```
Status: 200 OK
```

**Resposta sem Token:**
```json
Status: 401 Unauthorized
{
  "error": "host token required"
}
```

**Resposta com Token Inválido:**
```json
Status: 403 Forbidden
{
  "error": "only room host can perform this action"
}
```

### 4. Outras Ações (Sem Restrição)

Qualquer pessoa pode:

```bash
# Enviar mensagem
curl -X POST http://localhost:8080/api/rooms/550e8400-e29b-41d4-a716-446655440000/messages \
  -H "Content-Type: application/json" \
  -d '{"message": "Olá pessoal!"}'

# Reagir a mensagem
curl -X PATCH http://localhost:8080/api/rooms/550e8400-e29b-41d4-a716-446655440000/messages/abc-123/react

# Ver mensagens
curl -X GET http://localhost:8080/api/rooms/550e8400-e29b-41d4-a716-446655440000/messages
```

## 🔐 Detalhes Técnicos

### **Como o Token é Gerado:**
- UUID aleatório (formato: `123e4567-e89b-12d3-a456-426614174001`)
- Único por sala
- Gerado no momento da criação da sala

### **Onde o Token é Armazenado:**
- **Servidor**: Memória (map thread-safe)
- **Cliente**: Você deve salvar (localStorage, cookie, etc.)

### **Expiração:**
- **24 horas** após criação
- Limpeza automática a cada 30 minutos

### **Segurança:**
- ✅ Token único por sala
- ✅ Validação de UUID
- ✅ Verificação de expiração
- ✅ Logs de ações de host
- ✅ Thread-safe

## 🎯 Permissões

| Ação | Host | Usuário Comum |
|------|------|---------------|
| Criar sala | ✅ | ✅ |
| Ver salas | ✅ | ✅ |
| Enviar mensagem | ✅ | ✅ |
| Reagir à mensagem | ✅ | ✅ |
| **Marcar como respondida** | ✅ | ❌ |
| WebSocket (tempo real) | ✅ | ✅ |

## 🔄 Cenários Comuns

### **F5 na Página:**
1. ✅ Token permanece válido
2. ✅ Reconecta WebSocket
3. ✅ Mantém status de host

### **Token Expirado:**
1. ❌ Host perde privilégios
2. ❌ Não consegue marcar mensagens
3. ✅ Ainda pode usar como usuário comum

### **Múltiplos Hosts:**
1. ❌ Apenas 1 token por sala
2. ❌ Token único para o criador
3. ✅ Sistema simples e direto

## 📱 Implementação no Frontend

```javascript
// Salvar token ao criar sala
const createRoom = async (theme) => {
  const response = await fetch('/api/rooms', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ theme })
  });
  
  const { id, host_token } = await response.json();
  
  // Salvar token para usar depois
  localStorage.setItem(`host_token_${id}`, host_token);
  
  return { id, host_token };
};

// Usar token em ações de host
const markAsAnswered = async (roomId, messageId) => {
  const token = localStorage.getItem(`host_token_${roomId}`);
  
  if (!token) {
    alert('Você não é o host desta sala');
    return;
  }
  
  const response = await fetch(`/api/rooms/${roomId}/messages/${messageId}/answer`, {
    method: 'PATCH',
    headers: { 'X-Host-Token': token }
  });
  
  if (response.status === 403) {
    alert('Apenas o host pode marcar mensagens como respondidas');
  }
};

// Verificar se é host
const checkHostStatus = async (roomId) => {
  const token = localStorage.getItem(`host_token_${roomId}`);
  
  if (!token) return false;
  
  const response = await fetch(`/api/rooms/${roomId}/host-status`, {
    headers: { 'X-Host-Token': token }
  });
  
  const { is_host } = await response.json();
  return is_host;
};
```

## ✅ Sistema Simples e Eficaz!

- **Sem JWT** - Tokens UUID simples
- **Sem usuários** - Identificação apenas por token
- **Sem senhas** - Token gerado automaticamente
- **Sem banco** - Armazenamento em memória
- **Funcional** - Controla apenas o essencial

Perfeito para um sistema de chat simples! 🎉
