import Link from "next/link";
import {
  ArrowRight,
  Building2,
  CheckCircle2,
  GraduationCap,
  ShieldCheck,
  Sparkles,
  Users,
  Utensils,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

const roleCards = [
  {
    title: "Admin MBG",
    description: "Pantau data sekolah, mitra, dan laporan keamanan pangan.",
    className: "border-blue-100 bg-blue-50 text-blue-700",
    icon: ShieldCheck,
  },
  {
    title: "Sekolah",
    description: "Kelola ringkasan siswa, alergi, dan status distribusi.",
    className: "border-emerald-100 bg-emerald-50 text-emerald-700",
    icon: GraduationCap,
  },
  {
    title: "Mitra/SPPG",
    description: "Monitoring produksi, menu, dan distribusi makanan.",
    className: "border-sky-100 bg-sky-50 text-sky-700",
    icon: Building2,
  },
  {
    title: "Umum/Siswa",
    description: "Cek profil, alergi, dan catatan keamanan makanan.",
    className: "border-orange-100 bg-orange-50 text-orange-700",
    icon: Users,
  },
];

export default function LandingPage() {
  return (
    <main className="min-h-screen overflow-hidden bg-[#f8fbff] text-slate-950">
      <section className="relative">
        <div className="absolute left-[-120px] top-[-120px] h-80 w-80 rounded-full bg-blue-200/50 blur-3xl" />
        <div className="absolute right-[-120px] top-40 h-96 w-96 rounded-full bg-orange-200/50 blur-3xl" />

        <nav className="relative z-10 mx-auto flex max-w-7xl items-center justify-between px-10 py-7">
          <Link href="/" className="flex items-center gap-3">
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

          <div className="flex items-center gap-8 rounded-full border border-white/70 bg-white/70 px-6 py-3 text-sm font-medium text-slate-600 shadow-sm backdrop-blur">
            <a href="#fitur" className="hover:text-blue-600">
              Fitur
            </a>
            <a href="#role" className="hover:text-blue-600">
              Role
            </a>
            <a href="#demo" className="hover:text-blue-600">
              Demo
            </a>
          </div>

          <div className="flex items-center gap-3">
            <Button variant="ghost" asChild>
              <Link href="/login">Masuk</Link>
            </Button>

            <Button className="rounded-full bg-blue-600 px-6" asChild>
              <Link href="/register">Daftar</Link>
            </Button>
          </div>
        </nav>

        <div className="relative z-10 mx-auto grid max-w-7xl grid-cols-[1.05fr_0.95fr] items-center gap-14 px-10 pb-20 pt-10">
          <div>
            <Badge className="rounded-full bg-white px-4 py-2 text-blue-700 shadow-sm hover:bg-white">
              <Sparkles className="mr-2 h-4 w-4" />
              Platform Monitoring Keamanan Makanan MBG
            </Badge>

            <h1 className="mt-7 max-w-4xl text-7xl font-black leading-[0.95] tracking-[-0.05em] text-slate-950">
              Food safety monitoring untuk MBG yang lebih modern.
            </h1>

            <p className="mt-7 max-w-2xl text-lg leading-8 text-slate-600">
              NutriSafe membantu sekolah, mitra SPPG, admin, dan siswa dalam
              memantau distribusi makanan, mencatat alergi, serta mengurangi
              risiko keamanan pangan.
            </p>

            <div className="mt-9 flex items-center gap-4">
              <Button
                size="lg"
                className="rounded-full bg-blue-600 px-7"
                asChild
              >
                <Link href="/login">
                  Masuk Dashboard
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>

              <Button
                size="lg"
                variant="outline"
                className="rounded-full px-7"
                asChild
              >
                <Link href="/register">Buat Akun</Link>
              </Button>
            </div>

            <div className="mt-12 grid max-w-2xl grid-cols-3 gap-4">
              <div className="rounded-3xl border border-white bg-white/80 p-5 shadow-sm">
                <p className="text-3xl font-black">4</p>
                <p className="mt-1 text-sm text-slate-500">Role utama</p>
              </div>

              <div className="rounded-3xl border border-white bg-white/80 p-5 shadow-sm">
                <p className="text-3xl font-black">API</p>
                <p className="mt-1 text-sm text-slate-500">Backend aktif</p>
              </div>

              <div className="rounded-3xl border border-white bg-white/80 p-5 shadow-sm">
                <p className="text-3xl font-black">Safe</p>
                <p className="mt-1 text-sm text-slate-500">Food tracking</p>
              </div>
            </div>
          </div>

          <div id="demo" className="relative">
            <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-orange-300/40 blur-2xl" />

            <Card className="relative overflow-hidden rounded-[2rem] border-white bg-white/85 shadow-2xl shadow-blue-100 backdrop-blur">
              <CardContent className="p-6">
                <div className="rounded-[1.5rem] bg-slate-950 p-5 text-white">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-xs text-slate-400">NutriSafe Live</p>
                      <h2 className="mt-1 text-2xl font-black">
                        Monitoring Center
                      </h2>
                    </div>

                    <Badge className="rounded-full bg-emerald-400/20 text-emerald-300 hover:bg-emerald-400/20">
                      Online
                    </Badge>
                  </div>

                  <div className="mt-6 grid grid-cols-2 gap-3">
                    <div className="rounded-2xl bg-white/10 p-4">
                      <p className="text-xs text-slate-400">Sekolah</p>
                      <p className="mt-2 text-3xl font-black">24</p>
                    </div>

                    <div className="rounded-2xl bg-white/10 p-4">
                      <p className="text-xs text-slate-400">Mitra</p>
                      <p className="mt-2 text-3xl font-black">12</p>
                    </div>
                  </div>
                </div>

                <div className="mt-4 space-y-3">
                  <div className="flex items-center gap-3 rounded-2xl border border-blue-100 bg-blue-50 p-4">
                    <CheckCircle2 className="h-5 w-5 text-blue-600" />
                    <div>
                      <p className="text-sm font-bold">Backend Terhubung</p>
                      <p className="text-xs text-slate-500">
                        GET /api/ping berhasil digunakan.
                      </p>
                    </div>
                  </div>

                  <div className="flex items-center gap-3 rounded-2xl border border-orange-100 bg-orange-50 p-4">
                    <Utensils className="h-5 w-5 text-orange-600" />
                    <div>
                      <p className="text-sm font-bold">Data Alergi</p>
                      <p className="text-xs text-slate-500">
                        UI siap, integrasi menunggu route backend.
                      </p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      <section id="role" className="mx-auto max-w-7xl px-10 pb-20">
        <div className="mb-8 flex items-end justify-between">
          <div>
            <p className="text-sm font-bold text-blue-600">Role Dashboard</p>
            <h2 className="mt-2 text-4xl font-black tracking-tight">
              Satu platform, banyak kebutuhan.
            </h2>
          </div>

          <p className="max-w-md text-sm leading-6 text-slate-500">
            Setiap role memakai identitas warna sesuai referensi Figma, tapi
            tetap berada dalam satu sistem desain NutriSafe.
          </p>
        </div>

        <div className="grid grid-cols-4 gap-5">
          {roleCards.map((role) => {
            const Icon = role.icon;

            return (
              <Card
                key={role.title}
                className={`rounded-[1.75rem] shadow-sm transition hover:-translate-y-1 hover:shadow-xl ${role.className}`}
              >
                <CardContent className="p-6">
                  <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/80">
                    <Icon className="h-6 w-6" />
                  </div>

                  <h3 className="mt-5 text-lg font-black text-slate-950">
                    {role.title}
                  </h3>

                  <p className="mt-2 text-sm leading-6 text-slate-600">
                    {role.description}
                  </p>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </section>
    </main>
  );
}
