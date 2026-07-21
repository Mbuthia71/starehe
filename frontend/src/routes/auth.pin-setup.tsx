import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { ShieldCheck } from "lucide-react";

export const Route = createFileRoute("/auth/pin-setup")({
  component: PinSetup,
});

function PinSetup() {
  const navigate = useNavigate();
  const [pin, setPin] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState("");

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <div className="w-full max-w-sm">
        <div className="grid size-12 place-items-center rounded-2xl bg-primary/10 text-primary">
          <ShieldCheck className="size-5" />
        </div>
        <h1 className="mt-6 display text-3xl font-semibold tracking-tight">Set your PIN.</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Choose a 4-digit PIN to unlock the Old Starehian Society app on this device.
        </p>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            if (pin.length !== 4) { setError("PIN must be 4 digits."); return; }
            if (pin !== confirm) { setError("PINs don't match."); return; }
            navigate({ to: "/" });
          }}
          className="mt-8 space-y-4"
        >
          <input
            inputMode="numeric"
            maxLength={4}
            value={pin}
            onChange={(e) => setPin(e.target.value.replace(/\D/g, ""))}
            placeholder="New PIN"
            className="h-14 w-full rounded-2xl border border-input bg-card text-center text-2xl tracking-[0.4em] focus:outline-none focus:ring-2 focus:ring-ring/40"
          />
          <input
            inputMode="numeric"
            maxLength={4}
            value={confirm}
            onChange={(e) => setConfirm(e.target.value.replace(/\D/g, ""))}
            placeholder="Confirm PIN"
            className="h-14 w-full rounded-2xl border border-input bg-card text-center text-2xl tracking-[0.4em] focus:outline-none focus:ring-2 focus:ring-ring/40"
          />
          {error && <div className="text-xs text-danger">{error}</div>}
          <button type="submit" className="w-full rounded-2xl bg-primary py-3 text-sm font-semibold text-primary-foreground hover:brightness-110">
            Save PIN
          </button>
        </form>
      </div>
    </div>
  );
}
