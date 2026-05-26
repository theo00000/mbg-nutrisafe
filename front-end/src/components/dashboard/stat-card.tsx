import type { LucideIcon } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";

type DashboardStatCardProps = {
  label: string;
  value: string;
  note: string;
  trend: string;
  icon: LucideIcon;
};

export function DashboardStatCard({
  label,
  value,
  note,
  trend,
  icon: Icon,
}: DashboardStatCardProps) {
  return (
    <Card className="border-blue-100 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
      <CardContent className="p-5">
        <div className="flex items-start justify-between gap-4">
          <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-50 text-blue-700">
            <Icon className="h-5 w-5" />
          </div>

          <span className="rounded-full bg-blue-50 px-2.5 py-1 text-xs font-semibold text-blue-700">
            {trend}
          </span>
        </div>

        <p className="mt-5 text-sm text-slate-500">{label}</p>

        <p className="mt-2 text-3xl font-bold tracking-tight text-slate-950">
          {value}
        </p>

        <p className="mt-1 text-xs text-slate-500">{note}</p>
      </CardContent>
    </Card>
  );
}
