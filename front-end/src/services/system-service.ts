import { apiFetch } from "@/lib/api";

export type PingResponse = {
  message?: string;
  status?: string;
};

export function getPing() {
  return apiFetch<PingResponse>("/api/ping");
}
