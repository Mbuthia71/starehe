import { createFileRoute } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Building2, Users, MapPin, Calendar } from "lucide-react";

export const Route = createFileRoute("/_app/groups")({
  component: Chapters,
});

const chapters = [
  { id: "nbo", name: "OSS Nairobi Chapter", members: 2140, city: "Nairobi", nextEvent: "Chapter dinner · Sat 21 Nov" },
  { id: "mba", name: "OSS Mombasa Chapter", members: 312, city: "Mombasa", nextEvent: "Beach clean-up · Sat 05 Dec" },
  { id: "lon", name: "OSS London Chapter", members: 187, city: "London, UK", nextEvent: "Networking mixer · Fri 27 Nov" },
  { id: "yr08", name: "Class of 2008", members: 246, city: "Global", nextEvent: "20-year reunion planning call" },
  { id: "yr13", name: "Class of 2013", members: 231, city: "Global", nextEvent: "Career panel · online" },
  { id: "grf", name: "Griffin House Alumni", members: 1120, city: "Global", nextEvent: "House quiz night · Dec 12" },
];

function Chapters() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Alumni chapters</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">
          Find your circle.
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Chapters by city, class year and house. Join one to see events and posts.
        </p>
      </header>

      <div className="grid gap-3 md:grid-cols-2">
        {chapters.map((c, i) => (
          <motion.div
            key={c.id}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, delay: i * 0.04 }}
            className="card-elev p-5"
          >
            <div className="flex items-start gap-3">
              <div className="grid size-11 shrink-0 place-items-center rounded-2xl bg-primary/10 text-primary">
                <Building2 className="size-5" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-semibold text-foreground">{c.name}</div>
                <div className="mt-0.5 flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[11px] text-muted-foreground">
                  <span className="inline-flex items-center gap-1"><Users className="size-3" /> {c.members.toLocaleString()} Old Starehians</span>
                  <span className="inline-flex items-center gap-1"><MapPin className="size-3" /> {c.city}</span>
                </div>
              </div>
            </div>
            <div className="mt-4 flex items-center justify-between rounded-xl bg-muted/50 px-3 py-2 text-[11px] text-foreground">
              <span className="inline-flex items-center gap-1.5 text-muted-foreground">
                <Calendar className="size-3.5" /> Next
              </span>
              <span className="truncate font-medium">{c.nextEvent}</span>
            </div>
            <button className="mt-4 w-full rounded-xl bg-primary py-2 text-xs font-semibold text-primary-foreground hover:brightness-110">
              Join chapter
            </button>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
