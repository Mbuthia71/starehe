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

export const Route = createFileRoute("/_app/tenders")({
  component: Tenders,
});

const mockTenders = [
  {
    id: "1",
    title: "Supply of Office Equipment",
    organization: "Ministry of Education",
    location: "Nairobi",
    budgetRange: "KES 5,000,000 - 7,000,000",
    submissionDeadline: "2024-08-15",
    description: "Supply and delivery of office equipment for various educational institutions...",
    postedAt: "1 week ago",
    bidsCount: 8,
    category: "Procurement",
  },
  {
    id: "2",
    title: "Construction of Classrooms",
    organization: "County Government of Kiambu",
    location: "Kiambu",
    budgetRange: "KES 15,000,000 - 20,000,000",
    submissionDeadline: "2024-09-01",
    description: "Construction of 10 classrooms in selected primary schools...",
    postedAt: "2 weeks ago",
    bidsCount: 12,
    category: "Construction",
  },
  {
    id: "3",
    title: "IT Infrastructure Upgrade",
    organization: "National Treasury",
    location: "Nairobi",
    budgetRange: "KES 8,000,000 - 12,000,000",
    submissionDeadline: "2024-08-30",
    description: "Upgrade of IT infrastructure including servers and networking equipment...",
    postedAt: "3 days ago",
    bidsCount: 5,
    category: "IT Services",
  },
];

function Tenders() {
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
        {mockTenders.map((tender, i) => (
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
                    <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                      {tender.category}
                    </span>
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Building2 className="size-3.5" />
                      {tender.organization}
                    </div>
                    <div className="flex items-center gap-1">
                      <MapPin className="size-3.5" />
                      {tender.location}
                    </div>
                    <div className="flex items-center gap-1">
                      <DollarSign className="size-3.5" />
                      {tender.budgetRange}
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="size-3.5" />
                      Deadline: {tender.submissionDeadline}
                    </div>
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {tender.description}
                  </p>
                  <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
                    <FileText className="size-3.5" />
                    {tender.bidsCount} bids submitted
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

      {/* Create Tender Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Post a Tender</h2>
            <form className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Tender Title</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Supply of Office Equipment"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Organization</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Ministry of Education"
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
                  <label className="mb-1 block text-sm font-medium">Category</label>
                  <select className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary">
                    <option value="procurement">Procurement</option>
                    <option value="construction">Construction</option>
                    <option value="it">IT Services</option>
                    <option value="consulting">Consulting</option>
                  </select>
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Budget Range</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., KES 5,000,000 - 7,000,000"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Submission Deadline</label>
                <input
                  type="date"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the tender requirements..."
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
                  Post Tender
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
