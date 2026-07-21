import { createFileRoute } from "@tanstack/react-router";
import { BarChart3 } from "lucide-react";

export const Route = createFileRoute("/_admin/admin/dashboard")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <div>
        <div className="text-xs uppercase tracking-[0.2em] text-muted-foreground">Admin</div>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight">Society overview.</h1>
        <p className="mt-2 max-w-2xl text-sm text-muted-foreground">Headline metrics for the Old Starehian Society: total alumni, weekly sign-ins, active chapters and pending approvals.</p>
      </div>
      <div className="rounded-2xl border border-border/60 bg-card p-10 text-center">
        <BarChart3 className="mx-auto size-6 text-muted-foreground" />
        <div className="mt-3 text-sm font-semibold">Awaiting live data</div>
        <p className="mt-1 text-xs text-muted-foreground">The Go backend will populate this view once wired up.</p>
      </div>
    </div>
  );
}
