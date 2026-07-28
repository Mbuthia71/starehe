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
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_app/offers")({
  component: Offers,
});

type Offer = {
  id: string;
  title: string;
  description: string;
  business_id?: string;
  discount_percentage?: number;
  original_price?: number;
  offer_price?: number;
  valid_from: string;
  valid_until: string;
  terms_conditions?: string;
  image_url?: string;
  is_exclusive: boolean;
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

function Offers() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    business_id: "",
    discount_percentage: 0,
    original_price: 0,
    offer_price: 0,
    valid_from: "",
    valid_until: "",
    terms_conditions: "",
    image_url: "",
    is_exclusive: false,
  });

  // Fetch offers from API
  const { data: offers = [], isLoading, refetch } = useQuery<Offer[]>({
    queryKey: ['offers'],
    queryFn: async () => {
      const res = await fetch('/api/business/offers?limit=50&offset=0', {
        headers: authHeaders(),
      });
      if (!res.ok) return [];
      const data = await res.json();
      return data.offers || data || [];
    },
    enabled: !!getToken(),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError("");

    try {
      const response = await fetch('/api/business/offers', {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to create offer");
      }

      const result = await response.json();
      setShowCreateModal(false);
      setFormData({
        title: "",
        description: "",
        business_id: "",
        discount_percentage: 0,
        original_price: 0,
        offer_price: 0,
        valid_from: "",
        valid_until: "",
        terms_conditions: "",
        image_url: "",
        is_exclusive: false,
      });
      refetch();
      alert("Offer posted successfully!");
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
        {isLoading ? (
          <div className="p-4 text-center text-sm text-muted-foreground">
            Loading offers...
          </div>
        ) : offers.length === 0 ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            No offers posted yet. Be the first to post one!
          </div>
        ) : (
          offers.map((offer, i) => (
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
                      {offer.is_exclusive && (
                        <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-amber-500">
                          Exclusive
                        </span>
                      )}
                    </div>
                    <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Calendar className="size-3.5" />
                        Valid until {new Date(offer.valid_until).toLocaleDateString()}
                      </div>
                    </div>
                    {offer.discount_percentage && offer.original_price && offer.offer_price && (
                      <div className="mt-3 flex items-center gap-3">
                        <div className="flex items-center gap-1 text-lg font-bold text-primary">
                          <Percent className="size-4" />
                          {offer.discount_percentage}% OFF
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                          <span className="line-through text-muted-foreground">
                            KES {offer.original_price.toLocaleString()}
                          </span>
                          <span className="font-semibold text-foreground">
                            KES {offer.offer_price.toLocaleString()}
                          </span>
                        </div>
                      </div>
                    )}
                    <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                      {offer.description}
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

      {/* Create Offer Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">Post Merchant Offer</h2>
            {submitError && (
              <div className="mb-4 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {submitError}
              </div>
            )}
            <form onSubmit={handleSubmit} className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Offer Title</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., 20% Off Tech Products"
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
                  <label className="mb-1 block text-sm font-medium">Original Price (KES)</label>
                  <input
                    type="number"
                    value={formData.original_price}
                    onChange={(e) => setFormData({ ...formData, original_price: parseFloat(e.target.value) || 0 })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., 5000"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Offer Price (KES)</label>
                  <input
                    type="number"
                    value={formData.offer_price}
                    onChange={(e) => setFormData({ ...formData, offer_price: parseFloat(e.target.value) || 0 })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., 4000"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Discount Percentage (Optional)</label>
                <input
                  type="number"
                  value={formData.discount_percentage}
                  onChange={(e) => setFormData({ ...formData, discount_percentage: parseInt(e.target.value) || 0 })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., 20"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Valid From</label>
                  <input
                    type="date"
                    value={formData.valid_from}
                    onChange={(e) => setFormData({ ...formData, valid_from: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Valid Until</label>
                  <input
                    type="date"
                    value={formData.valid_until}
                    onChange={(e) => setFormData({ ...formData, valid_until: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={3}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe the offer..."
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Terms & Conditions (Optional)</label>
                <textarea
                  rows={2}
                  value={formData.terms_conditions}
                  onChange={(e) => setFormData({ ...formData, terms_conditions: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Any terms and conditions..."
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
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="exclusive"
                  checked={formData.is_exclusive}
                  onChange={(e) => setFormData({ ...formData, is_exclusive: e.target.checked })}
                  className="rounded border-border"
                />
                <label htmlFor="exclusive" className="text-sm">Make this offer exclusive</label>
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
                  {isSubmitting ? "Posting..." : "Post Offer"}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
