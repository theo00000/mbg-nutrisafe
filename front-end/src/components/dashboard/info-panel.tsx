import { ReactNode } from "react";

import { Card, CardContent } from "@/components/ui/card";

type InfoPanelProps = {
  title: string;
  description: string;
  children: ReactNode;
};

export function InfoPanel({ title, description, children }: InfoPanelProps) {
  return (
    <Card className="rounded-[1.75rem] border-white bg-white/85 shadow-sm backdrop-blur">
      <CardContent className="p-6">
        <h2 className="text-lg font-bold text-slate-950">{title}</h2>

        <p className="mt-2 text-sm leading-6 text-slate-500">{description}</p>

        <div className="mt-5">{children}</div>
      </CardContent>
    </Card>
  );
}
