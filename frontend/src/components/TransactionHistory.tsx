import { useQuery } from '@tanstack/react-query';
import { ArrowUpRight, ArrowDownLeft, Clock, CheckCircle, XCircle } from 'lucide-react';
import { pointsApi, type PointsTransaction } from '../lib/pointsApi';
import { useState } from 'react';

export function TransactionHistory() {
  const [cursor, setCursor] = useState<string | undefined>();
  const { data, isLoading, error } = useQuery<{
    transactions: PointsTransaction[];
    next_cursor?: string;
  }>({
    queryKey: ['points-transactions', cursor],
    queryFn: () => pointsApi.getTransactions(20, cursor),
  });

  const transactions = data?.transactions || [];
  const nextCursor = data?.next_cursor;

  const getTransactionIcon = (type: string) => {
    if (type === 'earn' || type === 'referral' || type === 'campaign') {
      return <ArrowUpRight className="size-4 text-green-500" />;
    }
    return <ArrowDownLeft className="size-4 text-red-500" />;
  };

  const getTransactionColor = (type: string) => {
    if (type === 'earn' || type === 'referral' || type === 'campaign') {
      return 'text-green-500';
    }
    return 'text-red-500';
  };

  if (isLoading && transactions.length === 0) {
    return (
      <div className="card-elev p-4">
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="animate-pulse h-16 bg-muted rounded-xl" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card-elev p-4">
        <p className="text-sm text-destructive">Failed to load transactions</p>
      </div>
    );
  }

  return (
    <div className="card-elev p-4 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold">Transaction History</h3>
        <Clock className="size-4 text-muted-foreground" />
      </div>

      {transactions.length === 0 ? (
        <div className="text-center py-8 text-sm text-muted-foreground">
          No transactions yet
        </div>
      ) : (
        <div className="space-y-3">
          {transactions.map((tx) => (
            <div
              key={tx.id}
              className="flex items-center gap-3 p-3 bg-muted/50 rounded-xl"
            >
              <div className="grid size-10 place-items-center rounded-xl bg-background">
                {getTransactionIcon(tx.type)}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">{tx.description}</p>
                <p className="text-xs text-muted-foreground">
                  {new Date(tx.created_at).toLocaleDateString()}
                </p>
              </div>
              <div className="text-right">
                <p className={`text-sm font-semibold ${getTransactionColor(tx.type)}`}>
                  {tx.type === 'earn' || tx.type === 'referral' || tx.type === 'campaign' ? '+' : '-'}
                  {Math.abs(tx.amount).toLocaleString()}
                </p>
                <p className="text-xs text-muted-foreground">
                  Balance: {tx.balance_after.toLocaleString()}
                </p>
              </div>
            </div>
          ))}

          {nextCursor && (
            <button
              onClick={() => setCursor(nextCursor)}
              className="w-full py-2 text-sm text-primary hover:underline"
            >
              Load more
            </button>
          )}
        </div>
      )}
    </div>
  );
}
