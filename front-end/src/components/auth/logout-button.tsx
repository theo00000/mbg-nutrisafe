"use client";

import { useRouter } from "next/navigation";

import { clearAuth, getAuthToken } from "@/lib/auth-storage";
import { logout } from "@/services/auth-service";
import { Button } from "@/components/ui/button";

export function LogoutButton() {
  const router = useRouter();

  async function handleLogout() {
    const token = getAuthToken();

    try {
      if (token) {
        await logout(token);
      }
    } catch (error) {
      console.error("Logout API gagal:", error);
    } finally {
      clearAuth();
      router.push("/login");
    }
  }

  return (
    <Button variant="outline" onClick={handleLogout}>
      Keluar
    </Button>
  );
}
