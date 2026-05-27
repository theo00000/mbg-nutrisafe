import type { LucideIcon } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";

type RoleStatCardProps = {
  label: string;
  value: string;
  note: string;
  icon: LucideIcon;
  accentClassName: string;
  iconClassName: string;
};

export function RoleStatCard({
  label,
  value,
  note,
  icon: Icon,
  accentClassName,
  iconClassName,
}: RoleStatCardProps) {
  return (
    <Card className="rounded-[1.75rem] border-white bg-white/85 shadow-sm backdrop-blur transition hover:-translate-y-1 hover:shadow-xl">
      <CardContent className="p-5">
        <div
          className={`flex h-12 w-12 items-center justify-center rounded-2xl ${iconClassName}`}
        >
          <Icon className="h-5 w-5" />
        </div>

        <p className="mt-5 text-sm text-slate-500">{label}</p>

        <p className="mt-2 text-4xl font-black tracking-tight text-slate-950">
          {value}
        </p>

        <p className={`mt-1 text-xs font-medium ${accentClassName}`}>{note}</p>
      </CardContent>
    </Card>
  );
}
