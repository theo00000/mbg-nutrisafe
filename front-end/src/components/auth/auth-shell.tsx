import { ReactNode } from "react";
import Link from "next/link";
import { ArrowLeft, CheckCircle2, ShieldCheck, Sparkles } from "lucide-react";

type AuthShellProps = {
  title: string;
  description: string;
  children: ReactNode;
};

export function AuthShell({ title, description, children }: AuthShellProps) {
  return (
    <main className="grid min-h-screen grid-cols-[0.95fr_1.05fr] overflow-hidden bg-[#f8fbff] text-slate-950">
      <section className="relative flex flex-col justify-between px-12 py-10">
        <div className="absolute left-[-120px] top-[-120px] h-80 w-80 rounded-full bg-blue-200/50 blur-3xl" />
        <div className="absolute bottom-[-140px] right-[-80px] h-96 w-96 rounded-full bg-orange-200/50 blur-3xl" />

        <div className="relative z-10">
          <Link href="/" className="inline-flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-blue-600 text-white shadow-lg shadow-blue-200">
              <ShieldCheck className="h-6 w-6" />
            </div>

            <div>
              <p className="text-lg font-black tracking-tight">NutriSafe</p>
              <p className="text-xs font-medium text-slate-500">
                MBG Safety Platform
              </p>
            </div>
          </Link>
        </div>

        <div className="relative z-10 max-w-xl">
          <div className="inline-flex items-center rounded-full border border-white bg-white/80 px-4 py-2 text-sm font-semibold text-blue-700 shadow-sm backdrop-blur">
            <Sparkles className="mr-2 h-4 w-4" />
            Platform Keamanan Makanan MBG
          </div>

          <h1 className="mt-7 text-6xl font-black leading-[0.95] tracking-[-0.05em] text-slate-950">
            Satu akses untuk monitoring makanan yang lebih aman.
          </h1>

          <p className="mt-6 text-base leading-8 text-slate-600">
            Masuk atau daftar untuk mengakses dashboard NutriSafe sesuai role:
            sekolah, mitra SPPG, umum/siswa, dan admin.
          </p>

          <div className="mt-8 grid max-w-lg grid-cols-2 gap-3">
            <div className="rounded-3xl border border-white bg-white/80 p-4 shadow-sm backdrop-blur">
              <CheckCircle2 className="h-5 w-5 text-blue-600" />
              <p className="mt-3 text-sm font-bold text-slate-900">
                Role-based dashboard
              </p>
              <p className="mt-1 text-xs leading-5 text-slate-500">
                Tampilan menyesuaikan kebutuhan pengguna.
              </p>
            </div>

            <div className="rounded-3xl border border-white bg-white/80 p-4 shadow-sm backdrop-blur">
              <ShieldCheck className="h-5 w-5 text-orange-600" />
              <p className="mt-3 text-sm font-bold text-slate-900">
                Food safety first
              </p>
              <p className="mt-1 text-xs leading-5 text-slate-500">
                Fokus pada alergi, distribusi, dan keamanan makanan.
              </p>
            </div>
          </div>
        </div>

        <div className="relative z-10">
          <Link
            href="/"
            className="inline-flex items-center text-sm font-semibold text-slate-500 transition hover:text-blue-600"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Kembali ke landing page
          </Link>
        </div>
      </section>

      <section className="flex max-h-screen items-start justify-center overflow-y-auto px-12 py-10">
        <div className="w-full max-w-xl rounded-[2rem] border border-white bg-white/85 p-8 shadow-2xl shadow-blue-100 backdrop-blur">
          <div className="mb-7">
            <p className="text-sm font-bold text-blue-600">MBG NutriSafe</p>

            <h2 className="mt-3 text-4xl font-black tracking-tight text-slate-950">
              {title}
            </h2>

            <p className="mt-3 text-sm leading-6 text-slate-600">
              {description}
            </p>
          </div>

          {children}
        </div>
      </section>
    </main>
  );
}
