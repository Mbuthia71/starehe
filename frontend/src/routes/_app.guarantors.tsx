import { createFileRoute } from "@tanstack/react-router";
import { Award } from "lucide-react";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/guarantors")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Coming soon</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">Endorsements.</h1>
        <p className="mt-2 text-sm text-muted-foreground">Endorse fellow Old Starehians for mentorship or professional roles.</p>
      </header>
      <div className="card-elev">
        <EmptyState
          icon={Award}
          title="We're building this next"
          description="This part of the Old Starehian Society network is on the roadmap. Watch the Activity feed for updates."
        />
      </div>
    </div>
  );
}
