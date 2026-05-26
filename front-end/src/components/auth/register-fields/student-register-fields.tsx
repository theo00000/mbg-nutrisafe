import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type StudentRegisterFieldsProps = {
  studentName: string;
  studentSchool: string;
  studentClass: string;
  initialAllergy: string;
  foodNote: string;
  setStudentName: (value: string) => void;
  setStudentSchool: (value: string) => void;
  setStudentClass: (value: string) => void;
  setInitialAllergy: (value: string) => void;
  setFoodNote: (value: string) => void;
};

export function StudentRegisterFields({
  studentName,
  studentSchool,
  studentClass,
  initialAllergy,
  foodNote,
  setStudentName,
  setStudentSchool,
  setStudentClass,
  setInitialAllergy,
  setFoodNote,
}: StudentRegisterFieldsProps) {
  return (
    <div className="rounded-3xl border border-orange-100 bg-orange-50/80 p-5">
      <div className="rounded-2xl bg-orange-500 px-4 py-3 text-center text-sm font-semibold text-white shadow-sm">
        Formulir Data Umum/Siswa
      </div>

      <p className="mt-3 rounded-2xl bg-white px-4 py-3 text-xs leading-5 text-slate-500">
        TODO: data tambahan siswa dan alergi awal belum dikirim ke backend
        karena endpoint register saat ini hanya menerima data akun dasar.
      </p>

      <div className="mt-5 space-y-4 rounded-2xl bg-white p-4">
        <div className="space-y-2 h-11 rounded-2xl">
          <Label>
            Nama Siswa <span className="text-red-500">*</span>
          </Label>
          <Input
            placeholder="Contoh: Andi Pratama"
            value={studentName}
            onChange={(event) => setStudentName(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>
            Sekolah Asal <span className="text-red-500">*</span>
          </Label>
          <Input
            placeholder="Contoh: SD Negeri 01 Jakarta"
            value={studentSchool}
            onChange={(event) => setStudentSchool(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>Kelas</Label>
          <Input
            placeholder="Contoh: 5A"
            value={studentClass}
            onChange={(event) => setStudentClass(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>Alergi Awal</Label>
          <Input
            placeholder="Contoh: kacang, susu, seafood"
            value={initialAllergy}
            onChange={(event) => setInitialAllergy(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>Catatan Makanan</Label>
          <Input
            placeholder="Contoh: tidak suka pedas, perlu menu rendah gula"
            value={foodNote}
            onChange={(event) => setFoodNote(event.target.value)}
          />
        </div>
      </div>
    </div>
  );
}
