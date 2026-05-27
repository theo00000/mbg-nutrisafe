"use client";

import { FormEvent, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type AllergyItem = {
  id: number;
  name: string;
  description: string;
};

const initialAllergies: AllergyItem[] = [
  {
    id: 1,
    name: "Kacang",
    description: "Hindari menu dengan kacang tanah atau olahan kacang.",
  },
  {
    id: 2,
    name: "Susu",
    description: "Perlu alternatif menu tanpa produk susu.",
  },
];

export function AllergyCard() {
  const [allergies, setAllergies] = useState(initialAllergies);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const newAllergy: AllergyItem = {
      id: Date.now(),
      name,
      description,
    };

    setAllergies((current) => [newAllergy, ...current]);
    setName("");
    setDescription("");
  }

  return (
    <Card className="rounded-[1.75rem] border-white bg-white/85 shadow-sm backdrop-blur">
      {" "}
      <CardContent className="p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold text-slate-950">
              Data Alergi
            </h2>

            <p className="mt-2 text-sm text-slate-500">
              TODO: sementara memakai dummy data karena route /api/allergies
              belum terdaftar di backend saat ini.
            </p>
          </div>

          <Badge className="bg-orange-100 text-orange-700 hover:bg-orange-100">
            Umum/Siswa
          </Badge>
        </div>

        <form onSubmit={handleSubmit} className="mt-5 space-y-3">
          <div className="space-y-2">
            <Label>Nama Alergi</Label>
            <Input
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder="Contoh: telur, seafood, kacang"
              className="h-11 rounded-2xl"
              required
            />
          </div>

          <div className="space-y-2">
            <Label>Catatan</Label>
            <Input
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              placeholder="Catatan singkat alergi"
              className="h-11 rounded-2xl"
              required
            />
          </div>

          <Button
            type="submit"
            className="h-11 rounded-full bg-orange-500 px-6 font-semibold hover:bg-orange-600"
          >
            Tambah Alergi
          </Button>
        </form>

        <div className="mt-6 space-y-3">
          {allergies.map((item) => (
            <div
              key={item.id}
              className="rounded-2xl border border-orange-100 bg-orange-50/80 p-4 transition hover:bg-orange-100/70"
            >
              <p className="text-sm font-semibold text-slate-900">
                {item.name}
              </p>
              <p className="mt-1 text-sm text-slate-600">{item.description}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
