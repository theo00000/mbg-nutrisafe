import { ReactNode } from "react";

type DashboardHeaderProps = {
  eyebrow: string;
  title: string;
  description: string;
  accentClassName: string;
  action?: ReactNode;
};

export function DashboardHeader({
  eyebrow,
  title,
  description,
  accentClassName,
  action,
}: DashboardHeaderProps) {
  return (
    <div className="mb-8 flex items-start justify-between rounded-[2rem] border border-white bg-white/80 p-6 shadow-sm backdrop-blur">
      <div>
        <p className={`text-sm font-bold ${accentClassName}`}>{eyebrow}</p>

        <h1 className="mt-3 text-4xl font-black tracking-tight text-slate-950">
          {title}
        </h1>

        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-600">
          {description}
        </p>
      </div>

      {action ? <div>{action}</div> : null}
    </div>
  );
}
