import { create, visitable, text } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/secrets/:backend/configuration'),
  defaultTTL: text('[data-test-row-value="Default Lease TTL"]'),
  maxTTL: text('[data-test-row-value="Max Lease TTL"]'),
});
