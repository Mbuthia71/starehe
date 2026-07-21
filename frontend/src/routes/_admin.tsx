import { createFileRoute, Outlet, Link, useRouterState, useNavigate } from "@tanstack/react-router";
import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import { Home, Users, UserPlus, Activity, Settings, BarChart3, CheckCircle, LogOut, Bell } from "lucide-react";
import { isAuthenticated, getStoredMember, clearAuth, hasPermission } from "@/lib/auth";
import logoAsset from "@/assets/siohioma-logo.png?url";

export const Route = createFileRoute("/_admin")({
  component: AdminLayout,
});

const adminTabs = [
  { to: "/admin/dashboard", label: "Overview", icon: BarChart3, exact: true, permission: "view_all" },
  { to: "/admin/loans", label: "Connections", icon: UserPlus, permission: "approve" },
  { to: "/admin/members", label: "Alumni", icon: Users, permission: "view_all" },
  { to: "/admin/transactions", label: "Activity log", icon: Activity, permission: "read" },
  { to: "/admin/approvals", label: "Approvals", icon: CheckCircle, permission: "approve" },
] as const;

function AdminSideTab({ to, label, icon: Icon, exact, permission }: { to: string; label: string; icon: typeof Home; exact?: boolean; permission?: string; }) {
  const hasAccess = !permission || hasPermission(permission);
  if (!hasAccess) return null;
  return (
    <Link
      to={to}
      activeOptions={{ exact: !!exact }}
      className="group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground data-[status=active]:bg-primary/10 data-[status=active]:text-primary"
    >
      <Icon className="size-4" strokeWidth={1.75} />
      <span>{label}</span>
    </Link>
  );
}

function AdminLayout() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const navigate = useNavigate();
  const [booting, setBooting] = useState(true);
  const [member, setMember] = useState<ReturnType<typeof getStoredMember>>(null);

  useEffect(() => {
    setMember(getStoredMember());
    if (!isAuthenticated()) { navigate({ to: "/auth/login" }); return; }
    if (!hasPermission("view_all") && !hasPermission("approve")) { navigate({ to: "/" }); return; }
    const t = setTimeout(() => setBooting(false), 600);
    return () => clearTimeout(t);
  }, []);

  return (
    <div className="relative min-h-screen w-full bg-background lg:grid lg:grid-cols-[280px_minmax(0,1fr)]">
      {booting && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-background">
          <div className="text-center">
            <img src={logoAsset} alt="OSS" className="mx-auto h-12 w-auto animate-pulse" />
            <div className="mt-4 text-sm text-muted-foreground">Loading OSS admin…</div>
          </div>
        </div>
      )}

      <aside className="hidden lg:flex lg:sticky lg:top-0 lg:h-screen lg:flex-col lg:border-r lg:border-border/60 lg:bg-card/40 lg:px-5 lg:py-6">
        <Link to="/admin/dashboard" className="flex items-center gap-2 px-2">
          <img src={logoAsset} alt="OSS" className="h-8 w-auto" />
          <span className="text-sm font-semibold">Admin</span>
        </Link>
        <nav className="mt-10 flex flex-col gap-1">
          {adminTabs.map((t) => (<AdminSideTab key={t.to} {...t} />))}
        </nav>
        <div className="mt-auto rounded-2xl border border-border/60 bg-background/60 p-4">
          <div className="text-xs text-muted-foreground">Signed in as</div>
          <div className="mt-1 truncate text-sm font-semibold">{member?.displayName || "Admin"}</div>
          <div className="mt-1 text-xs text-muted-foreground">{member?.role || "Council"}</div>
          <button onClick={() => { clearAuth(); navigate({ to: "/auth/login" }); }} className="mt-2 flex items-center gap-2 text-xs text-danger hover:underline">
            <LogOut className="size-3" /> Sign out
          </button>
        </div>
      </aside>

      <div className="mx-auto flex w-full max-w-7xl flex-col px-5 py-6 lg:px-8">
        <header className="flex items-center justify-between gap-4 border-b border-border/60 pb-4">
          <div>
            <div className="text-xs uppercase tracking-[0.2em] text-muted-foreground">Old Starehian Society</div>
            <h1 className="text-2xl font-semibold tracking-tight">Admin console</h1>
          </div>
          <div className="flex items-center gap-3">
            <button className="relative grid size-10 place-items-center rounded-full border border-border bg-card text-muted-foreground active:scale-95">
              <Bell className="size-4" />
              <span className="absolute right-2.5 top-2.5 size-1.5 rounded-full bg-warning" />
            </button>
          </div>
        </header>

        <main className="flex-1 py-6">
          <AnimatePresence mode="wait">
            <motion.div
              key={pathname}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.28, ease: [0.22, 1, 0.36, 1] }}
            >
              <Outlet />
            </motion.div>
          </AnimatePresence>
        </main>
      </div>
    </div>
  );
}
