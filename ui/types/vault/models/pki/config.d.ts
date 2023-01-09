import Model from '@ember-data/model';

export default class PkiConfigModel extends Model {
  secretMountPath: unknown;
  backend: string;
  formType: string;
  // import fields:
  pemBundle: string;
}
