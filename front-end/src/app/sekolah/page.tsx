import { roleTheme } from "@/lib/role-theme";

export default function SekolahPage() {
  const theme = roleTheme.sekolah;

  return (
    <main className={`min-h-screen ${theme.soft} p-10`}>
      <div className="rounded-3xl bg-white p-8 shadow-sm">
        <p className={`text-sm font-semibold ${theme.text}`}>{theme.name}</p>

        <h1 className="mt-2 text-3xl font-bold text-slate-950">
          Dashboard Sekolah MBG
        </h1>

        <p className="mt-3 text-sm text-slate-600">
          Kerangka awal halaman sekolah. Nanti kita isi dengan sidebar,
          statistik, tabel, dan data dari API dashboard.
        </p>
      </div>
    </main>
  );
}
