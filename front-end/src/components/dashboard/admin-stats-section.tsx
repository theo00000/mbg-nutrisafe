"use client";

import { useEffect, useState } from "react";
import {
  Building2,
  ClipboardList,
  School,
  ShieldAlert,
  Users,
} from "lucide-react";

import { DashboardStatCard } from "@/components/dashboard/stat-card";
import { getDashboardStats } from "@/services/dashboard-service";

const fallbackAdminStats = [
  {
    label: "Total Sekolah",
    value: "24",
    note: "Fallback dummy",
    trend: "+8%",
    icon: School,
  },
  {
    label: "Total Mitra",
    value: "12",
    note: "Fallback dummy",
    trend: "+3%",
    icon: Building2,
  },
  {
    label: "Total Siswa",
    value: "0",
    note: "Fallback dummy",
    trend: "+0%",
    icon: Users,
  },
  {
    label: "Kasus Prioritas",
    value: "4",
    note: "Belum tersedia dari API",
    trend: "Urgent",
    icon: ShieldAlert,
  },
];

function getNumberValue(
  data: Record<string, unknown> | undefined,
  keys: string[],
) {
  if (!data) return undefined;

  for (const key of keys) {
    const value = data[key];

    if (typeof value === "number") {
      return String(value);
    }

    if (typeof value === "string" && value.trim()) {
      return value;
    }
  }

  return undefined;
}

function mapStatsFromApi(data: Record<string, unknown> | undefined) {
  return [
    {
      ...fallbackAdminStats[0],
      value:
        getNumberValue(data, [
          "total_school",
          "total_schools",
          "totalSchools",
        ]) ?? fallbackAdminStats[0].value,
      note: "Dari Backend API",
    },
    {
      ...fallbackAdminStats[1],
      value:
        getNumberValue(data, ["total_spgg", "total_partners", "partners"]) ??
        fallbackAdminStats[1].value,
      note: "Dari Backend API",
    },
    {
      ...fallbackAdminStats[2],
      label: "Total Siswa",
      value:
        getNumberValue(data, ["total_student", "total_students", "students"]) ??
        fallbackAdminStats[2].value,
      note: "Dari Backend API",
    },
    {
      ...fallbackAdminStats[3],
      value:
        getNumberValue(data, [
          "priority_cases",
          "urgent_cases",
          "kasus_prioritas",
        ]) ?? fallbackAdminStats[3].value,
      note: "Fallback dummy",
    },
  ];
}

export function AdminStatsSection() {
  const [stats, setStats] = useState(fallbackAdminStats);
  const [isLoading, setIsLoading] = useState(true);
  const [source, setSource] = useState<"api" | "fallback">("fallback");

  useEffect(() => {
    getDashboardStats()
      .then((response) => {
        console.log("DASHBOARD STATS RESPONSE:", response);

        setSource("api");
        setStats(mapStatsFromApi(response.data));
      })
      .catch((error) => {
        console.error("Gagal mengambil dashboard stats:", error);

        setSource("fallback");
        setStats(fallbackAdminStats);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  return (
    <div>
      <p className="mb-3 text-xs font-medium text-slate-500">
        Sumber data: {source === "api" ? "Backend API" : "Fallback dummy"}
      </p>

      {isLoading ? (
        <p className="mb-3 text-xs font-medium text-slate-500">
          Sumber data:{" "}
          {source === "api" ? "GET /api/dashboard/stats" : "Fallback dummy"}
        </p>
      ) : null}

      <div className="grid grid-cols-2 gap-5 xl:grid-cols-4">
        {" "}
        {stats.map((item) => (
          <DashboardStatCard
            key={item.label}
            label={item.label}
            value={item.value}
            note={item.note}
            trend={item.trend}
            icon={item.icon}
          />
        ))}
      </div>
    </div>
  );
}
