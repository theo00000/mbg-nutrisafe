import { AdminSidebar } from "@/components/layout/admin-sidebar";
import { AdminHeader } from "@/components/layout/admin-header";
import { DashboardStatCard } from "@/components/dashboard/stat-card";
import { Badge } from "@/components/ui/badge";
import { RecentReports } from "@/components/dashboard/recent-reports";
import { QuickSummary } from "@/components/dashboard/quick-summary";
import { Building2, ClipboardList, School, ShieldAlert } from "lucide-react";

const stats = [
  {
    label: "Total Sekolah",
    value: "24",
    note: "Terdaftar aktif",
    trend: "+8%",
    icon: School,
  },
  {
    label: "Total Mitra",
    value: "12",
    note: "SPPG terhubung",
    trend: "+3%",
    icon: Building2,
  },
  {
    label: "Laporan Masuk",
    value: "38",
    note: "Bulan ini",
    trend: "+12%",
    icon: ClipboardList,
  },
  {
    label: "Kasus Prioritas",
    value: "4",
    note: "Perlu ditinjau",
    trend: "Urgent",
    icon: ShieldAlert,
  },
];

export default function AdminPage() {
  return (
    <main className="flex min-h-screen bg-blue-50">
      <AdminSidebar />

      <section className="flex-1 p-8">
        <AdminHeader />
        <div className="mb-8 flex items-start justify-between">
          <div>
            <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">
              Admin MBG
            </Badge>

            <h1 className="mt-3 text-3xl font-bold tracking-tight text-slate-950">
              Dashboard Admin
            </h1>

            <p className="mt-2 text-sm text-slate-600">
              Pantau data sekolah, mitra SPPG, laporan makanan, dan status
              keamanan distribusi MBG.
            </p>
          </div>

          <div className="rounded-2xl bg-white px-4 py-3 text-right shadow-sm">
            <p className="text-xs text-slate-500">Status Sistem</p>
            <p className="text-sm font-semibold text-blue-700">Online</p>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-5">
          {stats.map((item) => (
            <DashboardStatCard
              key={item.label}
              label={item.label}
              value={item.value}
              note={item.note}
              trend={item.trend}
              icon={item.icon}
            />
          ))}
        </div>

        <div className="mt-6 grid grid-cols-[1.4fr_0.8fr] gap-6">
          <RecentReports />

          <QuickSummary />
        </div>
      </section>
    </main>
  );
}
