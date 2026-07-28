import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import { Users, MapPin, Briefcase, Search, UserPlus, X } from "lucide-react";
import { useState, useMemo } from "react";

export const Route = createFileRoute("/_app/accounts")({
  component: Directory,
});

interface Alumnus {
  id: string;
  full_name: string;
  house: string;
  class_year: number;
  career: string;
  location: string;
  avatar_url?: string;
}

function Directory() {
  const [q, setQ] = useState("");
  const [house, setHouse] = useState<string>("");
  const [decade, setDecade] = useState<string>("");

  const getToken = () => {
    if (typeof window === 'undefined') return '';
    return localStorage.getItem('oss_token') || '';
  };

  const authHeaders = () => ({
    'Content-Type': 'application/json',
    Authorization: `Bearer ${getToken()}`,
  });

  // Fetch alumni profiles from API
  const { data: alumni = [], isLoading } = useQuery({
    queryKey: ['alumni-directory'],
    queryFn: async () => {
      const res = await fetch('/api/profiles', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.profiles || [];
    },
    enabled: !!getToken(),
  });

  const houses = useMemo(() => Array.from(new Set(alumni.map((a) => a.house).filter(Boolean))).sort(), [alumni]);
  const decades = useMemo(() => {
    const ds = Array.from(new Set(alumni.map((a) => a.class_year ? `${Math.floor(a.class_year / 10) * 10}s` : null).filter(Boolean)));
    return ds.sort();
  }, [alumni]);

  const filtered = useMemo(() => {
    const term = q.trim().toLowerCase();
    return alumni.filter((a) => {
      if (house && a.house !== house) return false;
      if (decade && a.class_year && `${Math.floor(a.class_year / 10) * 10}s` !== decade) return false;
      if (!term) return true;
      return [a.full_name, a.house, String(a.class_year), a.career, a.location].some((f) =>
        f?.toLowerCase().includes(term),
      );
    });
  }, [q, house, decade, alumni]);

  const clearFilters = () => { setHouse(""); setDecade(""); setQ(""); };
  const anyFilter = q || house || decade;

  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Alumni directory</div>
        <h1 className="display mt-1 text-3xl font-black tracking-tight md:text-4xl">
          Find an Old Starehian.
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          {filtered.length.toLocaleString()} of {alumni.length.toLocaleString()} alumni · filter by house, class year & city.
        </p>
      </header>

      <div className="relative">
        <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Search by name, house, class year, city…"
          className="h-11 w-full rounded-full border border-border bg-card pl-10 pr-4 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring/40"
        />
      </div>

      {/* Filter chips */}
      <div className="flex flex-wrap items-center gap-2">
        <span className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">House</span>
        {houses.map((h) => (
          <button
            key={h}
            onClick={() => setHouse((cur) => (cur === h ? "" : h))}
            className={`rounded-full border px-3 py-1 text-[11px] font-semibold transition-colors ${
              house === h
                ? "border-primary bg-primary text-primary-foreground"
                : "border-border bg-card text-foreground hover:bg-muted"
            }`}
          >
            {h}
          </button>
        ))}
        <span className="ml-2 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">Class</span>
        {decades.map((d) => (
          <button
            key={d}
            onClick={() => setDecade((cur) => (cur === d ? "" : d))}
            className={`rounded-full border px-3 py-1 text-[11px] font-semibold transition-colors ${
              decade === d
                ? "border-secondary bg-secondary text-secondary-foreground"
                : "border-border bg-card text-foreground hover:bg-muted"
            }`}
          >
            {d}
          </button>
        ))}
        {anyFilter && (
          <button
            onClick={clearFilters}
            className="inline-flex items-center gap-1 rounded-full border border-border bg-muted px-3 py-1 text-[11px] font-semibold text-foreground hover:bg-muted/70"
          >
            <X className="size-3" /> Clear
          </button>
        )}
      </div>

      <div className="space-y-2">
        {isLoading ? (
          <div className="card-elev flex items-center justify-center p-10 text-center text-sm text-muted-foreground">
            Loading alumni directory...
          </div>
        ) : filtered.length === 0 ? (
          <div className="card-elev flex flex-col items-center justify-center gap-2 p-10 text-center text-sm text-muted-foreground">
            <Users className="size-5" />
            No Old Starehians match that search.
          </div>
        ) : (
          filtered.map((a, i) => (
            <motion.div
              key={a.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: i * 0.03 }}
              className="card-elev flex items-center gap-4 p-4"
            >
              <div className="grid size-12 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary">
                {a.avatar_url ? (
                  <img src={a.avatar_url} alt={a.full_name} className="h-full w-full rounded-full object-cover" />
                ) : (
                  a.full_name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()
                )}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <Link
                    to="/p/$id"
                    params={{ id: a.id }}
                    className="truncate text-sm font-bold text-foreground hover:text-primary"
                  >
                    {a.full_name}
                  </Link>
                  {a.house && (
                    <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
                      {a.house} · {a.class_year ? `'${String(a.class_year).slice(2)}` : ''}
                    </span>
                  )}
                </div>
                <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[11px] text-muted-foreground">
                  {a.career && <span className="inline-flex items-center gap-1"><Briefcase className="size-3" /> {a.career}</span>}
                  {a.location && <span className="inline-flex items-center gap-1"><MapPin className="size-3" /> {a.location}</span>}
                </div>
              </div>
              <Link
                to="/loans"
                className="inline-flex shrink-0 items-center gap-1 rounded-full border border-border bg-card px-3 py-1.5 text-[11px] font-medium text-foreground transition-colors hover:bg-muted"
              >
                <UserPlus className="size-3" /> Connect
              </Link>
            </motion.div>
          ))
        )}
      </div>
    </div>
  );
}
