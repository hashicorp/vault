import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';

export default class PkiConfigModel extends Model {
  @service secretMountPath;
  @attr('string') formType;

  @attr('string') pemBundle;
  @attr importedIssuers;
  @attr importedKeys;

  get backend() {
    return this.secretMountPath.currentPath;
  }
}
