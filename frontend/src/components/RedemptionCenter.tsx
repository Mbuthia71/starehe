import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ShoppingBag, Gift, ExternalLink, CheckCircle, AlertCircle } from 'lucide-react';
import { pointsApi, type Partner, type Redemption } from '../lib/pointsApi';
import { useState } from 'react';

export function RedemptionCenter() {
  const queryClient = useQueryClient();
  const [selectedPartner, setSelectedPartner] = useState<Partner | null>(null);
  const [redeeming, setRedeeming] = useState(false);

  const { data: partners = [], isLoading: partnersLoading } = useQuery<Partner[]>({
    queryKey: ['points-partners'],
    queryFn: pointsApi.getPartners,
  });

  const { data: redemptions = [] } = useQuery<Redemption[]>({
    queryKey: ['points-redemptions'],
    queryFn: () => pointsApi.getRedemptions(20).then(d => d.redemptions),
  });

  const redeemMutation = useMutation({
    mutationFn: ({ partnerId, pointsCost }: { partnerId: string; pointsCost: number }) =>
      pointsApi.redeemPoints(partnerId, pointsCost),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['points-balance'] });
      queryClient.invalidateQueries({ queryKey: ['points-redemptions'] });
      setSelectedPartner(null);
      setRedeeming(false);
    },
    onError: () => {
      setRedeeming(false);
    },
  });

  const handleRedeem = (partner: Partner) => {
    const pointsCost = 1000; // Default cost, should come from partner config
    if (window.confirm(`Redeem ${pointsCost} points for ${partner.name}?`)) {
      setRedeeming(true);
      redeemMutation.mutate({ partnerId: partner.id, pointsCost });
    }
  };

  return (
    <div className="space-y-4">
      <div className="card-elev p-4">
        <div className="flex items-center gap-2 mb-4">
          <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
            <ShoppingBag className="size-5" />
          </div>
          <h3 className="font-semibold">Redeem Points</h3>
        </div>

        {partnersLoading ? (
          <div className="space-y-3">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="animate-pulse h-20 bg-muted rounded-xl" />
            ))}
          </div>
        ) : partners.length === 0 ? (
          <div className="text-center py-8 text-sm text-muted-foreground">
            No redemption partners available
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {partners.map((partner) => (
              <div
                key={partner.id}
                className="p-4 bg-muted/50 rounded-xl space-y-3"
              >
                <div className="flex items-start gap-3">
                  {partner.logo_url ? (
                    <img
                      src={partner.logo_url}
                      alt={partner.name}
                      className="size-12 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="grid size-12 place-items-center rounded-xl bg-primary/10 text-primary">
                      <Gift className="size-6" />
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <h4 className="font-semibold truncate">{partner.name}</h4>
                    <p className="text-xs text-muted-foreground line-clamp-2">
                      {partner.description}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => handleRedeem(partner)}
                  disabled={redeeming}
                  className="w-full py-2 text-sm font-semibold bg-primary text-primary-foreground rounded-xl hover:brightness-110 disabled:opacity-50"
                >
                  {redeeming ? 'Processing...' : 'Redeem 1,000 pts'}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {redemptions.length > 0 && (
        <div className="card-elev p-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
              <Gift className="size-5" />
            </div>
            <h3 className="font-semibold">Recent Redemptions</h3>
          </div>

          <div className="space-y-3">
            {redemptions.slice(0, 5).map((redemption) => (
              <div
                key={redemption.id}
                className="flex items-center gap-3 p-3 bg-muted/50 rounded-xl"
              >
                <div className="grid size-10 place-items-center rounded-xl bg-background">
                  {redemption.status === 'completed' ? (
                    <CheckCircle className="size-4 text-green-500" />
                  ) : redemption.status === 'pending' ? (
                    <AlertCircle className="size-4 text-amber-500" />
                  ) : (
                    <AlertCircle className="size-4 text-red-500" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium">
                    {redemption.voucher_code || 'Voucher pending'}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {new Date(redemption.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold">
                    -{redemption.points_cost.toLocaleString()} pts
                  </p>
                  <p className="text-xs text-muted-foreground capitalize">
                    {redemption.status}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
