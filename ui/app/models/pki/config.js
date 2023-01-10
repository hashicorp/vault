import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default class PkiConfigModel extends Model {
  @service secretMountPath;

  @attr('string') pemBundle;
  @attr('string') type;

  get backend() {
    return this.secretMountPath.currentPath;
  }

  @lazyCapabilities(apiPath`${'backend'}/config/ca`, 'backend') configCaPath;
  @lazyCapabilities(apiPath`${'backend'}/root/generate/${'type'}`, 'backend', 'type') generateRootPath;
  @lazyCapabilities(apiPath`${'backend'}/intermediate/generate/${'type'}`, 'backend', 'type') generateCsrPath;
  @lazyCapabilities(apiPath`${'backend'}/intermediate/cross-sign`, 'backend') crossSignPath;

  get canConfigCa() {
    return this.configCaPath.get('canCreate') !== false;
  }
  get canGenerateRoot() {
    return this.generateRootPath.get('canCreate') !== false;
  }
  get canGenerateIntermediate() {
    return this.generateCsrPath.get('canCreate') !== false;
  }
  get canCrossSign() {
    return this.crossSignPath.get('canCreate') !== false;
  }
}
