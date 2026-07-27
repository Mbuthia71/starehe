import { createFileRoute, useRouter } from '@tanstack/react-router';
import { useEffect, useState, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import chatApi, { Conversation, Message } from '../../lib/chat';

export const Route = createFileRoute('/dashboard/messages')();

function MessagesPage() {
  const router = useRouter();
  const [selectedConversation, setSelectedConversation] =
    useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Get auth token from local storage or session
  const getToken = () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.navigate({ to: '/auth/login' });
    }
    return token || '';
  };

  // Fetch conversations
  const { data: conversationsData, isLoading: conversationsLoading } =
    useQuery({
      queryKey: ['conversations'],
      queryFn: async () => {
        try {
          return await chatApi.getConversations(getToken());
        } catch (error) {
          console.error('Failed to fetch conversations:', error);
          return { conversations: [] };
        }
      },
      enabled: !!getToken(),
      refetchInterval: 5000,
    });

  // Fetch messages for selected conversation
  const { data: messagesData, isLoading: messagesLoading } = useQuery({
    queryKey: ['messages', selectedConversation?.id],
    queryFn: async () => {
      if (!selectedConversation) return null;
      try {
        return await chatApi.getMessages(getToken(), selectedConversation.id);
      } catch (error) {
        console.error('Failed to fetch messages:', error);
        return { messages: [] };
      }
    },
    enabled: !!selectedConversation && !!getToken(),
    refetchInterval: 2000,
  });

  useEffect(() => {
    if (messagesData?.messages) {
      setMessages(messagesData.messages);
      // Auto-scroll to bottom
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messagesData]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !selectedConversation) return;

    setLoading(true);
    try {
      await chatApi.sendMessage(
        getToken(),
        selectedConversation.id,
        newMessage
      );
      setNewMessage('');
      // Refetch messages after a short delay
      await new Promise((resolve) => setTimeout(resolve, 500));
    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Failed to send message. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const conversations = conversationsData?.conversations || [];

  return (
    <div className="flex h-full gap-4">
      {/* Conversations List */}
      <div className="w-64 border-r flex flex-col">
        <div className="p-4 border-b">
          <h2 className="text-lg font-semibold">Messages</h2>
        </div>
        <div className="divide-y overflow-y-auto flex-1">
          {conversationsLoading ? (
            <div className="p-4 text-center text-gray-500">Loading...</div>
          ) : conversations.length === 0 ? (
            <div className="p-4 text-center text-gray-500">No conversations</div>
          ) : (
            conversations.map((conv: Conversation) => (
              <button
                key={conv.id}
                onClick={() => setSelectedConversation(conv)}
                className={`w-full p-4 text-left hover:bg-gray-100 transition-colors ${
                  selectedConversation?.id === conv.id ? 'bg-blue-50' : ''
                }`}
              >
                <div className="font-medium">{conv.name || 'Direct Message'}</div>
                <div className="text-sm text-gray-500">
                  {new Date(conv.updated_at).toLocaleDateString()}
                </div>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 flex flex-col">
        {selectedConversation ? (
          <>
            {/* Header */}
            <div className="border-b p-4">
              <h1 className="font-semibold">{selectedConversation.name || 'Direct Message'}</h1>
            </div>

            {/* Messages Container */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messagesLoading ? (
                <div className="flex items-center justify-center h-full">
                  <div className="text-gray-500">Loading messages...</div>
                </div>
              ) : messages.length === 0 ? (
                <div className="flex items-center justify-center h-full">
                  <div className="text-gray-500">No messages yet</div>
                </div>
              ) : (
                <>
                  {messages.map((msg: Message) => (
                    <div
                      key={msg.id}
                      className="flex gap-2"
                    >
                      <div className="flex-1">
                        <div className="bg-gray-200 rounded-lg p-3 max-w-xs">
                          {msg.content && <p>{msg.content}</p>}
                          {msg.media_url && (
                            <img
                              src={msg.media_url}
                              alt="Media"
                              className="max-w-sm rounded mt-2"
                            />
                          )}
                        </div>
                        <div className="text-xs text-gray-400 mt-1">
                          {new Date(msg.created_at).toLocaleTimeString()}
                        </div>
                      </div>
                    </div>
                  ))}
                  <div ref={messagesEndRef} />
                </>
              )}
            </div>

            {/* Message Input */}
            <form onSubmit={handleSendMessage} className="border-t p-4 flex gap-2">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="Type a message..."
                disabled={loading}
                className="flex-1 rounded-lg border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="submit"
                disabled={loading || !newMessage.trim()}
                className="bg-blue-600 text-white rounded-lg px-4 py-2 hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? 'Sending...' : 'Send'}
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Select a conversation to start messaging
          </div>
        )}
      </div>
    </div>
  );
}
