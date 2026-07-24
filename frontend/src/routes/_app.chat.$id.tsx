import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect, useRef } from "react";
import { Send, Users, MoreVertical } from "lucide-react";

export const Route = createFileRoute("/_app/chat/$id")({
  component: GroupChat,
});

type Message = {
  id: string;
  group_id?: string;
  recipient_id?: string;
  sender_id: string;
  content?: string;
  media_url?: string;
  created_at: string;
  sender_name?: string;
  sender_avatar?: string;
};

type Group = {
  id: string;
  name: string;
  type: 'chapter' | 'career' | 'cohort' | 'custom';
  member_count: number;
  is_member: boolean;
  user_role?: string;
};

function GroupChat() {
  const { id } = Route.useParams();
  const [messages, setMessages] = useState<Message[]>([]);
  const [group, setGroup] = useState<Group | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [sending, setSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [onlineCount, setOnlineCount] = useState(0);

  useEffect(() => {
    fetchGroup();
    fetchMessages();
    // TODO: Setup Centrifugo WebSocket connection
  }, [id]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const fetchGroup = async () => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`/api/groups/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      setGroup(data);
    } catch (error) {
      console.error('Failed to fetch group:', error);
    }
  };

  const fetchMessages = async () => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`/api/chat/group/${id}?limit=50&offset=0`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      setMessages(data.messages || []);
    } catch (error) {
      console.error('Failed to fetch messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!message.trim() || sending) return;

    setSending(true);
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`/api/chat/conversations/${id}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: message,
          group_id: id,
        }),
      });

      if (response.ok) {
        const newMessage = await response.json();
        setMessages([...messages, newMessage]);
        setMessage("");
      }
    } catch (error) {
      console.error('Failed to send message:', error);
    } finally {
      setSending(false);
    }
  };

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-muted-foreground">Loading chat...</div>
      </div>
    );
  }

  return (
    <div className="flex h-screen flex-col">
      {/* Chat header */}
      <div className="border-b bg-card p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <h1 className="text-lg font-semibold">{group?.name || 'Group Chat'}</h1>
              <span className="text-xs text-muted-foreground">
                {group?.member_count || 0} members
              </span>
            </div>
            {onlineCount > 0 && (
              <span className="flex items-center gap-1 text-xs text-green-600">
                <span className="size-2 rounded-full bg-green-500" />
                {onlineCount} online
              </span>
            )}
          </div>
          <button className="p-2 hover:bg-muted rounded-lg">
            <MoreVertical className="size-5" />
          </button>
        </div>
      </div>

      {/* Messages list */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 ? (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            No messages yet. Start the conversation!
          </div>
        ) : (
          messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.sender_id === 'current-user' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                  msg.sender_id === 'current-user'
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted'
                }`}
              >
                {msg.sender_name && msg.sender_id !== 'current-user' && (
                  <div className="text-xs font-semibold mb-1 opacity-70">
                    {msg.sender_name}
                  </div>
                )}
                <div className="text-sm">{msg.content}</div>
                <div className="text-xs opacity-50 mt-1">
                  {new Date(msg.created_at).toLocaleTimeString([], {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </div>
              </div>
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Message composer */}
      <div className="border-t bg-card p-4">
        <form onSubmit={sendMessage} className="flex gap-2">
          <input
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type a message..."
            className="flex-1 rounded-xl border border-input bg-background px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            disabled={sending}
          />
          <button
            type="submit"
            disabled={!message.trim() || sending}
            className="rounded-xl bg-primary px-4 py-2 text-primary-foreground hover:brightness-110 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Send className="size-4" />
          </button>
        </form>
      </div>
    </div>
  );
}
