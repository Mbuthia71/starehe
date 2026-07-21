import { createFileRoute } from "@tanstack/react-router";
import { HeartHandshake } from "lucide-react";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/withdraw")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Coming soon</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">Chapter contributions.</h1>
        <p className="mt-2 text-sm text-muted-foreground">Contribute to chapter events, reunions and bursary drives.</p>
      </header>
      <div className="card-elev">
        <EmptyState
          icon={HeartHandshake}
          title="We're building this next"
          description="This part of the Old Starehian Society network is on the roadmap. Watch the Activity feed for updates."
        />
      </div>
    </div>
  );
}
