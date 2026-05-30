import { ReactNode } from "react";

import { RequireAuth } from "@/components/auth/require-auth";
import type { UserRole } from "@/lib/auth-storage";

type DashboardShellProps = {
  allowedRoles: UserRole[];
  sidebar: ReactNode;
  children: ReactNode;
};

export function DashboardShell({
  allowedRoles,
  sidebar,
  children,
}: DashboardShellProps) {
  return (
    <RequireAuth allowedRoles={allowedRoles}>
      <main className="flex min-h-screen bg-[#f8fbff]">
        {sidebar}

        <section className="flex-1 overflow-y-auto p-8">{children}</section>
      </main>
    </RequireAuth>
  );
}
