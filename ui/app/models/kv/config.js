import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { duration } from 'core/helpers/format-duration';

// This model is used only for display only - configuration happens via secret-engine model when an engine is mounted
@withFormFields(['casRequired', 'deleteVersionAfter', 'maxVersions'])
export default class KvConfigModel extends Model {
  @attr backend;
  @attr('number', { label: 'Maximum number of versions' }) maxVersions;

  @attr('boolean', { label: 'Require check and set' }) casRequired;

  @attr({ label: 'Automate secret deletion' }) deleteVersionAfter;

  @lazyCapabilities(apiPath`${'backend'}/config`, 'backend') configPath;

  get canRead() {
    return this.configPath.get('canRead') !== false;
  }

  // used in template to render using this model instead of secret-engine (where these attrs also exist)
  get displayFields() {
    return ['casRequired', 'deleteVersionAfter', 'maxVersions'];
  }

  get displayDeleteTtl() {
    if (this.deleteVersionAfter === '0s') return 'Never delete';
    return duration([this.deleteVersionAfter]);
  }
}
