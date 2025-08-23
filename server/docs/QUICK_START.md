# 🚀 Ask Me Anything - Guia de Início Rápido

## ⚡ Setup Rápido

### 1. Iniciar a API

```bash
# Clonar o repositório
git clone <seu-repo>
cd ask-me-anything

# Iniciar com Docker
make docker-reload

# OU executar localmente
make run
```

### 2. Testar a API

```bash
# Listar salas
curl http://localhost:8080/api/rooms

# Criar uma sala
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"theme": "Minha primeira sala"}'
```

## 📱 Exemplos Práticos

### Frontend Vanilla JavaScript

```html
<!DOCTYPE html>
<html>
<head>
    <title>Ask Me Anything</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .room { border: 1px solid #ddd; margin: 10px; padding: 15px; border-radius: 8px; }
        .message { background: #f5f5f5; margin: 10px 0; padding: 10px; border-radius: 5px; }
        .reaction-btn { background: #007bff; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; }
        .answered { background: #d4edda; border-color: #c3e6cb; }
        input, textarea { width: 100%; padding: 10px; margin: 5px 0; border: 1px solid #ddd; border-radius: 4px; }
        button { background: #28a745; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>🎯 Ask Me Anything</h1>
    
    <!-- Formulário para criar sala -->
    <div>
        <h2>Criar Nova Sala</h2>
        <input type="text" id="roomTheme" placeholder="Tema da sala...">
        <button onclick="createRoom()">Criar Sala</button>
    </div>
    
    <!-- Lista de salas -->
    <div>
        <h2>Salas Disponíveis</h2>
        <div id="roomsList"></div>
        <button onclick="loadRooms()">Atualizar Salas</button>
    </div>
    
    <!-- Sala atual -->
    <div id="currentRoom" style="display:none;">
        <h2 id="currentRoomTitle"></h2>
        <p><strong>Status:</strong> <span id="connectionStatus">Desconectado</span></p>
        
        <!-- Mensagens -->
        <div id="messagesList"></div>
        
        <!-- Enviar mensagem -->
        <div>
            <textarea id="newMessage" placeholder="Digite sua pergunta..." rows="3"></textarea>
            <button onclick="sendMessage()">Enviar Pergunta</button>
        </div>
    </div>

    <script>
        const API_BASE = 'http://localhost:8080';
        let currentRoomData = null;
        let hostToken = null;
        let websocket = null;

        // Utilitário para requisições
        async function apiRequest(endpoint, options = {}) {
            const response = await fetch(`${API_BASE}${endpoint}`, {
                headers: { 'Content-Type': 'application/json', ...options.headers },
                ...options
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return await response.json();
            }
            return response;
        }

        // Carregar salas
        async function loadRooms() {
            try {
                const rooms = await apiRequest('/api/rooms');
                const roomsList = document.getElementById('roomsList');
                
                if (rooms.length === 0) {
                    roomsList.innerHTML = '<p>Nenhuma sala encontrada. Crie a primeira!</p>';
                    return;
                }
                
                roomsList.innerHTML = rooms.map(room => `
                    <div class="room">
                        <h3>${room.theme}</h3>
                        <button onclick="joinRoom('${room.id}', '${room.theme}')">Entrar na Sala</button>
                    </div>
                `).join('');
            } catch (error) {
                alert('Erro ao carregar salas: ' + error.message);
            }
        }

        // Criar sala
        async function createRoom() {
            const theme = document.getElementById('roomTheme').value.trim();
            if (!theme) {
                alert('Digite um tema para a sala!');
                return;
            }

            try {
                const result = await apiRequest('/api/rooms', {
                    method: 'POST',
                    body: JSON.stringify({ theme })
                });
                
                // Salvar token de host
                localStorage.setItem(`host_token_${result.id}`, result.host_token);
                
                alert(`Sala criada com sucesso!\\nID: ${result.id}\\nVocê é o host desta sala.`);
                document.getElementById('roomTheme').value = '';
                loadRooms();
                
            } catch (error) {
                alert('Erro ao criar sala: ' + error.message);
            }
        }

        // Entrar em uma sala
        async function joinRoom(roomId, roomTitle) {
            try {
                currentRoomData = { id: roomId, theme: roomTitle };
                hostToken = localStorage.getItem(`host_token_${roomId}`);
                
                // Mostrar área da sala
                document.getElementById('currentRoom').style.display = 'block';
                document.getElementById('currentRoomTitle').textContent = roomTitle;
                
                if (hostToken) {
                    document.getElementById('currentRoomTitle').textContent += ' (Você é o host)';
                }
                
                // Carregar mensagens
                await loadMessages(roomId);
                
                // Conectar WebSocket
                connectWebSocket(roomId);
                
            } catch (error) {
                alert('Erro ao entrar na sala: ' + error.message);
            }
        }

        // Carregar mensagens
        async function loadMessages(roomId) {
            try {
                const messages = await apiRequest(`/api/rooms/${roomId}/messages`);
                displayMessages(messages);
            } catch (error) {
                console.error('Erro ao carregar mensagens:', error);
            }
        }

        // Exibir mensagens
        function displayMessages(messages) {
            const messagesList = document.getElementById('messagesList');
            
            if (messages.length === 0) {
                messagesList.innerHTML = '<p>Nenhuma pergunta ainda. Seja o primeiro a perguntar!</p>';
                return;
            }
            
            messagesList.innerHTML = messages.map(msg => `
                <div class="message ${msg.answered ? 'answered' : ''}" id="message-${msg.id}">
                    <p><strong>Pergunta:</strong> ${msg.message}</p>
                    <div>
                        <button class="reaction-btn" onclick="reactToMessage('${msg.id}')">
                            👍 ${msg.reaction_count}
                        </button>
                        ${hostToken && !msg.answered ? 
                            `<button class="reaction-btn" onclick="markAsAnswered('${msg.id}')" style="background: #dc3545;">
                                ✅ Marcar como Respondida
                            </button>` : ''
                        }
                        ${msg.answered ? '<span style="color: green;">✅ Respondida</span>' : ''}
                    </div>
                </div>
            `).join('');
        }

        // Conectar WebSocket
        function connectWebSocket(roomId) {
            if (websocket) {
                websocket.close();
            }
            
            websocket = new WebSocket(`ws://localhost:8080/subscribe/${roomId}`);
            
            websocket.onopen = () => {
                document.getElementById('connectionStatus').textContent = 'Conectado (tempo real)';
                document.getElementById('connectionStatus').style.color = 'green';
            };
            
            websocket.onmessage = (event) => {
                const data = JSON.parse(event.data);
                handleRealtimeMessage(data);
            };
            
            websocket.onclose = () => {
                document.getElementById('connectionStatus').textContent = 'Desconectado';
                document.getElementById('connectionStatus').style.color = 'red';
            };
            
            websocket.onerror = (error) => {
                console.error('Erro WebSocket:', error);
                document.getElementById('connectionStatus').textContent = 'Erro de conexão';
                document.getElementById('connectionStatus').style.color = 'red';
            };
        }

        // Manipular mensagens em tempo real
        function handleRealtimeMessage(data) {
            switch (data.kind) {
                case 'message_created':
                    // Recarregar mensagens para mostrar a nova
                    loadMessages(currentRoomData.id);
                    break;
                    
                case 'message_reaction_increased':
                    updateReactionCount(data.value.id, data.value.count);
                    break;
                    
                case 'message_answered':
                    markMessageAsAnswered(data.value.id);
                    break;
            }
        }

        // Atualizar contador de reações em tempo real
        function updateReactionCount(messageId, count) {
            const messageElement = document.getElementById(`message-${messageId}`);
            if (messageElement) {
                const reactionBtn = messageElement.querySelector('.reaction-btn');
                if (reactionBtn) {
                    reactionBtn.innerHTML = `👍 ${count}`;
                }
            }
        }

        // Marcar mensagem como respondida visualmente
        function markMessageAsAnswered(messageId) {
            const messageElement = document.getElementById(`message-${messageId}`);
            if (messageElement) {
                messageElement.classList.add('answered');
                // Remover botão de marcar como respondida
                const answerBtn = messageElement.querySelector('button[onclick*="markAsAnswered"]');
                if (answerBtn) {
                    answerBtn.remove();
                }
                // Adicionar indicador visual
                const indicator = document.createElement('span');
                indicator.innerHTML = '<span style="color: green;">✅ Respondida</span>';
                messageElement.querySelector('div').appendChild(indicator);
            }
        }

        // Enviar mensagem
        async function sendMessage() {
            const messageText = document.getElementById('newMessage').value.trim();
            if (!messageText || !currentRoomData) {
                alert('Digite uma pergunta!');
                return;
            }

            try {
                await apiRequest(`/api/rooms/${currentRoomData.id}/messages`, {
                    method: 'POST',
                    body: JSON.stringify({ message: messageText })
                });
                
                document.getElementById('newMessage').value = '';
                // A nova mensagem aparecerá via WebSocket
                
            } catch (error) {
                alert('Erro ao enviar mensagem: ' + error.message);
            }
        }

        // Reagir a mensagem
        async function reactToMessage(messageId) {
            try {
                await apiRequest(`/api/rooms/${currentRoomData.id}/messages/${messageId}/react`, {
                    method: 'PATCH'
                });
                // A atualização aparecerá via WebSocket
                
            } catch (error) {
                alert('Erro ao reagir: ' + error.message);
            }
        }

        // Marcar como respondida (apenas hosts)
        async function markAsAnswered(messageId) {
            if (!hostToken) {
                alert('Apenas o host pode marcar mensagens como respondidas!');
                return;
            }

            try {
                await apiRequest(`/api/rooms/${currentRoomData.id}/messages/${messageId}/answer`, {
                    method: 'PATCH',
                    headers: { 'X-Host-Token': hostToken }
                });
                // A atualização aparecerá via WebSocket
                
            } catch (error) {
                alert('Erro ao marcar como respondida: ' + error.message);
            }
        }

        // Carregar salas ao carregar a página
        window.onload = () => {
            loadRooms();
        };

        // Limpar WebSocket ao sair da página
        window.onbeforeunload = () => {
            if (websocket) {
                websocket.close();
            }
        };
    </script>
</body>
</html>
```

### Exemplo React Hook

```javascript
// hooks/useAskMeAnything.js
import { useState, useEffect, useCallback } from 'react';

const API_BASE = 'http://localhost:8080';

export const useAskMeAnything = () => {
  const [rooms, setRooms] = useState([]);
  const [currentRoom, setCurrentRoom] = useState(null);
  const [messages, setMessages] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const [ws, setWs] = useState(null);

  // Utility function for API requests
  const apiRequest = useCallback(async (endpoint, options = {}) => {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      headers: { 'Content-Type': 'application/json', ...options.headers },
      ...options
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }
    return response;
  }, []);

  // Load rooms
  const loadRooms = useCallback(async () => {
    try {
      const roomsData = await apiRequest('/api/rooms');
      setRooms(roomsData);
    } catch (error) {
      console.error('Error loading rooms:', error);
      throw error;
    }
  }, [apiRequest]);

  // Create room
  const createRoom = useCallback(async (theme) => {
    try {
      const result = await apiRequest('/api/rooms', {
        method: 'POST',
        body: JSON.stringify({ theme })
      });
      
      // Save host token
      localStorage.setItem(`host_token_${result.id}`, result.host_token);
      
      // Reload rooms
      await loadRooms();
      
      return result;
    } catch (error) {
      console.error('Error creating room:', error);
      throw error;
    }
  }, [apiRequest, loadRooms]);

  // Join room
  const joinRoom = useCallback(async (room) => {
    try {
      setCurrentRoom(room);
      
      // Load messages
      const messagesData = await apiRequest(`/api/rooms/${room.id}/messages`);
      setMessages(messagesData);
      
      // Connect WebSocket
      const websocket = new WebSocket(`ws://localhost:8080/subscribe/${room.id}`);
      
      websocket.onopen = () => setIsConnected(true);
      websocket.onclose = () => setIsConnected(false);
      websocket.onerror = () => setIsConnected(false);
      
      websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        
        switch (data.kind) {
          case 'message_created':
            setMessages(prev => [...prev, {
              id: data.value.id,
              message: data.value.message,
              reaction_count: 0,
              answered: false,
              room_id: room.id
            }]);
            break;
            
          case 'message_reaction_increased':
            setMessages(prev => prev.map(msg => 
              msg.id === data.value.id 
                ? { ...msg, reaction_count: data.value.count }
                : msg
            ));
            break;
            
          case 'message_answered':
            setMessages(prev => prev.map(msg => 
              msg.id === data.value.id 
                ? { ...msg, answered: true }
                : msg
            ));
            break;
        }
      };
      
      setWs(websocket);
      
    } catch (error) {
      console.error('Error joining room:', error);
      throw error;
    }
  }, [apiRequest]);

  // Send message
  const sendMessage = useCallback(async (message) => {
    if (!currentRoom) return;
    
    try {
      await apiRequest(`/api/rooms/${currentRoom.id}/messages`, {
        method: 'POST',
        body: JSON.stringify({ message })
      });
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  }, [apiRequest, currentRoom]);

  // React to message
  const reactToMessage = useCallback(async (messageId) => {
    if (!currentRoom) return;
    
    try {
      await apiRequest(`/api/rooms/${currentRoom.id}/messages/${messageId}/react`, {
        method: 'PATCH'
      });
    } catch (error) {
      console.error('Error reacting to message:', error);
      throw error;
    }
  }, [apiRequest, currentRoom]);

  // Mark as answered
  const markAsAnswered = useCallback(async (messageId) => {
    if (!currentRoom) return;
    
    const hostToken = localStorage.getItem(`host_token_${currentRoom.id}`);
    if (!hostToken) {
      throw new Error('Host token required');
    }
    
    try {
      await apiRequest(`/api/rooms/${currentRoom.id}/messages/${messageId}/answer`, {
        method: 'PATCH',
        headers: { 'X-Host-Token': hostToken }
      });
    } catch (error) {
      console.error('Error marking as answered:', error);
      throw error;
    }
  }, [apiRequest, currentRoom]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [ws]);

  // Check if user is host
  const isHost = useCallback(() => {
    if (!currentRoom) return false;
    return !!localStorage.getItem(`host_token_${currentRoom.id}`);
  }, [currentRoom]);

  return {
    // State
    rooms,
    currentRoom,
    messages,
    isConnected,
    
    // Actions
    loadRooms,
    createRoom,
    joinRoom,
    sendMessage,
    reactToMessage,
    markAsAnswered,
    
    // Helpers
    isHost
  };
};
```

### Uso do Hook no React

```jsx
// components/AskMeAnything.jsx
import React, { useState, useEffect } from 'react';
import { useAskMeAnything } from '../hooks/useAskMeAnything';

const AskMeAnything = () => {
  const {
    rooms,
    currentRoom,
    messages,
    isConnected,
    loadRooms,
    createRoom,
    joinRoom,
    sendMessage,
    reactToMessage,
    markAsAnswered,
    isHost
  } = useAskMeAnything();

  const [newRoomTheme, setNewRoomTheme] = useState('');
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    loadRooms();
  }, [loadRooms]);

  const handleCreateRoom = async (e) => {
    e.preventDefault();
    if (!newRoomTheme.trim()) return;
    
    try {
      await createRoom(newRoomTheme);
      setNewRoomTheme('');
    } catch (error) {
      alert('Erro ao criar sala: ' + error.message);
    }
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;
    
    try {
      await sendMessage(newMessage);
      setNewMessage('');
    } catch (error) {
      alert('Erro ao enviar mensagem: ' + error.message);
    }
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <h1>🎯 Ask Me Anything</h1>
      
      {!currentRoom ? (
        <div>
          {/* Create Room Form */}
          <form onSubmit={handleCreateRoom} style={{ marginBottom: '30px' }}>
            <h2>Criar Nova Sala</h2>
            <input
              type="text"
              value={newRoomTheme}
              onChange={(e) => setNewRoomTheme(e.target.value)}
              placeholder="Tema da sala..."
              style={{ width: '70%', padding: '10px', marginRight: '10px' }}
            />
            <button type="submit" style={{ padding: '10px 20px' }}>
              Criar Sala
            </button>
          </form>

          {/* Rooms List */}
          <div>
            <h2>Salas Disponíveis</h2>
            {rooms.length === 0 ? (
              <p>Nenhuma sala encontrada. Crie a primeira!</p>
            ) : (
              rooms.map(room => (
                <div key={room.id} style={{ 
                  border: '1px solid #ddd', 
                  margin: '10px 0', 
                  padding: '15px', 
                  borderRadius: '8px' 
                }}>
                  <h3>{room.theme}</h3>
                  <button onClick={() => joinRoom(room)}>
                    Entrar na Sala
                  </button>
                </div>
              ))
            )}
            <button onClick={loadRooms} style={{ marginTop: '10px' }}>
              Atualizar Salas
            </button>
          </div>
        </div>
      ) : (
        <div>
          {/* Current Room */}
          <div style={{ marginBottom: '20px' }}>
            <h2>{currentRoom.theme} {isHost() && '(Host)'}</h2>
            <p>
              Status: <span style={{ color: isConnected ? 'green' : 'red' }}>
                {isConnected ? 'Conectado (tempo real)' : 'Desconectado'}
              </span>
            </p>
            <button onClick={() => window.location.reload()}>
              Voltar às Salas
            </button>
          </div>

          {/* Messages */}
          <div style={{ marginBottom: '20px' }}>
            <h3>Perguntas</h3>
            {messages.length === 0 ? (
              <p>Nenhuma pergunta ainda. Seja o primeiro a perguntar!</p>
            ) : (
              messages.map(message => (
                <div 
                  key={message.id} 
                  style={{ 
                    background: message.answered ? '#d4edda' : '#f5f5f5',
                    margin: '10px 0', 
                    padding: '15px', 
                    borderRadius: '8px',
                    border: message.answered ? '1px solid #c3e6cb' : '1px solid #ddd'
                  }}
                >
                  <p><strong>Pergunta:</strong> {message.message}</p>
                  <div>
                    <button 
                      onClick={() => reactToMessage(message.id)}
                      style={{ 
                        background: '#007bff', 
                        color: 'white', 
                        border: 'none', 
                        padding: '5px 10px', 
                        borderRadius: '3px', 
                        marginRight: '10px' 
                      }}
                    >
                      👍 {message.reaction_count}
                    </button>
                    
                    {isHost() && !message.answered && (
                      <button 
                        onClick={() => markAsAnswered(message.id)}
                        style={{ 
                          background: '#dc3545', 
                          color: 'white', 
                          border: 'none', 
                          padding: '5px 10px', 
                          borderRadius: '3px' 
                        }}
                      >
                        ✅ Marcar como Respondida
                      </button>
                    )}
                    
                    {message.answered && (
                      <span style={{ color: 'green', marginLeft: '10px' }}>
                        ✅ Respondida
                      </span>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>

          {/* Send Message Form */}
          <form onSubmit={handleSendMessage}>
            <h3>Fazer uma Pergunta</h3>
            <textarea
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Digite sua pergunta..."
              rows="3"
              style={{ 
                width: '100%', 
                padding: '10px', 
                borderRadius: '4px', 
                border: '1px solid #ddd' 
              }}
            />
            <br />
            <button 
              type="submit" 
              style={{ 
                background: '#28a745', 
                color: 'white', 
                padding: '10px 20px', 
                border: 'none', 
                borderRadius: '4px', 
                marginTop: '10px' 
              }}
            >
              Enviar Pergunta
            </button>
          </form>
        </div>
      )}
    </div>
  );
};

export default AskMeAnything;
```

## 🎯 Comandos Úteis

```bash
# Iniciar API
make docker-reload

# Ver logs
make docker-logs

# Parar tudo
make dev-stop

# Executar migrações
make migrate-up

# Verificar status
docker compose ps

# Testar API
curl http://localhost:8080/api/rooms
```

## 🔗 URLs Importantes

- **API**: http://localhost:8080
- **WebSocket**: ws://localhost:8080
- **pgAdmin**: http://localhost:8081 (admin@admin.com / password)
- **Documentação**: [docs/API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

---

**🎉 Agora você tem exemplos completos para integrar com sua API!**
