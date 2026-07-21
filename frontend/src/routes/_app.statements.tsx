import { createFileRoute } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Activity, UserPlus, MessageSquare, ThumbsUp, Award, Calendar } from "lucide-react";

export const Route = createFileRoute("/_app/statements")({
  component: ActivityFeed,
});

interface FeedItem {
  id: string;
  actor: string;
  verb: string;
  target?: string;
  when: string;
  kind: "connect" | "comment" | "endorse" | "post" | "event";
}

const feed: FeedItem[] = [
  { id: "1", actor: "Kevin Mwangi (Griffin, '08)", verb: "sent you a connection request", when: "12 min ago", kind: "connect" },
  { id: "2", actor: "Class of 2010", verb: "posted", target: "Our 15-year reunion venue is confirmed", when: "1 h ago", kind: "post" },
  { id: "3", actor: "Prof. Wanjiku (Faculty)", verb: "commented on your post about", target: "engineering mentorship", when: "3 h ago", kind: "comment" },
  { id: "4", actor: "Brian Otieno (Livingstone, '11)", verb: "endorsed you for", target: "software engineering mentorship", when: "Yesterday", kind: "endorse" },
  { id: "5", actor: "OSS Nairobi Chapter", verb: "created an event", target: "Chapter dinner · Sat 21 Nov", when: "Yesterday", kind: "event" },
  { id: "6", actor: "OSS Secretariat", verb: "shared", target: "the November newsletter", when: "2 d ago", kind: "post" },
  { id: "7", actor: "Achieng Nyongo (Griffin, '13)", verb: "connected with", target: "Mary Wanjiku", when: "3 d ago", kind: "connect" },
];

const iconFor: Record<FeedItem["kind"], typeof Activity> = {
  connect: UserPlus,
  comment: MessageSquare,
  endorse: ThumbsUp,
  post: Activity,
  event: Calendar,
};

function ActivityFeed() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Community activity</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">
          What's happening.
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          A live feed of connections, posts, endorsements and events from your Starehe network.
        </p>
      </header>

      <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
        {feed.map((item, i) => {
          const Icon = iconFor[item.kind];
          return (
            <motion.div
              key={item.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: i * 0.03 }}
              className="flex items-start gap-3 px-4 py-4"
            >
              <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary/10 text-primary">
                <Icon className="size-4" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="text-sm text-foreground">
                  <span className="font-semibold">{item.actor}</span> {item.verb}
                  {item.target && <span className="font-medium text-primary"> · {item.target}</span>}
                </div>
                <div className="mt-0.5 text-[11px] text-muted-foreground">{item.when}</div>
              </div>
            </motion.div>
          );
        })}
      </div>

      <div className="card-elev flex items-center gap-3 p-4">
        <Award className="size-5 text-primary" />
        <div className="min-w-0 flex-1 text-xs text-muted-foreground">
          Your activity is visible to fellow Old Starehians only. Manage visibility in Profile.
        </div>
      </div>
    </div>
  );
}
