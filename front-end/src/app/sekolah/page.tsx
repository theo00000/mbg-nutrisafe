import { LogoutButton } from "@/components/auth/logout-button";
import { SekolahSidebar } from "@/components/layout/sekolah-sidebar";

import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { AlertCircle, ClipboardList, School, Utensils } from "lucide-react";
import { RoleStatCard } from "@/components/dashboard/role-stat-card";
import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { SchoolAllergyCard } from "@/components/allergies/school-allergy-card";
import { InfoPanel } from "@/components/dashboard/info-panel";

const schoolStats = [
  {
    label: "Total Siswa",
    value: "428",
    note: "Dummy sementara",
    icon: School,
  },
  {
    label: "Menu Hari Ini",
    value: "3",
    note: "Paket makanan aktif",
    icon: Utensils,
  },
  {
    label: "Data Alergi",
    value: "18",
    note: "Perlu perhatian",
    icon: AlertCircle,
  },
  {
    label: "Laporan",
    value: "5",
    note: "Bulan ini",
    icon: ClipboardList,
  },
];

export default function SekolahPage() {
  return (
    <DashboardShell allowedRoles={["school"]} sidebar={<SekolahSidebar />}>
      <DashboardHeader
        eyebrow="Dashboard Sekolah"
        title="Monitoring Keamanan Makanan Sekolah"
        description="Pantau ringkasan siswa, menu harian, alergi, dan laporan terkait distribusi MBG."
        accentClassName="text-emerald-700"
        action={<LogoutButton />}
      />

      <div className="grid grid-cols-4 gap-5">
        {schoolStats.map((item) => (
          <RoleStatCard
            key={item.label}
            label={item.label}
            value={item.value}
            note={item.note}
            icon={item.icon}
            accentClassName="text-emerald-700"
            iconClassName="bg-emerald-50 text-emerald-700"
          />
        ))}
      </div>

      <div className="mt-6 grid grid-cols-[1.2fr_0.8fr] gap-6">
        <InfoPanel
          title="Status Distribusi Hari Ini"
          description="TODO: sementara dummy karena endpoint distribusi sekolah belum tersedia."
        >
          <div className="rounded-[1.75rem] bg-emerald-600 p-5 text-white shadow-lg shadow-emerald-100">
            <p className="text-sm text-emerald-50">Status</p>
            <p className="mt-2 text-5xl font-black tracking-tight">Aman</p>
            <p className="mt-2 text-xs leading-5 text-emerald-50">
              Paket makanan hari ini tercatat diterima dan belum ada laporan
              prioritas tinggi.
            </p>
          </div>
        </InfoPanel>

        <InfoPanel
          title="Catatan Sekolah"
          description="Ringkasan alergi siswa, laporan menu, dan status penerimaan makanan."
        >
          <div className="space-y-3">
            <div className="rounded-[1.75rem] border border-emerald-100 bg-emerald-50 p-4">
              <p className="text-sm font-medium text-slate-900">
                18 siswa memiliki catatan alergi
              </p>
              <p className="mt-1 text-xs text-slate-500">
                Perlu validasi menu sebelum distribusi.
              </p>
            </div>

            <div className="rounded-[1.75rem] border border-emerald-100 bg-white p-4">
              <p className="text-sm font-medium text-slate-900">
                Tidak ada laporan darurat
              </p>
              <p className="mt-1 text-xs text-slate-500">
                Monitoring hari ini stabil.
              </p>
            </div>
          </div>
        </InfoPanel>
      </div>

      <div className="mt-6">
        <SchoolAllergyCard />
      </div>
    </DashboardShell>
  );
}
