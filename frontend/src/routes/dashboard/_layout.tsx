import { Outlet, createFileRoute, useRouter } from '@tanstack/react-router';
import { useEffect } from 'react';

export const Route = createFileRoute('/dashboard/_layout')();

function DashboardLayout() {
  const router = useRouter();

  useEffect(() => {
    // Check if user is authenticated
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.navigate({ to: '/auth/login' });
    }
  }, [router]);

  return (
    <div className="flex h-screen">
      {/* Placeholder for Sidebar - ensure it exists or create it */}
      <aside className="w-64 border-r bg-slate-50">
        <div className="p-4 border-b font-bold">Starehian</div>
      </aside>
      <main className="flex-1 overflow-hidden">
        <Outlet />
      </main>
    </div>
  );
}
