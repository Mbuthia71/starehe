import { createFileRoute, Link, useParams } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { ArrowLeft, BadgeCheck, Briefcase, GraduationCap, Linkedin, MapPin, Twitter, Instagram, Globe, UserPlus } from "lucide-react";
import logoAsset from "@/assets/siohioma-logo.png?url";

export const Route = createFileRoute("/p/$id")({
  component: PublicProfile,
  head: ({ params }) => ({
    meta: [
      { title: `Old Starehian · Member ${params.id}` },
      { name: "description", content: "Public member profile on the Old Starehian Society alumni network." },
    ],
  }),
});

// Minimal directory (matches _app.accounts.tsx). Later this will hit the Go API.
const alumni = [
  { id: "1", name: "Kevin Mwangi", house: "Griffin", classYear: 2008, role: "Product Manager · Safaricom", city: "Nairobi", initials: "KM" },
  { id: "2", name: "Brian Otieno", house: "Livingstone", classYear: 2011, role: "Senior Engineer · Andela", city: "Kampala", initials: "BO" },
  { id: "3", name: "Achieng Nyongo", house: "Griffin", classYear: 2013, role: "Doctor · Aga Khan Hospital", city: "Nairobi", initials: "AN" },
  { id: "4", name: "James Kariuki", house: "Livingstone", classYear: 2005, role: "Founder · Kilimanjaro Ventures", city: "Dar es Salaam", initials: "JK" },
  { id: "5", name: "Mary Wanjiku", house: "Griffin", classYear: 2015, role: "Lawyer · TripleOKLaw", city: "Mombasa", initials: "MW" },
  { id: "6", name: "Peter Kimani", house: "Livingstone", classYear: 2009, role: "Data Scientist · IBM", city: "London", initials: "PK" },
  { id: "7", name: "Sarah Njeri", house: "Griffin", classYear: 2012, role: "Journalist · NTV", city: "Nairobi", initials: "SN" },
  { id: "8", name: "David Ochieng", house: "Livingstone", classYear: 2007, role: "Architect · Planning Systems", city: "Kigali", initials: "DO" },
];

function PublicProfile() {
  const { id } = useParams({ from: "/p/$id" });
  const a = alumni.find((x) => x.id === id);

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border/60 bg-background/90 backdrop-blur">
        <div className="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
          <Link to="/" className="flex items-center gap-2">
            <div className="grid size-9 place-items-center rounded-xl bg-anchor">
              <img src={logoAsset} alt="OSS" className="h-6 w-auto" />
            </div>
            <div>
              <div className="text-[10px] font-bold uppercase tracking-[0.22em] text-muted-foreground">Old Starehian</div>
              <div className="display text-sm font-black">Society</div>
            </div>
          </Link>
          <Link to="/accounts" className="inline-flex items-center gap-1 text-[11px] font-semibold text-muted-foreground hover:text-foreground">
            <ArrowLeft className="size-3" /> Directory
          </Link>
        </div>
      </header>

      <main className="mx-auto max-w-4xl px-6 py-10">
        {!a ? (
          <div className="card-elev p-10 text-center">
            <div className="display text-2xl font-black">Member not found</div>
            <p className="mt-2 text-sm text-muted-foreground">This profile may be private or the link is invalid.</p>
            <Link to="/accounts" className="mt-6 inline-flex items-center gap-1 rounded-full bg-primary px-4 py-2 text-xs font-bold text-primary-foreground">
              Browse directory
            </Link>
          </div>
        ) : (
          <motion.article
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="grid gap-6 md:grid-cols-[300px_minmax(0,1fr)]"
          >
            <div className="surface-splash relative overflow-hidden rounded-2xl border border-white/10 p-4 text-anchor-foreground">
              <div className="starehe-stripe absolute inset-x-0 top-0 h-[3px]" />
              <div className="grid aspect-[4/5] w-full place-items-center rounded-xl bg-gradient-to-br from-primary/40 via-secondary/40 to-[oklch(0.70_0.18_45/0.5)] text-5xl font-black text-white">
                {a.initials}
              </div>
              <div className="mt-3 flex items-center gap-1.5">
                <div className="display truncate text-base font-black">{a.name}</div>
                <BadgeCheck className="size-4 shrink-0 text-primary" />
              </div>
              <div className="mt-0.5 text-[11px] text-anchor-foreground/60">{a.house} House · Class of {a.classYear}</div>
              <button className="mt-3 inline-flex w-full items-center justify-center gap-1.5 rounded-xl bg-white py-2.5 text-xs font-bold uppercase tracking-wider text-slate-900 hover:bg-white/90">
                <UserPlus className="size-3.5" /> Request to connect
              </button>
            </div>

            <div className="space-y-4">
              <div className="card-elev p-5">
                <div className="label-eyebrow">About</div>
                <p className="mt-3 text-sm leading-relaxed text-foreground">
                  Old Starehian, {a.house} House, Class of {a.classYear}. Currently working as {a.role.toLowerCase()} — based in {a.city}.
                </p>
                <div className="mt-4 grid grid-cols-1 gap-2 text-[12px] sm:grid-cols-2">
                  <span className="inline-flex items-center gap-2 text-muted-foreground">
                    <GraduationCap className="size-3.5" /> Starehe Boys' Centre · '{String(a.classYear).slice(2)}
                  </span>
                  <span className="inline-flex items-center gap-2 text-muted-foreground">
                    <Briefcase className="size-3.5" /> {a.role}
                  </span>
                  <span className="inline-flex items-center gap-2 text-muted-foreground">
                    <MapPin className="size-3.5" /> {a.city}
                  </span>
                </div>
              </div>

              <div className="card-elev p-5">
                <div className="label-eyebrow">Reach out</div>
                <div className="mt-3 flex flex-wrap gap-2 text-[12px]">
                  <SocialChip icon={Linkedin} label="LinkedIn" />
                  <SocialChip icon={Twitter} label="Twitter" />
                  <SocialChip icon={Instagram} label="Instagram" />
                  <SocialChip icon={Globe} label="Website" />
                </div>
                <p className="mt-3 text-[11px] text-muted-foreground">
                  Contact details are only shown to signed-in Old Starehians the member has allowed.
                </p>
              </div>
            </div>
          </motion.article>
        )}
      </main>
    </div>
  );
}

function SocialChip({ icon: Icon, label }: { icon: typeof Linkedin; label: string }) {
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full border border-border bg-card px-3 py-1 font-semibold text-muted-foreground">
      <Icon className="size-3" /> {label}
    </span>
  );
}