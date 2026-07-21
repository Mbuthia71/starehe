import { createFileRoute } from "@tanstack/react-router";
import { motion } from "framer-motion";
import {
  Building2,
  MapPin,
  Phone,
  Globe,
  Instagram,
  Facebook,
  Star,
  ArrowLeft,
  Shield,
  Lock,
  Clock,
  CheckCircle,
  XCircle,
  DollarSign,
} from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

export const Route = createFileRoute("/_app/business/$id")({
  component: BusinessDetail,
});

const mockBusiness = {
  id: "1",
  businessName: "TechVentures Kenya",
  phoneNumber: "+254 712 345 678",
  location: "Nairobi",
  website: "https://techventures.co.ke",
  description: "Leading software development company specializing in custom solutions for businesses across East Africa. We offer web development, mobile apps, cloud services, and IT consulting.",
  instagramHandle: "@techventureske",
  facebookHandle: "TechVentures Kenya",
  logoUrl: null,
  isVerified: true,
  isFeatured: true,
  status: "active",
  createdAt: "3 months ago",
  userId: "user-123",
};

const mockEscrowTransactions = [
  {
    id: "1",
    businessId: "1",
    buyerId: "user-456",
    sellerId: "user-123",
    amount: 50000,
    description: "Website development project",
    status: "funded",
    releaseConditions: "Project completion and delivery",
    fundedAt: "2024-07-15T10:00:00Z",
    completedAt: null,
    cancelledAt: null,
    createdAt: "2024-07-15T09:00:00Z",
  },
  {
    id: "2",
    businessId: "1",
    buyerId: "user-789",
    sellerId: "user-123",
    amount: 25000,
    description: "Mobile app consultation",
    status: "completed",
    releaseConditions: "Consultation session delivered",
    fundedAt: "2024-07-10T14:00:00Z",
    completedAt: "2024-07-12T16:00:00Z",
    cancelledAt: null,
    createdAt: "2024-07-10T13:00:00Z",
  },
];

function BusinessDetail() {
  const { id } = Route.useParams();
  const [showEscrowModal, setShowEscrowModal] = useState(false);
  const [escrowAmount, setEscrowAmount] = useState("");
  const [escrowDescription, setEscrowDescription] = useState("");
  const [escrowConditions, setEscrowConditions] = useState("");

  const handleCreateEscrow = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Call API to create escrow transaction
    console.log("Creating escrow:", { amount: escrowAmount, description: escrowDescription, conditions: escrowConditions });
    setShowEscrowModal(false);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "funded":
        return <Clock className="size-4 text-amber-500" />;
      case "completed":
        return <CheckCircle className="size-4 text-green-500" />;
      case "cancelled":
        return <XCircle className="size-4 text-red-500" />;
      default:
        return <Lock className="size-4 text-muted-foreground" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "funded":
        return "bg-amber-500/10 text-amber-500";
      case "completed":
        return "bg-green-500/10 text-green-500";
      case "cancelled":
        return "bg-red-500/10 text-red-500";
      default:
        return "bg-muted text-muted-foreground";
    }
  };

  return (
    <div className="space-y-6">
      {/* Back button */}
      <Link to="/business" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="size-4" />
        Back to Directory
      </Link>

      {/* Business Header */}
      <motion.div
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        className="card-elev p-6"
      >
        <div className="flex items-start justify-between gap-6">
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <h1 className="display text-3xl font-semibold tracking-tight md:text-4xl">
                {mockBusiness.businessName}
              </h1>
              {mockBusiness.isVerified && (
                <Star className="size-5 fill-amber-500 text-amber-500" />
              )}
              {mockBusiness.isFeatured && (
                <span className="rounded-full bg-amber-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-amber-500">
                  Featured
                </span>
              )}
            </div>
            <div className="mt-4 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <MapPin className="size-4" />
                {mockBusiness.location}
              </div>
              <div className="flex items-center gap-2">
                <Phone className="size-4" />
                {mockBusiness.phoneNumber}
              </div>
              {mockBusiness.website && (
                <a
                  href={mockBusiness.website}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 hover:text-primary"
                >
                  <Globe className="size-4" />
                  <span className="truncate max-w-[200px]">{mockBusiness.website}</span>
                </a>
              )}
            </div>
            <p className="mt-4 text-sm text-muted-foreground">{mockBusiness.description}</p>
            <div className="mt-4 flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
              {mockBusiness.instagramHandle && (
                <div className="flex items-center gap-2">
                  <Instagram className="size-4" />
                  {mockBusiness.instagramHandle}
                </div>
              )}
              {mockBusiness.facebookHandle && (
                <div className="flex items-center gap-2">
                  <Facebook className="size-4" />
                  <span className="truncate max-w-[150px]">{mockBusiness.facebookHandle}</span>
                </div>
              )}
            </div>
          </div>
          <button
            onClick={() => setShowEscrowModal(true)}
            className="inline-flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
          >
            <Shield className="size-4" />
            Create Escrow Transaction
          </button>
        </div>
      </motion.div>

      {/* Escrow Transactions Section */}
      <motion.div
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold">Escrow Transactions</h2>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Lock className="size-4" />
            Secure payments with escrow protection
          </div>
        </div>
        <div className="space-y-3">
          {mockEscrowTransactions.map((transaction, i) => (
            <motion.div
              key={transaction.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.15 + 0.05 * i }}
              className="card-elev p-5"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-foreground">{transaction.description}</h3>
                    <span className={`rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider ${getStatusColor(transaction.status)}`}>
                      {transaction.status}
                    </span>
                  </div>
                  <div className="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <DollarSign className="size-3.5" />
                      KES {transaction.amount.toLocaleString()}
                    </div>
                    <div className="flex items-center gap-1">
                      {getStatusIcon(transaction.status)}
                      {transaction.status === "funded" && "Awaiting completion"}
                      {transaction.status === "completed" && "Completed and released"}
                      {transaction.status === "cancelled" && "Transaction cancelled"}
                    </div>
                  </div>
                  {transaction.releaseConditions && (
                    <div className="mt-2 text-xs text-muted-foreground">
                      <span className="font-medium">Release conditions:</span> {transaction.releaseConditions}
                    </div>
                  )}
                  <div className="mt-2 text-xs text-muted-foreground">
                    Created: {new Date(transaction.createdAt).toLocaleDateString()}
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
          {mockEscrowTransactions.length === 0 && (
            <div className="card-elev p-8 text-center">
              <Shield className="mx-auto size-12 text-muted-foreground" />
              <p className="mt-2 text-sm text-muted-foreground">No escrow transactions yet</p>
            </div>
          )}
        </div>
      </motion.div>

      {/* Create Escrow Modal */}
      {showEscrowModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="card-elev w-full max-w-lg p-6"
          >
            <div className="mb-4 flex items-center gap-2">
              <Shield className="size-5 text-primary" />
              <h2 className="text-xl font-semibold">Create Escrow Transaction</h2>
            </div>
            <div className="mb-4 rounded-lg bg-primary/5 p-3">
              <p className="text-xs text-muted-foreground">
                <strong>How escrow works:</strong> Your payment is held securely until the seller delivers the agreed-upon work or product. Once you're satisfied, the funds are released.
              </p>
            </div>
            <form onSubmit={handleCreateEscrow} className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium">Amount (KES)</label>
                <input
                  type="number"
                  value={escrowAmount}
                  onChange={(e) => setEscrowAmount(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., 50000"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Description</label>
                <input
                  type="text"
                  value={escrowDescription}
                  onChange={(e) => setEscrowDescription(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Website development project"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium">Release Conditions</label>
                <textarea
                  rows={3}
                  value={escrowConditions}
                  onChange={(e) => setEscrowConditions(e.target.value)}
                  className="w-full rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
                  placeholder="e.g., Project completion and delivery of source code"
                  required
                />
              </div>
              <div className="flex items-center gap-2 rounded-lg bg-amber-500/10 p-3">
                <Lock className="size-4 text-amber-500" />
                <p className="text-xs text-muted-foreground">
                  Your payment will be held securely until the release conditions are met.
                </p>
              </div>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowEscrowModal(false)}
                  className="flex-1 rounded-xl border border-border/60 bg-card px-4 py-2.5 text-sm font-medium hover:bg-muted"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 rounded-xl bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
                >
                  Create Transaction
                </button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
