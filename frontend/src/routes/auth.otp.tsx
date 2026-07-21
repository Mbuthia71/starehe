import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { AlertCircle, KeyRound } from "lucide-react";

export const Route = createFileRoute("/auth/otp")({
  component: OtpPage,
});

function OtpPage() {
  const navigate = useNavigate();
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <div className="w-full max-w-sm">
        <div className="grid size-12 place-items-center rounded-2xl bg-primary/10 text-primary">
          <KeyRound className="size-5" />
        </div>
        <h1 className="mt-6 display text-3xl font-semibold tracking-tight">Enter your code.</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          We sent a six-digit code to your Old Starehian contact number. Enter it below to continue.
        </p>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            if (code.length !== 6) { setError("Enter the full 6-digit code."); return; }
            navigate({ to: "/" });
          }}
          className="mt-8 space-y-4"
        >
          <input
            inputMode="numeric"
            maxLength={6}
            value={code}
            onChange={(e) => setCode(e.target.value.replace(/\D/g, ""))}
            placeholder="••••••"
            className="h-14 w-full rounded-2xl border border-input bg-card text-center text-2xl tracking-[0.4em] focus:outline-none focus:ring-2 focus:ring-ring/40"
          />
          {error && (
            <div className="flex items-center gap-2 rounded-xl border border-danger/30 bg-danger/10 px-3 py-2 text-xs text-danger">
              <AlertCircle className="size-4" /> {error}
            </div>
          )}
          <button type="submit" className="w-full rounded-2xl bg-primary py-3 text-sm font-semibold text-primary-foreground hover:brightness-110">
            Verify
          </button>
        </form>
      </div>
    </div>
  );
}
