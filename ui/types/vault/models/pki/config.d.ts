import Model from '@ember-data/model';

export default class PkiConfigModel extends Model {
  secretMountPath: unknown;
  pemBundle: string;
  type: string;
  get backend(): string;
  // apiPaths for capabilities
  importBundlePath: string;
  generateIssuerRootPath: string;
  generateIssuerCsrPath: string;
  crossSignPath: string;
  // Capabilities
  get canImportBundle(): boolean;
  get canGenerateIssuerRoot(): boolean;
  get canGenerateIssuerIntermediate(): boolean;
  get canCrossSign(): boolean;
}
