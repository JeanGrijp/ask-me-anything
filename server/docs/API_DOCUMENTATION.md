# 📚 Ask Me Anything - Documentação da API

## 🎯 Visão Geral

A **Ask Me Anything API** é uma API REST construída em Go que permite criar salas de perguntas e respostas em tempo real. A API suporta criação de salas, envio de mensagens, sistema de reações e notificações via WebSocket.

## 🚀 Como Conectar seu Frontend

### 📋 Configuração Base

```javascript
// Configuração base da API
const API_BASE_URL = 'http://localhost:8080';
const WS_BASE_URL = 'ws://localhost:8080';

// Headers padrão para requisições
const defaultHeaders = {
  'Content-Type': 'application/json',
  'Accept': 'application/json',
};

// Para requisições que precisam de autenticação de host
const hostHeaders = (hostToken) => ({
  ...defaultHeaders,
  'X-Host-Token': hostToken,
});
```

### 🔧 Utilitário de Requisições

```javascript
// Função helper para fazer requisições
async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const config = {
    headers: defaultHeaders,
    ...options,
  };
  
  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    // Verificar se há conteúdo para parsear
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }
    
    return response;
  } catch (error) {
    console.error('API Request Error:', error);
    throw error;
  }
}
```

## 📍 Endpoints da API

### 🏠 **Salas (Rooms)**

#### **GET /api/rooms**
Lista todas as salas disponíveis.

```javascript
// Listar todas as salas
const getRooms = async () => {
  return await apiRequest('/api/rooms');
};

// Exemplo de uso
const rooms = await getRooms();
console.log(rooms);
// Resposta: [{"id": "uuid", "theme": "Tema da sala"}]
```

**Resposta:**
```json
[
  {
    "id": "ec99fdaa-92cf-4b85-883f-8599fd9d4df1",
    "theme": "Discussão sobre tecnologia"
  }
]
```

---

#### **POST /api/rooms**
Cria uma nova sala.

```javascript
// Criar uma nova sala
const createRoom = async (theme) => {
  return await apiRequest('/api/rooms', {
    method: 'POST',
    body: JSON.stringify({ theme }),
  });
};

// Exemplo de uso
const newRoom = await createRoom('Minha nova sala');
console.log(newRoom);
// Resposta: {"id": "uuid", "host_token": "token"}
```

**Body:**
```json
{
  "theme": "Tema da sua sala"
}
```

**Resposta:**
```json
{
  "id": "ec99fdaa-92cf-4b85-883f-8599fd9d4df1",
  "host_token": "002b39c3-5e70-4f94-8223-2f16341df7df"
}
```

---

#### **GET /api/rooms/{room_id}**
Obtém detalhes de uma sala específica.

```javascript
// Obter detalhes de uma sala
const getRoom = async (roomId) => {
  return await apiRequest(`/api/rooms/${roomId}`);
};

// Exemplo de uso
const room = await getRoom('ec99fdaa-92cf-4b85-883f-8599fd9d4df1');
```

**Resposta:**
```json
{
  "id": "ec99fdaa-92cf-4b85-883f-8599fd9d4df1",
  "theme": "Discussão sobre tecnologia"
}
```

---

#### **GET /api/rooms/{room_id}/host-status**
Verifica se o usuário atual é o host da sala.

```javascript
// Verificar status de host
const getHostStatus = async (roomId, hostToken = null) => {
  const headers = hostToken ? hostHeaders(hostToken) : defaultHeaders;
  
  return await apiRequest(`/api/rooms/${roomId}/host-status`, {
    headers,
  });
};

// Exemplo de uso
const status = await getHostStatus('room-id', 'host-token');
console.log(status);
// Resposta: {"is_host": true, "room_id": "room-id"}
```

**Resposta:**
```json
{
  "is_host": true,
  "room_id": "ec99fdaa-92cf-4b85-883f-8599fd9d4df1"
}
```

---

### 💬 **Mensagens (Messages)**

#### **GET /api/rooms/{room_id}/messages**
Lista todas as mensagens de uma sala.

```javascript
// Listar mensagens de uma sala
const getRoomMessages = async (roomId) => {
  return await apiRequest(`/api/rooms/${roomId}/messages`);
};

// Exemplo de uso
const messages = await getRoomMessages('ec99fdaa-92cf-4b85-883f-8599fd9d4df1');
```

**Resposta:**
```json
[
  {
    "id": "b01f60db-9b7d-4081-b339-947a23909505",
    "room_id": "ec99fdaa-92cf-4b85-883f-8599fd9d4df1",
    "message": "Qual é a sua linguagem favorita?",
    "reaction_count": 5,
    "answered": false
  }
]
```

---

#### **POST /api/rooms/{room_id}/messages**
Cria uma nova mensagem em uma sala.

```javascript
// Criar uma nova mensagem
const createMessage = async (roomId, message) => {
  return await apiRequest(`/api/rooms/${roomId}/messages`, {
    method: 'POST',
    body: JSON.stringify({ message }),
  });
};

// Exemplo de uso
const newMessage = await createMessage(
  'ec99fdaa-92cf-4b85-883f-8599fd9d4df1', 
  'Como você aprendeu programação?'
);
```

**Body:**
```json
{
  "message": "Sua pergunta aqui"
}
```

**Resposta:**
```json
{
  "id": "b01f60db-9b7d-4081-b339-947a23909505"
}
```

---

#### **GET /api/rooms/{room_id}/messages/{message_id}**
Obtém uma mensagem específica.

```javascript
// Obter uma mensagem específica
const getMessage = async (roomId, messageId) => {
  return await apiRequest(`/api/rooms/${roomId}/messages/${messageId}`);
};
```

---

#### **PATCH /api/rooms/{room_id}/messages/{message_id}/react**
Adiciona uma reação a uma mensagem.

```javascript
// Adicionar reação a uma mensagem
const reactToMessage = async (roomId, messageId) => {
  return await apiRequest(`/api/rooms/${roomId}/messages/${messageId}/react`, {
    method: 'PATCH',
  });
};

// Exemplo de uso
const reaction = await reactToMessage(
  'ec99fdaa-92cf-4b85-883f-8599fd9d4df1',
  'b01f60db-9b7d-4081-b339-947a23909505'
);
console.log(reaction);
// Resposta: {"count": 6}
```

**Resposta:**
```json
{
  "count": 6
}
```

---

#### **DELETE /api/rooms/{room_id}/messages/{message_id}/react**
Remove uma reação de uma mensagem.

```javascript
// Remover reação de uma mensagem
const removeReactionFromMessage = async (roomId, messageId) => {
  return await apiRequest(`/api/rooms/${roomId}/messages/${messageId}/react`, {
    method: 'DELETE',
  });
};
```

---

#### **PATCH /api/rooms/{room_id}/messages/{message_id}/answer** 🔐
Marca uma mensagem como respondida (apenas hosts).

```javascript
// Marcar mensagem como respondida (requer host token)
const markMessageAsAnswered = async (roomId, messageId, hostToken) => {
  return await apiRequest(`/api/rooms/${roomId}/messages/${messageId}/answer`, {
    method: 'PATCH',
    headers: hostHeaders(hostToken),
  });
};

// Exemplo de uso
await markMessageAsAnswered(
  'ec99fdaa-92cf-4b85-883f-8599fd9d4df1',
  'b01f60db-9b7d-4081-b339-947a23909505',
  '002b39c3-5e70-4f94-8223-2f16341df7df'
);
```

---

## 🔌 **WebSocket - Tempo Real**

### **Conectar ao WebSocket**

```javascript
// Conectar ao WebSocket de uma sala
const connectToRoom = (roomId) => {
  const ws = new WebSocket(`${WS_BASE_URL}/subscribe/${roomId}`);
  
  ws.onopen = () => {
    console.log(`Conectado à sala ${roomId}`);
  };
  
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    handleRealtimeMessage(data);
  };
  
  ws.onclose = () => {
    console.log('Conexão WebSocket fechada');
  };
  
  ws.onerror = (error) => {
    console.error('Erro no WebSocket:', error);
  };
  
  return ws;
};

// Manipular mensagens em tempo real
const handleRealtimeMessage = (data) => {
  switch (data.kind) {
    case 'message_created':
      console.log('Nova mensagem:', data.value);
      // Adicionar mensagem à UI
      break;
      
    case 'message_reaction_increased':
      console.log('Reação adicionada:', data.value);
      // Atualizar contador de reações
      break;
      
    case 'message_reaction_decreased':
      console.log('Reação removida:', data.value);
      // Atualizar contador de reações
      break;
      
    case 'message_answered':
      console.log('Mensagem respondida:', data.value);
      // Marcar mensagem como respondida
      break;
  }
};
```

### **Tipos de Mensagens WebSocket**

#### **message_created**
```json
{
  "kind": "message_created",
  "value": {
    "id": "message-id",
    "message": "Conteúdo da mensagem"
  }
}
```

#### **message_reaction_increased**
```json
{
  "kind": "message_reaction_increased",
  "value": {
    "id": "message-id",
    "count": 5
  }
}
```

#### **message_reaction_decreased**
```json
{
  "kind": "message_reaction_decreased",
  "value": {
    "id": "message-id",
    "count": 4
  }
}
```

#### **message_answered**
```json
{
  "kind": "message_answered",
  "value": {
    "id": "message-id"
  }
}
```

---

## 🛠️ **Exemplo Completo - React**

```jsx
import React, { useState, useEffect } from 'react';

const AskMeAnything = () => {
  const [rooms, setRooms] = useState([]);
  const [currentRoom, setCurrentRoom] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [ws, setWs] = useState(null);

  // Carregar salas ao montar o componente
  useEffect(() => {
    loadRooms();
  }, []);

  // Conectar ao WebSocket quando sala mudar
  useEffect(() => {
    if (currentRoom) {
      connectToWebSocket(currentRoom.id);
      loadMessages(currentRoom.id);
    }
    
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [currentRoom]);

  const loadRooms = async () => {
    try {
      const roomsData = await getRooms();
      setRooms(roomsData);
    } catch (error) {
      console.error('Erro ao carregar salas:', error);
    }
  };

  const loadMessages = async (roomId) => {
    try {
      const messagesData = await getRoomMessages(roomId);
      setMessages(messagesData);
    } catch (error) {
      console.error('Erro ao carregar mensagens:', error);
    }
  };

  const connectToWebSocket = (roomId) => {
    const websocket = connectToRoom(roomId);
    
    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      if (data.kind === 'message_created') {
        setMessages(prev => [...prev, {
          id: data.value.id,
          message: data.value.message,
          reaction_count: 0,
          answered: false,
        }]);
      } else if (data.kind === 'message_reaction_increased') {
        setMessages(prev => prev.map(msg => 
          msg.id === data.value.id 
            ? { ...msg, reaction_count: data.value.count }
            : msg
        ));
      }
    };
    
    setWs(websocket);
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !currentRoom) return;

    try {
      await createMessage(currentRoom.id, newMessage);
      setNewMessage('');
    } catch (error) {
      console.error('Erro ao enviar mensagem:', error);
    }
  };

  const handleReact = async (messageId) => {
    try {
      await reactToMessage(currentRoom.id, messageId);
    } catch (error) {
      console.error('Erro ao reagir:', error);
    }
  };

  return (
    <div>
      <h1>Ask Me Anything</h1>
      
      {/* Lista de Salas */}
      <div>
        <h2>Salas</h2>
        {rooms.map(room => (
          <button 
            key={room.id} 
            onClick={() => setCurrentRoom(room)}
          >
            {room.theme}
          </button>
        ))}
      </div>

      {/* Mensagens da Sala Atual */}
      {currentRoom && (
        <div>
          <h2>{currentRoom.theme}</h2>
          
          <div>
            {messages.map(message => (
              <div key={message.id}>
                <p>{message.message}</p>
                <button onClick={() => handleReact(message.id)}>
                  👍 {message.reaction_count}
                </button>
                {message.answered && <span>✅ Respondida</span>}
              </div>
            ))}
          </div>

          {/* Formulário de Nova Mensagem */}
          <form onSubmit={handleSendMessage}>
            <input
              type="text"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Digite sua pergunta..."
            />
            <button type="submit">Enviar</button>
          </form>
        </div>
      )}
    </div>
  );
};

export default AskMeAnything;
```

---

## 🔐 **Autenticação**

### **Host Token**
- Obtido ao criar uma sala
- Necessário para ações de host (marcar como respondida)
- Enviar no header `X-Host-Token`

### **Exemplo de Uso**
```javascript
// Criar sala e salvar token
const { id, host_token } = await createRoom('Minha Sala');
localStorage.setItem(`host_token_${id}`, host_token);

// Usar token posteriormente
const hostToken = localStorage.getItem(`host_token_${id}`);
await markMessageAsAnswered(id, messageId, hostToken);
```

---

## ⚠️ **Tratamento de Erros**

```javascript
// Wrapper com tratamento de erro
const safeApiRequest = async (apiCall, errorMessage = 'Erro na API') => {
  try {
    return await apiCall();
  } catch (error) {
    console.error(errorMessage, error);
    
    // Exibir erro para usuário
    alert(`${errorMessage}: ${error.message}`);
    
    // Ou usar sistema de notificação
    showNotification(errorMessage, 'error');
    
    throw error;
  }
};

// Uso
const rooms = await safeApiRequest(
  () => getRooms(),
  'Erro ao carregar salas'
);
```

---

## 🚀 **Dicas de Performance**

1. **Cache de Dados**: Cache salas e mensagens localmente
2. **Debounce**: Use debounce para reações múltiplas
3. **Reconnect**: Implemente reconexão automática do WebSocket
4. **Lazy Loading**: Carregue mensagens por páginas
5. **Throttle**: Limite frequência de requisições

---

## 🌐 **URLs de Acesso**

- **API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080
- **pgAdmin**: http://localhost:8081 (admin@admin.com / password)

---

## 📝 **CORS**

A API está configurada para aceitar requisições de qualquer origem em desenvolvimento:

```javascript
// Headers permitidos
'Accept', 'Authorization', 'Content-Type', 'X-CSRF-Token', 'X-Host-Token'

// Métodos permitidos  
'GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'
```

---

## 🔄 **Fluxo Típico de Uso**

1. **Listar salas** disponíveis
2. **Criar nova sala** (opcional)
3. **Conectar ao WebSocket** da sala
4. **Carregar mensagens** existentes
5. **Enviar nova mensagem**
6. **Reagir a mensagens**
7. **Marcar como respondida** (se for host)

---

**🎉 Agora você tem tudo que precisa para integrar seu frontend com a Ask Me Anything API!**
