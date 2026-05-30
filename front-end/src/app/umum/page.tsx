import { DashboardShell } from "@/components/dashboard/dashboard-shell";
import { DashboardHeader } from "@/components/dashboard/dashboard-header";
import { LogoutButton } from "@/components/auth/logout-button";
import { UmumSidebar } from "@/components/layout/umum-sidebar";
import { ProfileCard } from "@/components/profile/profile-card";
import { ProfileEditForm } from "@/components/profile/profile-edit-form";
import { AllergyCard } from "@/components/allergies/allergy-card";

export default function UmumPage() {
  return (
    <DashboardShell allowedRoles={["umum"]} sidebar={<UmumSidebar />}>
      <DashboardHeader
        eyebrow="Dashboard Umum"
        title="Kelola Informasi Akun"
        description="Pantau profil, data alergi, dan informasi keamanan makanan pribadi."
        accentClassName="text-orange-600"
        action={<LogoutButton />}
      />

      <div className="grid grid-cols-[1fr_0.85fr] gap-6">
        <ProfileCard />
        <ProfileEditForm />
      </div>

      <div className="mt-6">
        <AllergyCard />
      </div>
    </DashboardShell>
  );
}
