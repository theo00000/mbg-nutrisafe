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
  const [severity, setSeverity] = useState("ringan");
  const [description, setDescription] = useState("");
  const [actionRequired, setActionRequired] = useState("");

  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    const token = getAuthToken();

    if (!token) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setIsLoading(false);
      return;
    }

    getAllergies(token)
      .then((data) => {
        setAllergies(data);
      })
      .catch((error) => {
        const message =
          error instanceof Error ? error.message : "Gagal memuat data alergi.";

        setErrorMessage(message);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const token = getAuthToken();

    if (!token) {
      setErrorMessage("Token tidak ditemukan. Silakan login ulang.");
      return;
    }

    if (!studentName.trim() || !className.trim() || !allergyType.trim()) {
      setErrorMessage("Nama siswa, kelas, dan jenis alergi wajib diisi.");
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
        severity,
        description,
        action_required: actionRequired,
      });

      setAllergies((current) => [created, ...current]);

      setStudentName("");
      setClassName("");
      setAllergyType("");
      setSeverity("ringan");
      setDescription("");
      setActionRequired("");

      setSuccessMessage("Data alergi berhasil ditambahkan.");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Gagal menyimpan data alergi.";

      setErrorMessage(message);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Card className="rounded-[1.75rem] border-white bg-white/85 shadow-sm backdrop-blur">
      <CardContent className="p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-lg font-bold text-slate-950">
              Data Alergi Siswa
            </h2>

            <p className="mt-2 text-sm leading-6 text-slate-500">
              Data ini terhubung ke endpoint GET dan POST /api/allergies untuk
              role Sekolah.
            </p>
          </div>

          <div className="rounded-full bg-emerald-50 px-4 py-2 text-xs font-bold text-emerald-700">
            Backend API
          </div>
        </div>

        <form
          onSubmit={handleSubmit}
          className="mt-6 grid grid-cols-[1fr_1fr] gap-4"
        >
          <div className="space-y-2">
            <Label>Nama Siswa</Label>
            <Input
              value={studentName}
              onChange={(event) => setStudentName(event.target.value)}
              placeholder="Contoh: Andi Pratama"
              className="h-11 rounded-2xl"
            />
          </div>

          <div className="space-y-2">
            <Label>Kelas</Label>
            <Input
              value={className}
              onChange={(event) => setClassName(event.target.value)}
              placeholder="Contoh: 5A"
              className="h-11 rounded-2xl"
            />
          </div>

          <div className="space-y-2">
            <Label>Jenis Alergi</Label>
            <Input
              value={allergyType}
              onChange={(event) => setAllergyType(event.target.value)}
              placeholder="Contoh: kacang, susu, seafood"
              className="h-11 rounded-2xl"
            />
          </div>

          <div className="space-y-2">
            <Label>Tingkat Keparahan</Label>
            <select
              value={severity}
              onChange={(event) => setSeverity(event.target.value)}
              className="h-11 w-full rounded-2xl border border-input bg-background px-3 text-sm outline-none transition focus:border-emerald-300 focus:ring-4 focus:ring-emerald-100"
            >
              <option value="ringan">Ringan</option>
              <option value="sedang">Sedang</option>
              <option value="berat">Berat</option>
            </select>
          </div>

          <div className="space-y-2">
            <Label>Deskripsi</Label>
            <Input
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              placeholder="Contoh: muncul ruam setelah konsumsi susu"
              className="h-11 rounded-2xl"
            />
          </div>

          <div className="space-y-2">
            <Label>Tindakan yang Diperlukan</Label>
            <Input
              value={actionRequired}
              onChange={(event) => setActionRequired(event.target.value)}
              placeholder="Contoh: sediakan menu pengganti"
              className="h-11 rounded-2xl"
            />
          </div>

          <div className="col-span-2">
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
              className="h-11 rounded-full bg-emerald-600 px-6 font-semibold hover:bg-emerald-700"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Menyimpan..." : "Tambah Data Alergi"}
            </Button>
          </div>
        </form>

        <div className="mt-7">
          <h3 className="text-sm font-bold text-slate-950">
            Daftar Alergi Siswa
          </h3>

          {isLoading ? (
            <p className="mt-3 text-sm text-slate-500">Memuat data alergi...</p>
          ) : null}

          {!isLoading && allergies.length === 0 ? (
            <p className="mt-3 rounded-2xl border border-dashed border-emerald-200 bg-emerald-50/60 px-4 py-5 text-sm text-slate-500">
              Belum ada data alergi siswa.
            </p>
          ) : null}

          <div className="mt-4 grid grid-cols-2 gap-4">
            {allergies.map((item) => (
              <div
                key={item.id}
                className="rounded-[1.5rem] border border-emerald-100 bg-emerald-50/70 p-4"
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
                  {item.description || "Tidak ada deskripsi tambahan."}
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
