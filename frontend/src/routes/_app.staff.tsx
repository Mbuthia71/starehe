import { createFileRoute } from "@tanstack/react-router";
import { Users } from "lucide-react";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/staff")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Coming soon</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">Society leadership.</h1>
        <p className="mt-2 text-sm text-muted-foreground">Meet the Old Starehian Society council and secretariat.</p>
      </header>
      <div className="card-elev">
        <EmptyState
          icon={Users}
          title="We're building this next"
          description="This part of the Old Starehian Society network is on the roadmap. Watch the Activity feed for updates."
        />
      </div>
    </div>
  );
}
