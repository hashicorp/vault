import Model from '@ember-data/model';
import { FormField, ModelValidations } from 'vault/app-types';

export default class PkiActionModel extends Model {
  secretMountPath: unknown;
  pemBundle: string;
  type: string;
  actionType: string | null;
  get backend(): string;
  // apiPaths for capabilities
  importBundlePath: string;
  generateIssuerRootPath: string;
  generateIssuerCsrPath: string;
  crossSignPath: string;
  allFields: Array<FormField>;
  validate(): ModelValidations;
  // Capabilities
  get canImportBundle(): boolean;
  get canGenerateIssuerRoot(): boolean;
  get canGenerateIssuerIntermediate(): boolean;
  get canCrossSign(): boolean;
}
