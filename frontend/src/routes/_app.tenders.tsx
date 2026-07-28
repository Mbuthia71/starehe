import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  FileText,
  MapPin,
  Calendar,
  DollarSign,
  Building2,
  ArrowUpRight,
  Plus,
  Search,
  Filter,
} from "lucide-react";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_app/tenders")({
  component: Tenders,
});

type Tender = {
  id: string;
  title: string;
  organization_name: string;
  location?: string;
  budget_range?: string;
  submission_deadline: string;
  description: string;
  requirements?: string;
  tender_number?: string;
  category?: string;
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

function Tenders() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [formData, setFormData] = useState({
    organization_name: "",
    title: "",
    description: "",
    requirements: "",
    budget_range: "",
    tender_number: "",
    category: "",
    submission_deadline: "",
  });

  // Fetch tenders from API
  const { data: tenders = [], isLoading, refetch } = useQuery<Tender[]>({
    queryKey: ['tenders'],
    queryFn: async () => {
      const res = await fetch('/api/business/tenders?limit=50&offset=0', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.tenders || data || [];
    },
    enabled: !!getToken(),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError("");

    try {
      const response = await fetch('/api/business/tenders', {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to submit tender");
      }

      const result = await response.json();
      setShowCreateModal(false);
      setFormData({
        organization_name: "",
        title: "",
        description: "",
        requirements: "",
        budget_range: "",
        tender_number: "",
        category: "",
        submission_deadline: "",
      });
      refetch();
      alert("Tender submitted successfully!");
    } catch (error: any) {
      setSubmitError(error.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
            <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
            Opportunities
          </div>
          <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
            Tenders
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Find and post tender opportunities within the Starehe community
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          Post Tender
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search tenders..."
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

      {/* Tenders list */}
      <div className="space-y-3">
        {isLoading ? (
          <div className="p-4 text-center text-sm text-muted-foreground">
            Loading tenders...
          </div>
        ) : tenders.length === 0 ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            No tenders posted yet. Be the first to post one!
          </div>
        ) : (
          tenders.map((tender, i) => (
            <motion.div
              key={tender.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.35, delay: 0.05 * i }}
            >
              <Link
                to="/tenders/$id"
                params={{ id: tender.id }}
                className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className="text-lg font-semibold text-foreground">{tender.title}</h3>
                      {tender.category && (
                        <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                          {tender.category}
                        </span>
                      )}
                    </div>
                    <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Building2 className="size-3.5" />
                        {tender.organization_name}
                      </div>
                      {tender.location && (
                        <div className="flex items-center gap-1">
                          <MapPin className="size-3.5" />
                          {tender.location}
                        </div>
                      )}
                      {tender.budget_range && (
                        <div className="flex items-center gap-1">
                          <DollarSign className="size-3.5" />
                          {tender.budget_range}
                        </div>
                      )}
                      <div className="flex items-center gap-1">
                        <Calendar className="size-3.5" />
                        Deadline: {new Date(tender.submission_deadline).toLocaleDateString()}
                      </div>
                    </div>
                    <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                      {tender.description}
                    </p>
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

      {/* Create Tender Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Post a Tender</h2>
            {submitError && (
              <div className="mb-4 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {submitError}
              </div>
            )}
            <form onSubmit={handleSubmit} className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Tender Title</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Supply of Office Equipment"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Organization</label>
                <input
                  type="text"
                  value={formData.organization_name}
                  onChange={(e) => setFormData({ ...formData, organization_name: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Ministry of Education"
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Category</label>
                  <select
                    value={formData.category}
                    onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="">Select category</option>
                    <option value="procurement">Procurement</option>
                    <option value="construction">Construction</option>
                    <option value="it">IT Services</option>
                    <option value="consulting">Consulting</option>
                  </select>
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Tender Number (Optional)</label>
                  <input
                    type="text"
                    value={formData.tender_number}
                    onChange={(e) => setFormData({ ...formData, tender_number: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., TND/2024/001"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Budget Range (Optional)</label>
                <input
                  type="text"
                  value={formData.budget_range}
                  onChange={(e) => setFormData({ ...formData, budget_range: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., KES 5,000,000 - 7,000,000"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Submission Deadline</label>
                <input
                  type="date"
                  value={formData.submission_deadline}
                  onChange={(e) => setFormData({ ...formData, submission_deadline: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the tender requirements..."
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Requirements (Optional)</label>
                <textarea
                  rows={3}
                  value={formData.requirements}
                  onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="List the requirements for this tender..."
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
                  {isSubmitting ? "Submitting..." : "Post Tender"}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
