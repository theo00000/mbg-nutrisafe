import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type MitraRegisterFieldsProps = {
  institutionName: string;
  identityNumber: string;
  sppgName: string;
  sppgAddress: string;
  productionCapacity: string;
  setInstitutionName: (value: string) => void;
  setIdentityNumber: (value: string) => void;
  setSppgName: (value: string) => void;
  setSppgAddress: (value: string) => void;
  setProductionCapacity: (value: string) => void;
};

export function MitraRegisterFields({
  institutionName,
  identityNumber,
  sppgName,
  sppgAddress,
  productionCapacity,
  setInstitutionName,
  setIdentityNumber,
  setSppgName,
  setSppgAddress,
  setProductionCapacity,
}: MitraRegisterFieldsProps) {
  return (
    <div className="rounded-3xl border border-blue-100 bg-blue-50/80 p-5">
      <div className="rounded-2xl bg-blue-600 px-4 py-3 text-center text-sm font-semibold text-white shadow-sm">
        Formulir Pendaftaran Mitra SPPG Baru
      </div>

      <p className="mt-3 rounded-2xl bg-white px-4 py-3 text-xs leading-5 text-slate-500">
        TODO: data tambahan Mitra/SPPG dan upload dokumen belum dikirim ke
        backend karena endpoint register saat ini hanya menerima data akun
        dasar.
      </p>

      <div className="mt-5 space-y-6">
        <section>
          <h3 className="text-sm font-bold text-slate-950">Data Mitra</h3>

          <div className="mt-3 space-y-3 rounded-2xl bg-white p-4">
            <div className="space-y-2">
              <Label>
                Nama Instansi <span className="text-red-500">*</span>
              </Label>
              <Input
                placeholder="Contoh: Yayasan Pangan Sehat"
                value={institutionName}
                onChange={(event) => setInstitutionName(event.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>
                NIK / NPWP <span className="text-red-500">*</span>
              </Label>
              <Input
                placeholder="Masukkan NIK atau NPWP"
                value={identityNumber}
                onChange={(event) => setIdentityNumber(event.target.value)}
              />
            </div>
          </div>
        </section>

        <section>
          <h3 className="text-sm font-bold text-slate-950">Data SPPG</h3>

          <div className="mt-3 space-y-3 rounded-2xl bg-white p-4">
            <div className="space-y-2">
              <Label>
                Nama SPPG <span className="text-red-500">*</span>
              </Label>
              <Input
                placeholder="Contoh: SPPG Sehat Mandiri"
                value={sppgName}
                onChange={(event) => setSppgName(event.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Alamat SPPG</Label>
              <Input
                placeholder="Alamat lengkap dapur/SPPG"
                value={sppgAddress}
                onChange={(event) => setSppgAddress(event.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Kapasitas Produksi</Label>
              <Input
                placeholder="Contoh: 1000 porsi/hari"
                value={productionCapacity}
                onChange={(event) => setProductionCapacity(event.target.value)}
              />
            </div>
          </div>
        </section>

        <section>
          <h3 className="text-sm font-bold text-slate-950">Upload Dokumen</h3>

          <div className="mt-3 space-y-3 rounded-2xl bg-white p-4">
            <div className="space-y-2">
              <Label>Upload Proposal PDF</Label>
              <Input type="file" accept=".pdf" />
            </div>

            <div className="space-y-2">
              <Label>Upload Foto Dapur</Label>
              <Input type="file" accept="image/*" />
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
