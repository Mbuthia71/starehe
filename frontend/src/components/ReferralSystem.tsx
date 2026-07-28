import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Users, Copy, Check, Share2, Gift } from 'lucide-react';
import { pointsApi, type Referral } from '../lib/pointsApi';
import { useState } from 'react';

export function ReferralSystem() {
  const queryClient = useQueryClient();
  const [copied, setCopied] = useState(false);

  const { data: referral, isLoading, error } = useQuery<Referral>({
    queryKey: ['points-referral'],
    queryFn: pointsApi.getReferral,
  });

  const createMutation = useMutation({
    mutationFn: () => pointsApi.createReferral(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['points-referral'] });
    },
  });

  const handleCopyCode = () => {
    if (referral?.referral_code) {
      navigator.clipboard.writeText(referral.referral_code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleShare = async () => {
    if (referral?.referral_code) {
      const shareText = `Join the Starehe Society Platform! Use my referral code: ${referral.referral_code}`;
      if (navigator.share) {
        try {
          await navigator.share({
            title: 'Join Starehe Society',
            text: shareText,
          });
        } catch (err) {
          console.log('Share failed:', err);
        }
      } else {
        handleCopyCode();
      }
    }
  };

  if (isLoading) {
    return (
      <div className="card-elev p-4">
        <div className="animate-pulse h-24 bg-muted rounded-xl" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="card-elev p-4">
        <p className="text-sm text-destructive">Failed to load referral info</p>
      </div>
    );
  }

  if (!referral) {
    return (
      <div className="card-elev p-4 space-y-4">
        <div className="flex items-center gap-2">
          <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
            <Users className="size-5" />
          </div>
          <h3 className="font-semibold">Refer & Earn</h3>
        </div>
        <div className="text-center py-4">
          <div className="grid size-16 place-items-center rounded-2xl bg-muted mx-auto mb-3">
            <Gift className="size-8 text-primary" />
          </div>
          <p className="text-sm text-muted-foreground mb-4">
            Invite friends and earn 500 points when they join!
          </p>
          <button
            onClick={() => createMutation.mutate()}
            disabled={createMutation.isPending}
            className="w-full py-3 text-sm font-semibold bg-primary text-primary-foreground rounded-xl hover:brightness-110 disabled:opacity-50"
          >
            {createMutation.isPending ? 'Creating...' : 'Get Referral Code'}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="card-elev p-4 space-y-4">
      <div className="flex items-center gap-2">
        <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
          <Users className="size-5" />
        </div>
        <h3 className="font-semibold">Your Referral Code</h3>
      </div>

      <div className="bg-muted/50 rounded-xl p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-2xl font-bold tracking-wider">{referral.referral_code}</p>
            <p className="text-xs text-muted-foreground">
              Status: <span className="capitalize">{referral.status}</span>
            </p>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleCopyCode}
              className="grid size-10 place-items-center rounded-xl bg-background hover:bg-muted transition-colors"
            >
              {copied ? <Check className="size-4 text-green-500" /> : <Copy className="size-4" />}
            </button>
            <button
              onClick={handleShare}
              className="grid size-10 place-items-center rounded-xl bg-background hover:bg-muted transition-colors"
            >
              <Share2 className="size-4" />
            </button>
          </div>
        </div>

        <div className="pt-3 border-t border-border">
          <div className="grid grid-cols-2 gap-3">
            <div className="text-center">
              <p className="text-lg font-semibold text-primary">500</p>
              <p className="text-xs text-muted-foreground">Points per referral</p>
            </div>
            <div className="text-center">
              <p className="text-lg font-semibold">{referral.points_awarded}</p>
              <p className="text-xs text-muted-foreground">Points earned</p>
            </div>
          </div>
        </div>
      </div>

      <div className="text-xs text-muted-foreground text-center">
        Share your code with friends to earn points when they join the platform
      </div>
    </div>
  );
}
