import { createFileRoute, Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Building2,
  MapPin,
  Phone,
  Globe,
  Instagram,
  Facebook,
  Star,
  ArrowUpRight,
  Plus,
  Search,
  Filter,
  Shield,
} from "lucide-react";
import { useState } from "react";
import { API_CONFIG } from "../lib/api";

export const Route = createFileRoute("/_app/business")({
  component: Business,
});

const mockBusinesses = [
  {
    id: "1",
    businessName: "TechVentures Kenya",
    phoneNumber: "+254 712 345 678",
    location: "Nairobi",
    website: "https://techventures.co.ke",
    description: "Leading software development company specializing in custom solutions for businesses across East Africa...",
    instagramHandle: "@techventureske",
    facebookHandle: "TechVentures Kenya",
    logoUrl: null,
    isVerified: true,
    isFeatured: true,
    status: "active",
    createdAt: "3 months ago",
  },
  {
    id: "2",
    businessName: "Savanna Safaris Ltd",
    phoneNumber: "+254 723 456 789",
    location: "Nairobi",
    website: "https://savannasafaris.co.ke",
    description: "Premium safari tours and wildlife experiences across Kenya's national parks...",
    instagramHandle: "@savannasafaris",
    facebookHandle: "Savanna Safaris",
    logoUrl: null,
    isVerified: true,
    isFeatured: false,
    status: "active",
    createdAt: "6 months ago",
  },
  {
    id: "3",
    businessName: "Urban Design Studio",
    phoneNumber: "+254 734 567 890",
    location: "Mombasa",
    website: "https://urbandesign.co.ke",
    description: "Creative design agency offering branding, web design, and marketing services...",
    instagramHandle: "@urbandesignstudio",
    facebookHandle: "Urban Design Studio",
    logoUrl: null,
    isVerified: false,
    isFeatured: false,
    status: "active",
    createdAt: "2 months ago",
  },
];

function Business() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [formData, setFormData] = useState({
    business_name: "",
    phone_number: "",
    location: "",
    website: "",
    description: "",
    instagram_handle: "",
    facebook_handle: "",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitError("");

    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_CONFIG.baseUrl}/business/listings`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to create business listing");
      }

      const result = await response.json();
      setShowCreateModal(false);
      setFormData({
        business_name: "",
        phone_number: "",
        location: "",
        website: "",
        description: "",
        instagram_handle: "",
        facebook_handle: "",
      });
      alert("Business listed successfully!");
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
            Alumni Businesses
          </div>
          <h1 className="display mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
            Business Directory
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Discover and support businesses owned by Starehe alumni
          </p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="size-4" />
          List Business
        </button>
      </div>

      {/* Search and filters */}
      <div className="flex gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search businesses..."
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

      {/* Businesses list */}
      <div className="grid gap-3 sm:grid-cols-2">
        {mockBusinesses.map((business, i) => (
          <motion.div
            key={business.id}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, delay: 0.05 * i }}
          >
            <Link
              to="/business/$id"
              params={{ id: business.id }}
              className="card-elev block p-5 transition-all hover:-translate-y-0.5 hover:border-primary/30"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className="text-lg font-semibold text-foreground">{business.businessName}</h3>
                    {business.isVerified && (
                      <Star className="size-4 fill-amber-500 text-amber-500" />
                    )}
                    {business.isFeatured && (
                      <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-amber-500">
                        Featured
                      </span>
                    )}
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <MapPin className="size-3.5" />
                      {business.location}
                    </div>
                    <div className="flex items-center gap-1">
                      <Phone className="size-3.5" />
                      {business.phoneNumber}
                    </div>
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {business.description}
                  </p>
                  <div className="mt-3 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    {business.website && (
                      <div className="flex items-center gap-1">
                        <Globe className="size-3.5" />
                        <span className="truncate max-w-[150px]">{business.website}</span>
                      </div>
                    )}
                    {business.instagramHandle && (
                      <div className="flex items-center gap-1">
                        <Instagram className="size-3.5" />
                        {business.instagramHandle}
                      </div>
                    )}
                    {business.facebookHandle && (
                      <div className="flex items-center gap-1">
                        <Facebook className="size-3.5" />
                        <span className="truncate max-w-[100px]">{business.facebookHandle}</span>
                      </div>
                    )}
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

      {/* Create Business Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <h2 className="text-xl font-semibold">List Your Business</h2>
            {submitError && (
              <div className="mb-4 rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                {submitError}
              </div>
            )}
            <form onSubmit={handleSubmit} className="mt-4 space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Business Name</label>
                <input
                  type="text"
                  value={formData.business_name}
                  onChange={(e) => setFormData({ ...formData, business_name: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., TechVentures Kenya"
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Phone Number</label>
                  <input
                    type="tel"
                    value={formData.phone_number}
                    onChange={(e) => setFormData({ ...formData, phone_number: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., +254 712 345 678"
                    required
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Location</label>
                  <input
                    type="text"
                    value={formData.location}
                    onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., Nairobi"
                    required
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Website (Optional)</label>
                <input
                  type="url"
                  value={formData.website}
                  onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., https://yourbusiness.co.ke"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="mb-1 block text-sm font-medium">Instagram Handle (Optional)</label>
                  <input
                    type="text"
                    value={formData.instagram_handle}
                    onChange={(e) => setFormData({ ...formData, instagram_handle: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., @yourbusiness"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium">Facebook Handle (Optional)</label>
                  <input
                    type="text"
                    value={formData.facebook_handle}
                    onChange={(e) => setFormData({ ...formData, facebook_handle: e.target.value })}
                    className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                    placeholder="e.g., Your Business"
                  />
                </div>
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <textarea
                  rows={4}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="Describe your business..."
                  required
                />
              </div>
              <div className="flex items-center gap-2 rounded-lg bg-primary/5 p-3">
                <Shield className="size-4 text-primary" />
                <p className="text-xs text-muted-foreground">
                  Business listings are free. Escrow integration available for secure transactions.
                </p>
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
                  {isSubmitting ? "Listing..." : "List Business"}
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
