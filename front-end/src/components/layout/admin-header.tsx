import { Bell, Search, Wifi } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function AdminHeader() {
  return (
    <header className="mb-8 flex items-center justify-between gap-6">
      <div>
        <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">
          Admin MBG
        </Badge>

        <h1 className="mt-3 text-3xl font-bold tracking-tight text-slate-950">
          Dashboard Admin
        </h1>

        <p className="mt-2 text-sm text-slate-600">
          Pantau data sekolah, mitra SPPG, laporan makanan, dan status keamanan
          distribusi MBG.
        </p>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative w-72">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
          <Input
            placeholder="Cari laporan..."
            className="h-11 rounded-2xl border-blue-100 bg-white pl-9"
          />
        </div>

        <div className="flex h-11 items-center gap-2 rounded-2xl bg-white px-4 shadow-sm">
          <Wifi className="h-4 w-4 text-blue-600" />
          <span className="text-sm font-medium text-slate-700">Online</span>
        </div>

        <Button size="icon" variant="outline" className="h-11 w-11 rounded-2xl">
          <Bell className="h-4 w-4" />
        </Button>
      </div>
    </header>
  );
}
