"use client";

import { useEffect, useState } from "react";

import { getAuthToken } from "@/lib/auth-storage";
import { getProfileDetail, ProfileDetail } from "@/services/profile-service";
import { Card, CardContent } from "@/components/ui/card";

export function ProfileCard() {
  const [profile, setProfile] = useState<ProfileDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const token = getAuthToken();

    if (!token) {
      // Auth berbasis localStorage perlu update state setelah komponen mount.
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");

      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsLoading(false);

      return;
    }

    getProfileDetail(token)
      .then((data) => {
        console.log("PROFILE RESPONSE:", data);
        setProfile(data);
      })
      .catch((error) => {
        const message =
          error instanceof Error ? error.message : "Gagal mengambil profil.";

        setErrorMessage(message);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  return (
    <Card className="border-orange-100 shadow-sm">
      <CardContent className="p-6">
        <h2 className="text-lg font-semibold text-slate-950">
          Profil Pengguna
        </h2>

        <p className="mt-2 text-sm text-slate-500">
          Data ini diambil dari endpoint GET /api/profile/detail.
        </p>

        {isLoading ? (
          <p className="mt-5 text-sm text-slate-500">Memuat profil...</p>
        ) : null}

        {errorMessage ? (
          <p className="mt-5 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-600">
            {errorMessage}
          </p>
        ) : null}

        {profile ? (
          <div className="mt-5 grid grid-cols-2 gap-4">
            <div className="rounded-2xl bg-orange-50 p-4">
              <p className="text-xs text-slate-500">Nama</p>
              <p className="mt-1 text-sm font-semibold text-slate-900">
                {profile.name}
              </p>
            </div>

            <div className="rounded-2xl bg-orange-50 p-4">
              <p className="text-xs text-slate-500">Email</p>
              <p className="mt-1 text-sm font-semibold text-slate-900">
                {profile.email}
              </p>
            </div>

            <div className="rounded-2xl bg-orange-50 p-4">
              <p className="text-xs text-slate-500">Telepon</p>
              <p className="mt-1 text-sm font-semibold text-slate-900">
                {profile.phone || "-"}
              </p>
            </div>

            <div className="rounded-2xl bg-orange-50 p-4">
              <p className="text-xs text-slate-500">Role</p>
              <p className="mt-1 text-sm font-semibold text-slate-900">
                {profile.role || profile.role_name || "-"}
              </p>
            </div>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
