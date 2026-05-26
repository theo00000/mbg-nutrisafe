import { LogoutButton } from "@/components/auth/logout-button";
import { RequireAuth } from "@/components/auth/require-auth";
import { MitraSidebar } from "@/components/layout/mitra-sidebar";
import { Card, CardContent } from "@/components/ui/card";

const mitraStats = [
  {
    label: "Paket Diproduksi",
    value: "1.280",
    note: "Dummy hari ini",
  },
  {
    label: "Sekolah Tujuan",
    value: "12",
    note: "Distribusi aktif",
  },
  {
    label: "Menu Aktif",
    value: "4",
    note: "Paket makan siang",
  },
  {
    label: "Laporan Kualitas",
    value: "2",
    note: "Perlu ditinjau",
  },
];

export default function MitraPage() {
  return (
    <RequireAuth allowedRoles={["spgg"]}>
      <main className="flex min-h-screen bg-sky-50">
        <MitraSidebar />

        <section className="flex-1 p-8">
          <div className="mb-8 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-sky-700">
                Dashboard Mitra/SPPG
              </p>

              <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-950">
                Monitoring Produksi dan Distribusi
              </h1>

              <p className="mt-2 text-sm text-slate-600">
                Pantau status produksi makanan, distribusi ke sekolah, dan
                laporan kualitas harian.
              </p>
            </div>

            <LogoutButton />
          </div>

          <div className="grid grid-cols-4 gap-5">
            {mitraStats.map((item) => (
              <Card key={item.label} className="border-sky-100 shadow-sm">
                <CardContent className="p-5">
                  <p className="text-sm text-slate-500">{item.label}</p>
                  <p className="mt-3 text-3xl font-bold text-slate-950">
                    {item.value}
                  </p>
                  <p className="mt-1 text-xs text-slate-500">{item.note}</p>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="mt-6 grid grid-cols-[1.2fr_0.8fr] gap-6">
            <Card className="border-sky-100 shadow-sm">
              <CardContent className="p-6">
                <h2 className="text-lg font-semibold text-slate-950">
                  Status Produksi Hari Ini
                </h2>

                <p className="mt-2 text-sm text-slate-500">
                  TODO: sementara dummy karena endpoint produksi/menu belum
                  tersedia.
                </p>

                <div className="mt-5 rounded-3xl bg-sky-700 p-5 text-white">
                  <p className="text-sm text-sky-100">Progress produksi</p>
                  <p className="mt-2 text-3xl font-bold">82%</p>
                  <p className="mt-2 text-xs leading-5 text-sky-100">
                    Produksi makanan berjalan normal dan siap masuk tahap
                    distribusi.
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card className="border-sky-100 shadow-sm">
              <CardContent className="p-6">
                <h2 className="text-lg font-semibold text-slate-950">
                  Catatan Distribusi
                </h2>

                <p className="mt-2 text-sm text-slate-500">
                  Area ini nanti bisa diisi status pengiriman, validasi
                  kualitas, dan laporan sekolah tujuan.
                </p>

                <div className="mt-5 space-y-3">
                  <div className="rounded-2xl border border-sky-100 bg-sky-50 p-4">
                    <p className="text-sm font-medium text-slate-900">
                      8 sekolah sudah menerima paket
                    </p>
                    <p className="mt-1 text-xs text-slate-500">
                      4 sekolah masih dalam proses distribusi.
                    </p>
                  </div>

                  <div className="rounded-2xl border border-sky-100 bg-white p-4">
                    <p className="text-sm font-medium text-slate-900">
                      Tidak ada laporan kualitas tinggi
                    </p>
                    <p className="mt-1 text-xs text-slate-500">
                      Monitoring kualitas hari ini stabil.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </section>
      </main>
    </RequireAuth>
  );
}
