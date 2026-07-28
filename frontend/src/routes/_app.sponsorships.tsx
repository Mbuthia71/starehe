import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Heart,
  MapPin,
  Calendar,
  DollarSign,
  Users,
  ArrowUpRight,
  Plus,
  Search,
  Filter,
} from "lucide-react";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_app/sponsorships")({
  component: Sponsorships,
});

type Sponsorship = {
  id: string;
  title: string;
  description: string;
  sponsorship_type: string;
  target_amount: number;
  start_date: string;
  end_date?: string;
  beneficiary?: string;
  image_url?: string;
  created_at: string;
};

const getToken = () => {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem('oss_token') || '';
};

const authHeaders = () => ({
  'Content-Type': 'application/json',
  Authorization: `Bearer ${getToken()}`,
});

function Sponsorships() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    sponsorship_type: "",
    target_amount: 0,
    start_date: "",
    end_date: "",
    beneficiary: "",
    image_url: "",
  });

  // Fetch sponsorships from API
  const { data: sponsorships = [], isLoading, refetch } = useQuery<Sponsorship[]>({
    queryKey: ['sponsorships'],
    queryFn: async () => {
      const res = await fetch('/api/business/sponsorships?limit=50&offset=0', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.sponsorships || data || [];
    },
    enabled: !!getToken(),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError("");

    try {
      const response = await fetch('/api/business/sponsorships', {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to create sponsorship");
      }

      const result = await response.json();
      setShowCreateModal(false);
      setFormData({
        title: "",
        description: "",
        sponsorship_type: "",
        target_amount: 0,
        start_date: "",
        end_date: "",
        beneficiary: "",
        image_url: "",
      });
      refetch();
      alert("Sponsorship created successfully!");
    } catch (error: any) {
      setSubmitError(error.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const progressPercentage = (current: number, target: number) =>
    Math.round((current / target) * 100);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
            <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
            Giving Back
          </div>
          <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
            Sponsorships
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Support Starehe through sponsorships and endowments
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          Create Sponsorship
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search sponsorships..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full rounded-xl border border-border/60 bg-card pl-10 pr-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
        <button className="inline-flex items-center gap-2 rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm font-medium hover:bg-muted">
          <Filter className="size-4" />
          Filters
        </button>
      </div>

      {/* Sponsorships list */}
      <div className="space-y-3">
        {isLoading ? (
          <div className="p-4 text-center text-sm text-muted-foreground">
            Loading sponsorships...
          </div>
        ) : sponsorships.length === 0 ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            No sponsorships posted yet. Be the first to create one!
          </div>
        ) : (
          sponsorships.map((sponsorship, i) => (
            <motion.div
              key={sponsorship.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.35, delay: 0.05 * i }}
            >
              <Link
                to="/sponsorships/$id"
                params={{ id: sponsorship.id }}
                className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className="text-lg font-semibold text-foreground">{sponsorship.title}</h3>
                      <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                        {sponsorship.sponsorship_type}
                      </span>
                    </div>
                    <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                      {sponsorship.beneficiary && (
                        <div className="flex items-center gap-1">
                          <Users className="size-3.5" />
                          {sponsorship.beneficiary}
                        </div>
                      )}
                      {sponsorship.end_date && (
                        <div className="flex items-center gap-1">
                          <Calendar className="size-3.5" />
                          Until {new Date(sponsorship.end_date).toLocaleDateString()}
                        </div>
                      )}
                    </div>
                    <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                      {sponsorship.description}
                    </p>
                    
                    {/* Target amount display */}
                    <div className="mt-4">
                      <div className="mb-2 flex items-center justify-between text-sm">
                        <span className="font-medium text-foreground">
                          Target: KES {sponsorship.target_amount.toLocaleString()}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="grid size-8 shrink-0 place-items-center rounded-full bg-muted text-muted-foreground">
                    <ArrowUpRight className="size-4" />
                  </div>
                </div>
              </Link>
            </motion.div>
          ))
        )}
      </div>

      {/* Create Sponsorship Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Create Sponsorship</h2>
            {submitError && (
              <div className="mb-4 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {submitError}
              </div>
            )}
            <form onSubmit={handleSubmit} className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Title</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Education Support Fund"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Type</label>
                <select
                  value={formData.sponsorship_type}
                  onChange={(e) => setFormData({ ...formData, sponsorship_type: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  required
                >
                  <option value="">Select type</option>
                  <option value="education">Education</option>
                  <option value="sports">Sports</option>
                  <option value="infrastructure">Infrastructure</option>
                  <option value="events">Events</option>
                  <option value="other">Other</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Target Amount (KES)</label>
                  <input
                    type="number"
                    value={formData.target_amount}
                    onChange={(e) => setFormData({ ...formData, target_amount: parseFloat(e.target.value) || 0 })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., 5000000"
                    required
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Start Date</label>
                  <input
                    type="date"
                    value={formData.start_date}
                    onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">End Date</label>
                <input
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Beneficiary</label>
                <input
                  type="text"
                  value={formData.beneficiary}
                  onChange={(e) => setFormData({ ...formData, beneficiary: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Starehe Students"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Image URL (Optional)</label>
                <input
                  type="url"
                  value={formData.image_url}
                  onChange={(e) => setFormData({ ...formData, image_url: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="https://..."
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the sponsorship purpose..."
                  required
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm font-medium hover:bg-muted"
                  disabled={isSubmitting}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 rounded-xl bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Creating..." : "Create Sponsorship"}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
