import { ApiStatusCard } from "@/components/dashboard/api-status-card";
import { AdminStatsSection } from "@/components/dashboard/admin-stats-section";
import { QuickSummary } from "@/components/dashboard/quick-summary";
import { RecentReports } from "@/components/dashboard/recent-reports";
import { AdminHeader } from "@/components/layout/admin-header";
import { AdminSidebar } from "@/components/layout/admin-sidebar";

export default function AdminPage() {
  return (
    <main className="flex min-h-screen bg-[#f8fbff]">
      <AdminSidebar />

      <section className="flex-1 overflow-y-auto p-8">
        <AdminHeader />

        <div className="mb-5 flex justify-end">
          <ApiStatusCard />
        </div>

        <AdminStatsSection />

        <div className="mt-6 grid grid-cols-[1.4fr_0.8fr] gap-6">
          <RecentReports />
          <QuickSummary />
        </div>
      </section>
    </main>
  );
}
