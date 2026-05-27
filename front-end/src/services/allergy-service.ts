import { apiFetch } from "@/lib/api";

export type AllergyItem = {
  id: number;
  school: {
    id: number;
    name_school: string;
  };
  student_name: string;
  class_name: string;
  allergy_type: string;
  description: string;
  severity: string;
  action_required: string;
};

type AllergyListResponse = {
  status: string;
  data: AllergyItem[];
};

type AllergyCreateResponse = {
  status: string;
  message: string;
  data: AllergyItem;
};

export type CreateAllergyPayload = {
  student_name: string;
  class_name: string;
  allergy_type: string;
  description?: string;
  severity: string;
  action_required?: string;
};

export function getAllergies(token: string) {
  return apiFetch<AllergyListResponse>("/api/allergies", {
    method: "GET",
    token,
  }).then((response) => response.data);
}

export function createAllergy(token: string, payload: CreateAllergyPayload) {
  return apiFetch<AllergyCreateResponse>("/api/allergies", {
    method: "POST",
    token,
    body: payload,
  }).then((response) => response.data);
}
