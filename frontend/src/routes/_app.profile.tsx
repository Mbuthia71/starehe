import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Fingerprint, KeyRound, LogOut, Phone, Mail, MapPin, ShieldCheck, ChevronRight,
  Building2, GraduationCap, Camera, Linkedin, Twitter, Instagram, Globe, Check, Award, BadgeCheck, Briefcase, Users,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { clearAuth, getStoredMember } from "@/lib/auth";
import { StatusPill } from "@/components/StatusPill";
import { toast } from "sonner";

export const Route = createFileRoute("/_app/profile")({
  component: Profile,
});

function Profile() {
  const navigate = useNavigate();
  const member = getStoredMember();
  const [biometric, setBiometric] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);
  const [avatar, setAvatar] = useState<string | null>(null);
  const [socials, setSocials] = useState({ linkedin: "", twitter: "", instagram: "", website: "" });
  const [privacy, setPrivacy] = useState({
    linkedin: true, twitter: true, instagram: true, website: true,
    email: true, phone: false, city: true, employer: true,
  });
  const [savedFlash, setSavedFlash] = useState(false);
  const [careerGroups, setCareerGroups] = useState<any[]>([]);
  const [loadingCareerGroups, setLoadingCareerGroups] = useState(false);

  useEffect(() => {
    try {
      const a = localStorage.getItem("oss.profile.avatar");
      if (a) setAvatar(a);
      const s = localStorage.getItem("oss.profile.socials");
      if (s) setSocials(JSON.parse(s));
      const p = localStorage.getItem("oss.profile.privacy");
      if (p) setPrivacy((prev) => ({ ...prev, ...JSON.parse(p) }));
    } catch {}
    
    fetchCareerGroups();
  }, []);

  const fetchCareerGroups = async () => {
    setLoadingCareerGroups(true);
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch('/api/groups?limit=10&offset=0', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      const careerGroupsOnly = (data.groups || []).filter((g: any) => g.type === 'career' && !g.is_member);
      setCareerGroups(careerGroupsOnly.slice(0, 3));
    } catch (error) {
      console.error('Failed to fetch career groups:', error);
    } finally {
      setLoadingCareerGroups(false);
    }
  };

  const joinCareerGroup = async (groupId: string) => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`/api/groups/${groupId}/join`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      if (response.ok) {
        toast.success("Joined career group!");
        fetchCareerGroups();
      }
    } catch (error) {
      console.error('Failed to join career group:', error);
      toast.error("Failed to join group");
    }
  };

  const onFile = (file: File) => {
    const MAX_MB = 5;
    if (!file.type.startsWith("image/")) {
      toast.error("That file isn't an image. Please choose a JPG, PNG or WebP.");
      return;
    }
    if (!/^image\/(jpe?g|png|webp|gif)$/i.test(file.type)) {
      toast.error("Unsupported format. Use JPG, PNG, WebP or GIF.");
      return;
    }
    if (file.size > MAX_MB * 1024 * 1024) {
      toast.error(`Image is too large (${(file.size / 1024 / 1024).toFixed(1)} MB). Max ${MAX_MB} MB.`);
      return;
    }
    const reader = new FileReader();
    reader.onload = () => {
      const url = String(reader.result);
      setAvatar(url);
      try { localStorage.setItem("oss.profile.avatar", url); } catch {}
      toast.success("Profile photo updated");
    };
    reader.onerror = () => toast.error("Could not read that image. Try another file.");
    reader.readAsDataURL(file);
  };

  const saveSocials = () => {
    try { localStorage.setItem("oss.profile.socials", JSON.stringify(socials)); } catch {}
    try { localStorage.setItem("oss.profile.privacy", JSON.stringify(privacy)); } catch {}
    setSavedFlash(true);
    setTimeout(() => setSavedFlash(false), 1600);
  };

  const togglePrivacy = (k: keyof typeof privacy) =>
    setPrivacy((p) => ({ ...p, [k]: !p[k] }));

  const shareUrl = typeof window !== "undefined" && member
    ? `${window.location.origin}/p/${member.accountNo ?? member.id}`
    : "";

  const copyShareUrl = async () => {
    if (!shareUrl) return;
    try {
      await navigator.clipboard.writeText(shareUrl);
      toast.success("Profile link copied");
    } catch {
      toast.error("Could not copy link");
    }
  };

  const initials = (member?.displayName ?? "OS").split(" ").map((p) => p[0]).slice(0, 2).join("");

  return (
    <div className="space-y-7">
      <header>
        <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
          <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
          Your profile
        </div>
        <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
          {member?.displayName ?? "Old Starehian"}.
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Public handle · <span className="font-mono">{member?.accountNo ?? "OSS-0000"}</span>
        </p>
      </header>

      {/* Shareable member card (image-4 style) */}
      <motion.section
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.45 }}
        className="grid gap-5 md:grid-cols-[300px_minmax(0,1fr)]"
      >
        {/* Card */}
        <div className="surface-splash relative overflow-hidden rounded-2xl border border-white/10 p-4 text-anchor-foreground shadow-[0_30px_60px_-30px_oklch(0.15_0.06_265/0.7)]">
          <div className="starehe-stripe absolute inset-x-0 top-0 h-[3px]" />
          {/* Photo */}
          <div className="relative mx-auto aspect-[4/5] w-full overflow-hidden rounded-xl">
            {avatar ? (
              <img src={avatar} alt={member?.displayName ?? "avatar"} className="h-full w-full object-cover" />
            ) : (
              <div className="grid h-full w-full place-items-center bg-gradient-to-br from-primary/40 via-secondary/40 to-[oklch(0.70_0.18_45/0.5)] text-4xl font-semibold text-white">
                {initials}
              </div>
            )}
            <button
              onClick={() => fileRef.current?.click()}
              className="absolute bottom-2 right-2 inline-flex items-center gap-1.5 rounded-full bg-black/60 px-3 py-1.5 text-[11px] font-semibold text-white backdrop-blur hover:bg-black/70"
            >
              <Camera className="size-3.5" /> {avatar ? "Change" : "Upload"}
            </button>
            <input
              ref={fileRef}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={(e) => e.target.files?.[0] && onFile(e.target.files[0])}
            />
          </div>
          {/* Name row */}
          <div className="mt-3 flex items-center gap-1.5">
            <div className="display truncate text-base font-semibold">{member?.displayName ?? "Old Starehian"}</div>
            <BadgeCheck className="size-4 shrink-0 text-primary" />
          </div>
          <div className="mt-0.5 text-[11px] text-anchor-foreground/60">
            {member?.officeName ?? "Griffin House · Class of 2008"}
          </div>
          {/* Stats */}
          <div className="mt-3 grid grid-cols-3 gap-2 text-center">
            <div className="rounded-lg bg-white/5 py-2">
              <div className="text-[10px] uppercase tracking-wider text-anchor-foreground/50">Connections</div>
              <div className="mt-0.5 text-sm font-semibold">184</div>
            </div>
            <div className="rounded-lg bg-white/5 py-2">
              <div className="text-[10px] uppercase tracking-wider text-anchor-foreground/50">Mentees</div>
              <div className="mt-0.5 text-sm font-semibold">7</div>
            </div>
            <div className="rounded-lg bg-white/5 py-2">
              <div className="text-[10px] uppercase tracking-wider text-anchor-foreground/50">Chapters</div>
              <div className="mt-0.5 text-sm font-semibold">3</div>
            </div>
          </div>
          <button
            onClick={copyShareUrl}
            className="mt-3 w-full rounded-xl bg-white text-slate-900 py-2.5 text-xs font-bold uppercase tracking-wider hover:bg-white/90"
          >
            Copy shareable link
          </button>
          {shareUrl && (
            <div className="mt-2 truncate text-center text-[10px] text-white/50">{shareUrl}</div>
          )}
        </div>

        {/* Bio + about */}
        <div className="space-y-4">
          <div className="card-elev p-5">
            <div className="flex items-center justify-between">
              <div className="label-eyebrow">About</div>
              <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                <Award className="size-3" /> Life member
              </span>
            </div>
            <p className="mt-3 text-sm leading-relaxed text-foreground">
              Old Starehian, Griffin House. Product manager currently building fintech for East Africa.
              Happy to mentor Form-4 leavers and Class of '24–'26 on tech and product careers.
            </p>
            <div className="mt-4 grid grid-cols-1 gap-2 text-[12px] sm:grid-cols-2">
              <span className="inline-flex items-center gap-2 text-muted-foreground">
                <GraduationCap className="size-3.5" /> Starehe Boys' Centre · 2004–2008
              </span>
              <span className="inline-flex items-center gap-2 text-muted-foreground">
                <Building2 className="size-3.5" /> Product · Independent
              </span>
              <span className="inline-flex items-center gap-2 text-muted-foreground">
                <MapPin className="size-3.5" /> Nairobi, Kenya
              </span>
              <span className="inline-flex items-center gap-2 text-muted-foreground">
                <ShieldCheck className="size-3.5" /> Verified alumnus
              </span>
            </div>
          </div>

          {/* Socials editor */}
          <div className="card-elev p-5">
            <div className="flex items-center justify-between">
              <div className="label-eyebrow">Socials</div>
              <button
                onClick={saveSocials}
                className="inline-flex items-center gap-1 rounded-full bg-primary px-3 py-1 text-[11px] font-semibold text-primary-foreground hover:brightness-110"
              >
                {savedFlash ? (<><Check className="size-3.5" /> Saved</>) : "Save"}
              </button>
            </div>
            <div className="mt-3 grid gap-2 sm:grid-cols-2">
              <SocialInput icon={Linkedin} label="LinkedIn" value={socials.linkedin} onChange={(v) => setSocials((s) => ({ ...s, linkedin: v }))} placeholder="linkedin.com/in/…" />
              <SocialInput icon={Twitter} label="X / Twitter" value={socials.twitter} onChange={(v) => setSocials((s) => ({ ...s, twitter: v }))} placeholder="@handle" />
              <SocialInput icon={Instagram} label="Instagram" value={socials.instagram} onChange={(v) => setSocials((s) => ({ ...s, instagram: v }))} placeholder="@handle" />
              <SocialInput icon={Globe} label="Website" value={socials.website} onChange={(v) => setSocials((s) => ({ ...s, website: v }))} placeholder="you.com" />
            </div>
          </div>

          {/* Privacy settings */}
          <div className="card-elev p-5">
            <div className="flex items-center justify-between">
              <div className="label-eyebrow">Privacy · what other alumni can see</div>
            </div>
            <p className="mt-2 text-[11px] text-muted-foreground">
              Toggle each field to control visibility on your public member card.
            </p>
            <div className="mt-3 grid gap-2 sm:grid-cols-2">
              <PrivacyRow label="LinkedIn" on={privacy.linkedin} onToggle={() => togglePrivacy("linkedin")} />
              <PrivacyRow label="X / Twitter" on={privacy.twitter} onToggle={() => togglePrivacy("twitter")} />
              <PrivacyRow label="Instagram" on={privacy.instagram} onToggle={() => togglePrivacy("instagram")} />
              <PrivacyRow label="Website" on={privacy.website} onToggle={() => togglePrivacy("website")} />
              <PrivacyRow label="Email" on={privacy.email} onToggle={() => togglePrivacy("email")} />
              <PrivacyRow label="Phone number" on={privacy.phone} onToggle={() => togglePrivacy("phone")} />
              <PrivacyRow label="Home city" on={privacy.city} onToggle={() => togglePrivacy("city")} />
              <PrivacyRow label="Employer" on={privacy.employer} onToggle={() => togglePrivacy("employer")} />
            </div>
          </div>

          {/* Career group suggestions */}
          {!loadingCareerGroups && careerGroups.length > 0 && (
            <div className="card-elev p-5">
              <div className="label-eyebrow">Suggested career groups</div>
              <p className="mt-2 text-[11px] text-muted-foreground">
                Join groups matching your profession to network with peers.
              </p>
              <div className="mt-3 space-y-2">
                {careerGroups.map((group) => (
                  <div key={group.id} className="flex items-center justify-between rounded-lg bg-muted/50 p-3">
                    <div className="flex items-center gap-2">
                      <div className="grid size-8 shrink-0 place-items-center rounded-lg bg-primary/10 text-primary">
                        <Briefcase className="size-4" />
                      </div>
                      <div>
                        <div className="text-sm font-medium">{group.name}</div>
                        <div className="text-[11px] text-muted-foreground">
                          <Users className="inline size-3 mr-1" />
                          {group.member_count} members
                        </div>
                      </div>
                    </div>
                    <button
                      onClick={() => joinCareerGroup(group.id)}
                      className="rounded-lg bg-primary px-3 py-1.5 text-[11px] font-semibold text-primary-foreground hover:brightness-110"
                    >
                      Join
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </motion.section>

      {/* Contact */}
      <section>
        <div className="label-eyebrow mb-3">Contact</div>
        <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
          <Row icon={Mail} label="Email" value={member?.emailAddress ?? "—"} />
          <Row icon={Phone} label="Phone" value={member?.mobileNo ?? "—"} />
          <Row icon={MapPin} label="Home city" value="Nairobi, Kenya" />
          <Row icon={Building2} label="Employer" value="Independent consultant" />
        </div>
      </section>

      {/* Security */}
      <section>
        <div className="label-eyebrow mb-3">Security</div>
        <div className="card-elev divide-y divide-border/60 overflow-hidden p-0">
          <button className="flex w-full items-center gap-3 px-4 py-3 text-left transition-colors hover:bg-muted/50">
            <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary/10 text-primary">
              <KeyRound className="size-4" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="text-sm font-medium text-foreground">Change password</div>
              <div className="text-[11px] text-muted-foreground">Recommended every 90 days</div>
            </div>
            <ChevronRight className="size-4 text-muted-foreground" />
          </button>

          <div className="flex items-center gap-3 px-4 py-3">
            <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary/10 text-primary">
              <Fingerprint className="size-4" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="text-sm font-medium text-foreground">Biometric sign-in</div>
              <div className="text-[11px] text-muted-foreground">Use Face ID or fingerprint on this device</div>
            </div>
            <button
              onClick={() => setBiometric((b) => !b)}
              className={`relative h-6 w-11 rounded-full transition-colors ${biometric ? "bg-primary" : "bg-muted"}`}
              aria-label="Toggle biometric"
            >
              <span className={`absolute top-0.5 size-5 rounded-full bg-white transition-transform ${biometric ? "translate-x-5" : "translate-x-0.5"}`} />
            </button>
          </div>

          <div className="flex items-center gap-3 px-4 py-3">
            <div className="grid size-10 shrink-0 place-items-center rounded-full bg-primary/10 text-primary">
              <ShieldCheck className="size-4" />
            </div>
            <div className="min-w-0 flex-1">
              <div className="text-sm font-medium text-foreground">Two-step verification</div>
              <div className="text-[11px] text-muted-foreground">On · SMS code on new devices</div>
            </div>
            <StatusPill status="Active" />
          </div>
        </div>
      </section>

      <button
        onClick={() => { clearAuth(); navigate({ to: "/auth/login" }); }}
        className="flex w-full items-center justify-center gap-2 rounded-2xl border border-danger/30 bg-danger/5 py-3 text-sm font-semibold text-danger transition-colors hover:bg-danger/10"
      >
        <LogOut className="size-4" /> Sign out
      </button>
    </div>
  );
}

function SocialInput({
  icon: Icon, label, value, onChange, placeholder,
}: { icon: typeof Mail; label: string; value: string; onChange: (v: string) => void; placeholder: string }) {
  return (
    <label className="block">
      <div className="mb-1 flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
        <Icon className="size-3" /> {label}
      </div>
      <input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="h-10 w-full rounded-lg border border-border bg-background px-3 text-sm text-foreground placeholder:text-muted-foreground/60 focus:outline-none focus:ring-2 focus:ring-primary/30"
      />
    </label>
  );
}

function Row({ icon: Icon, label, value }: { icon: typeof Mail; label: string; value: string }) {
  return (
    <div className="flex items-center gap-3 px-4 py-3">
      <div className="grid size-10 shrink-0 place-items-center rounded-full bg-muted text-muted-foreground">
        <Icon className="size-4" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="text-[11px] uppercase tracking-wider text-muted-foreground">{label}</div>
        <div className="truncate text-sm font-medium text-foreground">{value}</div>
      </div>
    </div>
  );
}

function PrivacyRow({ label, on, onToggle }: { label: string; on: boolean; onToggle: () => void }) {
  return (
    <div className="flex items-center justify-between rounded-lg border border-border/60 bg-background/50 px-3 py-2">
      <span className="text-xs font-semibold text-foreground">{label}</span>
      <button
        onClick={onToggle}
        className={`relative h-5 w-9 rounded-full transition-colors ${on ? "bg-primary" : "bg-muted"}`}
        aria-label={`Toggle ${label} visibility`}
      >
        <span className={`absolute top-0.5 size-4 rounded-full bg-white transition-transform ${on ? "translate-x-4" : "translate-x-0.5"}`} />
      </button>
    </div>
  );
}
