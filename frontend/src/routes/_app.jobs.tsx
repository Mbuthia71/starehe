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

export const Route = createFileRoute("/_app/jobs")({
  component: Jobs,
});

const mockJobs = [
  {
    id: "1",
    title: "Senior Software Engineer",
    company: "TechCorp Kenya",
    location: "Nairobi",
    jobType: "full-time",
    salaryRange: "KES 150,000 - 250,000",
    description: "We are looking for an experienced software engineer to join our growing team...",
    postedAt: "2 days ago",
    applicationsCount: 12,
  },
  {
    id: "2",
    title: "Marketing Manager",
    company: "Brand Solutions Ltd",
    location: "Remote",
    jobType: "remote",
    salaryRange: "KES 80,000 - 120,000",
    description: "Lead our marketing efforts and help grow our brand presence...",
    postedAt: "1 week ago",
    applicationsCount: 8,
  },
  {
    id: "3",
    title: "Financial Analyst",
    company: "Invest Africa",
    location: "Mombasa",
    jobType: "full-time",
    salaryRange: "KES 100,000 - 180,000",
    description: "Analyze financial data and provide insights for investment decisions...",
    postedAt: "3 days ago",
    applicationsCount: 15,
  },
];

function Jobs() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);

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
        {mockJobs.map((job, i) => (
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
                    <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                      {job.jobType}
                    </span>
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Building2 className="size-3.5" />
                      {job.company}
                    </div>
                    <div className="flex items-center gap-1">
                      <MapPin className="size-3.5" />
                      {job.location}
                    </div>
                    <div className="flex items-center gap-1">
                      <DollarSign className="size-3.5" />
                      {job.salaryRange}
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="size-3.5" />
                      {job.postedAt}
                    </div>
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {job.description}
                  </p>
                  <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
                    <Briefcase className="size-3.5" />
                    {job.applicationsCount} applications
                  </div>
                </div>
                <div className="grid size-8 shrink-0 place-items-center rounded-full bg-muted text-muted-foreground">
                  <ArrowUpRight className="size-4" />
                </div>
              </div>
            </Link>
          </motion.div>
        ))}
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
            <form className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Job Title</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Senior Software Engineer"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Company</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., TechCorp Kenya"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Location</label>
                  <input
                    type="text"
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., Nairobi"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Job Type</label>
                  <select className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary">
                    <option value="full-time">Full-time</option>
                    <option value="part-time">Part-time</option>
                    <option value="contract">Contract</option>
                    <option value="remote">Remote</option>
                  </select>
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Salary Range</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., KES 150,000 - 250,000"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the role and requirements..."
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm font-medium hover:bg-muted"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 rounded-xl bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
                >
                  Post Job
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
