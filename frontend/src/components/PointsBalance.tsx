import { useQuery } from '@tanstack/react-query';
import { Coins, TrendingUp, Award } from 'lucide-react';
import { pointsApi, type PointsWallet } from '../lib/pointsApi';

export function PointsBalance() {
  const { data: wallet, isLoading, error } = useQuery<PointsWallet>({
    queryKey: ['points-balance'],
    queryFn: pointsApi.getBalance,
    refetchInterval: 30000, // Refresh every 30 seconds
  });

  if (isLoading) {
    return (
      <div className="card-elev p-4">
        <div className="animate-pulse h-20 bg-muted rounded-xl" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="card-elev p-4">
        <p className="text-sm text-destructive">Failed to load balance</p>
      </div>
    );
  }

  const availableBalance = wallet?.balance || 0;
  const heldBalance = wallet?.held_balance || 0;
  const totalBalance = availableBalance + heldBalance;

  return (
    <div className="card-elev p-4 space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
            <Coins className="size-5" />
          </div>
          <div>
            <h3 className="font-semibold">Points Balance</h3>
            <p className="text-xs text-muted-foreground">Current tier: {wallet?.tier || 'Bronze'}</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-primary">{totalBalance.toLocaleString()}</p>
          <p className="text-xs text-muted-foreground">total points</p>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="bg-muted/50 rounded-xl p-3">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp className="size-4 text-green-500" />
            <span className="text-xs text-muted-foreground">Available</span>
          </div>
          <p className="text-lg font-semibold">{availableBalance.toLocaleString()}</p>
        </div>
        <div className="bg-muted/50 rounded-xl p-3">
          <div className="flex items-center gap-2 mb-1">
            <Award className="size-4 text-amber-500" />
            <span className="text-xs text-muted-foreground">Pending</span>
          </div>
          <p className="text-lg font-semibold">{heldBalance.toLocaleString()}</p>
        </div>
      </div>
    </div>
  );
}
