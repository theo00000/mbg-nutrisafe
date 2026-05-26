import { apiFetch } from "@/lib/api";

export type ProfileDetail = {
  id?: number;
  user_id?: number;
  name: string;
  email: string;
  phone?: string;
  role?: string;
  role_name?: string;
};

export type UpdateProfilePayload = {
  name?: string;
  phone?: string;
};

type ProfileDetailResponse = {
  status: string;
  message: string;
  data: ProfileDetail;
};

export async function getProfileDetail(token: string) {
  const response = await apiFetch<ProfileDetailResponse>(
    "/api/profile/detail",
    {
      method: "GET",
      token,
    },
  );

  return response.data;
}

export async function updateProfileDetail(
  token: string,
  payload: UpdateProfilePayload,
) {
  const response = await apiFetch<ProfileDetailResponse>(
    "/api/profile/detail",
    {
      method: "PUT",
      token,
      body: payload,
    },
  );

  return response.data;
}
