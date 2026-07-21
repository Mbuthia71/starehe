import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Users, MapPin, Briefcase, Search, UserPlus, X } from "lucide-react";
import { useState, useMemo } from "react";

export const Route = createFileRoute("/_app/accounts")({
  component: Directory,
});

interface Alumnus {
  id: string;
  name: string;
  house: string;
  classYear: number;
  role: string;
  city: string;
  initials: string;
}

const alumni: Alumnus[] = [
  { id: "1", name: "Kevin Mwangi", house: "Griffin", classYear: 2008, role: "Product Manager · Safaricom", city: "Nairobi", initials: "KM" },
  { id: "2", name: "Brian Otieno", house: "Livingstone", classYear: 2011, role: "Senior Engineer · Andela", city: "Kampala", initials: "BO" },
  { id: "3", name: "Achieng Nyongo", house: "Griffin", classYear: 2013, role: "Doctor · Aga Khan Hospital", city: "Nairobi", initials: "AN" },
  { id: "4", name: "James Kariuki", house: "Livingstone", classYear: 2005, role: "Founder · Kilimanjaro Ventures", city: "Dar es Salaam", initials: "JK" },
  { id: "5", name: "Mary Wanjiku", house: "Griffin", classYear: 2015, role: "Lawyer · TripleOKLaw", city: "Mombasa", initials: "MW" },
  { id: "6", name: "Peter Kimani", house: "Livingstone", classYear: 2009, role: "Data Scientist · IBM", city: "London", initials: "PK" },
  { id: "7", name: "Sarah Njeri", house: "Griffin", classYear: 2012, role: "Journalist · NTV", city: "Nairobi", initials: "SN" },
  { id: "8", name: "David Ochieng", house: "Livingstone", classYear: 2007, role: "Architect · Planning Systems", city: "Kigali", initials: "DO" },
];

function Directory() {
  const [q, setQ] = useState("");
  const [house, setHouse] = useState<string>("");
  const [decade, setDecade] = useState<string>("");

  const houses = useMemo(() => Array.from(new Set(alumni.map((a) => a.house))).sort(), []);
  const decades = useMemo(() => {
    const ds = Array.from(new Set(alumni.map((a) => `${Math.floor(a.classYear / 10) * 10}s`)));
    return ds.sort();
  }, []);

  const filtered = useMemo(() => {
    const term = q.trim().toLowerCase();
    return alumni.filter((a) => {
      if (house && a.house !== house) return false;
      if (decade && `${Math.floor(a.classYear / 10) * 10}s` !== decade) return false;
      if (!term) return true;
      return [a.name, a.house, String(a.classYear), a.role, a.city].some((f) =>
        f.toLowerCase().includes(term),
      );
    });
  }, [q, house, decade]);

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
          {filtered.length.toLocaleString()} of 10,000+ alumni · filter by house, class year & city.
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
        {filtered.map((a, i) => (
          <motion.div
            key={a.id}
            initial={{ opacity: 0, y: 6 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: i * 0.03 }}
            className="card-elev flex items-center gap-4 p-4"
          >
            <div className="grid size-12 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary">
              {a.initials}
            </div>
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <Link
                  to="/p/$id"
                  params={{ id: a.id }}
                  className="truncate text-sm font-bold text-foreground hover:text-primary"
                >
                  {a.name}
                </Link>
                <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
                  {a.house} · '{String(a.classYear).slice(2)}
                </span>
              </div>
              <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[11px] text-muted-foreground">
                <span className="inline-flex items-center gap-1"><Briefcase className="size-3" /> {a.role}</span>
                <span className="inline-flex items-center gap-1"><MapPin className="size-3" /> {a.city}</span>
              </div>
            </div>
            <Link
              to="/loans"
              className="inline-flex shrink-0 items-center gap-1 rounded-full border border-border bg-card px-3 py-1.5 text-[11px] font-medium text-foreground transition-colors hover:bg-muted"
            >
              <UserPlus className="size-3" /> Connect
            </Link>
          </motion.div>
        ))}
        {filtered.length === 0 && (
          <div className="card-elev flex flex-col items-center justify-center gap-2 p-10 text-center text-sm text-muted-foreground">
            <Users className="size-5" />
            No Old Starehians match that search.
          </div>
        )}
      </div>
    </div>
  );
}
