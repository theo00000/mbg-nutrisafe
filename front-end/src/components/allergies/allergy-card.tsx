"use client";

import { FormEvent, useEffect, useState } from "react";

import { getAuthToken } from "@/lib/auth-storage";
import {
  AllergyItem,
  createAllergy,
  getAllergies,
} from "@/services/allergy-service";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export function SchoolAllergyCard() {
  const [allergies, setAllergies] = useState<AllergyItem[]>([]);

  const [studentName, setStudentName] = useState("");
  const [className, setClassName] = useState("");
  const [allergyType, setAllergyType] = useState("");
  const [severity, setSeverity] = useState("Ringan");
  const [description, setDescription] = useState("");
  const [actionRequired, setActionRequired] = useState("");

  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  async function loadAllergies() {
    const token = getAuthToken();

    if (!token) {
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      setIsLoading(false);
      return;
    }

    try {
      setErrorMessage("");
      const data = await getAllergies(token);
      setAllergies(data);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Gagal memuat data alergi.";

      setErrorMessage(message);
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    loadAllergies();
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const token = getAuthToken();

    if (!token) {
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      return;
    }

    if (!studentName || !className || !allergyType || !severity) {
      setErrorMessage(
        "Nama siswa, kelas, jenis alergi, dan tingkat keparahan wajib diisi.",
      );
      return;
    }

    setIsSubmitting(true);
    setErrorMessage("");
    setSuccessMessage("");

    try {
      const created = await createAllergy(token, {
        student_name: studentName,
        class_name: className,
        allergy_type: allergyType,
        description,
        severity,
        action_required: actionRequired,
      });

      setAllergies((current) => [created, ...current]);

      setStudentName("");
      setClassName("");
      setAllergyType("");
      setSeverity("Ringan");
      setDescription("");
      setActionRequired("");

      setSuccessMessage("Data alergi siswa berhasil disimpan.");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Gagal menyimpan data alergi.";

      setErrorMessage(message);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Card className="rounded-[1.75rem] border-white/80 bg-white/90 shadow-sm">
      <CardContent className="p-6">
        <div className="mb-5 flex items-start justify-between gap-4">
          <div>
            <h2 className="font-heading text-xl font-black text-slate-950">
              Data Alergi Siswa
            </h2>
            <p className="mt-1 text-sm text-slate-500">
              Terhubung ke endpoint GET dan POST /school/allergy-data.
            </p>
          </div>

          <span className="rounded-full bg-emerald-100 px-4 py-2 text-xs font-bold text-emerald-700">
            Backend API
          </span>
        </div>

        <form
          onSubmit={handleSubmit}
          className="grid grid-cols-1 gap-4 md:grid-cols-2"
        >
          <div>
            <Label>Nama Siswa</Label>
            <Input
              value={studentName}
              onChange={(event) => setStudentName(event.target.value)}
              className="mt-2 h-11 rounded-2xl"
              placeholder="Contoh: Vina"
            />
          </div>

          <div>
            <Label>Kelas</Label>
            <Input
              value={className}
              onChange={(event) => setClassName(event.target.value)}
              className="mt-2 h-11 rounded-2xl"
              placeholder="Contoh: 6A"
            />
          </div>

          <div>
            <Label>Jenis Alergi</Label>
            <Input
              value={allergyType}
              onChange={(event) => setAllergyType(event.target.value)}
              className="mt-2 h-11 rounded-2xl"
              placeholder="Contoh: Alergi Kacang"
            />
          </div>

          <div>
            <Label>Tingkat Keparahan</Label>
            <select
              value={severity}
              onChange={(event) => setSeverity(event.target.value)}
              className="mt-2 h-11 w-full rounded-2xl border border-input bg-background px-3 text-sm"
            >
              <option value="Ringan">Ringan</option>
              <option value="Sedang">Sedang</option>
              <option value="Berat">Berat</option>
            </select>
          </div>

          <div>
            <Label>Deskripsi</Label>
            <Input
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              className="mt-2 h-11 rounded-2xl"
              placeholder="Contoh: Gatal-gatal setelah konsumsi kacang"
            />
          </div>

          <div>
            <Label>Tindakan yang Diperlukan</Label>
            <Input
              value={actionRequired}
              onChange={(event) => setActionRequired(event.target.value)}
              className="mt-2 h-11 rounded-2xl"
              placeholder="Contoh: Pisahkan menu kacang"
            />
          </div>

          <div className="md:col-span-2">
            {errorMessage ? (
              <p className="mb-3 rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-600">
                {errorMessage}
              </p>
            ) : null}

            {successMessage ? (
              <p className="mb-3 rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
                {successMessage}
              </p>
            ) : null}

            <Button
              type="submit"
              disabled={isSubmitting}
              className="rounded-xl bg-emerald-600 px-8 hover:bg-emerald-700"
            >
              {isSubmitting ? "Menyimpan..." : "Simpan Data Alergi"}
            </Button>
          </div>
        </form>

        <div className="mt-7">
          <h3 className="font-heading text-base font-black text-slate-950">
            Daftar Alergi Siswa
          </h3>

          {isLoading ? (
            <p className="mt-3 text-sm text-slate-500">Memuat data alergi...</p>
          ) : null}

          {!isLoading && allergies.length === 0 ? (
            <p className="mt-3 rounded-2xl border border-dashed border-emerald-200 bg-emerald-50 px-4 py-5 text-sm text-slate-500">
              Belum ada data alergi siswa.
            </p>
          ) : null}

          <div className="mt-4 grid grid-cols-1 gap-3 md:grid-cols-2">
            {allergies.map((item) => (
              <div
                key={item.id}
                className="rounded-2xl border border-emerald-100 bg-emerald-50 p-4"
              >
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-sm font-bold text-slate-950">
                      {item.student_name}
                    </p>
                    <p className="mt-1 text-xs text-slate-500">
                      Kelas {item.class_name}
                    </p>
                  </div>

                  <span className="rounded-full bg-white px-3 py-1 text-xs font-bold text-emerald-700">
                    {item.severity}
                  </span>
                </div>

                <p className="mt-3 text-sm font-semibold text-slate-800">
                  {item.allergy_type}
                </p>

                <p className="mt-1 text-xs leading-5 text-slate-500">
                  {item.description || "Tidak ada deskripsi."}
                </p>

                <p className="mt-3 text-xs font-medium text-slate-600">
                  Tindakan: {item.action_required || "-"}
                </p>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
