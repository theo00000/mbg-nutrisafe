import { AlertTriangle, Clock } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";

const reports = [
  {
    title: "Dugaan makanan basi",
    school: "SD Negeri 01",
    category: "Kualitas makanan",
    time: "10 menit lalu",
    priority: "Sedang",
  },
  {
    title: "Keterlambatan distribusi",
    school: "SMP Harapan Bangsa",
    category: "Distribusi",
    time: "35 menit lalu",
    priority: "Rendah",
  },
  {
    title: "Keluhan alergi siswa",
    school: "SMA Cendekia",
    category: "Alergi",
    time: "1 jam lalu",
    priority: "Tinggi",
  },
];

function getPriorityClass(priority: string) {
  if (priority === "Tinggi") {
    return "bg-red-100 text-red-700 hover:bg-red-100";
  }

  if (priority === "Sedang") {
    return "bg-yellow-100 text-yellow-700 hover:bg-yellow-100";
  }

  return "bg-blue-100 text-blue-700 hover:bg-blue-100";
}

export function RecentReports() {
  return (
    <Card className="border-blue-100 shadow-sm">
      <CardContent className="p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold text-slate-950">
              Laporan Terbaru
            </h2>
            <p className="mt-2 text-sm text-slate-500">
              TODO: data laporan masih dummy karena endpoint laporan belum
              tersedia.
            </p>
          </div>

          <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100">
            Live Monitoring
          </Badge>
        </div>

        <div className="mt-5 space-y-3">
          {reports.map((report) => (
            <div
              key={report.title}
              className="rounded-2xl border border-blue-100 bg-white p-4 transition hover:bg-blue-50"
            >
              <div className="flex items-start justify-between gap-4">
                <div>
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-4 w-4 text-blue-600" />
                    <p className="text-sm font-semibold text-slate-900">
                      {report.title}
                    </p>
                  </div>

                  <p className="mt-1 text-xs text-slate-500">
                    {report.school} • {report.category}
                  </p>

                  <div className="mt-3 flex items-center gap-1 text-xs text-slate-400">
                    <Clock className="h-3.5 w-3.5" />
                    <span>{report.time}</span>
                  </div>
                </div>

                <Badge className={getPriorityClass(report.priority)}>
                  {report.priority}
                </Badge>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
