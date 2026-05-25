export type UserRole = "school" | "spgg" | "umum";

const TOKEN_KEY = "nutrisafe_token";
const ROLE_KEY = "nutrisafe_role";

export function saveAuth(token: string, role: UserRole) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(ROLE_KEY, role);
}

export function getAuthToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function getAuthRole() {
  return localStorage.getItem(ROLE_KEY) as UserRole | null;
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(ROLE_KEY);
}

export function getDashboardPathByRole(role: UserRole) {
  if (role === "school") return "/sekolah";
  if (role === "spgg") return "/mitra";
  if (role === "umum") return "/umum";

  return "/";
}
