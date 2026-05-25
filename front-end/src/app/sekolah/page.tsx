import { LogoutButton } from "@/components/auth/logout-button";
import { RequireAuth } from "@/components/auth/require-auth";
import { SekolahSidebar } from "@/components/layout/sekolah-sidebar";
import { Card, CardContent } from "@/components/ui/card";

const schoolStats = [
  {
    label: "Total Siswa",
    value: "428",
    note: "Dummy sementara",
  },
  {
    label: "Menu Hari Ini",
    value: "3",
    note: "Paket makanan aktif",
  },
  {
    label: "Data Alergi",
    value: "18",
    note: "Perlu perhatian",
  },
  {
    label: "Laporan",
    value: "5",
    note: "Bulan ini",
  },
];

export default function SekolahPage() {
  return (
    <RequireAuth allowedRoles={["school"]}>
      <main className="flex min-h-screen bg-emerald-50">
        <SekolahSidebar />

        <section className="flex-1 p-8">
          <div className="mb-8 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-emerald-700">
                Dashboard Sekolah
              </p>

              <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-950">
                Monitoring Keamanan Makanan Sekolah
              </h1>

              <p className="mt-2 text-sm text-slate-600">
                Pantau ringkasan siswa, menu harian, alergi, dan laporan terkait
                distribusi MBG.
              </p>
            </div>

            <LogoutButton />
          </div>

          <div className="grid grid-cols-4 gap-5">
            {schoolStats.map((item) => (
              <Card key={item.label} className="border-emerald-100 shadow-sm">
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
            <Card className="border-emerald-100 shadow-sm">
              <CardContent className="p-6">
                <h2 className="text-lg font-semibold text-slate-950">
                  Status Distribusi Hari Ini
                </h2>

                <p className="mt-2 text-sm text-slate-500">
                  TODO: sementara dummy karena endpoint distribusi sekolah belum
                  tersedia.
                </p>

                <div className="mt-5 rounded-3xl bg-emerald-600 p-5 text-white">
                  <p className="text-sm text-emerald-50">Status</p>
                  <p className="mt-2 text-3xl font-bold">Aman</p>
                  <p className="mt-2 text-xs leading-5 text-emerald-50">
                    Paket makanan hari ini tercatat diterima dan belum ada
                    laporan prioritas tinggi.
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card className="border-emerald-100 shadow-sm">
              <CardContent className="p-6">
                <h2 className="text-lg font-semibold text-slate-950">
                  Catatan Sekolah
                </h2>

                <p className="mt-2 text-sm text-slate-500">
                  Area ini nanti bisa diisi ringkasan alergi siswa, laporan
                  menu, atau status penerimaan makanan.
                </p>

                <div className="mt-5 space-y-3">
                  <div className="rounded-2xl border border-emerald-100 bg-emerald-50 p-4">
                    <p className="text-sm font-medium text-slate-900">
                      18 siswa memiliki catatan alergi
                    </p>
                    <p className="mt-1 text-xs text-slate-500">
                      Perlu validasi menu sebelum distribusi.
                    </p>
                  </div>

                  <div className="rounded-2xl border border-emerald-100 bg-white p-4">
                    <p className="text-sm font-medium text-slate-900">
                      Tidak ada laporan darurat
                    </p>
                    <p className="mt-1 text-xs text-slate-500">
                      Monitoring hari ini stabil.
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
