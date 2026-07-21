import { createFileRoute } from "@tanstack/react-router";
import { Send } from "lucide-react";
import { EmptyState } from "@/components/EmptyState";

export const Route = createFileRoute("/_app/transfer")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <header>
        <div className="label-eyebrow">Coming soon</div>
        <h1 className="display mt-1 text-3xl font-semibold tracking-tight">Message an alumnus.</h1>
        <p className="mt-2 text-sm text-muted-foreground">Direct-message any Old Starehian in the network.</p>
      </header>
      <div className="card-elev">
        <EmptyState
          icon={Send}
          title="We're building this next"
          description="This part of the Old Starehian Society network is on the roadmap. Watch the Activity feed for updates."
        />
      </div>
    </div>
  );
}
