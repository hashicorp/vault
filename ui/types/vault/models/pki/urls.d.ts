import Model from '@ember-data/model';

export default class PkiUrlsModel extends Model {
  get useOpenAPI(): boolean;
  getHelpUrl(backendPath: string): string;
  issuingCertificates: array;
  crlDistributionPoints: array;
  ocspServers: array;
  urlsPath: string;
  get canSet(): boolean;
}
