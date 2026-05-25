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
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const token = getAuthToken();
    const role = getAuthRole();

    if (!token || !role) {
      router.push("/login");
      return;
    }

    if (allowedRoles && !allowedRoles.includes(role)) {
      router.push(getDashboardPathByRole(role));
      return;
    }

    setIsChecking(false);
  }, [allowedRoles, router]);

  if (isChecking) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50">
        <p className="text-sm text-slate-500">Memeriksa sesi login...</p>
      </main>
    );
  }

  return <>{children}</>;
}
