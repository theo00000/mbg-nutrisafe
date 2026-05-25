import { ReactNode } from "react";
import Link from "next/link";
import { ShieldCheck } from "lucide-react";

type AuthShellProps = {
  title: string;
  description: string;
  children: ReactNode;
};

export function AuthShell({ title, description, children }: AuthShellProps) {
  return (
    <main className="grid min-h-screen grid-cols-[0.95fr_1.05fr] bg-gradient-to-br from-sky-50 via-white to-orange-50">
      <section className="flex flex-col justify-between px-12 py-10">
        <Link href="/" className="flex items-center gap-3">
          <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-600 text-white shadow-sm">
            <ShieldCheck className="h-6 w-6" />
          </div>

          <div>
            <p className="text-lg font-bold text-slate-950">NutriSafe</p>
            <p className="text-xs text-slate-500">MBG Food Safety Platform</p>
          </div>
        </Link>

        <div className="max-w-xl">
          <p className="text-sm font-semibold text-blue-600">
            Platform Keamanan Makanan MBG
          </p>

          <h1 className="mt-5 text-5xl font-bold tracking-tight text-slate-950">
            Monitoring makanan lebih aman untuk sekolah dan siswa.
          </h1>

          <p className="mt-5 text-base leading-7 text-slate-600">
            NutriSafe membantu mencatat profil, memantau alergi, dan mendukung
            proses distribusi makanan yang lebih aman.
          </p>
        </div>

        <p className="text-sm text-slate-500">
          © 2026 MBG NutriSafe. Hackathon Project.
        </p>
      </section>

      <section className="flex max-h-screen items-start justify-center overflow-y-auto px-12 py-10">
        <div className="w-full max-w-xl rounded-3xl border border-blue-100 bg-white/90 p-8 shadow-xl backdrop-blur">
          <p className="text-sm font-semibold text-blue-600">MBG NutriSafe</p>

          <h2 className="mt-3 text-3xl font-bold tracking-tight text-slate-950">
            {title}
          </h2>

          <p className="mt-2 text-sm leading-6 text-slate-600">{description}</p>

          <div className="mt-6">{children}</div>
        </div>
      </section>
    </main>
  );
}
