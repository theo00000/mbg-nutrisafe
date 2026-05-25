"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { register } from "@/services/auth-service";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type RoleName = "school" | "spgg" | "umum";

export default function RegisterPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [roleName, setRoleName] = useState<RoleName>("school");

  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setErrorMessage("");
    setSuccessMessage("");
    setIsLoading(true);

    try {
      const response = await register({
        name,
        email,
        phone,
        password,
        role_name: roleName,
      });

      setSuccessMessage(response.message || "Akun berhasil dibuat.");

      setTimeout(() => {
        router.push("/login");
      }, 1000);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Registrasi gagal. Coba lagi.";

      setErrorMessage(message);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-sky-50 px-10">
      <Card className="w-full max-w-md border-blue-100 shadow-sm">
        <CardContent className="p-8">
          <p className="text-sm font-semibold text-blue-600">MBG NutriSafe</p>

          <h1 className="mt-3 text-3xl font-bold tracking-tight text-slate-950">
            Daftar Akun
          </h1>

          <p className="mt-2 text-sm text-slate-600">
            Buat akun sesuai peran pengguna di ekosistem MBG NutriSafe.
          </p>

          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
            <div className="space-y-2">
              <Label>Nama</Label>
              <Input
                placeholder="Nama lengkap"
                value={name}
                onChange={(event) => setName(event.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Email</Label>
              <Input
                type="email"
                placeholder="nama@email.com"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Nomor Telepon</Label>
              <Input
                placeholder="08xxxxxxxxxx"
                value={phone}
                onChange={(event) => setPhone(event.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Password</Label>
              <Input
                type="password"
                placeholder="Minimal 6 karakter"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>Role</Label>
              <select
                value={roleName}
                onChange={(event) =>
                  setRoleName(event.target.value as RoleName)
                }
                className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
              >
                <option value="school">Sekolah</option>
                <option value="spgg">Mitra/SPGG</option>
                <option value="umum">Umum/Siswa</option>
              </select>
            </div>

            {errorMessage ? (
              <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-600">
                {errorMessage}
              </p>
            ) : null}

            {successMessage ? (
              <p className="rounded-xl bg-emerald-50 px-3 py-2 text-sm text-emerald-700">
                {successMessage}
              </p>
            ) : null}

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Mendaftarkan..." : "Daftar Akun"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
