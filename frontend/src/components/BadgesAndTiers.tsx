import { useQuery } from '@tanstack/react-query';
import { Award, Trophy, Star, Crown } from 'lucide-react';
import { pointsApi, type UserBadge, type UserTier } from '../lib/pointsApi';

export function BadgesAndTiers() {
  const { data: badges = [], isLoading: badgesLoading } = useQuery<UserBadge[]>({
    queryKey: ['points-badges'],
    queryFn: pointsApi.getBadges,
  });

  const { data: tier, isLoading: tierLoading } = useQuery<UserTier>({
    queryKey: ['points-tier'],
    queryFn: pointsApi.getTier,
  });

  const getTierIcon = (tierName: string) => {
    switch (tierName?.toLowerCase()) {
      case 'platinum':
        return <Crown className="size-6 text-purple-500" />;
      case 'gold':
        return <Trophy className="size-6 text-amber-500" />;
      case 'silver':
        return <Star className="size-6 text-gray-400" />;
      default:
        return <Award className="size-6 text-amber-700" />;
    }
  };

  const getTierColor = (tierName: string) => {
    switch (tierName?.toLowerCase()) {
      case 'platinum':
        return 'bg-purple-500/10 text-purple-500 border-purple-500/20';
      case 'gold':
        return 'bg-amber-500/10 text-amber-500 border-amber-500/20';
      case 'silver':
        return 'bg-gray-400/10 text-gray-400 border-gray-400/20';
      default:
        return 'bg-amber-700/10 text-amber-700 border-amber-700/20';
    }
  };

  if (badgesLoading || tierLoading) {
    return (
      <div className="card-elev p-4">
        <div className="space-y-4">
          <div className="animate-pulse h-24 bg-muted rounded-xl" />
          <div className="animate-pulse h-32 bg-muted rounded-xl" />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Current Tier */}
      <div className={`card-elev p-4 border-2 ${getTierColor(tier?.tier || 'bronze')}`}>
        <div className="flex items-center gap-3 mb-4">
          <div className="grid size-12 place-items-center rounded-xl bg-background">
            {getTierIcon(tier?.tier || 'bronze')}
          </div>
          <div>
            <h3 className="font-semibold capitalize">{tier?.tier || 'Bronze'} Tier</h3>
            <p className="text-xs text-muted-foreground">
              Achieved {new Date(tier?.achieved_at || '').toLocaleDateString()}
            </p>
          </div>
        </div>

        {tier?.benefits && tier.benefits.length > 0 && (
          <div className="space-y-2">
            <p className="text-xs font-medium text-muted-foreground">Benefits:</p>
            <ul className="space-y-1">
              {tier.benefits.map((benefit, idx) => (
                <li key={idx} className="text-xs flex items-center gap-2">
                  <Star className="size-3 fill-current" />
                  {benefit}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* Badges */}
      <div className="card-elev p-4">
        <div className="flex items-center gap-2 mb-4">
          <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
            <Award className="size-5" />
          </div>
          <h3 className="font-semibold">Your Badges</h3>
        </div>

        {badges.length === 0 ? (
          <div className="text-center py-6 text-sm text-muted-foreground">
            No badges earned yet. Complete activities to earn badges!
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
            {badges.map((badge) => (
              <div
                key={badge.id}
                className="p-3 bg-muted/50 rounded-xl text-center space-y-2"
              >
                {badge.badge_icon ? (
                  <img
                    src={badge.badge_icon}
                    alt={badge.badge_name}
                    className="size-12 mx-auto rounded-lg object-cover"
                  />
                ) : (
                  <div className="grid size-12 place-items-center rounded-xl bg-primary/10 text-primary mx-auto">
                    <Award className="size-6" />
                  </div>
                )}
                <div>
                  <p className="text-sm font-medium truncate">{badge.badge_name}</p>
                  <p className="text-[10px] text-muted-foreground">
                    {new Date(badge.earned_at).toLocaleDateString()}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
