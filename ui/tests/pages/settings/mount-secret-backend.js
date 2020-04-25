import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import mountForm from 'vault/tests/pages/components/mount-backend-form';

export default create({
  visit: visitable('/vault/settings/mount-secret-backend'),
  ...mountForm,
  version: fillable('[data-test-input="options.version"]'),
  enableMaxTtl: clickable('[data-test-toggle-input="Max Lease TTL"]'),
  maxTTLVal: fillable('[data-test-ttl-value="Max Lease TTL"]'),
  maxTTLUnit: fillable('[data-test-ttl-unit="Max Lease TTL"] [data-test-select="ttl-unit"]'),
  enableDefaultTtl: clickable('[data-test-toggle-input="Default Lease TTL"]'),
  defaultTTLVal: fillable('input[data-test-ttl-value="Default Lease TTL"]'),
  defaultTTLUnit: fillable('[data-test-ttl-unit="Default Lease TTL"] [data-test-select="ttl-unit"]'),
  enable: async function(type, path) {
    await this.visit();
    await this.mount(type, path);
  },
});
