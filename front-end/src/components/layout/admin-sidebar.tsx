import Link from "next/link";
import {
  LayoutDashboard,
  ClipboardList,
  School,
  Building2,
  User,
  Settings,
  ShieldCheck,
} from "lucide-react";

const menuItems = [
  {
    label: "Dashboard",
    href: "/admin",
    icon: LayoutDashboard,
    active: true,
  },
  {
    label: "Pelaporan",
    href: "/admin/laporan",
    icon: ClipboardList,
  },
  {
    label: "Data Sekolah",
    href: "/admin/sekolah",
    icon: School,
  },
  {
    label: "Data Mitra",
    href: "/admin/mitra",
    icon: Building2,
  },
  {
    label: "Akun",
    href: "/admin/akun",
    icon: User,
  },
  {
    label: "Pengaturan",
    href: "/admin/pengaturan",
    icon: Settings,
  },
  {
    label: "Keamanan",
    href: "/admin/keamanan",
    icon: ShieldCheck,
  },
];

export function AdminSidebar() {
  return (
    <aside className="flex h-screen w-72 flex-col border-r border-white/10 bg-slate-950 px-5 py-5 text-white">
      {" "}
      <div className="mb-8 flex items-center gap-3 rounded-3xl bg-white/10 p-4 ring-1 ring-white/15">
        <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white text-blue-700">
          <ShieldCheck className="h-6 w-6" />
        </div>

        <div>
          <p className="text-sm font-bold leading-none">NutriSafe</p>
          <p className="mt-1 text-xs text-blue-100">Admin MBG Panel</p>
        </div>
      </div>
      <nav className="space-y-2">
        {menuItems.map((item) => {
          const Icon = item.icon;

          return (
            <Link
              key={item.label}
              href={item.href}
              className={`group flex items-center gap-3 rounded-2xl px-3.5 py-3 text-sm font-semibold transition ${
                item.active
                  ? "bg-blue-600 text-white shadow-lg shadow-blue-900/30"
                  : "text-slate-300 hover:bg-white/10 hover:text-white"
              }`}
            >
              <Icon className="h-4 w-4 shrink-0" />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>
      <div className="mt-auto rounded-[1.75rem] bg-white/10 p-4 ring-1 ring-white/10">
        <p className="text-xs text-slate-400">Login sebagai</p>
        <p className="mt-1 text-sm font-bold">Admin MBG</p>
        <p className="mt-1 text-xs text-slate-400">admin@nutrisafe.id</p>
      </div>
    </aside>
  );
}
