import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

export default function LandingPage() {
  return (
    <main className="min-h-screen bg-sky-50 px-10 py-8">
      <section className="mx-auto flex max-w-6xl items-center justify-between gap-10 rounded-3xl bg-white p-10 shadow-sm">
        <div className="max-w-xl">
          <p className="text-sm font-semibold text-blue-600">MBG NutriSafe</p>

          <h1 className="mt-3 text-5xl font-bold tracking-tight text-slate-950">
            Platform Monitoring Keamanan Makanan MBG
          </h1>

          <p className="mt-5 text-base leading-7 text-slate-600">
            Membantu sekolah, mitra, admin, dan siswa dalam pelaporan,
            pemantauan, serta pengelolaan data keamanan makanan.
          </p>

          <div className="mt-8 flex gap-3">
            <Button>Masuk Akun</Button>
            <Button variant="outline">Lihat Dashboard</Button>
          </div>
        </div>

        <Card className="w-[360px] border-blue-100 bg-blue-50">
          <CardContent className="p-6">
            <p className="text-sm font-semibold text-blue-700">Status Sistem</p>
            <h2 className="mt-3 text-3xl font-bold text-slate-900">
              Siap Dipakai
            </h2>
            <p className="mt-2 text-sm text-slate-600">
              Frontend sudah siap untuk mulai dikembangkan per role.
            </p>
          </CardContent>
        </Card>
      </section>
    </main>
  );
}
