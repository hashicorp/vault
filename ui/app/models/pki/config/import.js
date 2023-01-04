import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';

export default class PkiConfigImportModel extends Model {
  @service secretMountPath;

  @attr('string') pemBundle;
  @attr importedIssuers;
  @attr importedKeys;

  get backend() {
    return this.secretMountPath.currentPath;
  }
}
