export const KES = (n: number | undefined | null, opts: { compact?: boolean; decimals?: number } = {}) => {
  const { compact, decimals = 0 } = opts;
  const num = n ?? 0;
  if (compact && Math.abs(num) >= 1_000_000) {
    return `KES ${(num / 1_000_000).toFixed(1)}M`;
  }
  if (compact && Math.abs(num) >= 1_000) {
    return `KES ${(num / 1_000).toFixed(1)}K`;
  }
  return `KES ${num.toLocaleString("en-KE", {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })}`;
};

export const Num = (n: number, decimals = 0) =>
  n.toLocaleString("en-KE", { minimumFractionDigits: decimals, maximumFractionDigits: decimals });

export const formatDate = (iso: string, withTime = false) => {
  const d = new Date(iso);
  const date = d.toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  if (!withTime) return date;
  const time = d.toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit" });
  return `${date} · ${time}`;
};