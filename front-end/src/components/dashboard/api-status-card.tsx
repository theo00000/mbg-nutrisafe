"use client";

import { useEffect, useState } from "react";
import { Wifi, WifiOff } from "lucide-react";

import { getPing } from "@/services/system-service";

type ApiStatus = "checking" | "online" | "offline";

export function ApiStatusCard() {
  const [status, setStatus] = useState<ApiStatus>("checking");

  useEffect(() => {
    getPing()
      .then(() => setStatus("online"))
      .catch(() => setStatus("offline"));
  }, []);

  const isOnline = status === "online";
  const isChecking = status === "checking";

  return (
    <div className="rounded-2xl bg-white px-4 py-3 shadow-sm">
      <div className="flex items-center gap-2">
        {isOnline ? (
          <Wifi className="h-4 w-4 text-blue-600" />
        ) : (
          <WifiOff className="h-4 w-4 text-red-500" />
        )}

        <div>
          <p className="text-xs text-slate-500">Backend API</p>
          <p className="text-sm font-semibold text-slate-800">
            {isChecking
              ? "Mengecek..."
              : isOnline
                ? "Terhubung"
                : "Tidak terhubung"}
          </p>
        </div>
      </div>
    </div>
  );
}
