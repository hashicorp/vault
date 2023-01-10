import Model from '@ember-data/model';

export default class PkiConfigModel extends Model {
  secretMountPath: unknown;
  pemBundle: string;
  type: string;
  get backend(): string;
  // apiPaths for capabilities
  configCaPath: string;
  generateRootPath: string;
  generateCsrPath: string;
  crossSignPath: string;
  // Capabilities
  get canConfigCa(): boolean;
  get canGenerateRoot(): boolean;
  get canGenerateIntermediate(): boolean;
  get canCrossSign(): boolean;
}
