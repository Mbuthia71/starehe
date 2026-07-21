import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  GraduationCap,
  Users,
  Calendar,
  ArrowUpRight,
  Plus,
  Search,
  Filter,
} from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/_app/class-groups")({
  component: ClassGroups,
});

const mockClassGroups = [
  {
    id: "1",
    schoolType: "SBC",
    yearOfCompletion: 2008,
    className: "Griffin House",
    description: "Griffin House Class of 2008",
    memberCount: 45,
    isActive: true,
  },
  {
    id: "2",
    schoolType: "SBC",
    yearOfCompletion: 2010,
    className: "Livingstone House",
    description: "Livingstone House Class of 2010",
    memberCount: 52,
    isActive: true,
  },
  {
    id: "3",
    schoolType: "SGC",
    yearOfCompletion: 2012,
    className: "Form 4 - 2012",
    description: "Starehe Girls Centre Class of 2012",
    memberCount: 38,
    isActive: true,
  },
  {
    id: "4",
    schoolType: "SBC",
    yearOfCompletion: 2015,
    className: "Patel House",
    description: "Patel House Class of 2015",
    memberCount: 48,
    isActive: true,
  },
];

function ClassGroups() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedSchool, setSelectedSchool] = useState<"all" | "SBC" | "SGC">("all");
  const [showCreateModal, setShowCreateModal] = useState(false);

  const filteredGroups = mockClassGroups.filter(
    (group) =>
      (selectedSchool === "all" || group.schoolType === selectedSchool) &&
      (group.className.toLowerCase().includes(searchQuery.toLowerCase()) ||
        group.description?.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
            <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
            Community
          </div>
          <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
            Class Groups
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Connect with your classmates from Starehe Boys' Centre and Starehe Girls' Centre
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          Create Group
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search class groups..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full rounded-xl border border-border/60 bg-card pl-10 pr-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setSelectedSchool("all")}
            className={`rounded-xl border px-4 py-2.5 text-sm font-medium ${
              selectedSchool === "all"
                ? "border-primary bg-primary text-primary-foreground"
                : "border-border/60 bg-card hover:bg-muted"
            }`}
          >
            All
          </button>
          <button
            onClick={() => setSelectedSchool("SBC")}
            className={`rounded-xl border px-4 py-2.5 text-sm font-medium ${
              selectedSchool === "SBC"
                ? "border-primary bg-primary text-primary-foreground"
                : "border-border/60 bg-card hover:bg-muted"
            }`}
          >
            SBC
          </button>
          <button
            onClick={() => setSelectedSchool("SGC")}
            className={`rounded-xl border px-4 py-2.5 text-sm font-medium ${
              selectedSchool === "SGC"
                ? "border-primary bg-primary text-primary-foreground"
                : "border-border/60 bg-card hover:bg-muted"
            }`}
          >
            SGC
          </button>
        </div>
      </div>

      {/* Class groups list */}
      <div className="grid gap-3 sm:grid-cols-2">
        {filteredGroups.map((group, i) => (
          <motion.div
            key={group.id}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, delay: 0.05 * i }}
          >
            <Link
              to="/class-groups/$id"
              params={{ id: group.id }}
              className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                      {group.schoolType}
                    </span>
                    <h3 className="text-lg font-semibold text-foreground">{group.className}</h3>
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Calendar className="size-3.5" />
                      Class of {group.yearOfCompletion}
                    </div>
                    <div className="flex items-center gap-1">
                      <Users className="size-3.5" />
                      {group.memberCount} members
                    </div>
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {group.description}
                  </p>
                </div>
                <div className="grid size-8 shrink-0 place-items-center rounded-full bg-muted text-muted-foreground">
                  <ArrowUpRight className="size-4" />
                </div>
              </div>
            </Link>
          </motion.div>
        ))}
      </div>

      {/* Create Class Group Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Create Class Group</h2>
            <form className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">School Type</label>
                <select className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary">
                  <option value="SBC">Starehe Boys' Centre (SBC)</option>
                  <option value="SGC">Starehe Girls' Centre (SGC)</option>
                </select>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Year of Completion</label>
                <input
                  type="number"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., 2008"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Class Name</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Griffin House"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={3}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe your class group..."
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
                  Create Group
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
