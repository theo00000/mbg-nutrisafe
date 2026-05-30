import Image from "next/image";
import Link from "next/link";
import {
  ArrowRight,
  Building2,
  FileText,
  Megaphone,
  School,
  ShieldCheck,
  Users,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

const stats = [
  {
    label: "Sekolah Terdaftar",
    value: "443.703",
    icon: School,
    className: "text-emerald-600 bg-emerald-50",
  },
  {
    label: "Mitra SPPG",
    value: "24.000",
    icon: Building2,
    className: "text-blue-600 bg-blue-50",
  },
  {
    label: "Siswa",
    value: "59,86 Juta",
    icon: Users,
    className: "text-orange-500 bg-orange-50",
  },
];

const mitraSteps = [
  {
    number: "01",
    title: "Registrasi Mitra",
    description:
      "Daftarkan diri melalui halaman pendaftaran mitra dengan mengisi data diri dan informasi usaha secara lengkap.",
  },
  {
    number: "02",
    title: "Pengajuan Data SPPG",
    description:
      "Lengkapi nama SPPG, alamat, kapasitas produksi, dan dokumen pendukung yang dibutuhkan.",
  },
  {
    number: "03",
    title: "Proses Verifikasi",
    description:
      "Tim MBG melakukan verifikasi terhadap data dan dokumen yang telah diajukan.",
  },
  {
    number: "04",
    title: "Persetujuan Mitra",
    description:
      "Mitra yang lolos verifikasi akan ditetapkan sebagai mitra SPPG resmi.",
  },
  {
    number: "05",
    title: "Penugasan Distribusi",
    description:
      "Mitra aktif menerima penugasan distribusi ke sekolah dan memantau kegiatan melalui sistem.",
  },
];

const accessCards = [
  {
    title: "Untuk Mitra",
    description:
      "Monitoring distribusi, sekolah tujuan, dan pelaporan menu harian.",
    image: "/images/landing/Mitra.svg",
    actions: [
      {
        label: "Monitoring Allergy Siswa",
        href: "/mitra",
        icon: ShieldCheck,
        className: "bg-blue-600 hover:bg-blue-700",
      },
      {
        label: "Data Sekolah",
        href: "/mitra",
        icon: School,
        className: "bg-blue-600 hover:bg-blue-700",
      },
      {
        label: "Pelaporan Menu Harian",
        href: "/mitra",
        icon: FileText,
        className: "bg-blue-600 hover:bg-blue-700",
      },
    ],
  },
  {
    title: "Untuk Umum",
    description: "Laporkan kendala makanan dan pantau informasi keamanan MBG.",
    image: "/images/landing/Umum.svg",
    actions: [
      {
        label: "Pelaporan Makanan",
        href: "/umum",
        icon: Megaphone,
        className: "bg-orange-500 hover:bg-orange-600",
      },
    ],
  },
  {
    title: "Untuk Sekolah",
    description: "Kelola data alergi siswa dan laporan makanan dari sekolah.",
    image: "/images/landing/Sekolah.svg",
    actions: [
      {
        label: "Menambahkan Data Alergi Siswa",
        href: "/sekolah",
        icon: ShieldCheck,
        className: "bg-emerald-600 hover:bg-emerald-700",
      },
      {
        label: "Pelaporan Makanan",
        href: "/sekolah",
        icon: FileText,
        className: "bg-emerald-600 hover:bg-emerald-700",
      },
    ],
  },
];

export default function LandingPage() {
  return (
    <main className="min-h-screen bg-sky-100 text-slate-950">
      <section className="mx-auto min-h-screen w-full max-w-[1440px] overflow-hidden bg-sky-100 shadow-none lg:shadow-2xl">
        <nav className="relative z-20 flex flex-col gap-5 bg-sky-100/90 px-5 py-5 backdrop-blur md:px-8 lg:flex-row lg:items-center lg:justify-between lg:px-12">
          <Link href="/" className="flex items-center gap-3">
            <div className="relative h-12 w-12 overflow-hidden rounded-2xl bg-white shadow-sm">
              <Image
                src="/images/nutrisafe-logo.svg"
                alt="NutriSafe MBG"
                fill
                className="object-cover"
              />
            </div>

            <div>
              <p className="font-heading text-xl font-black tracking-tight text-blue-950">
                NutriSafe MBG
              </p>
              <p className="text-xs font-medium text-blue-700">
                Food Safety Monitoring
              </p>
            </div>
          </Link>
          <div className="flex w-full flex-wrap items-center justify-center gap-2 rounded-3xl bg-white/70 px-4 py-3 text-xs font-bold text-blue-950 shadow-sm backdrop-blur sm:text-sm lg:w-auto lg:gap-7 lg:rounded-full lg:px-7">
            <a
              href="#beranda"
              className="rounded-full bg-blue-600 px-4 py-2 text-white sm:px-5"
            >
              Beranda
            </a>
            <a href="#tentang" className="transition hover:text-blue-600">
              Tentang MBG
            </a>
            <a href="#mitra" className="transition hover:text-blue-600">
              Cara Menjadi Mitra
            </a>
            <Link href="/login" className="transition hover:text-blue-600">
              Masuk/daftar
            </Link>
          </div>
        </nav>
        <section id="beranda" className="relative">
          <div className="relative h-[620px] overflow-hidden sm:h-[680px] lg:h-[640px]">
            <Image
              src="/images/landing/hero-section.svg"
              alt="Program makan bergizi untuk siswa"
              fill
              priority
              className="object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-b from-sky-100/95 via-sky-100/60 to-sky-100/10 lg:bg-gradient-to-r lg:from-sky-100/95 lg:via-sky-100/45 lg:to-transparent" />
            <div className="absolute left-5 right-5 top-10 max-w-xl sm:left-8 lg:left-12 lg:right-auto lg:top-14">
              <p className="inline-flex rounded-full bg-white/80 px-4 py-2 text-sm font-bold text-blue-700 shadow-sm backdrop-blur">
                Platform Keamanan Makanan Siswa
              </p>
              <h1 className="font-heading mt-5 text-4xl font-black leading-tight tracking-tight text-blue-950 sm:text-5xl">
                Monitoring & Pelaporan Keamanan Makanan Siswa MBG.
              </h1>
              <p className="mt-5 max-w-lg text-base leading-7 text-slate-700">
                NutriSafe membantu sekolah, mitra SPPG, dan pengguna umum
                memantau distribusi makanan, alergi siswa, serta laporan
                keamanan pangan secara lebih terstruktur.
              </p>
              <div className="mt-7 flex flex-col gap-3 sm:flex-row sm:items-center">
                <Button
                  asChild
                  className="h-12 rounded-full bg-blue-600 px-7 font-bold hover:bg-blue-700"
                >
                  <Link href="/register">
                    Daftar Menjadi Mitra
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
                <Button
                  asChild
                  variant="outline"
                  className="h-12 rounded-full bg-white/80 px-7 font-bold"
                >
                  <Link href="/login">Masuk Dashboard</Link>
                </Button>
              </div>
            </div>
          </div>
        </section>
        <section className="grid grid-cols-1 gap-4 border-y border-slate-200 bg-white px-5 py-6 sm:grid-cols-3 sm:gap-0 md:px-8 lg:px-12">
          {stats.map((item, index) => {
            const Icon = item.icon;

            return (
              <div
                key={item.label}
                className={`flex items-center justify-center gap-5 ${
                  index !== stats.length - 1
                    ? "sm:border-r sm:border-slate-300"
                    : ""
                }`}
              >
                <div
                  className={`flex h-16 w-16 items-center justify-center rounded-2xl ${item.className}`}
                >
                  <Icon className="h-8 w-8" />
                </div>

                <div>
                  <p className="font-heading text-3xl font-black tracking-tight">
                    {item.value}
                  </p>
                  <p className="text-sm font-bold text-slate-700">
                    {item.label}
                  </p>
                </div>
              </div>
            );
          })}
        </section>
        <section
          id="tentang"
          className="bg-sky-200 px-5 py-12 md:px-8 lg:px-12"
        >
          <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="text-sm font-black text-blue-700">
                Akses Berdasarkan Peran
              </p>
              <h2 className="font-heading mt-2 text-3xl font-black tracking-tight text-slate-950 sm:text-4xl">
                Pilih kebutuhan monitoring NutriSafe.
              </h2>
            </div>
          </div>
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3 lg:gap-8">
            {accessCards.map((card) => (
              <Card
                key={card.title}
                className="group overflow-hidden rounded-[1.75rem] border-white bg-white/95 shadow-sm transition hover:-translate-y-1 hover:shadow-2xl"
              >
                <CardContent className="flex min-h-[390px] flex-col p-6 lg:min-h-[420px]">
                  <div className="relative h-44 overflow-hidden rounded-[1.5rem] bg-sky-50">
                    <Image
                      src={card.image}
                      alt={card.title}
                      fill
                      className="object-contain p-4 transition duration-500 group-hover:scale-105"
                    />
                  </div>
                  <h3 className="font-heading mt-6 text-center text-xl font-black text-slate-950">
                    {card.title}
                  </h3>
                  <p className="mt-2 min-h-[48px] text-center text-sm leading-6 text-slate-500">
                    {card.description}
                  </p>
                  <div className="mt-6 space-y-3">
                    {card.actions.map((action) => {
                      const ActionIcon = action.icon;

                      return (
                        <Button
                          key={action.label}
                          asChild
                          className={`h-10 w-full justify-start rounded-xl text-sm font-bold text-white ${action.className}`}
                        >
                          <Link href={action.href}>
                            <ActionIcon className="mr-2 h-4 w-4" />
                            {action.label}
                          </Link>
                        </Button>
                      );
                    })}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
        <section
          id="mitra"
          className="relative overflow-hidden bg-gradient-to-br from-white via-sky-50 to-blue-100 px-5 py-14 md:px-8 lg:px-12 lg:py-16"
        >
          <div className="absolute -right-24 -top-24 h-72 w-72 rounded-full bg-blue-200/50 blur-3xl" />
          <div className="absolute -bottom-24 left-20 h-72 w-72 rounded-full bg-emerald-200/40 blur-3xl" />

          <div className="relative z-10 grid grid-cols-1 items-center gap-10 lg:grid-cols-[1fr_0.9fr] lg:gap-12">
            <div>
              <p className="inline-flex rounded-full bg-blue-100 px-4 py-2 text-sm font-black text-blue-700">
                Cara Menjadi Mitra
              </p>

              <h2 className="font-heading mt-5 max-w-2xl text-3xl font-black leading-tight tracking-tight text-slate-950 sm:text-4xl">
                Daftarkan SPPG dan mulai kelola distribusi makanan dengan aman.
              </h2>

              <p className="mt-4 max-w-2xl text-sm leading-7 text-slate-600 sm:text-base">
                Mitra dapat menggunakan NutriSafe untuk memantau produksi
                makanan, melihat sekolah tujuan, mencatat distribusi, dan
                memastikan makanan aman sebelum diterima oleh siswa.
              </p>

              <div className="mt-8 grid gap-4">
                {mitraSteps.map((step) => (
                  <div
                    key={step.number}
                    className="flex gap-4 rounded-[1.5rem] border border-blue-100 bg-white/80 p-4 shadow-sm backdrop-blur transition hover:-translate-y-0.5 hover:shadow-md"
                  >
                    <div className="font-heading flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-blue-600 text-sm font-black text-white">
                      {step.number}
                    </div>

                    <div>
                      <h3 className="font-heading text-base font-black text-slate-950">
                        {step.title}
                      </h3>
                      <p className="mt-1 text-sm leading-6 text-slate-500">
                        {step.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>

              <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                <Button
                  asChild
                  className="h-12 rounded-full bg-blue-600 px-7 font-bold hover:bg-blue-700"
                >
                  <Link href="/register">
                    Daftar Mitra
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>

                <Button
                  asChild
                  variant="outline"
                  className="h-12 rounded-full bg-white/80 px-7 font-bold"
                >
                  <Link href="/login">Masuk Dashboard</Link>
                </Button>
              </div>
            </div>

            <div className="relative">
              <div className="relative h-[280px] overflow-hidden rounded-[2rem] border border-white bg-white shadow-2xl sm:h-[360px] lg:h-[420px] lg:rounded-[2.25rem]">
                <Image
                  src="/images/landing/delivery-mbg.svg"
                  alt="Distribusi makanan MBG"
                  fill
                  className="object-contain bg-sky-50 p-6"
                />
              </div>

              <div className="mt-4 rounded-[1.75rem] border border-white bg-white/90 p-5 shadow-xl backdrop-blur lg:absolute lg:-bottom-6 lg:left-8 lg:right-8 lg:mt-0">
                <div className="grid grid-cols-3 gap-4 text-center">
                  <div>
                    <p className="font-heading text-2xl font-black text-blue-700">
                      24K
                    </p>
                    <p className="mt-1 text-xs font-semibold text-slate-500">
                      Mitra SPPG
                    </p>
                  </div>

                  <div>
                    <p className="font-heading text-2xl font-black text-emerald-600">
                      Aman
                    </p>
                    <p className="mt-1 text-xs font-semibold text-slate-500">
                      Distribusi
                    </p>
                  </div>

                  <div>
                    <p className="font-heading text-2xl font-black text-orange-500">
                      MBG
                    </p>
                    <p className="mt-1 text-xs font-semibold text-slate-500">
                      Monitoring
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </section>
    </main>
  );
}
