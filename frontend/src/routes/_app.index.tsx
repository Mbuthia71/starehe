import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Users,
  UserPlus,
  Building2,
  Activity,
  Calendar,
  MessageSquare,
  ChevronRight,
  ArrowUpRight,
} from "lucide-react";
import { getStoredMember } from "@/lib/auth";
import heroSplash from "@/assets/hero-splash.jpg?url";

export const Route = createFileRoute("/_app/")({
  component: Home,
});

const quickActions = [
  { to: "/accounts", label: "Directory", icon: Users },
  { to: "/loans", label: "Connections", icon: UserPlus },
  { to: "/groups", label: "Chapters", icon: Building2 },
  { to: "/statements", label: "Activity", icon: Activity },
] as const;

const upcomingEvents = [
  { title: "Class of 2008 · 20-year reunion planning call", when: "Thu 12 Nov · 7:00 pm EAT" },
  { title: "OSS Nairobi chapter dinner", when: "Sat 21 Nov · Muthaiga Country Club" },
  { title: "Mentorship drive · Form 4 leavers", when: "Sat 05 Dec · Starehe Boys' Centre" },
];

const recentActivity = [
  { who: "Kevin Mwangi (Griffin, '08)", what: "sent you a connection request", when: "12 min ago", tone: "primary" as const },
  { who: "Class of 2010", what: "posted a chapter announcement", when: "1 h ago", tone: "muted" as const },
  { who: "Prof. Wanjiku (Faculty)", what: "commented on your post", when: "3 h ago", tone: "muted" as const },
  { who: "Brian Otieno (Livingstone, '11)", what: "endorsed you for mentorship", when: "Yesterday", tone: "muted" as const },
  { who: "OSS Secretariat", what: "shared the November newsletter", when: "2 d ago", tone: "muted" as const },
];

function Home() {
  const storedMember = getStoredMember();
  const greeting = (() => {
    const h = new Date().getHours();
    if (h < 12) return "Good morning";
    if (h < 17) return "Good afternoon";
    return "Good evening";
  })();

  const firstName = storedMember?.firstName || storedMember?.displayName?.split(" ")?.[0] || "Old Starehian";

  return (
    <div className="space-y-8">
      <div>
        <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
          <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
          {greeting}
        </div>
        <h1 className="display mt-2 text-4xl font-semibold tracking-tight md:text-5xl">
          {firstName}.
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          <span className="font-semibold text-foreground">Nec aspera terrent.</span> Difficulties be damned — welcome back to the brotherhood.
        </p>
      </div>

      {/* Dramatic hero — Starehe crest of colour */}
      <motion.section
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
        className="relative overflow-hidden rounded-3xl border border-white/10 p-8 text-white shadow-[0_40px_80px_-40px_oklch(0.15_0.06_265/0.8)]"
      >
        {/* Crimson × navy paint-splash background */}
        <motion.img
          aria-hidden
          src={heroSplash}
          alt=""
          initial={{ scale: 1.1, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 1.4, ease: [0.22, 1, 0.36, 1] }}
          className="pointer-events-none absolute inset-0 h-full w-full object-cover"
        />
        {/* Dark legibility scrim */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              "linear-gradient(135deg, rgba(8,8,20,0.55) 0%, rgba(8,8,20,0.15) 45%, rgba(8,8,20,0.75) 100%)",
          }}
        />
        <div className="starehe-stripe absolute inset-x-0 top-0 z-10 h-[3px]" />

        <div className="relative">
          <div className="flex items-center justify-between">
            <div className="text-[10px] font-bold uppercase tracking-[0.28em] text-white/75">
              Your Starehe network
            </div>
            <span className="glass-panel rounded-full px-2.5 py-1 text-[10px] font-bold uppercase tracking-widest text-white/85">
              Griffin · '08
            </span>
          </div>

          {/* Glass content block over the splash */}
          <div className="glass-panel mt-6 flex items-end gap-6 rounded-2xl p-5">
            <div>
              <div className="num display text-6xl font-black tracking-tight text-white md:text-7xl">184</div>
              <div className="mt-1 text-xs font-medium text-white/70">
                {storedMember?.officeName ?? "Class of 2008 · Griffin House"}
              </div>
            </div>
            <div className="mb-1 flex-1 border-l border-white/20 pl-6">
              <div className="text-[10px] font-bold uppercase tracking-[0.22em] text-white/60">This week</div>
              <div className="mt-1 text-2xl font-black text-white">+3</div>
              <div className="text-[11px] text-white/60">new connections</div>
            </div>
          </div>

          <div className="mt-8 flex flex-wrap items-center gap-2">
            <Link
              to="/accounts"
              className="inline-flex items-center gap-1.5 rounded-full bg-white text-slate-900 px-4 py-2 text-xs font-bold uppercase tracking-wider hover:bg-white/90"
            >
              Open directory <ArrowUpRight className="size-3.5" />
            </Link>
            <Link
              to="/groups"
              className="glass-panel inline-flex items-center gap-1.5 rounded-full px-4 py-2 text-xs font-bold uppercase tracking-wider text-white hover:bg-white/15"
            >
              Find your chapter
            </Link>
          </div>
        </div>
      </motion.section>

      {/* Quick actions */}
      <section>
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {quickActions.map((a, i) => (
            <motion.div
              key={a.to}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.35, delay: 0.05 * i, ease: [0.22, 1, 0.36, 1] }}
            >
              <Link
                to={a.to}
                className="group flex h-full flex-col items-start gap-3 rounded-2xl border border-border/60 bg-card p-4 transition-all hover:-translate-y-0.5 hover:border-primary/30 hover:shadow-[0_18px_40px_-24px_oklch(0.52_0.20_25/0.35)]"
              >
                <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary ring-1 ring-primary/15 transition-transform group-hover:scale-105">
                  <a.icon className="size-5" strokeWidth={1.75} />
                </div>
                <span className="text-sm font-semibold leading-tight text-foreground">
                  {a.label}
                </span>
              </Link>
            </motion.div>
          ))}
        </div>
      </section>

      {/* Pending connection */}
      <Link to="/loans" className="block">
        <motion.section
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.15 }}
          className="card-elev flex items-center justify-between gap-4 p-5"
        >
          <div className="min-w-0">
            <div className="label-eyebrow">Pending connection</div>
            <div className="display mt-1 text-xl font-semibold">
              Kevin Mwangi
            </div>
            <div className="mt-1 truncate text-xs text-muted-foreground">
              Griffin House · Class of 2008 · wants to connect
            </div>
          </div>
          <div className="grid size-10 shrink-0 place-items-center rounded-full bg-muted text-muted-foreground">
            <ChevronRight className="size-4" />
          </div>
        </motion.section>
      </Link>

      {/* Upcoming events */}
      <section>
        <div className="mb-3 flex items-baseline justify-between">
          <div className="label-eyebrow">Upcoming</div>
          <Link to="/groups" className="text-xs font-medium text-primary">
            All chapters
          </Link>
        </div>
        <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
          {upcomingEvents.map((e) => (
            <div key={e.title} className="flex items-center gap-3 px-4 py-3">
              <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary/10 text-primary">
                <Calendar className="size-4" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-medium text-foreground">{e.title}</div>
                <div className="text-[11px] text-muted-foreground">{e.when}</div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Recent activity */}
      <section>
        <div className="mb-3 flex items-baseline justify-between">
          <div className="label-eyebrow">Recent activity</div>
          <Link to="/statements" className="text-xs font-medium text-primary">
            See all
          </Link>
        </div>
        <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
          {recentActivity.map((r, i) => (
            <div key={i} className="flex items-center gap-3 px-4 py-3">
              <div
                className={`grid size-10 shrink-0 place-items-center rounded-full ${
                  r.tone === "primary" ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"
                }`}
              >
                <MessageSquare className="size-4" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-medium text-foreground">
                  <span className="font-semibold">{r.who}</span> {r.what}
                </div>
                <div className="text-[11px] text-muted-foreground">{r.when}</div>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
