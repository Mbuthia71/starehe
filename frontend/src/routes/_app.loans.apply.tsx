import { createFileRoute } from "@tanstack/react-router";
import { UserPlus } from "lucide-react";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/loans/apply")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Coming soon</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">Request an introduction.</h1>
        <p className="mt-2 text-sm text-muted-foreground">Ask a mutual Old Starehian to introduce you to a fellow alumnus.</p>
      </header>
      <div className="card-elev">
        <EmptyState
          icon={UserPlus}
          title="We're building this next"
          description="This part of the Old Starehian Society network is on the roadmap. Watch the Activity feed for updates."
        />
      </div>
    </div>
  );
}
