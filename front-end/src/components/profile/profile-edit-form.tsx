"use client";

import { FormEvent, useEffect, useState } from "react";

import { getAuthToken } from "@/lib/auth-storage";
import {
  getProfileDetail,
  updateProfileDetail,
} from "@/services/profile-service";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export function ProfileEditForm() {
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const token = getAuthToken();

    if (!token) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsLoading(false);
      return;
    }

    getProfileDetail(token)
      .then((profile) => {
        setName(profile.name || "");
        setPhone(profile.phone || "");
      })
      .catch((error) => {
        const text =
          error instanceof Error ? error.message : "Gagal memuat profil.";

        setErrorMessage(text);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const token = getAuthToken();

    if (!token) {
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      return;
    }

    setMessage("");
    setErrorMessage("");
    setIsSubmitting(true);

    try {
      await updateProfileDetail(token, {
        name,
        phone,
      });

      setMessage("Profil berhasil diperbarui.");
    } catch (error) {
      const text =
        error instanceof Error ? error.message : "Gagal memperbarui profil.";

      setErrorMessage(text);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Card className="border-orange-100 shadow-sm">
      <CardContent className="p-6">
        <h2 className="text-lg font-semibold text-slate-950">Edit Profil</h2>

        <p className="mt-2 text-sm text-slate-500">
          Perbarui informasi dasar akun pengguna.
        </p>

        {isLoading ? (
          <p className="mt-5 text-sm text-slate-500">Memuat data...</p>
        ) : (
          <form onSubmit={handleSubmit} className="mt-5 space-y-4">
            <div className="space-y-2">
              <Label>Nama</Label>
              <Input
                value={name}
                onChange={(event) => setName(event.target.value)}
                placeholder="Nama lengkap"
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Nomor Telepon</Label>
              <Input
                value={phone}
                onChange={(event) => setPhone(event.target.value)}
                placeholder="08xxxxxxxxxx"
              />
            </div>

            {message ? (
              <p className="rounded-xl bg-emerald-50 px-3 py-2 text-sm text-emerald-700">
                {message}
              </p>
            ) : null}

            {errorMessage ? (
              <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-600">
                {errorMessage}
              </p>
            ) : null}

            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Menyimpan..." : "Simpan Perubahan"}
            </Button>
          </form>
        )}
      </CardContent>
    </Card>
  );
}
