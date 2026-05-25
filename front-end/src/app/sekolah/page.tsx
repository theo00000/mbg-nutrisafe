import { LogoutButton } from "@/components/auth/logout-button";
import { roleTheme } from "@/lib/role-theme";
import { RequireAuth } from "@/components/auth/require-auth";

export default function SekolahPage() {
  const theme = roleTheme.sekolah;

  return (
    <RequireAuth allowedRoles={["school"]}>
      <main className={`min-h-screen ${theme.soft} p-10`}>
        <div className="mb-6 flex items-center justify-between">
          <div>
            <p className={`text-sm font-semibold ${theme.text}`}>
              {theme.name}
            </p>

            <h1 className="mt-2 text-3xl font-bold text-slate-950">
              Dashboard Sekolah
            </h1>
          </div>

          <LogoutButton />
        </div>

        <div className="rounded-3xl bg-white p-8 shadow-sm">
          <p className="text-sm text-slate-600">
            Halaman sekolah sudah siap untuk integrasi profil dan data alergi.
          </p>
        </div>
      </main>
    </RequireAuth>
  );
}
