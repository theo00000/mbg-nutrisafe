import { LogoutButton } from "@/components/auth/logout-button";
import { RequireAuth } from "@/components/auth/require-auth";
import { roleTheme } from "@/lib/role-theme";
import { ProfileCard } from "@/components/profile/profile-card";
import { ProfileEditForm } from "@/components/profile/profile-edit-form";
import { AllergyCard } from "@/components/allergies/allergy-card";
import { UmumSidebar } from "@/components/layout/umum-sidebar";

export default function UmumPage() {
  const theme = roleTheme.umum;

  return (
    <RequireAuth allowedRoles={["umum"]}>
      <main className="flex min-h-screen bg-orange-50">
        <UmumSidebar />

        <section className="flex-1 p-8">
          <div className="mb-6 flex items-center justify-between">
            <div>
              <p className={`text-sm font-semibold ${theme.text}`}>
                {theme.name}
              </p>

              <h1 className="mt-2 text-3xl font-bold text-slate-950">
                Dashboard Umum/Siswa
              </h1>
            </div>

            <LogoutButton />
          </div>

          <div className="grid grid-cols-[1fr_0.8fr] gap-6">
            <ProfileCard />
            <ProfileEditForm />
          </div>

          <div className="mt-6">
            <AllergyCard />
          </div>
        </section>
      </main>
    </RequireAuth>
  );
}
