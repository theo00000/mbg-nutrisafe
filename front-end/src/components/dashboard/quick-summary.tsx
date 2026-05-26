import { Activity, CheckCircle2, Database } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

const summaryItems = [
  {
    label: "Keamanan hari ini",
    value: "Stabil",
    icon: CheckCircle2,
  },
  {
    label: "Sinkronisasi data",
    value: "Aktif",
    icon: Database,
  },
  {
    label: "Monitoring sistem",
    value: "Berjalan",
    icon: Activity,
  },
];

export function QuickSummary() {
  return (
    <Card className="border-blue-100 shadow-sm">
      <CardContent className="p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold text-slate-950">
              Ringkasan Cepat
            </h2>
            <p className="mt-2 text-sm text-slate-500">
              TODO: sebagian data nanti disambungkan ke GET
              /api/dashboard/stats.
            </p>
          </div>

          <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">
            Overview
          </Badge>
        </div>

        <div className="mt-5 rounded-3xl bg-blue-600 p-5 text-white">
          <p className="text-sm text-blue-100">Status operasional</p>
          <p className="mt-2 text-3xl font-bold">Aman</p>
          <p className="mt-2 text-xs leading-5 text-blue-100">
            Tidak ada anomali besar pada distribusi dan laporan prioritas.
          </p>
        </div>

        <div className="mt-5 space-y-3">
          {summaryItems.map((item) => {
            const Icon = item.icon;

            return (
              <div
                key={item.label}
                className="flex items-center justify-between rounded-2xl border border-blue-100 bg-blue-50 px-4 py-3"
              >
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-white text-blue-700">
                    <Icon className="h-4 w-4" />
                  </div>

                  <p className="text-sm font-medium text-slate-700">
                    {item.label}
                  </p>
                </div>

                <p className="text-sm font-semibold text-blue-700">
                  {item.value}
                </p>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
