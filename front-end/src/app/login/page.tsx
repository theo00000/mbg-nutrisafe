"use client";

import { getDashboardPathByRole, saveAuth } from "@/lib/auth-storage";
import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { login } from "@/services/auth-service";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function LoginPage() {
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setErrorMessage("");
    setIsLoading(true);

    try {
      const response = await login({
        email,
        password,
      });

      saveAuth(response.token, response.role);
      router.push(getDashboardPathByRole(response.role));
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Login gagal. Coba lagi.";

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
            Masuk Akun
          </h1>

          <p className="mt-2 text-sm text-slate-600">
            Masuk untuk mengakses dashboard sesuai role pengguna.
          </p>

          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
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
              <Label>Password</Label>
              <Input
                type="password"
                placeholder="Masukkan password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>

            {errorMessage ? (
              <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-600">
                {errorMessage}
              </p>
            ) : null}

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Memproses..." : "Masuk"}
            </Button>

            <p className="text-center text-sm text-slate-600">
              Belum punya akun?{" "}
              <Link href="/register" className="font-semibold text-blue-600">
                Daftar di sini
              </Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
