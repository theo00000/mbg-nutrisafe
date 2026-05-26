"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

import { AuthShell } from "@/components/auth/auth-shell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { register } from "@/services/auth-service";
import { MitraRegisterFields } from "@/components/auth/register-fields/mitra-register-fields";
import { SchoolRegisterFields } from "@/components/auth/register-fields/school-register-fields";
import { StudentRegisterFields } from "@/components/auth/register-fields/student-register-fields";

type RoleName = "school" | "spgg" | "umum";

const rolePreview = {
  school: {
    title: "Sekolah",
    description:
      "Untuk sekolah yang ingin memantau siswa, alergi, dan distribusi MBG.",
    className: "border-emerald-100 bg-emerald-50 text-emerald-700",
  },
  spgg: {
    title: "Mitra/SPPG",
    description:
      "Untuk mitra penyedia makanan yang mengelola produksi dan distribusi.",
    className: "border-blue-100 bg-blue-50 text-blue-700",
  },
  umum: {
    title: "Umum/Siswa",
    description:
      "Untuk siswa atau pengguna umum yang ingin mencatat profil dan alergi.",
    className: "border-orange-100 bg-orange-50 text-orange-700",
  },
};

export default function RegisterPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [roleName, setRoleName] = useState<RoleName>("school");

  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  // State tambahan Mitra/SPPG
  const [institutionName, setInstitutionName] = useState("");
  const [identityNumber, setIdentityNumber] = useState("");
  const [sppgName, setSppgName] = useState("");
  const [sppgAddress, setSppgAddress] = useState("");
  const [productionCapacity, setProductionCapacity] = useState("");

  // State tambahan Sekolah
  const [schoolName, setSchoolName] = useState("");
  const [npsn, setNpsn] = useState("");
  const [schoolAddress, setSchoolAddress] = useState("");
  const [schoolPIC, setSchoolPIC] = useState("");

  // State tambahan Umum/Siswa
  const [studentName, setStudentName] = useState("");
  const [studentSchool, setStudentSchool] = useState("");
  const [studentClass, setStudentClass] = useState("");
  const [initialAllergy, setInitialAllergy] = useState("");
  const [foodNote, setFoodNote] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setErrorMessage("");
    setSuccessMessage("");
    setIsLoading(true);

    if (
      roleName === "spgg" &&
      (!institutionName || !identityNumber || !sppgName)
    ) {
      setErrorMessage(
        "Untuk Mitra/SPPG, Nama Instansi, NIK/NPWP, dan Nama SPPG wajib diisi.",
      );
      setIsLoading(false);
      return;
    }

    if (roleName === "school" && (!schoolName.trim() || !npsn.trim())) {
      setErrorMessage("Untuk Sekolah, Nama Sekolah dan NPSN wajib diisi.");
      setIsLoading(false);
      return;
    }

    if (roleName === "umum" && (!studentName.trim() || !studentSchool.trim())) {
      setErrorMessage(
        "Untuk Umum/Siswa, Nama Siswa dan Sekolah Asal wajib diisi.",
      );
      setIsLoading(false);
      return;
    }

    try {
      const response = await register({
        name,
        email,
        phone,
        password,
        role_name: roleName,
      });

      setSuccessMessage(response.message || "Akun berhasil dibuat.");

      setTimeout(() => {
        router.push("/login");
      }, 1000);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Registrasi gagal. Coba lagi.";

      setErrorMessage(message);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <AuthShell
      title="Daftar Akun"
      description="Buat akun sesuai peran pengguna di ekosistem MBG NutriSafe."
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label>Nama</Label>
          <Input
            placeholder="Nama lengkap"
            value={name}
            onChange={(event) => setName(event.target.value)}
            className="h-11 rounded-2xl"
            required
          />
        </div>

        <div className="space-y-2">
          <Label>Email</Label>
          <Input
            type="email"
            placeholder="nama@email.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            className="h-11 rounded-2xl"
            required
          />
        </div>

        <div className="space-y-2">
          <Label>Nomor Telepon</Label>
          <Input
            placeholder="08xxxxxxxxxx"
            value={phone}
            onChange={(event) => setPhone(event.target.value)}
            className="h-11 rounded-2xl"
            required
          />
        </div>

        <div className="space-y-2">
          <Label>Password</Label>
          <Input
            type="password"
            placeholder="Minimal 6 karakter"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            className="h-11 rounded-2xl"
            required
          />
        </div>

        <div className="space-y-2">
          <Label>Role</Label>

          <select
            value={roleName}
            onChange={(event) => setRoleName(event.target.value as RoleName)}
            className="h-11 w-full rounded-2xl border border-input bg-background px-3 text-sm outline-none transition focus:border-blue-300 focus:ring-4 focus:ring-blue-100"
          >
            <option value="school">Sekolah</option>
            <option value="spgg">Mitra/SPPG</option>
            <option value="umum">Umum/Siswa</option>
          </select>

          <div
            className={`rounded-2xl border px-4 py-3 ${rolePreview[roleName].className}`}
          >
            <p className="text-sm font-bold">
              Mendaftar sebagai {rolePreview[roleName].title}
            </p>
            <p className="mt-1 text-xs leading-5 text-slate-600">
              {rolePreview[roleName].description}
            </p>
          </div>
        </div>

        {roleName === "school" ? (
          <SchoolRegisterFields
            schoolName={schoolName}
            npsn={npsn}
            schoolAddress={schoolAddress}
            schoolPIC={schoolPIC}
            setSchoolName={setSchoolName}
            setNpsn={setNpsn}
            setSchoolAddress={setSchoolAddress}
            setSchoolPIC={setSchoolPIC}
          />
        ) : null}

        {roleName === "spgg" ? (
          <MitraRegisterFields
            institutionName={institutionName}
            identityNumber={identityNumber}
            sppgName={sppgName}
            sppgAddress={sppgAddress}
            productionCapacity={productionCapacity}
            setInstitutionName={setInstitutionName}
            setIdentityNumber={setIdentityNumber}
            setSppgName={setSppgName}
            setSppgAddress={setSppgAddress}
            setProductionCapacity={setProductionCapacity}
          />
        ) : null}

        {roleName === "umum" ? (
          <StudentRegisterFields
            studentName={studentName}
            studentSchool={studentSchool}
            studentClass={studentClass}
            initialAllergy={initialAllergy}
            foodNote={foodNote}
            setStudentName={setStudentName}
            setStudentSchool={setStudentSchool}
            setStudentClass={setStudentClass}
            setInitialAllergy={setInitialAllergy}
            setFoodNote={setFoodNote}
          />
        ) : null}

        {errorMessage ? (
          <p className="rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-600">
            {errorMessage}
          </p>
        ) : null}

        {successMessage ? (
          <p className="rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
            {successMessage}
          </p>
        ) : null}

        <Button
          type="submit"
          className="h-11 w-full rounded-full bg-blue-600 font-semibold hover:bg-blue-700"
          disabled={isLoading}
        >
          {isLoading ? "Mendaftarkan..." : "Daftar Akun"}
        </Button>

        <p className="text-center text-sm text-slate-600">
          Sudah punya akun?{" "}
          <Link href="/login" className="font-semibold text-blue-600">
            Masuk di sini
          </Link>
        </p>
      </form>
    </AuthShell>
  );
}
