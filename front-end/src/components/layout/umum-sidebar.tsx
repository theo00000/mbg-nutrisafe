import { AlertCircle, Home, Settings, User } from "lucide-react";

const menuItems = [
  {
    label: "Dashboard Umum",
    icon: Home,
    active: true,
  },
  {
    label: "Profil",
    icon: User,
  },
  {
    label: "Data Alergi",
    icon: AlertCircle,
  },
  {
    label: "Pengaturan",
    icon: Settings,
  },
];

export function UmumSidebar() {
  return (
    <aside className="flex h-screen w-72 flex-col border-r border-white/10 bg-orange-950 px-5 py-5 text-white">
      <div className="mb-8 rounded-[1.75rem] bg-white/10 p-4 ring-1 ring-white/10">
        <p className="text-sm font-black tracking-tight">NutriSafe MBG</p>
        <p className="mt-1 text-xs text-orange-100/70">Student Safety Panel</p>
      </div>
      <nav className="space-y-2">
        {menuItems.map((item) => {
          const Icon = item.icon;

          return (
            <div
              key={item.label}
              className={`flex items-center gap-3 rounded-2xl px-3.5 py-3 text-sm font-semibold transition ${
                item.active
                  ? "bg-orange-500 text-white shadow-lg shadow-orange-900/30"
                  : "text-orange-50/80 hover:bg-white/10 hover:text-white"
              }`}
            >
              <Icon className="h-4 w-4 shrink-0" />
              <span>{item.label}</span>
            </div>
          );
        })}
      </nav>
      <div className="mt-auto rounded-[1.75rem] bg-white/10 p-4 ring-1 ring-white/10">
        <p className="text-xs text-orange-100/70">Role aktif</p>
        <p className="mt-1 text-sm font-bold">Umum/Siswa</p>
      </div>
    </aside>
  );
}
