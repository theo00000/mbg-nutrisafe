import { apiFetch } from "@/lib/api";

export type DashboardStatsResponse = {
  status?: string;
  message?: string;
  data?: Record<string, unknown>;
};

export function getDashboardStats() {
  return apiFetch<DashboardStatsResponse>("/api/dashboard/stats");
}
