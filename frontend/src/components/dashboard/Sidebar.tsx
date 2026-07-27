import { Link, useLocation } from '@tanstack/react-router';
import {
  MessageCircle,
  Home,
  Users,
  LogOut,
  Settings,
  Bell,
  Search,
} from 'lucide-react';
import { useState } from 'react';

export function Sidebar() {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(true);

  const isActive = (path: string) => location.pathname.includes(path);

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    window.location.href = '/auth/login';
  };

  return (
    <aside
      className={`bg-sidebar border-r transition-all duration-200 flex flex-col ${
        isOpen ? 'w-64' : 'w-20'
      }`}
    >
      {/* Header */}
      <div className="p-4 border-b flex items-center justify-between">
        {isOpen && (
          <h1 className="font-bold text-lg tracking-tight">Starehian</h1>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
        <NavItem
          icon={<Home className="w-5 h-5" />}
          label="Dashboard"
          href="/dashboard"
          isActive={isActive('/dashboard') && !isActive('/dashboard/messages')}
          isOpen={isOpen}
        />
        <NavItem
          icon={<Users className="w-5 h-5" />}
          label="Connections"
          href="/dashboard/connections"
          isActive={isActive('/connections')}
          isOpen={isOpen}
        />
        <NavItem
          icon={<MessageCircle className="w-5 h-5" />}
          label="Messages"
          href="/dashboard/messages"
          isActive={isActive('/messages')}
          isOpen={isOpen}
          badge="new"
        />
        <NavItem
          icon={<Bell className="w-5 h-5" />}
          label="Notifications"
          href="/dashboard/notifications"
          isActive={isActive('/notifications')}
          isOpen={isOpen}
        />
        <NavItem
          icon={<Search className="w-5 h-5" />}
          label="Directory"
          href="/dashboard/directory"
          isActive={isActive('/directory')}
          isOpen={isOpen}
        />
      </nav>

      {/* Footer */}
      <div className="border-t p-4 space-y-2">
        <NavItem
          icon={<Settings className="w-5 h-5" />}
          label="Settings"
          href="/dashboard/settings"
          isActive={isActive('/settings')}
          isOpen={isOpen}
        />
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-red-600 hover:bg-red-50 transition-colors"
        >
          <LogOut className="w-5 h-5" />
          {isOpen && <span>Logout</span>}
        </button>
      </div>
    </aside>
  );
}

interface NavItemProps {
  icon: React.ReactNode;
  label: string;
  href: string;
  isActive: boolean;
  isOpen: boolean;
  badge?: string;
}

function NavItem({
  icon,
  label,
  href,
  isActive,
  isOpen,
  badge,
}: NavItemProps) {
  return (
    <Link
      to={href}
      className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors relative ${
        isActive
          ? 'bg-primary text-primary-foreground'
          : 'text-foreground hover:bg-accent'
      }`}
      title={!isOpen ? label : undefined}
    >
      {icon}
      {isOpen && (
        <>
          <span>{label}</span>
          {badge && (
            <span className="ml-auto text-xs bg-red-500 text-white rounded-full px-2 py-0.5">
              {badge}
            </span>
          )}
        </>
      )}
    </Link>
  );
}
