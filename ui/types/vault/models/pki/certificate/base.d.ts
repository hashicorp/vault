import Model from '@ember-data/model';
export default class PkiCertificateBaseModel extends Model {
  secretMountPath: class;
  get useOpenAPI(): boolean;
  get backend(): string;
  getHelpUrl(): void;
  commonName: string;
  caChain: string;
  certificate: string;
  expiration: number;
  issuingCa: string;
  privateKey: string;
  privateKeyType: string;
  serialNumber: string;
  notValidAfter: date;
  notValidBefore: date;
  pemBundle: string;
  importedIssuers: string[];
  importedKeys: string[];
  revokePath: string;
  get canRevoke(): boolean;
}
