import { API_CONFIG } from './api';

export interface Conversation {
  id: string;
  type: 'direct' | 'group';
  name?: string;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  conversation_id?: string;
  group_id?: string;
  recipient_id?: string;
  sender_id: string;
  content?: string;
  media_url?: string;
  created_at: string;
  read_at?: string;
}

export interface ConversationMember {
  id: string;
  conversation_id: string;
  user_id: string;
  role: 'admin' | 'member';
  joined_at: string;
}

const chatApi = {
  // Conversations
  getConversations: async (token: string, limit = 20, offset = 0) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/?limit=${limit}&offset=${offset}`,
      {
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to fetch conversations');
    return response.json();
  },

  createDirectConversation: async (token: string, targetUserId: string) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/direct/${targetUserId}`,
      {
        method: 'POST',
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to create conversation');
    return response.json();
  },

  createGroupConversation: async (
    token: string,
    name: string,
    members: string[]
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/group`,
      {
        method: 'POST',
        headers: API_CONFIG.headers(token),
        body: JSON.stringify({
          name,
          members,
          type: 'group',
        }),
      }
    );
    if (!response.ok) throw new Error('Failed to create group conversation');
    return response.json();
  },

  // Messages
  getMessages: async (
    token: string,
    conversationId: string,
    limit = 50,
    offset = 0
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}/messages?limit=${limit}&offset=${offset}`,
      {
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to fetch messages');
    return response.json();
  },

  sendMessage: async (
    token: string,
    conversationId: string,
    content: string,
    mediaUrl?: string
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}/messages`,
      {
        method: 'POST',
        headers: API_CONFIG.headers(token),
        body: JSON.stringify({
          content,
          media_url: mediaUrl,
        }),
      }
    );
    if (!response.ok) throw new Error('Failed to send message');
    return response.json();
  },

  markAsRead: async (
    token: string,
    conversationId: string,
    messageId: string
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}/read`,
      {
        method: 'POST',
        headers: API_CONFIG.headers(token),
        body: JSON.stringify({
          message_id: messageId,
        }),
      }
    );
    if (!response.ok) throw new Error('Failed to mark as read');
    return response.json();
  },

  // Centrifugo tokens
  getConnectionToken: async (token: string) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/token`,
      {
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to get connection token');
    return response.json();
  },

  getChannelToken: async (token: string, channel: string) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/channel-token?channel=${channel}`,
      {
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to get channel token');
    return response.json();
  },

  // Member management
  addMember: async (
    token: string,
    conversationId: string,
    userId: string
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}/members`,
      {
        method: 'POST',
        headers: API_CONFIG.headers(token),
        body: JSON.stringify({
          user_id: userId,
        }),
      }
    );
    if (!response.ok) throw new Error('Failed to add member');
    return response.json();
  },

  removeMember: async (
    token: string,
    conversationId: string,
    memberId: string
  ) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}/members/${memberId}`,
      {
        method: 'DELETE',
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to remove member');
    return response.json();
  },

  leaveConversation: async (token: string, conversationId: string) => {
    const response = await fetch(
      `${API_CONFIG.baseUrl}/chat/${conversationId}`,
      {
        method: 'DELETE',
        headers: API_CONFIG.headers(token),
      }
    );
    if (!response.ok) throw new Error('Failed to leave conversation');
    return response.json();
  },
};

export default chatApi;
