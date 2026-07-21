import { createFileRoute } from "@tanstack/react-router";
import { Users } from "lucide-react";

export const Route = createFileRoute("/_admin/admin/members")({
  component: Page,
});

function Page() {
  return (
    <div className="space-y-6">
      <div>
        <div className="text-xs uppercase tracking-[0.2em] text-muted-foreground">Admin</div>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight">Alumni roster.</h1>
        <p className="mt-2 max-w-2xl text-sm text-muted-foreground">Search, verify and manage the 10,000+ Old Starehian profiles.</p>
      </div>
      <div className="rounded-2xl border border-border/60 bg-card p-10 text-center">
        <Users className="mx-auto size-6 text-muted-foreground" />
        <div className="mt-3 text-sm font-semibold">Awaiting live data</div>
        <p className="mt-1 text-xs text-muted-foreground">The Go backend will populate this view once wired up.</p>
      </div>
    </div>
  );
}
