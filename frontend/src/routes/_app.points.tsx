import { createFileRoute } from '@tanstack/react-router';
import { PointsBalance } from '../components/PointsBalance';
import { TransactionHistory } from '../components/TransactionHistory';
import { RedemptionCenter } from '../components/RedemptionCenter';
import { ReferralSystem } from '../components/ReferralSystem';
import { BadgesAndTiers } from '../components/BadgesAndTiers';
import { Campaigns } from '../components/Campaigns';
import { Coins, History, ShoppingBag, Users, Award, Flame } from 'lucide-react';
import { useState } from 'react';

export const Route = createFileRoute('/_app/points')({
  component: PointsPage,
});

type Tab = 'balance' | 'history' | 'redeem' | 'referral' | 'badges' | 'campaigns';

const tabs: { id: Tab; label: string; icon: any }[] = [
  { id: 'balance', label: 'Balance', icon: Coins },
  { id: 'history', label: 'History', icon: History },
  { id: 'redeem', label: 'Redeem', icon: ShoppingBag },
  { id: 'referral', label: 'Referral', icon: Users },
  { id: 'badges', label: 'Badges', icon: Award },
  { id: 'campaigns', label: 'Campaigns', icon: Flame },
];

function PointsPage() {
  const [activeTab, setActiveTab] = useState<Tab>('balance');

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Points & Rewards</h1>
        <p className="text-muted-foreground">Earn, redeem, and track your points</p>
      </div>

      {/* Mobile Tab Navigation */}
      <div className="md:hidden overflow-x-auto pb-2">
        <div className="flex gap-2 min-w-max">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
                  activeTab === tab.id
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted text-muted-foreground hover:bg-muted/80'
                }`}
              >
                <Icon className="size-4" />
                {tab.label}
              </button>
            );
          })}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Always show balance */}
        <div className="lg:col-span-1 space-y-4">
          <PointsBalance />
          <ReferralSystem />
        </div>

        {/* Right Column - Tab Content */}
        <div className="lg:col-span-2">
          {/* Desktop Tab Navigation */}
          <div className="hidden md:flex gap-2 mb-4">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
                    activeTab === tab.id
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-muted text-muted-foreground hover:bg-muted/80'
                  }`}
                >
                  <Icon className="size-4" />
                  {tab.label}
                </button>
              );
            })}
          </div>

          {/* Tab Content */}
          <div className="md:hidden">
            {activeTab === 'balance' && <PointsBalance />}
            {activeTab === 'history' && <TransactionHistory />}
            {activeTab === 'redeem' && <RedemptionCenter />}
            {activeTab === 'referral' && <ReferralSystem />}
            {activeTab === 'badges' && <BadgesAndTiers />}
            {activeTab === 'campaigns' && <Campaigns />}
          </div>

          <div className="hidden md:block">
            {activeTab === 'history' && <TransactionHistory />}
            {activeTab === 'redeem' && <RedemptionCenter />}
            {activeTab === 'badges' && <BadgesAndTiers />}
            {activeTab === 'campaigns' && <Campaigns />}
          </div>
        </div>
      </div>
    </div>
  );
}
