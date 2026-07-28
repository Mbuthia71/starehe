import { createFileRoute } from '@tanstack/react-router';
import { useEffect, useState, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Building2, Send, MessageSquare, ArrowLeft } from 'lucide-react';

export const Route = createFileRoute('/_app/messages')({
  component: MessagesPage,
});

type Chapter = {
  id: string;
  name: string;
  type: string;
  member_count: number;
  is_member: boolean;
};

type ChatMessage = {
  id: string;
  group_id?: string;
  sender_id: string;
  sender_name?: string;
  sender_avatar?: string;
  content?: string;
  media_url?: string;
  created_at: string;
  delivered?: boolean;
  read?: boolean;
};

const getToken = () => {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem('oss_token') || '';
};

// Decode the JWT payload to identify the current user (client-only).
const getMyUserId = (): string => {
  if (typeof window === 'undefined') return '';
  const stored = localStorage.getItem('oss_user_id');
  if (stored) return stored;
  const token = getToken();
  if (!token) return '';
  try {
    const payload = JSON.parse(atob(token.split('.')[1] || ''));
    return payload.sub || payload.user_id || payload.id || '';
  } catch {
    return '';
  }
};

const authHeaders = () => ({
  'Content-Type': 'application/json',
  Authorization: `Bearer ${getToken()}`,
});

function MessagesPage() {
  const [selectedGroup, setSelectedGroup] = useState<Chapter | null>(null);
  const [newMessage, setNewMessage] = useState('');
  const [sending, setSending] = useState(false);
  const [optimisticMessages, setOptimisticMessages] = useState<ChatMessage[]>([]);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const myUserId = getMyUserId();

  // Fetch the chapters/groups the user belongs to.
  const { data: groups = [], isLoading: groupsLoading } = useQuery<Chapter[]>({
    queryKey: ['chat-chapters'],
    queryFn: async () => {
      const res = await fetch('/api/groups?limit=100&offset=0', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return (data.groups || []).filter((g: Chapter) => g.is_member);
    },
    enabled: !!getToken(),
  });

  // Fetch messages for the selected chapter, polling for near real-time.
  const { data: messages = [], isLoading: messagesLoading } = useQuery<
    ChatMessage[]
  >({
    queryKey: ['chat-group-messages', selectedGroup?.id],
    queryFn: async () => {
      if (!selectedGroup) return [];
      const res = await fetch(
        `/api/chat/group/${selectedGroup.id}/messages?limit=100&offset=0`,
        { headers: authHeaders() }
      );
      if (!res.ok) return [];
      const data = await res.json();
      return data.messages || [];
    },
    enabled: !!selectedGroup && !!getToken(),
    refetchInterval: 2500,
  });

  // Merge optimistic messages with fetched messages
  const allMessages = [...messages];
  optimisticMessages.forEach((optMsg) => {
    if (!allMessages.find((m) => m.id === optMsg.id)) {
      allMessages.push(optMsg);
    }
  });
  // Sort by created_at
  allMessages.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [allMessages]);

  // Clear optimistic messages when they appear in the real messages
  useEffect(() => {
    const messageIds = new Set(messages.map((m) => m.id));
    setOptimisticMessages((prev) => prev.filter((m) => !messageIds.has(m.id)));
  }, [messages]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !selectedGroup) return;

    const tempId = `temp-${Date.now()}`;
    const optimisticMsg: ChatMessage = {
      id: tempId,
      sender_id: myUserId,
      content: newMessage.trim(),
      created_at: new Date().toISOString(),
      delivered: false,
      read: false,
    };

    // Add optimistic message immediately
    setOptimisticMessages((prev) => [...prev, optimisticMsg]);
    setNewMessage('');

    setSending(true);
    try {
      const res = await fetch(`/api/chat/group/${selectedGroup.id}/messages`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({
          content: optimisticMsg.content,
        }),
      });
      if (!res.ok) throw new Error('send failed');
      // The optimistic message will be cleared when the real message appears
    } catch (err) {
      console.error('Failed to send message:', err);
      // Remove optimistic message on error
      setOptimisticMessages((prev) => prev.filter((m) => m.id !== tempId));
      alert('Failed to send message. Please try again.');
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="flex h-[calc(100vh-5rem)] md:h-[calc(100vh-8rem)] gap-0 md:gap-4">
      {/* Chapter list */}
      <div className={`w-full md:w-72 shrink-0 card-elev flex flex-col overflow-hidden ${selectedGroup ? 'hidden md:flex' : 'flex'}`}>
        <div className="p-4 md:p-4 border-b border-border">
          <div className="label-eyebrow">Chapter chats</div>
          <h2 className="text-lg font-semibold">Messages</h2>
        </div>
        <div className="overflow-y-auto flex-1">
          {groupsLoading ? (
            <div className="p-4 text-center text-sm text-muted-foreground">
              Loading chapters...
            </div>
          ) : groups.length === 0 ? (
            <div className="p-6 text-center text-sm text-muted-foreground">
              You haven't joined any chapters yet. Visit{' '}
              <span className="font-semibold text-foreground">Chapters</span> to
              join one and start chatting.
            </div>
          ) : (
            groups.map((g) => (
              <button
                key={g.id}
                onClick={() => setSelectedGroup(g)}
                className={`w-full flex items-center gap-3 p-4 md:p-4 text-left transition-colors hover:bg-muted active:bg-muted/80 ${
                  selectedGroup?.id === g.id ? 'bg-muted' : ''
                }`}
              >
                <div className="grid size-10 md:size-10 shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
                  <Building2 className="size-5" />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="truncate text-sm font-semibold text-foreground">
                    {g.name}
                  </div>
                  <div className="text-[11px] text-muted-foreground">
                    {g.member_count} members
                  </div>
                </div>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Conversation panel */}
      <div className={`flex-1 card-elev flex flex-col overflow-hidden ${!selectedGroup ? 'hidden md:flex' : 'flex'}`}>
        {selectedGroup ? (
          <>
            <div className="border-b border-border p-3 md:p-4 flex items-center gap-3">
              <button
                onClick={() => setSelectedGroup(null)}
                className="md:hidden grid size-9 place-items-center rounded-xl bg-muted hover:bg-muted/80 active:bg-muted/60"
              >
                <ArrowLeft className="size-5" />
              </button>
              <div className="grid size-9 md:size-9 place-items-center rounded-2xl bg-primary/10 text-primary">
                <Building2 className="size-4" />
              </div>
              <div className="min-w-0 flex-1">
                <h1 className="font-semibold leading-tight truncate text-base md:text-base">{selectedGroup.name}</h1>
                <div className="text-[11px] text-muted-foreground">
                  {selectedGroup.member_count} members
                </div>
              </div>
            </div>

            <div className="flex-1 overflow-y-auto p-3 md:p-4 space-y-2 md:space-y-3">
              {messagesLoading ? (
                <div className="flex items-center justify-center h-full text-sm text-muted-foreground">
                  Loading messages...
                </div>
              ) : allMessages.length === 0 ? (
                <div className="flex items-center justify-center h-full text-sm text-muted-foreground">
                  No messages yet. Say hello!
                </div>
              ) : (
                allMessages.map((msg) => {
                  const mine = msg.sender_id === myUserId;
                  const isOptimistic = msg.id.startsWith('temp-');
                  return (
                    <div
                      key={msg.id}
                      className={`flex gap-2 md:gap-2 ${mine ? 'justify-end' : 'justify-start'}`}
                    >
                      {!mine && (
                        <div className="flex-shrink-0">
                          {msg.sender_avatar ? (
                            <img
                              src={msg.sender_avatar}
                              alt={msg.sender_name || 'User'}
                              className="size-8 md:size-8 rounded-full object-cover"
                            />
                          ) : (
                            <div className="size-8 md:size-8 rounded-full bg-muted flex items-center justify-center text-xs font-semibold">
                              {msg.sender_name?.charAt(0).toUpperCase() || '?'}
                            </div>
                          )}
                        </div>
                      )}
                      <div className="flex flex-col max-w-[85%] md:max-w-[70%]">
                        {!mine && msg.sender_name && (
                          <span className="text-xs font-medium text-muted-foreground mb-1 ml-1">
                            {msg.sender_name}
                          </span>
                        )}
                        <div
                          className={`rounded-2xl px-3 py-2 md:px-4 md:py-2 ${
                            mine
                              ? 'bg-primary text-primary-foreground'
                              : 'bg-muted text-foreground'
                          }`}
                        >
                          {msg.content && (
                            <p className="text-sm md:text-sm whitespace-pre-wrap break-words">
                              {msg.content}
                            </p>
                          )}
                          {msg.media_url && (
                            <img
                              src={msg.media_url}
                              alt="Media"
                              className="max-w-full md:max-w-xs rounded-lg mt-2"
                            />
                          )}
                          <div
                            className={`text-[10px] mt-1 flex items-center gap-1 ${
                              mine
                                ? 'text-primary-foreground/70'
                                : 'text-muted-foreground'
                            }`}
                          >
                            {new Date(msg.created_at).toLocaleTimeString([], {
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                            {mine && (
                              <span className="flex items-center gap-0.5">
                                {isOptimistic ? (
                                  <span className="opacity-50">⏳</span>
                                ) : msg.read ? (
                                  <span className="text-blue-300">✓✓</span>
                                ) : msg.delivered ? (
                                  <span>✓✓</span>
                                ) : (
                                  <span>✓</span>
                                )}
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })
              )}
              <div ref={messagesEndRef} />
            </div>

            <form
              onSubmit={handleSendMessage}
              className="border-t border-border p-3 md:p-3 flex gap-2"
            >
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder={`Message ${selectedGroup.name}...`}
                disabled={sending}
                className="flex-1 rounded-xl border border-border bg-background px-4 py-3 md:px-4 md:py-2 text-base md:text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              />
              <button
                type="submit"
                disabled={sending || !newMessage.trim()}
                className="inline-flex items-center justify-center gap-1.5 rounded-xl bg-primary px-4 py-3 md:px-4 md:py-2 text-sm font-semibold text-primary-foreground hover:brightness-110 active:brightness-120 disabled:opacity-50"
              >
                <Send className="size-4" />
                <span className="hidden md:inline">{sending ? 'Sending' : 'Send'}</span>
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center text-center text-muted-foreground gap-3 p-4">
            <div className="grid size-14 md:size-14 place-items-center rounded-3xl bg-muted">
              <MessageSquare className="size-7" />
            </div>
            <p className="text-sm md:text-sm">Select a chapter to start messaging</p>
          </div>
        )}
      </div>
    </div>
  );
}
