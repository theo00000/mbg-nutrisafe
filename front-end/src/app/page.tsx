import Link from "next/link";
import {
  AlertTriangle,
  CheckCircle2,
  School,
  ShieldCheck,
  Users,
  Utensils,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

const features = [
  {
    title: "Pantau Keamanan Makanan",
    description: "Monitoring laporan, alergi, dan status distribusi MBG.",
    icon: ShieldCheck,
  },
  {
    title: "Multi Role Dashboard",
    description: "Admin, sekolah, mitra SPPG, dan siswa punya tampilan khusus.",
    icon: Users,
  },
  {
    title: "Data Alergi Siswa",
    description:
      "Membantu mencegah risiko makanan yang tidak sesuai kondisi siswa.",
    icon: AlertTriangle,
  },
];

export default function LandingPage() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-sky-50 via-white to-orange-50">
      <nav className="mx-auto flex max-w-7xl items-center justify-between px-10 py-6">
        <div className="flex items-center gap-3">
          <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-600 text-white shadow-sm">
            <ShieldCheck className="h-6 w-6" />
          </div>

          <div>
            <p className="text-lg font-bold text-slate-950">NutriSafe</p>
            <p className="text-xs text-slate-500">MBG Food Safety Platform</p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <Button variant="ghost" asChild>
            <Link href="/login">Masuk</Link>
          </Button>

          <Button asChild>
            <Link href="/register">Daftar</Link>
          </Button>
        </div>
      </nav>

      <section className="mx-auto grid max-w-7xl grid-cols-[1.1fr_0.9fr] items-center gap-12 px-10 py-16">
        <div>
          <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">
            Platform Keamanan Makanan MBG
          </Badge>

          <h1 className="mt-6 max-w-3xl text-6xl font-bold tracking-tight text-slate-950">
            Monitoring keamanan makanan MBG yang lebih aman dan terintegrasi.
          </h1>

          <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600">
            NutriSafe membantu sekolah, mitra SPPG, admin, dan siswa dalam
            memantau distribusi makanan, mencatat alergi, serta mendukung
            pelaporan risiko keamanan pangan.
          </p>

          <div className="mt-8 flex gap-3">
            <Button size="lg" asChild>
              <Link href="/login">Masuk Dashboard</Link>
            </Button>

            <Button size="lg" variant="outline" asChild>
              <Link href="/register">Buat Akun</Link>
            </Button>
          </div>

          <div className="mt-10 grid max-w-xl grid-cols-3 gap-4">
            <div>
              <p className="text-3xl font-bold text-slate-950">4</p>
              <p className="mt-1 text-sm text-slate-500">Role pengguna</p>
            </div>

            <div>
              <p className="text-3xl font-bold text-slate-950">API</p>
              <p className="mt-1 text-sm text-slate-500">Backend terhubung</p>
            </div>

            <div>
              <p className="text-3xl font-bold text-slate-950">Safe</p>
              <p className="mt-1 text-sm text-slate-500">Food monitoring</p>
            </div>
          </div>
        </div>

        <Card className="border-blue-100 bg-white/80 shadow-xl backdrop-blur">
          <CardContent className="p-8">
            <div className="rounded-3xl bg-blue-600 p-6 text-white">
              <p className="text-sm text-blue-100">Status Sistem</p>
              <h2 className="mt-3 text-4xl font-bold">Aktif</h2>
              <p className="mt-3 text-sm leading-6 text-blue-100">
                Dashboard siap digunakan untuk monitoring role sekolah, mitra,
                umum/siswa, dan admin.
              </p>
            </div>

            <div className="mt-5 space-y-3">
              <div className="flex items-center gap-3 rounded-2xl border border-blue-100 bg-blue-50 p-4">
                <CheckCircle2 className="h-5 w-5 text-blue-600" />
                <div>
                  <p className="text-sm font-semibold text-slate-900">
                    Backend API Terhubung
                  </p>
                  <p className="text-xs text-slate-500">
                    GET /api/ping sudah berhasil digunakan.
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3 rounded-2xl border border-orange-100 bg-orange-50 p-4">
                <Utensils className="h-5 w-5 text-orange-600" />
                <div>
                  <p className="text-sm font-semibold text-slate-900">
                    Data Alergi
                  </p>
                  <p className="text-xs text-slate-500">
                    UI tersedia, integrasi menunggu route backend aktif.
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3 rounded-2xl border border-emerald-100 bg-emerald-50 p-4">
                <School className="h-5 w-5 text-emerald-600" />
                <div>
                  <p className="text-sm font-semibold text-slate-900">
                    Dashboard Sekolah
                  </p>
                  <p className="text-xs text-slate-500">
                    Monitoring dummy siap untuk presentasi UI.
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="mx-auto grid max-w-7xl grid-cols-3 gap-5 px-10 pb-16">
        {features.map((feature) => {
          const Icon = feature.icon;

          return (
            <Card key={feature.title} className="border-slate-100 shadow-sm">
              <CardContent className="p-6">
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-sky-50 text-blue-600">
                  <Icon className="h-6 w-6" />
                </div>

                <h3 className="mt-5 text-lg font-semibold text-slate-950">
                  {feature.title}
                </h3>

                <p className="mt-2 text-sm leading-6 text-slate-600">
                  {feature.description}
                </p>
              </CardContent>
            </Card>
          );
        })}
      </section>
    </main>
  );
}
