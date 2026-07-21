import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Tag,
  MapPin,
  Calendar,
  Percent,
  Building2,
  ArrowUpRight,
  Plus,
  Search,
  Filter,
} from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/_app/offers")({
  component: Offers,
});

const mockOffers = [
  {
    id: "1",
    title: "20% Off Tech Products",
    business: "TechStore Kenya",
    location: "Nairobi",
    discountPercentage: 20,
    originalPrice: 5000,
    offerPrice: 4000,
    validFrom: "2024-07-20",
    validUntil: "2024-08-20",
    description: "Get 20% off on all tech products including laptops, phones, and accessories...",
    isExclusive: false,
    postedAt: "2 days ago",
    redemptionsCount: 45,
  },
  {
    id: "2",
    title: "Buy 1 Get 1 Free Coffee",
    business: "Java House",
    location: "Multiple Locations",
    discountPercentage: 50,
    originalPrice: 500,
    offerPrice: 250,
    validFrom: "2024-07-15",
    validUntil: "2024-07-31",
    description: "Buy one coffee and get another one absolutely free at any Java House location...",
    isExclusive: true,
    postedAt: "1 week ago",
    redemptionsCount: 120,
  },
  {
    id: "3",
    title: "15% Off Gym Membership",
    business: "FitLife Gym",
    location: "Westlands",
    discountPercentage: 15,
    originalPrice: 8000,
    offerPrice: 6800,
    validFrom: "2024-07-10",
    validUntil: "2024-09-10",
    description: "Get 15% off on monthly gym membership plans at FitLife Gym Westlands...",
    isExclusive: false,
    postedAt: "3 days ago",
    redemptionsCount: 28,
  },
];

function Offers() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.22em] text-muted-foreground">
            <span className="inline-block h-[3px] w-6 rounded-full bg-primary" />
            Benefits
          </div>
          <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
            Merchant Offers
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Exclusive discounts and offers from alumni businesses
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          Post Offer
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search offers..."
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

      {/* Offers list */}
      <div className="space-y-3">
        {mockOffers.map((offer, i) => (
          <motion.div
            key={offer.id}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, delay: 0.05 * i }}
          >
            <Link
              to="/offers/$id"
              params={{ id: offer.id }}
              className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className="text-lg font-semibold text-foreground">{offer.title}</h3>
                    {offer.isExclusive && (
                      <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-amber-500">
                        Exclusive
                      </span>
                    )}
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Building2 className="size-3.5" />
                      {offer.business}
                    </div>
                    <div className="flex items-center gap-1">
                      <MapPin className="size-3.5" />
                      {offer.location}
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="size-3.5" />
                      Valid until {offer.validUntil}
                    </div>
                  </div>
                  <div className="mt-3 flex items-center gap-3">
                    <div className="flex items-center gap-1 text-lg font-bold text-primary">
                      <Percent className="size-4" />
                      {offer.discountPercentage}% OFF
                    </div>
                    <div className="flex items-center gap-2 text-sm">
                      <span className="line-through text-muted-foreground">
                        KES {offer.originalPrice.toLocaleString()}
                      </span>
                      <span className="font-semibold text-foreground">
                        KES {offer.offerPrice.toLocaleString()}
                      </span>
                    </div>
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {offer.description}
                  </p>
                  <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
                    <Tag className="size-3.5" />
                    {offer.redemptionsCount} redemptions
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

      {/* Create Offer Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Post Merchant Offer</h2>
            <form className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Offer Title</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., 20% Off Tech Products"
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Business</label>
                <input
                  type="text"
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., TechStore Kenya"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Original Price (KES)</label>
                  <input
                    type="number"
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., 5000"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Offer Price (KES)</label>
                  <input
                    type="number"
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., 4000"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Valid From</label>
                  <input
                    type="date"
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Valid Until</label>
                  <input
                    type="date"
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={3}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the offer..."
                />
              </div>
              <div className="flex items-center gap-2">
                <input type="checkbox" id="exclusive" className="rounded border-border" />
                <label htmlFor="exclusive" className="text-sm">Make this offer exclusive</label>
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
                  Post Offer
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
