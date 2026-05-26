import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type SchoolRegisterFieldsProps = {
  schoolName: string;
  npsn: string;
  schoolAddress: string;
  schoolPIC: string;
  setSchoolName: (value: string) => void;
  setNpsn: (value: string) => void;
  setSchoolAddress: (value: string) => void;
  setSchoolPIC: (value: string) => void;
};

export function SchoolRegisterFields({
  schoolName,
  npsn,
  schoolAddress,
  schoolPIC,
  setSchoolName,
  setNpsn,
  setSchoolAddress,
  setSchoolPIC,
}: SchoolRegisterFieldsProps) {
  return (
    <div className="rounded-3xl border border-emerald-100 bg-emerald-50/80 p-5">
      <div className="rounded-2xl bg-emerald-600 px-4 py-3 text-center text-sm font-semibold text-white shadow-sm">
        Formulir Data Sekolah
      </div>

      <p className="mt-3 rounded-2xl bg-white px-4 py-3 text-xs leading-5 text-slate-500">
        TODO: data tambahan sekolah belum dikirim ke backend karena endpoint
        register saat ini hanya menerima data akun dasar.
      </p>

      <div className="mt-5 space-y-4 rounded-2xl bg-white p-4">
        <div className="space-y-2 h-11 rounded-2xl">
          <Label>
            Nama Sekolah <span className="text-red-500">*</span>
          </Label>
          <Input
            placeholder="Contoh: SD Negeri 01 Jakarta"
            value={schoolName}
            onChange={(event) => setSchoolName(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>
            NPSN <span className="text-red-500">*</span>
          </Label>
          <Input
            placeholder="Nomor Pokok Sekolah Nasional"
            value={npsn}
            onChange={(event) => setNpsn(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>Alamat Sekolah</Label>
          <Input
            placeholder="Alamat lengkap sekolah"
            value={schoolAddress}
            onChange={(event) => setSchoolAddress(event.target.value)}
          />
        </div>

        <div className="space-y-2 h-11 rounded-2xl">
          <Label>Penanggung Jawab</Label>
          <Input
            placeholder="Nama guru/staf penanggung jawab"
            value={schoolPIC}
            onChange={(event) => setSchoolPIC(event.target.value)}
          />
        </div>
      </div>
    </div>
  );
}
