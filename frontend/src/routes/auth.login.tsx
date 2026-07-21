import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { AlertCircle, ArrowRight } from "lucide-react";
import logoAsset from "@/assets/siohioma-logo.png?url";
import stareheVideo from "@/assets/videostarehe.mp4?url";
import { login, signup } from "@/lib/auth";

export const Route = createFileRoute("/auth/login")({
  component: Login,
});

function Login() {
  const navigate = useNavigate();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [greeting, setGreeting] = useState("");
  const [mode, setMode] = useState<"login" | "signup">("login");

  useEffect(() => {
    const hour = new Date().getHours();
    if (hour >= 6 && hour < 12) setGreeting("Good morning");
    else if (hour >= 12 && hour < 15) setGreeting("Good afternoon");
    else if (hour >= 15 && hour < 18) setGreeting("Good evening");
    else setGreeting("Good evening");
  }, []);

  const isEmail = identifier.includes("@");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!identifier.trim() || !password.trim()) return;
    setLoading(true);
    setError("");
    try {
      if (mode === "signup") {
        if (!fullName.trim()) {
          setError("Full name is required for signup");
          setLoading(false);
          return;
        }
        await signup(identifier.trim(), password.trim(), fullName.trim());
      } else {
        await login(identifier.trim(), password.trim());
      }
      navigate({ to: "/" });
    } catch (err: any) {
      setError(err.message || "Authentication failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative min-h-screen w-full overflow-hidden bg-black">
      {/* Full-bleed video */}
      <motion.video
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 1.4 }}
        autoPlay
        muted
        loop
        playsInline
        className="absolute inset-0 h-full w-full object-cover"
      >
        <source src={stareheVideo} type="video/mp4" />
      </motion.video>

      {/* Neutral cinematic vignette — bottom shadow for legibility only */}
      <div
        aria-hidden
        className="absolute inset-0"
        style={{
          background:
            "linear-gradient(180deg, rgba(0,0,0,0.55) 0%, rgba(0,0,0,0.15) 35%, rgba(0,0,0,0.15) 55%, rgba(0,0,0,0.75) 100%)",
        }}
      />
      {/* Starehe stripe */}
      <div className="starehe-stripe absolute inset-x-0 top-0 z-30 h-[3px]" />

      {/* Top brand row */}
      <div className="relative z-20 flex items-center justify-between px-6 pt-6 md:px-12 md:pt-8">
        <div className="flex items-center gap-3">
          <img src={logoAsset} alt="Old Starehian Society" className="h-11 w-auto brightness-0 invert drop-shadow-lg" />
          <div className="hidden sm:block">
            <div className="text-[10px] font-bold uppercase tracking-[0.3em] text-white/80">Old Starehian Society</div>
            <div className="text-[9px] font-medium uppercase tracking-[0.28em] text-white/50">United to serve · Est. 1959</div>
          </div>
        </div>
        <div className="hidden items-center gap-2 md:flex">
          <span className="h-[2px] w-6 bg-white/70" />
          <span className="text-[10px] font-bold uppercase tracking-[0.3em] text-white/70">Nec aspera terrent</span>
        </div>
      </div>

      {/* Layered composition: giant headline left, glass sign-in card right */}
      <div className="relative z-10 mx-auto grid min-h-[calc(100vh-88px)] w-full max-w-[1400px] grid-cols-1 items-center gap-10 px-6 pb-12 pt-8 md:px-12 lg:grid-cols-[1.15fr_minmax(0,440px)]">
        {/* Hero copy */}
        <motion.div
          initial={{ opacity: 0, y: 24 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.9, ease: [0.22, 1, 0.36, 1] }}
          className="max-w-[720px] text-white"
        >
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-white/20 bg-white/5 px-3 py-1 backdrop-blur-md">
            <span className="size-1.5 rounded-full bg-primary" />
            <span className="text-[10px] font-bold uppercase tracking-[0.28em] text-white/80">Alumni Network</span>
          </div>

          <h1 className="display font-black leading-[0.92] tracking-tight text-white text-[clamp(3rem,9vw,7.5rem)]">
            {greeting || "Welcome"},
            <br />
            <span className="italic text-white/85">Old Starehian.</span>
          </h1>

          <p className="mt-6 max-w-[46ch] text-base font-light leading-relaxed text-white/75 md:text-lg">
            Ten thousand brothers. Every class. Every corner of the world. Reconnect with the men who wore the maroon and gold before you — and the ones who will after.
          </p>

          <div className="mt-8 hidden gap-8 md:flex">
            <Stat n="10,000+" label="Alumni worldwide" />
            <Stat n="65+" label="Years of tradition" />
            <Stat n="120+" label="Chapters & class years" />
          </div>
        </motion.div>

        {/* Glass sign-in card */}
        <motion.form
          initial={{ opacity: 0, y: 24, scale: 0.98 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ delay: 0.25, duration: 0.75, ease: [0.22, 1, 0.36, 1] }}
          onSubmit={handleSubmit}
          className="glass-panel relative w-full max-w-[440px] justify-self-end rounded-3xl p-7 md:p-9"
        >
          <div className="mb-6">
            <div className="mb-4 flex gap-2">
              <button
                type="button"
                onClick={() => setMode("login")}
                className={`flex-1 rounded-lg px-4 py-2 text-[11px] font-bold uppercase tracking-wider transition-all ${
                  mode === "login"
                    ? "bg-white text-black"
                    : "bg-white/10 text-white/60 hover:bg-white/20"
                }`}
              >
                Sign In
              </button>
              <button
                type="button"
                onClick={() => setMode("signup")}
                className={`flex-1 rounded-lg px-4 py-2 text-[11px] font-bold uppercase tracking-wider transition-all ${
                  mode === "signup"
                    ? "bg-white text-black"
                    : "bg-white/10 text-white/60 hover:bg-white/20"
                }`}
              >
                Sign Up
              </button>
            </div>
            <h2 className="display text-3xl font-black text-white">
              {mode === "login" ? "Enter the brotherhood." : "Join the brotherhood."}
            </h2>
            <p className="mt-2 text-[13px] text-white/60">
              {mode === "login"
                ? "Use the phone number or email your class rep registered."
                : "Enter your phone number to create an account."}
            </p>
          </div>

          <div className="space-y-4">
            <label className="block">
              <div className="mb-1.5 text-[10px] font-bold uppercase tracking-[0.2em] text-white/60">
                {isEmail ? "Email (admin)" : "Phone number"}
              </div>
              <input
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                autoComplete="username"
                placeholder={isEmail ? "admin@oldstarehian.org" : "+254712345678"}
                className="w-full rounded-xl border border-white/15 bg-white/[0.08] px-4 py-3.5 text-[15px] font-medium text-white placeholder:text-white/35 focus:border-white/40 focus:bg-white/[0.12] focus:outline-none focus:ring-2 focus:ring-white/20"
              />
            </label>

            <label className="block">
              <div className="mb-1.5 flex items-center justify-between">
                <span className="text-[10px] font-bold uppercase tracking-[0.2em] text-white/60">Password</span>
                <a href="#" className="text-[10px] font-semibold text-white/70 hover:text-white">Forgot?</a>
              </div>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete={mode === "signup" ? "new-password" : "current-password"}
                placeholder="••••••••"
                className="w-full rounded-xl border border-white/15 bg-white/[0.08] px-4 py-3.5 text-[15px] font-medium text-white placeholder:text-white/35 focus:border-white/40 focus:bg-white/[0.12] focus:outline-none focus:ring-2 focus:ring-white/20"
              />
            </label>

            {mode === "signup" && (
              <label className="block">
                <div className="mb-1.5 text-[10px] font-bold uppercase tracking-[0.2em] text-white/60">
                  Full name
                </div>
                <input
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="John Doe"
                  autoComplete="name"
                  className="w-full rounded-xl border border-white/15 bg-white/[0.08] px-4 py-3.5 text-[15px] font-medium text-white placeholder:text-white/35 focus:border-white/40 focus:bg-white/[0.12] focus:outline-none focus:ring-2 focus:ring-white/20"
                />
              </label>
            )}
          </div>

          <AnimatePresence>
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -6 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -6 }}
                className="mt-4 flex items-center gap-2 rounded-xl border border-red-400/30 bg-red-500/15 px-3.5 py-3 text-[13px] font-medium text-red-100 backdrop-blur"
              >
                <AlertCircle className="size-4 shrink-0" />
                {error}
              </motion.div>
            )}
          </AnimatePresence>

          <motion.button
            type="submit"
            disabled={loading || !identifier || !password}
            whileTap={{ scale: 0.985 }}
            className="mt-6 flex w-full items-center justify-center gap-2 rounded-xl bg-white px-6 py-3.5 text-sm font-bold uppercase tracking-[0.14em] text-slate-900 shadow-xl transition-all hover:bg-white/90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading ? "Processing…" : (mode === "login" ? "Sign In" : "Sign Up")}
            {!loading && <ArrowRight className="size-4" />}
          </motion.button>

          <p className="mt-5 text-center text-[11px] text-white/50">
            New Old Starehian? <a className="font-semibold text-white/80 underline-offset-2 hover:underline" href="#">Request an invite</a> from your class rep.
          </p>
        </motion.form>
      </div>
    </div>
  );
}

function Stat({ n, label }: { n: string; label: string }) {
  return (
    <div className="border-l border-white/20 pl-4">
      <div className="display text-2xl font-black text-white">{n}</div>
      <div className="mt-0.5 text-[10px] font-semibold uppercase tracking-[0.22em] text-white/50">{label}</div>
    </div>
  );
}
