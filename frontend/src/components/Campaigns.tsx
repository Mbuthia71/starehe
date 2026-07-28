import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Target, Calendar, CheckCircle, Flame } from 'lucide-react';
import { pointsApi, type PointsCampaign } from '../lib/pointsApi';

export function Campaigns() {
  const queryClient = useQueryClient();

  const { data: campaigns = [], isLoading } = useQuery<PointsCampaign[]>({
    queryKey: ['points-campaigns'],
    queryFn: pointsApi.getCampaigns,
  });

  const joinMutation = useMutation({
    mutationFn: (campaignId: string) => pointsApi.joinCampaign(campaignId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['points-campaigns'] });
      queryClient.invalidateQueries({ queryKey: ['points-balance'] });
    },
  });

  const isCampaignActive = (campaign: PointsCampaign) => {
    const now = new Date();
    const start = new Date(campaign.start_date);
    const end = new Date(campaign.end_date);
    return campaign.is_active && now >= start && now <= end;
  };

  const isCampaignEnded = (campaign: PointsCampaign) => {
    return new Date() > new Date(campaign.end_date);
  };

  if (isLoading) {
    return (
      <div className="card-elev p-4">
        <div className="space-y-3">
          {[...Array(2)].map((_, i) => (
            <div key={i} className="animate-pulse h-28 bg-muted rounded-xl" />
          ))}
        </div>
      </div>
    );
  }

  const activeCampaigns = campaigns.filter(isCampaignActive);
  const endedCampaigns = campaigns.filter(isCampaignEnded);

  return (
    <div className="card-elev p-4 space-y-4">
      <div className="flex items-center gap-2">
        <div className="grid size-10 place-items-center rounded-xl bg-primary/10 text-primary">
          <Flame className="size-5" />
        </div>
        <h3 className="font-semibold">Active Campaigns</h3>
      </div>

      {activeCampaigns.length === 0 && endedCampaigns.length === 0 ? (
        <div className="text-center py-8 text-sm text-muted-foreground">
          No active campaigns at the moment
        </div>
      ) : (
        <div className="space-y-3">
          {activeCampaigns.map((campaign) => (
            <div
              key={campaign.id}
              className="p-4 bg-gradient-to-br from-primary/10 to-primary/5 rounded-xl space-y-3"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h4 className="font-semibold">{campaign.name}</h4>
                  <p className="text-sm text-muted-foreground mt-1">{campaign.description}</p>
                </div>
                <div className="text-right">
                  <p className="text-lg font-bold text-primary">+{campaign.points_reward}</p>
                  <p className="text-xs text-muted-foreground">points</p>
                </div>
              </div>

              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Calendar className="size-3" />
                <span>
                  {new Date(campaign.start_date).toLocaleDateString()} - {new Date(campaign.end_date).toLocaleDateString()}
                </span>
              </div>

              <button
                onClick={() => joinMutation.mutate(campaign.id)}
                disabled={joinMutation.isPending}
                className="w-full py-2 text-sm font-semibold bg-primary text-primary-foreground rounded-xl hover:brightness-110 disabled:opacity-50"
              >
                {joinMutation.isPending ? 'Joining...' : 'Join Campaign'}
              </button>
            </div>
          ))}

          {endedCampaigns.length > 0 && (
            <div className="pt-4 border-t border-border">
              <p className="text-xs font-medium text-muted-foreground mb-3">Ended Campaigns</p>
              {endedCampaigns.map((campaign) => (
                <div
                  key={campaign.id}
                  className="p-3 bg-muted/50 rounded-xl opacity-60"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm font-medium">{campaign.name}</h4>
                      <p className="text-xs text-muted-foreground mt-1">
                        Ended {new Date(campaign.end_date).toLocaleDateString()}
                      </p>
                    </div>
                    <CheckCircle className="size-4 text-muted-foreground" />
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
