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

  @lazyCapabilities(apiPath`${'backend'}/issuers/import/bundle`, 'backend') importBundlePath;
  @lazyCapabilities(apiPath`${'backend'}/issuers/generate/root/${'type'}`, 'backend', 'type')
  generateIssuerRootPath;
  @lazyCapabilities(apiPath`${'backend'}/issuers/generate/intermediate/${'type'}`, 'backend', 'type')
  generateIssuerCsrPath;

  get canImportBundle() {
    return this.importBundlePath.get('canCreate') !== false;
  }
  get canGenerateIssuerRoot() {
    return this.generateIssuerRootPath.get('canCreate') !== false;
  }
  get canGenerateIssuerIntermediate() {
    return this.generateIssuerCsrPath.get('canCreate') !== false;
  }
  get canCrossSign() {
    return this.crossSignPath.get('canCreate') !== false;
  }
}
