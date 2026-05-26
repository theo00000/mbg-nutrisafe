"use client";

import { getDashboardPathByRole, saveAuth } from "@/lib/auth-storage";
import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { login } from "@/services/auth-service";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AuthShell } from "@/components/auth/auth-shell";

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
    <AuthShell
      title="Masuk Akun"
      description="Masuk untuk mengakses dashboard sesuai role pengguna."
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label>Email</Label>
          <Input
            type="email"
            placeholder="nama@email.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            className="h-11 rounded-2xl"
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
            className="h-11 rounded-2xl"
            required
          />
        </div>

        {errorMessage ? (
          <p className="rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-600">
            {errorMessage}
          </p>
        ) : null}

        <Button
          type="submit"
          className="h-11 w-full rounded-full bg-blue-600 font-semibold hover:bg-blue-700"
          disabled={isLoading}
        >
          {isLoading ? "Memproses..." : "Masuk"}
        </Button>

        <p className="text-center text-sm text-slate-600">
          Belum punya akun?{" "}
          <Link href="/register" className="font-semibold text-blue-600">
            Daftar di sini
          </Link>
        </p>
      </form>
    </AuthShell>
  );
}
