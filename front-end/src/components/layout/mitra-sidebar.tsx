import {
  Building2,
  ClipboardList,
  Home,
  PackageCheck,
  Settings,
} from "lucide-react";

const menuItems = [
  {
    label: "Dashboard",
    icon: Home,
    active: true,
  },
  {
    label: "Produksi",
    icon: PackageCheck,
  },
  {
    label: "Distribusi",
    icon: Building2,
  },
  {
    label: "Laporan",
    icon: ClipboardList,
  },
  {
    label: "Pengaturan",
    icon: Settings,
  },
];

export function MitraSidebar() {
  return (
    <aside className="flex h-screen w-72 flex-col bg-sky-700 px-5 py-5 text-white">
      <div className="mb-8 rounded-3xl bg-white/15 p-4 ring-1 ring-white/20">
        <p className="text-sm font-bold">NutriSafe MBG</p>
        <p className="mt-1 text-xs text-sky-100">Dashboard Mitra/SPPG</p>
      </div>

      <nav className="space-y-2">
        {menuItems.map((item) => {
          const Icon = item.icon;

          return (
            <div
              key={item.label}
              className={`flex items-center gap-3 rounded-2xl px-3.5 py-3 text-sm font-medium ${
                item.active
                  ? "bg-white text-sky-700 shadow-sm"
                  : "text-sky-50 hover:bg-white/10"
              }`}
            >
              <Icon className="h-4 w-4 shrink-0" />
              <span>{item.label}</span>
            </div>
          );
        })}
      </nav>

      <div className="mt-auto rounded-3xl bg-white/15 p-4 ring-1 ring-white/20">
        <p className="text-xs text-sky-100">Role aktif</p>
        <p className="mt-1 text-sm font-semibold">Mitra/SPPG</p>
      </div>
    </aside>
  );
}
