import { motion } from "framer-motion";
import logoAsset from "@/assets/siohioma-logo.png?url";

interface SiohiomaLoaderProps {
  /** "fullscreen" overlays the viewport; "inline" sits in flow. */
  variant?: "fullscreen" | "inline";
  label?: string;
  /** Background tone of the overlay. "cream" matches the app, "anchor" for dark surfaces, "transparent" for none. */
  surface?: "cream" | "anchor" | "transparent";
}

/**
 * Enterprise loading state — Siohioma mark with a measured breath animation
 * and a thin progress hairline. Designed to feel calm, not busy.
 */
export function SiohiomaLoader({
  variant = "fullscreen",
  label = "Preparing your workspace",
  surface = "cream",
}: SiohiomaLoaderProps) {
  const wrapClass =
    variant === "fullscreen"
      ? `fixed inset-0 z-[100] grid place-items-center ${
          surface === "cream"
            ? "bg-background"
            : surface === "anchor"
              ? "bg-anchor"
              : "bg-transparent"
        }`
      : "grid w-full place-items-center py-16";

  return (
    <div className={wrapClass} role="status" aria-live="polite" aria-busy="true">
      <div className="flex flex-col items-center gap-10">
        {/* Logo with breathing halo */}
        <div className="relative grid place-items-center">
          <motion.div
            aria-hidden
            className="absolute size-48 rounded-full"
            style={{
              background:
                "radial-gradient(closest-side, oklch(0.265 0.045 162 / 0.10), transparent 70%)",
            }}
            animate={{ scale: [1, 1.18, 1], opacity: [0.5, 0.9, 0.5] }}
            transition={{ duration: 2.8, repeat: Infinity, ease: "easeInOut" }}
          />
          <motion.div
            aria-hidden
            className="absolute size-32 rounded-full"
            style={{
              background:
                "radial-gradient(closest-side, oklch(0.62 0.18 55 / 0.18), transparent 70%)",
            }}
            animate={{ scale: [1, 1.08, 1], opacity: [0.4, 0.8, 0.4] }}
            transition={{ duration: 2.8, repeat: Infinity, ease: "easeInOut", delay: 0.4 }}
          />
          <motion.img
            src={logoAsset}
            alt="Siohioma"
            className="relative h-20 w-auto select-none"
            draggable={false}
            initial={{ opacity: 0, y: 8, scale: 0.96 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            transition={{ duration: 0.9, ease: [0.22, 1, 0.36, 1] }}
          />
        </div>

        {/* Progress hairline */}
        <div className="flex flex-col items-center gap-4">
          <div className="relative h-px w-56 overflow-hidden bg-foreground/10">
            <motion.div
              className="absolute inset-y-0 w-1/3 bg-foreground/70"
              initial={{ x: "-100%" }}
              animate={{ x: "320%" }}
              transition={{ duration: 1.8, repeat: Infinity, ease: [0.65, 0, 0.35, 1] }}
            />
          </div>
          <div className="text-[10px] font-semibold uppercase tracking-[0.32em] text-muted-foreground">
            {label}
          </div>
        </div>
      </div>
      <span className="sr-only">Loading</span>
    </div>
  );
}

/** Compact inline variant — for cards and panels. */
export function SiohiomaSpinner({ size = 28 }: { size?: number }) {
  return (
    <div className="inline-grid place-items-center" role="status" aria-busy="true">
      <motion.div
        className="relative grid place-items-center"
        animate={{ rotate: 360 }}
        transition={{ duration: 2.4, repeat: Infinity, ease: "linear" }}
        style={{ width: size, height: size }}
      >
        <svg viewBox="0 0 36 36" className="size-full">
          <circle cx="18" cy="18" r="15" fill="none" stroke="currentColor" strokeOpacity="0.12" strokeWidth="2" />
          <circle
            cx="18"
            cy="18"
            r="15"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeDasharray="24 70"
          />
        </svg>
      </motion.div>
      <span className="sr-only">Loading</span>
    </div>
  );
}