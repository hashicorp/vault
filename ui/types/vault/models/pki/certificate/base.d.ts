import Model from '@ember-data/model';
export default class PkiCertificateBaseModel extends Model {
  secretMountPath: class;
  get useOpenAPI(): boolean;
  get backend(): string;
  getHelpUrl(): void;
  altNames: string;
  commonName: string;
  caChain: string;
  certificate: string;
  excludeCnFromSans: boolean;
  expiration: number;
  ipSans: string;
  issuingCa: string;
  notValidAfter: date;
  notValidBefore: date;
  otherSans: string;
  privateKey: string;
  privateKeyType: string;
  revokePath: string;
  revocationTime: number;
  serialNumber: string;
  get canRevoke(): boolean;
}
