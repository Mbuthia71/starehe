import { Outlet, createFileRoute, useRouter } from '@tanstack/react-router';
import { Sidebar } from '../../components/dashboard/Sidebar';
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
      <Sidebar />
      <main className="flex-1 overflow-hidden">
        <Outlet />
      </main>
    </div>
  );
}
