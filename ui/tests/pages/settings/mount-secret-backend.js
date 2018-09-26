import { create, visitable, fillable } from 'ember-cli-page-object';
import mountForm from 'vault/tests/pages/components/mount-backend-form';
import withFlash from 'vault/tests/helpers/with-flash';

export default create({
  visit: visitable('/vault/settings/mount-secret-backend'),
  ...mountForm,
  version: fillable('[data-test-input="options.version"]'),
  maxTTLVal: fillable('[data-test-input="config.maxLeaseTtl"] [data-test-ttl-value]'),
  maxTTLUnit: fillable('[data-test-input="config.maxLeaseTtl"] [data-test-ttl-unit]'),
  defaultTTLVal: fillable('[data-test-input="config.defaultLeaseTtl"] [data-test-ttl-value]'),
  defaultTTLUnit: fillable('[data-test-input="config.defaultLeaseTtl"] [data-test-ttl-unit]'),
  enable: async function(type, path) {
    await this.visit();
    return withFlash(this.mount(type, path));
  },
});
