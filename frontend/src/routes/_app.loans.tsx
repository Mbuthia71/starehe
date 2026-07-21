import { createFileRoute } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { useState } from "react";
import { UserPlus, Check, X, Users } from "lucide-react";
import { StatusPill } from "@/components/StatusPill";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/loans")({
  component: Connections,
});

interface ConnectionRequest {
  id: string;
  name: string;
  detail: string;
  status: "Pending" | "Accepted" | "Declined";
  when: string;
}

const initial: ConnectionRequest[] = [
  { id: "c1", name: "Kevin Mwangi", detail: "Griffin · Class of 2008 · Product Manager", status: "Pending", when: "12 min ago" },
  { id: "c2", name: "Sarah Njeri", detail: "Griffin · Class of 2012 · Journalist", status: "Pending", when: "2 h ago" },
  { id: "c3", name: "Brian Otieno", detail: "Livingstone · Class of 2011 · Senior Engineer", status: "Accepted", when: "Yesterday" },
  { id: "c4", name: "Peter Kimani", detail: "Livingstone · Class of 2009 · Data Scientist", status: "Accepted", when: "3 d ago" },
];

const suggestions = [
  { id: "s1", name: "Mary Wanjiku", detail: "Griffin · Class of 2015 · Lawyer", reason: "2 mutual connections" },
  { id: "s2", name: "David Ochieng", detail: "Livingstone · Class of 2007 · Architect", reason: "Same chapter · Nairobi" },
  { id: "s3", name: "Achieng Nyongo", detail: "Griffin · Class of 2013 · Doctor", reason: "5 mutual connections" },
];

function Connections() {
  const [requests, setRequests] = useState(initial);
  const pending = requests.filter((r) => r.status === "Pending");
  const accepted = requests.filter((r) => r.status === "Accepted");

  const decide = (id: string, status: "Accepted" | "Declined") => {
    setRequests((rs) => rs.map((r) => (r.id === id ? { ...r, status } : r)));
  };

  return (
    <div className="space-y-8">
      <header>
        <div className="label-eyebrow">Alumni connections</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">
          Your network.
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Accept requests, discover fellow Old Starehians and grow your circle.
        </p>
      </header>

      <section>
        <div className="mb-3 flex items-baseline justify-between">
          <div className="label-eyebrow">Pending requests</div>
          <span className="text-xs text-muted-foreground">{pending.length} waiting</span>
        </div>
        {pending.length === 0 ? (
          <div className="card-elev">
            <EmptyState icon={UserPlus} title="No pending requests" description="You're all caught up." />
          </div>
        ) : (
          <div className="space-y-2">
            {pending.map((r, i) => (
              <motion.div
                key={r.id}
                initial={{ opacity: 0, y: 6 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: i * 0.05 }}
                className="card-elev flex items-center gap-4 p-4"
              >
                <div className="grid size-11 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary">
                  {r.name.split(" ").map((p) => p[0]).slice(0, 2).join("")}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="truncate text-sm font-semibold text-foreground">{r.name}</div>
                  <div className="truncate text-[11px] text-muted-foreground">{r.detail}</div>
                  <div className="text-[10px] text-muted-foreground">{r.when}</div>
                </div>
                <div className="flex shrink-0 gap-2">
                  <button
                    onClick={() => decide(r.id, "Accepted")}
                    className="inline-flex items-center gap-1 rounded-full bg-primary px-3 py-1.5 text-[11px] font-medium text-primary-foreground hover:brightness-110"
                  >
                    <Check className="size-3" /> Accept
                  </button>
                  <button
                    onClick={() => decide(r.id, "Declined")}
                    className="inline-flex items-center gap-1 rounded-full border border-border bg-card px-3 py-1.5 text-[11px] font-medium text-foreground hover:bg-muted"
                  >
                    <X className="size-3" /> Decline
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </section>

      <section>
        <div className="mb-3 flex items-baseline justify-between">
          <div className="label-eyebrow">Recently connected</div>
        </div>
        <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
          {accepted.map((r) => (
            <div key={r.id} className="flex items-center gap-3 px-4 py-3">
              <div className="grid size-10 shrink-0 place-items-center rounded-full bg-muted text-sm font-semibold text-muted-foreground">
                {r.name.split(" ").map((p) => p[0]).slice(0, 2).join("")}
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-medium text-foreground">{r.name}</div>
                <div className="truncate text-[11px] text-muted-foreground">{r.detail}</div>
              </div>
              <StatusPill status={r.status} />
            </div>
          ))}
        </div>
      </section>

      <section>
        <div className="mb-3 flex items-baseline justify-between">
          <div className="label-eyebrow">Suggested Old Starehians</div>
        </div>
        <div className="space-y-2">
          {suggestions.map((s, i) => (
            <motion.div
              key={s.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.1 + i * 0.05 }}
              className="card-elev flex items-center gap-4 p-4"
            >
              <div className="grid size-11 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary">
                {s.name.split(" ").map((p) => p[0]).slice(0, 2).join("")}
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-semibold text-foreground">{s.name}</div>
                <div className="truncate text-[11px] text-muted-foreground">{s.detail}</div>
                <div className="mt-0.5 inline-flex items-center gap-1 text-[10px] text-primary">
                  <Users className="size-3" /> {s.reason}
                </div>
              </div>
              <button className="inline-flex shrink-0 items-center gap-1 rounded-full border border-border bg-card px-3 py-1.5 text-[11px] font-medium text-foreground hover:bg-muted">
                <UserPlus className="size-3" /> Connect
              </button>
            </motion.div>
          ))}
        </div>
      </section>
    </div>
  );
}
