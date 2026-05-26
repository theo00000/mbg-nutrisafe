"use client";

import { ReactNode, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import {
  getAuthRole,
  getAuthToken,
  getDashboardPathByRole,
  UserRole,
} from "@/lib/auth-storage";

type RequireAuthProps = {
  children: ReactNode;
  allowedRoles?: UserRole[];
};

export function RequireAuth({ children, allowedRoles }: RequireAuthProps) {
  const router = useRouter();
  const [isAllowed, setIsAllowed] = useState(false);

  useEffect(() => {
    const token = getAuthToken();
    const role = getAuthRole();

    if (!token || !role) {
      router.replace("/login");
      return;
    }

    if (allowedRoles && !allowedRoles.includes(role)) {
      router.replace(getDashboardPathByRole(role));
      return;
    }

    // Auth guard sederhana berbasis localStorage perlu update state setelah mount.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setIsAllowed(true);
  }, [allowedRoles, router]);

  if (!isAllowed) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50">
        <p className="text-sm text-slate-500">Memeriksa sesi login...</p>
      </main>
    );
  }

  return <>{children}</>;
}
