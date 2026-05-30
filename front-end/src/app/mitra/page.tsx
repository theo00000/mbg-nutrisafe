import { LogoutButton } from "@/components/auth/logout-button";
import { MitraSidebar } from "@/components/layout/mitra-sidebar";
import { Card, CardContent } from "@/components/ui/card";
import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { Building2, ClipboardList, PackageCheck, Utensils } from "lucide-react";
import { RoleStatCard } from "@/components/dashboard/role-stat-card";
import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { InfoPanel } from "@/components/dashboard/info-panel";

const mitraStats = [
  {
    label: "Paket Diproduksi",
    value: "1.280",
    note: "Dummy hari ini",
    icon: PackageCheck,
  },
  {
    label: "Sekolah Tujuan",
    value: "12",
    note: "Distribusi aktif",
    icon: Building2,
  },
  {
    label: "Menu Aktif",
    value: "4",
    note: "Paket makan siang",
    icon: Utensils,
  },
  {
    label: "Laporan Kualitas",
    value: "2",
    note: "Perlu ditinjau",
    icon: ClipboardList,
  },
];

export default function MitraPage() {
  return (
    <DashboardShell allowedRoles={["spgg"]} sidebar={<MitraSidebar />}>
      <DashboardHeader
        eyebrow="Dashboard Mitra/SPPG"
        title="Monitoring Produksi dan Distribusi"
        description="Pantau status produksi makanan, distribusi ke sekolah, dan laporan kualitas harian."
        accentClassName="text-sky-700"
        action={<LogoutButton />}
      />

      <div className="grid grid-cols-4 gap-5">
        {mitraStats.map((item) => (
          <RoleStatCard
            key={item.label}
            label={item.label}
            value={item.value}
            note={item.note}
            icon={item.icon}
            accentClassName="text-sky-700"
            iconClassName="bg-sky-50 text-sky-700"
          />
        ))}
      </div>

      <div className="mt-6 grid grid-cols-[1.2fr_0.8fr] gap-6">
        <InfoPanel
          title="Status Produksi Hari Ini"
          description="TODO: sementara dummy karena endpoint produksi/menu belum tersedia."
        >
          <div className="rounded-[1.75rem] bg-sky-600 p-5 text-white shadow-lg shadow-sky-100">
            <p className="text-sm text-sky-100">Progress produksi</p>
            <p className="mt-2 text-5xl font-black tracking-tight">82%</p>
            <p className="mt-2 text-xs leading-5 text-sky-100">
              Produksi makanan berjalan normal dan siap masuk tahap distribusi.
            </p>
          </div>
        </InfoPanel>

        <InfoPanel
          title="Catatan Distribusi"
          description="Area ini nanti bisa diisi status pengiriman, validasi kualitas, dan laporan sekolah tujuan."
        >
          <div className="space-y-3">
            <div className="rounded-[1.75rem] border border-sky-100 bg-sky-50 p-4">
              <p className="text-sm font-medium text-slate-900">
                8 sekolah sudah menerima paket
              </p>
              <p className="mt-1 text-xs text-slate-500">
                4 sekolah masih dalam proses distribusi.
              </p>
            </div>

            <div className="rounded-[1.75rem] border border-sky-100 bg-white p-4">
              <p className="text-sm font-medium text-slate-900">
                Tidak ada laporan kualitas tinggi
              </p>
              <p className="mt-1 text-xs text-slate-500">
                Monitoring kualitas hari ini stabil.
              </p>
            </div>
          </div>
        </InfoPanel>
      </div>
    </DashboardShell>
  );
}
