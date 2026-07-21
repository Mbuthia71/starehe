import { createFileRoute } from "@tanstack/react-router";
import { Wrench } from "lucide-react";

export const Route = createFileRoute("/_app/dev/emails")({
  component: () => (
    <div className="p-8 text-sm text-muted-foreground">
      <Wrench className="mb-3 size-5" />
      Email template previews are being rebuilt for the Old Starehian Society.
    </div>
  ),
});
