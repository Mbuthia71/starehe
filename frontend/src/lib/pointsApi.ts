const getToken = () => {
  if (typeof window === 'undefined') return '';
  return localStorage.getItem('oss_token') || '';
};

const authHeaders = () => ({
  'Content-Type': 'application/json',
  Authorization: `Bearer ${getToken()}`,
});

// Types
export type PointsWallet = {
  user_id: string;
  balance: number;
  held_balance: number;
  tier: string;
  created_at: string;
  updated_at: string;
};

export type PointsTransaction = {
  id: string;
  user_id: string;
  type: string;
  amount: number;
  balance_after: number;
  description: string;
  source: string | null;
  metadata: Record<string, any> | null;
  created_at: string;
};

export type Partner = {
  id: string;
  name: string;
  description: string;
  logo_url: string | null;
  website: string | null;
  integration_mode: string;
  api_config: Record<string, any> | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type Redemption = {
  id: string;
  user_id: string;
  partner_id: string;
  points_cost: number;
  currency_value: number;
  status: string;
  voucher_code: string | null;
  partner_reference: string | null;
  metadata: Record<string, any> | null;
  created_at: string;
  updated_at: string;
};

export type Referral = {
  id: string;
  referrer_id: string;
  referred_user_id: string | null;
  referral_code: string;
  status: string;
  points_awarded: number;
  referred_at: string | null;
  completed_at: string | null;
  created_at: string;
};

export type UserBadge = {
  id: string;
  user_id: string;
  badge_id: string;
  badge_name: string;
  badge_description: string;
  badge_icon: string | null;
  earned_at: string;
};

export type UserTier = {
  id: string;
  user_id: string;
  tier: string;
  points_threshold: number;
  benefits: string[];
  achieved_at: string;
};

export type PointsCampaign = {
  id: string;
  name: string;
  description: string;
  points_reward: number;
  start_date: string;
  end_date: string;
  is_active: boolean;
  requirements: Record<string, any>;
};

// API Functions
export const pointsApi = {
  // Get user's points balance
  getBalance: async (): Promise<PointsWallet> => {
    const res = await fetch('/api/points/balance', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch balance');
    return res.json();
  },

  // Get transaction history
  getTransactions: async (limit = 20, cursor?: string): Promise<{ transactions: PointsTransaction[]; next_cursor?: string }> => {
    const url = cursor 
      ? `/api/points/transactions?limit=${limit}&cursor=${cursor}`
      : `/api/points/transactions?limit=${limit}`;
    const res = await fetch(url, {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch transactions');
    return res.json();
  },

  // Get available partners for redemption
  getPartners: async (): Promise<Partner[]> => {
    const res = await fetch('/api/points/partners', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch partners');
    const data = await res.json();
    return data.partners || [];
  },

  // Redeem points
  redeemPoints: async (partnerId: string, pointsCost: number): Promise<Redemption> => {
    const res = await fetch('/api/points/redeem', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({
        partner_id: partnerId,
        points_cost: pointsCost,
      }),
    });
    if (!res.ok) throw new Error('Failed to redeem points');
    return res.json();
  },

  // Get redemption history
  getRedemptions: async (limit = 20, cursor?: string): Promise<{ redemptions: Redemption[]; next_cursor?: string }> => {
    const url = cursor
      ? `/api/points/redemptions?limit=${limit}&cursor=${cursor}`
      : `/api/points/redemptions?limit=${limit}`;
    const res = await fetch(url, {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch redemptions');
    return res.json();
  },

  // Create referral
  createReferral: async (): Promise<Referral> => {
    const res = await fetch('/api/points/referrals', {
      method: 'POST',
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to create referral');
    return res.json();
  },

  // Get referral info
  getReferral: async (): Promise<Referral> => {
    const res = await fetch('/api/points/referrals', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch referral');
    return res.json();
  },

  // Get user badges
  getBadges: async (): Promise<UserBadge[]> => {
    const res = await fetch('/api/points/badges', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch badges');
    const data = await res.json();
    return data.badges || [];
  },

  // Get user tier
  getTier: async (): Promise<UserTier> => {
    const res = await fetch('/api/points/tier', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch tier');
    return res.json();
  },

  // Get available campaigns
  getCampaigns: async (): Promise<PointsCampaign[]> => {
    const res = await fetch('/api/points/campaigns', {
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch campaigns');
    const data = await res.json();
    return data.campaigns || [];
  },

  // Join campaign
  joinCampaign: async (campaignId: string): Promise<void> => {
    const res = await fetch(`/api/points/campaigns/${campaignId}/join`, {
      method: 'POST',
      headers: authHeaders(),
    });
    if (!res.ok) throw new Error('Failed to join campaign');
  },
};
