import { createFileRoute, Outlet, Link, useRouterState, useNavigate } from "@tanstack/react-router";
import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import {
  Home, Users, UserPlus, Building2, Activity, User, Bell, AlertTriangle,
  ChevronsLeft, ChevronsRight, LogOut, Search, Sun, Moon, Menu, X,
  Briefcase, FileText, GraduationCap, Tag, Heart, Store,
} from "lucide-react";
import { SiohiomaLoader } from "@/components/SiohiomaLoader";
import { isAuthenticated, getStoredMember, clearAuth, updateActivity, startSessionTimeoutCheck, detectDeviceSecurity, trustCurrentDevice, isTrustedDevice, type AuthMember } from "@/lib/auth";
import logoAsset from "@/assets/siohioma-logo.png?url";

export const Route = createFileRoute("/_app")({
  component: AppLayout,
});

const tabs = [
  { to: "/", label: "Home", icon: Home, exact: true },
  { to: "/accounts", label: "Directory", icon: Users },
  { to: "/jobs", label: "Jobs", icon: Briefcase },
  { to: "/tenders", label: "Tenders", icon: FileText },
  { to: "/class-groups", label: "Class Groups", icon: GraduationCap },
  { to: "/offers", label: "Offers", icon: Tag },
  { to: "/sponsorships", label: "Sponsorships", icon: Heart },
  { to: "/business", label: "Business", icon: Store },
  { to: "/groups", label: "Chapters", icon: Building2 },
  { to: "/profile", label: "Profile", icon: User },
] as const;

type SidebarMode = "full" | "rail";

function NavItem({
  to, label, icon: Icon, exact, compact, onNavigate,
}: { to: string; label: string; icon: typeof Home; exact?: boolean; compact?: boolean; onNavigate?: () => void }) {
  return (
    <Link
      to={to}
      activeOptions={{ exact: !!exact }}
      onClick={onNavigate}
      title={compact ? label : undefined}
      className={`group relative flex items-center gap-3 rounded-xl text-sm font-medium text-sidebar-foreground/70 transition-all hover:bg-sidebar-accent hover:text-sidebar-foreground data-[status=active]:bg-sidebar-accent data-[status=active]:text-sidebar-foreground ${
        compact ? "h-11 w-11 justify-center" : "px-3 py-2.5"
      }`}
    >
      <span className="absolute left-0 top-1/2 h-0 w-[3px] -translate-y-1/2 rounded-r-full bg-sidebar-primary transition-all group-data-[status=active]:h-6" />
      <Icon className="size-[18px] shrink-0" strokeWidth={1.75} />
      {!compact && <span className="truncate">{label}</span>}
    </Link>
  );
}

function SidebarBody({
  mode, onToggleMode, member, onSignOut, theme, onToggleTheme, onNavigate, showToggle = true,
}: {
  mode: SidebarMode;
  onToggleMode: () => void;
  member: AuthMember | null;
  onSignOut: () => void;
  theme: "light" | "dark";
  onToggleTheme: () => void;
  onNavigate?: () => void;
  showToggle?: boolean;
}) {
  const compact = mode === "rail";
  const initials = (member?.displayName ?? "OS").split(" ").map((p) => p[0]).slice(0, 2).join("");

  return (
    <div className="relative flex h-full flex-col bg-sidebar text-sidebar-foreground">
      {/* Starehe stripe accent */}
      <div className="starehe-stripe absolute inset-x-0 top-0 h-[3px] opacity-90" />

      {/* Header */}
      <div className={`flex items-center gap-3 px-4 pb-4 pt-6 ${compact ? "justify-center px-2" : ""}`}>
        <Link to="/" onClick={onNavigate} className="flex items-center gap-2.5">
          <div className="grid size-9 shrink-0 place-items-center rounded-xl bg-white/10 ring-1 ring-white/10">
            <img src={logoAsset} alt="OSS" className="h-6 w-auto" />
          </div>
          {!compact && (
            <div className="min-w-0">
              <div className="text-[10px] font-semibold uppercase tracking-[0.2em] text-sidebar-foreground/50">
                Old Starehian
              </div>
              <div className="display truncate text-sm font-semibold">Society</div>
            </div>
          )}
        </Link>
        {showToggle && !compact && (
          <button
            onClick={onToggleMode}
            aria-label="Collapse sidebar"
            className="ml-auto grid size-7 place-items-center rounded-md text-sidebar-foreground/60 hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <ChevronsLeft className="size-4" />
          </button>
        )}
      </div>

      {/* Search (hidden in rail) */}
      {!compact && (
        <div className="px-3 pb-3">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-1/2 size-3.5 -translate-y-1/2 text-sidebar-foreground/40" />
            <input
              placeholder="Search alumni…"
              className="h-9 w-full rounded-lg border border-white/10 bg-white/[0.04] pl-9 pr-3 text-xs text-sidebar-foreground placeholder:text-sidebar-foreground/40 focus:outline-none focus:ring-2 focus:ring-sidebar-primary/40"
            />
          </div>
        </div>
      )}

      {/* Menu label */}
      {!compact && (
        <div className="px-4 pb-1 pt-2 text-[10px] font-semibold uppercase tracking-[0.2em] text-sidebar-foreground/40">
          Menu
        </div>
      )}

      {/* Primary nav */}
      <nav className={`flex flex-col gap-1 ${compact ? "items-center px-2" : "px-2"}`}>
        {tabs.map((t) => (
          <NavItem key={t.to} {...t} compact={compact} onNavigate={onNavigate} />
        ))}
      </nav>

      {/* Expand button in rail */}
      {showToggle && compact && (
        <button
          onClick={onToggleMode}
          aria-label="Expand sidebar"
          className="mx-auto mt-3 grid size-8 place-items-center rounded-md text-sidebar-foreground/60 hover:bg-sidebar-accent hover:text-sidebar-foreground"
        >
          <ChevronsRight className="size-4" />
        </button>
      )}

      <div className="mt-auto flex flex-col gap-3 px-3 pb-4">
        {/* Theme toggle */}
        {!compact ? (
          <div className="flex items-center rounded-xl bg-white/[0.05] p-1 ring-1 ring-white/5">
            <button
              onClick={() => theme !== "light" && onToggleTheme()}
              className={`flex flex-1 items-center justify-center gap-1.5 rounded-lg py-1.5 text-[11px] font-semibold transition-colors ${
                theme === "light" ? "bg-white/10 text-sidebar-foreground" : "text-sidebar-foreground/50"
              }`}
            >
              <Sun className="size-3.5" /> Light
            </button>
            <button
              onClick={() => theme !== "dark" && onToggleTheme()}
              className={`flex flex-1 items-center justify-center gap-1.5 rounded-lg py-1.5 text-[11px] font-semibold transition-colors ${
                theme === "dark" ? "bg-white/10 text-sidebar-foreground" : "text-sidebar-foreground/50"
              }`}
            >
              <Moon className="size-3.5" /> Dark
            </button>
          </div>
        ) : (
          <button
            onClick={onToggleTheme}
            aria-label="Toggle theme"
            className="mx-auto grid size-9 place-items-center rounded-lg bg-white/[0.05] text-sidebar-foreground/70 ring-1 ring-white/5 hover:bg-white/10"
          >
            {theme === "dark" ? <Sun className="size-4" /> : <Moon className="size-4" />}
          </button>
        )}

        {/* Member card */}
        <div className={`flex items-center gap-3 rounded-xl bg-white/[0.04] p-2 ring-1 ring-white/5 ${compact ? "justify-center p-2" : ""}`}>
          <div className="grid size-9 shrink-0 place-items-center rounded-full bg-primary text-primary-foreground text-[11px] font-semibold">
            {initials}
          </div>
          {!compact && (
            <>
              <div className="min-w-0 flex-1">
                <div className="truncate text-xs font-semibold">{member?.displayName ?? "Old Starehian"}</div>
                <div className="truncate text-[10px] text-sidebar-foreground/50">{member?.emailAddress ?? "OSS alumni"}</div>
              </div>
              <button
                onClick={onSignOut}
                aria-label="Sign out"
                className="grid size-8 shrink-0 place-items-center rounded-md text-sidebar-foreground/50 hover:bg-white/10 hover:text-sidebar-foreground"
              >
                <LogOut className="size-4" />
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function AppLayout() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const navigate = useNavigate();
  const [booting, setBooting] = useState(true);
  const [showSecurityAlert, setShowSecurityAlert] = useState(false);
  const [securityMessage, setSecurityMessage] = useState("");
  const [member, setMember] = useState<AuthMember | null>(null);
  const [mode, setMode] = useState<SidebarMode>("full");
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [theme, setTheme] = useState<"light" | "dark">("light");

  useEffect(() => {
    const stored = typeof window !== "undefined" ? localStorage.getItem("oss.sidebar") : null;
    if (stored === "rail" || stored === "full") setMode(stored);
    const t = typeof window !== "undefined" ? localStorage.getItem("oss.theme") : null;
    if (t === "dark" || t === "light") {
      setTheme(t);
      document.documentElement.classList.toggle("dark", t === "dark");
    }
  }, []);

  const toggleMode = () => {
    const next: SidebarMode = mode === "full" ? "rail" : "full";
    setMode(next);
    try { localStorage.setItem("oss.sidebar", next); } catch {}
  };
  const toggleTheme = () => {
    const next = theme === "dark" ? "light" : "dark";
    setTheme(next);
    document.documentElement.classList.toggle("dark", next === "dark");
    try { localStorage.setItem("oss.theme", next); } catch {}
  };
  const signOut = () => { clearAuth(); navigate({ to: "/auth/login" }); };

  useEffect(() => {
    setMember(getStoredMember());
    if (!isAuthenticated()) {
      navigate({ to: "/auth/login" });
      return;
    }

    const security = detectDeviceSecurity();
    if (security.isRooted || security.isJailbroken) {
      setSecurityMessage("Your device appears to be modified. For your security, please use a standard device.");
      setShowSecurityAlert(true);
    } else if (!isTrustedDevice()) {
      setSecurityMessage("New device detected. Trust this device to keep it signed in?");
      setShowSecurityAlert(true);
    }

    const timeoutInterval = startSessionTimeoutCheck(() => {
      clearAuth();
      navigate({ to: "/auth/login" });
    });

    const handleActivity = () => updateActivity();
    window.addEventListener('mousemove', handleActivity);
    window.addEventListener('keydown', handleActivity);
    window.addEventListener('touchstart', handleActivity);

    const t = setTimeout(() => setBooting(false), 1400);
    return () => {
      clearTimeout(t);
      clearInterval(timeoutInterval);
      window.removeEventListener('mousemove', handleActivity);
      window.removeEventListener('keydown', handleActivity);
      window.removeEventListener('touchstart', handleActivity);
    };
  }, []);

  // close drawer on route change
  useEffect(() => { setDrawerOpen(false); }, [pathname]);

  const sidebarWidth = mode === "full" ? "260px" : "68px";

  return (
    <div
      className="relative min-h-screen w-full bg-background md:grid"
      style={{ gridTemplateColumns: `${sidebarWidth} minmax(0,1fr)` }}
    >
      <AnimatePresence>
        {booting && (
          <motion.div
            key="boot"
            initial={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
            className="fixed inset-0 z-[100]"
          >
            <SiohiomaLoader surface="cream" label={`Karibu, ${member?.firstName || member?.displayName || "Old Starehian"}`} />
          </motion.div>
        )}
      </AnimatePresence>

      <AnimatePresence>
        {showSecurityAlert && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4"
          >
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="w-full max-w-md rounded-2xl bg-background p-6 shadow-lg"
            >
              <div className="flex items-start gap-4">
                <div className="rounded-full bg-warning/20 p-2">
                  <AlertTriangle className="size-5 text-warning" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-semibold">Security check</h3>
                  <p className="mt-2 text-sm text-muted-foreground">{securityMessage}</p>
                  <div className="mt-4 flex gap-2">
                    <button
                      onClick={() => { trustCurrentDevice(); setShowSecurityAlert(false); }}
                      className="flex-1 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:brightness-110"
                    >
                      Trust device
                    </button>
                    <button
                      onClick={() => { clearAuth(); navigate({ to: "/auth/login" }); }}
                      className="flex-1 rounded-lg border border-input bg-background px-4 py-2 text-sm font-medium text-foreground hover:bg-muted"
                    >
                      Sign out
                    </button>
                  </div>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Desktop / tablet sidebar */}
      <aside
        className="sticky top-0 hidden h-screen border-r border-sidebar-border md:block"
        style={{ width: sidebarWidth }}
      >
        <SidebarBody
          mode={mode}
          onToggleMode={toggleMode}
          member={member}
          onSignOut={signOut}
          theme={theme}
          onToggleTheme={toggleTheme}
        />
      </aside>

      {/* Mobile off-canvas drawer */}
      <AnimatePresence>
        {drawerOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setDrawerOpen(false)}
              className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm md:hidden"
            />
            <motion.aside
              initial={{ x: "-100%" }}
              animate={{ x: 0 }}
              exit={{ x: "-100%" }}
              transition={{ type: "spring", damping: 26, stiffness: 260 }}
              className="fixed inset-y-0 left-0 z-50 w-[280px] md:hidden"
            >
              <SidebarBody
                mode="full"
                onToggleMode={() => setDrawerOpen(false)}
                member={member}
                onSignOut={signOut}
                theme={theme}
                onToggleTheme={toggleTheme}
                onNavigate={() => setDrawerOpen(false)}
                showToggle={false}
              />
              <button
                onClick={() => setDrawerOpen(false)}
                aria-label="Close menu"
                className="absolute right-3 top-5 grid size-8 place-items-center rounded-md text-sidebar-foreground/70 hover:bg-white/10"
              >
                <X className="size-4" />
              </button>
            </motion.aside>
          </>
        )}
      </AnimatePresence>

      <div className="flex w-full flex-col">
        <header className="sticky top-0 z-30 flex items-center justify-between gap-3 border-b border-border/60 bg-background/90 px-4 py-3 backdrop-blur md:px-8 md:py-5">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setDrawerOpen(true)}
              aria-label="Open menu"
              className="grid size-10 place-items-center rounded-xl border border-border bg-card text-foreground md:hidden"
            >
              <Menu className="size-4" />
            </button>
            <Link to="/" className="flex items-center gap-2 md:hidden">
              <img src={logoAsset} alt="Old Starehian Society" className="h-7 w-auto" />
            </Link>
            <div className="hidden md:block">
              <div className="text-[10px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
                Alumni network
              </div>
              <div className="display text-sm font-semibold">Old Starehian Society</div>
            </div>
          </div>
          <button
            aria-label="Notifications"
            className="relative grid size-10 place-items-center rounded-full border border-border bg-card text-muted-foreground active:scale-95"
          >
            <Bell className="size-4" />
            <span className="absolute right-2.5 top-2.5 size-1.5 rounded-full bg-primary" />
          </button>
        </header>

        <main className="mx-auto w-full max-w-[860px] flex-1 px-5 pb-16 pt-6 md:px-8">
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
