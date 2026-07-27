import { createFileRoute, useRouter } from '@tanstack/react-router';
import { useEffect, useState } from 'react';
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

  // Get auth token from local storage or session
  const getToken = () => localStorage.getItem('access_token') || '';

  // Fetch conversations
  const { data: conversationsData, isLoading: conversationsLoading } =
    useQuery({
      queryKey: ['conversations'],
      queryFn: () => chatApi.getConversations(getToken()),
      enabled: !!getToken(),
    });

  // Fetch messages for selected conversation
  const { data: messagesData, isLoading: messagesLoading } = useQuery({
    queryKey: ['messages', selectedConversation?.id],
    queryFn: () =>
      selectedConversation
        ? chatApi.getMessages(getToken(), selectedConversation.id)
        : Promise.resolve(null),
    enabled: !!selectedConversation && !!getToken(),
    refetchInterval: 3000, // Poll every 3 seconds
  });

  useEffect(() => {
    if (messagesData?.messages) {
      setMessages(messagesData.messages);
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
      // Refetch messages
      await new Promise((resolve) => setTimeout(resolve, 500));
    } catch (error) {
      console.error('Failed to send message:', error);
    } finally {
      setLoading(false);
    }
  };

  const conversations = conversationsData?.conversations || [];

  return (
    <div className="flex h-full gap-4">
      {/* Conversations List */}
      <div className="w-64 border-r">
        <div className="p-4 border-b">
          <h2 className="text-lg font-semibold">Messages</h2>
        </div>
        <div className="divide-y overflow-y-auto">
          {conversationsLoading ? (
            <div className="p-4 text-center text-muted-foreground">Loading...</div>
          ) : conversations.length === 0 ? (
            <div className="p-4 text-center text-muted-foreground">No conversations</div>
          ) : (
            conversations.map((conv: Conversation) => (
              <button
                key={conv.id}
                onClick={() => setSelectedConversation(conv)}
                className={`w-full p-4 text-left hover:bg-accent transition-colors ${
                  selectedConversation?.id === conv.id ? 'bg-accent' : ''
                }`}
              >
                <div className="font-medium">{conv.name || 'Direct Message'}</div>
                <div className="text-sm text-muted-foreground">
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
                  <div className="text-muted-foreground">Loading messages...</div>
                </div>
              ) : messages.length === 0 ? (
                <div className="flex items-center justify-center h-full">
                  <div className="text-muted-foreground">No messages yet</div>
                </div>
              ) : (
                messages.map((msg: Message) => (
                  <div
                    key={msg.id}
                    className="flex gap-2"
                  >
                    <div className="flex-1">
                      <div className="bg-secondary rounded-lg p-3">
                        {msg.content && <p>{msg.content}</p>}
                        {msg.media_url && (
                          <img
                            src={msg.media_url}
                            alt="Media"
                            className="max-w-sm rounded mt-2"
                          />
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        {new Date(msg.created_at).toLocaleTimeString()}
                      </div>
                    </div>
                  </div>
                ))
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
                className="flex-1 rounded-lg border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary"
              />
              <button
                type="submit"
                disabled={loading || !newMessage.trim()}
                className="bg-primary text-primary-foreground rounded-lg px-4 py-2 hover:bg-primary/90 disabled:opacity-50"
              >
                {loading ? 'Sending...' : 'Send'}
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-muted-foreground">
            Select a conversation to start messaging
          </div>
        )}
      </div>
    </div>
  );
}
