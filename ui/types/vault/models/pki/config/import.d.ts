import Model from '@ember-data/model';

export default class PkiConfigImportModel extends Model {
  backend: string;
  secretMountPath: unknown;
  importFile: string;
  pemBundle: string;
  certificate: string;
}
