import Model from '@ember-data/model';

export default class PkiConfigModel extends Model {
  secretMountPath: unknown;
  backend: string;
  type: string;
  // from config/import (TODO: cleanup)
  importFile: string;
  pemBundle: string;
  certificate: string;
}
