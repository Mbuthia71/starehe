import { createFileRoute } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Building2, Users, MapPin, Calendar, Briefcase, GraduationCap } from "lucide-react";
import { useState, useEffect } from "react";

export const Route = createFileRoute("/_app/groups")({
  component: Groups,
});

type Group = {
  id: string;
  name: string;
  type: 'chapter' | 'career' | 'cohort' | 'custom';
  join_policy: 'open' | 'approval_required' | 'auto';
  description?: string;
  member_count: number;
  is_member: boolean;
  user_role?: string;
  created_at: string;
};

function Groups() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<'all' | 'chapter' | 'career' | 'cohort'>('all');
  const [showCreate, setShowCreate] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({
    name: '',
    type: 'chapter' as 'chapter' | 'career' | 'cohort' | 'custom',
    description: '',
  });

  useEffect(() => {
    fetchGroups();
  }, [filter]);

  const createGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) return;
    setCreating(true);
    try {
      const token = localStorage.getItem('oss_token');
      const response = await fetch(`/api/groups`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: form.name.trim(),
          type: form.type,
          join_policy: 'open',
          description: form.description.trim(),
        }),
      });
      if (response.ok) {
        setForm({ name: '', type: 'chapter', description: '' });
        setShowCreate(false);
        fetchGroups();
      } else {
        const data = await response.json().catch(() => ({}));
        alert(data.error || 'Failed to create chapter');
      }
    } catch (error) {
      console.error('Failed to create chapter:', error);
      alert('Failed to create chapter. Please try again.');
    } finally {
      setCreating(false);
    }
  };

  const fetchGroups = async () => {
    try {
      const token = localStorage.getItem('oss_token');
      const response = await fetch(`/api/groups?limit=50&offset=0`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      setGroups(data.groups || []);
    } catch (error) {
      console.error('Failed to fetch groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const joinGroup = async (groupId: string) => {
    try {
      const token = localStorage.getItem('oss_token');
      const response = await fetch(`/api/groups/${groupId}/join`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      if (response.ok) {
        fetchGroups();
      } else {
        const data = await response.json().catch(() => ({}));
        alert(data.error || 'Failed to join group');
      }
    } catch (error) {
      console.error('Failed to join group:', error);
      alert('Failed to join group. Please try again.');
    }
  };

  const leaveGroup = async (groupId: string) => {
    try {
      const token = localStorage.getItem('oss_token');
      const response = await fetch(`/api/groups/${groupId}/leave`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      if (response.ok) {
        fetchGroups();
      }
    } catch (error) {
      console.error('Failed to leave group:', error);
    }
  };

  const filteredGroups = filter === 'all' ? groups : groups.filter(g => g.type === filter);

  const getGroupIcon = (type: string) => {
    switch (type) {
      case 'chapter': return Building2;
      case 'career': return Briefcase;
      case 'cohort': return GraduationCap;
      default: return Users;
    }
  };

  const getGroupTypeLabel = (type: string) => {
    switch (type) {
      case 'chapter': return 'Chapter';
      case 'career': return 'Career';
      case 'cohort': return 'Class Year';
      default: return 'Group';
    }
  };

  return (
    <div className="space-y-6">
      <header className="flex items-start justify-between gap-4">
        <div>
          <div className="label-eyebrow">Alumni groups</div>
          <h1 className="display mt-1 text-3xl font-semibold tracking-tight">
            Find your circle.
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Chapters, career groups, and class cohorts. Join to connect and chat.
          </p>
        </div>
        <button
          onClick={() => setShowCreate((v) => !v)}
          className="shrink-0 rounded-xl bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:brightness-110"
        >
          {showCreate ? 'Cancel' : 'New chapter'}
        </button>
      </header>

      {showCreate && (
        <form onSubmit={createGroup} className="card-elev p-5 space-y-3">
          <div className="grid gap-3 md:grid-cols-2">
            <div>
              <label className="text-xs font-semibold text-muted-foreground">Name</label>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. Nairobi Chapter"
                className="mt-1 w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              />
            </div>
            <div>
              <label className="text-xs font-semibold text-muted-foreground">Type</label>
              <select
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value as typeof form.type })}
                className="mt-1 w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
              >
                <option value="chapter">Chapter</option>
                <option value="career">Career</option>
                <option value="cohort">Class Year</option>
                <option value="custom">Custom</option>
              </select>
            </div>
          </div>
          <div>
            <label className="text-xs font-semibold text-muted-foreground">Description</label>
            <textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              placeholder="What is this group about?"
              rows={2}
              className="mt-1 w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/40"
            />
          </div>
          <button
            type="submit"
            disabled={creating || !form.name.trim()}
            className="rounded-xl bg-primary px-5 py-2 text-xs font-semibold text-primary-foreground hover:brightness-110 disabled:opacity-50"
          >
            {creating ? 'Creating...' : 'Create chapter'}
          </button>
        </form>
      )}

      <div className="flex gap-2">
        {(['all', 'chapter', 'career', 'cohort'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-lg text-xs font-semibold capitalize transition-all ${
              filter === f
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted text-muted-foreground hover:bg-muted/80'
            }`}
          >
            {f}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="text-center py-12 text-muted-foreground">Loading groups...</div>
      ) : (
        <div className="grid gap-3 md:grid-cols-2">
          {filteredGroups.map((group, i) => {
            const Icon = getGroupIcon(group.type);
            return (
              <motion.div
                key={group.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.35, delay: i * 0.04 }}
                className="card-elev p-5"
              >
                <div className="flex items-start gap-3">
                  <div className="grid size-11 shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
                    <Icon className="size-5" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <div className="truncate text-sm font-semibold text-foreground">{group.name}</div>
                      <span className="text-[10px] px-2 py-0.5 rounded-full bg-muted text-muted-foreground capitalize">
                        {getGroupTypeLabel(group.type)}
                      </span>
                    </div>
                    <div className="mt-0.5 flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[11px] text-muted-foreground">
                      <span className="inline-flex items-center gap-1"><Users className="size-3" /> {group.member_count} members</span>
                    </div>
                  </div>
                </div>
                {group.description && (
                  <p className="mt-3 text-xs text-muted-foreground line-clamp-2">{group.description}</p>
                )}
                {group.is_member ? (
                  <button
                    onClick={() => leaveGroup(group.id)}
                    className="mt-4 w-full rounded-xl border border-destructive/50 bg-destructive/10 py-2 text-xs font-semibold text-destructive hover:bg-destructive/20"
                  >
                    Leave group
                  </button>
                ) : (
                  <button
                    onClick={() => joinGroup(group.id)}
                    className="mt-4 w-full rounded-xl bg-primary py-2 text-xs font-semibold text-primary-foreground hover:brightness-110"
                  >
                    {group.join_policy === 'approval_required' ? 'Request to join' : 'Join group'}
                  </button>
                )}
              </motion.div>
            );
          })}
        </div>
      )}

      {!loading && filteredGroups.length === 0 && (
        <div className="text-center py-12 text-muted-foreground">
          No groups found for this filter.
        </div>
      )}
    </div>
  );
}
