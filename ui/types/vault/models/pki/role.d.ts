import Model from '@ember-data/model';
import { FModelValidations } from 'vault/app-types';

export default class PkiRoleModel extends Model {
  get useOpenAPI(): boolean;
  name: string;
  issuerRef: string;
  getHelpUrl(backendPath: string): string;
  validate(): ModelValidations;
  isNew: boolean;
}
