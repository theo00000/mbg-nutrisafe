import { apiFetch } from "@/lib/api";

export type LoginPayload = {
  email: string;
  password: string;
};

export type LoginResponse = {
  status: string;
  message: string;
  token: string;
  role: "school" | "spgg" | "umum";
};

export type RegisterPayload = {
  name: string;
  email: string;
  phone: string;
  password: string;
  role_name: "school" | "spgg" | "umum";
};

export type RegisterResponse = {
  status: string;
  message: string;
  data: {
    user_id: number;
    name: string;
    email: string;
    phone: string;
    role: string;
  };
};

export function login(payload: LoginPayload) {
  return apiFetch<LoginResponse>("/api/login", {
    method: "POST",
    body: payload,
  });
}

export function register(payload: RegisterPayload) {
  return apiFetch<RegisterResponse>("/api/register", {
    method: "POST",
    body: payload,
  });
}

export function logout(token: string) {
  return apiFetch<{ status: string; message: string }>("/api/logout", {
    method: "POST",
    token,
  });
}
