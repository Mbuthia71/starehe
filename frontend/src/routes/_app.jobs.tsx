import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Briefcase,
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

export const Route = createFileRoute("/_app/jobs")({
  component: Jobs,
});

type Job = {
  id: string;
  title: string;
  description: string;
  business_id?: string;
  requirements?: string;
  responsibilities?: string;
  location?: string;
  job_type?: string;
  salary_range?: string;
  application_deadline?: string;
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

function Jobs() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    business_id: "",
    requirements: "",
    responsibilities: "",
    location: "",
    job_type: "",
    salary_range: "",
    application_deadline: "",
  });

  // Fetch jobs from API
  const { data: jobs = [], isLoading, refetch } = useQuery<Job[]>({
    queryKey: ['jobs'],
    queryFn: async () => {
      const res = await fetch('/api/business/jobs?limit=50&offset=0', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.jobs || data || [];
    },
    enabled: !!getToken(),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError("");

    try {
      const response = await fetch('/api/business/jobs', {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to submit job");
      }

      const result = await response.json();
      setShowCreateModal(false);
      setFormData({
        title: "",
        description: "",
        business_id: "",
        requirements: "",
        responsibilities: "",
        location: "",
        job_type: "",
        salary_range: "",
        application_deadline: "",
      });
      refetch();
      alert("Job posted successfully!");
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
            Jobs
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Find and post job opportunities within the Starehe community
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          Post Job
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search jobs..."
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

      {/* Jobs list */}
      <div className="space-y-3">
        {isLoading ? (
          <div className="p-4 text-center text-sm text-muted-foreground">
            Loading jobs...
          </div>
        ) : jobs.length === 0 ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            No jobs posted yet. Be the first to post one!
          </div>
        ) : (
          jobs.map((job, i) => (
            <motion.div
              key={job.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.35, delay: 0.05 * i }}
            >
              <Link
                to="/jobs/$id"
                params={{ id: job.id }}
                className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className="text-lg font-semibold text-foreground">{job.title}</h3>
                      {job.job_type && (
                        <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                          {job.job_type}
                        </span>
                      )}
                    </div>
                    <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                      {job.location && (
                        <div className="flex items-center gap-1">
                          <MapPin className="size-3.5" />
                          {job.location}
                        </div>
                      )}
                      {job.salary_range && (
                        <div className="flex items-center gap-1">
                          <DollarSign className="size-3.5" />
                          {job.salary_range}
                        </div>
                      )}
                      <div className="flex items-center gap-1">
                        <Calendar className="size-3.5" />
                        {new Date(job.created_at).toLocaleDateString()}
                      </div>
                    </div>
                    <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                      {job.description}
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

      {/* Create Job Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Post a Job</h2>
            {submitError && (
              <div className="mb-4 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {submitError}
              </div>
            )}
            <form onSubmit={handleSubmit} className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Job Title</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Senior Software Engineer"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Business ID (Optional)</label>
                <input
                  type="text"
                  value={formData.business_id}
                  onChange={(e) => setFormData({ ...formData, business_id: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Link to your business listing"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Location</label>
                  <input
                    type="text"
                    value={formData.location}
                    onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., Nairobi"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Job Type</label>
                  <select
                    value={formData.job_type}
                    onChange={(e) => setFormData({ ...formData, job_type: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="">Select type</option>
                    <option value="full-time">Full-time</option>
                    <option value="part-time">Part-time</option>
                    <option value="contract">Contract</option>
                    <option value="remote">Remote</option>
                  </select>
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Salary Range (Optional)</label>
                <input
                  type="text"
                  value={formData.salary_range}
                  onChange={(e) => setFormData({ ...formData, salary_range: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., KES 150,000 - 250,000"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Application Deadline (Optional)</label>
                <input
                  type="date"
                  value={formData.application_deadline}
                  onChange={(e) => setFormData({ ...formData, application_deadline: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the role and requirements..."
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
                  placeholder="List the requirements for this role..."
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Responsibilities (Optional)</label>
                <textarea
                  rows={3}
                  value={formData.responsibilities}
                  onChange={(e) => setFormData({ ...formData, responsibilities: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="List the key responsibilities..."
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
                  {isSubmitting ? "Submitting..." : "Post Job"}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
