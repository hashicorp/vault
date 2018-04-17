import { create, visitable, fillable, clickable } from 'ember-cli-page-object';

export default create({
  visit: visitable('/vault/settings/mount-secret-backend'),
  type: fillable('[data-test-secret-backend-type]'),
  path: fillable('[data-test-secret-backend-path]'),
  submit: clickable('[data-test-secret-backend-submit]'),
  toggleOptions: clickable('[data-test-secret-backend-options]'),
  version: fillable('[data-test-secret-backend-version]'),
  maxTTLVal: fillable('[data-test-ttl-value]', { scope: '[data-test-secret-backend-max-ttl]' }),
  maxTTLUnit: fillable('[data-test-ttl-unit]', { scope: '[data-test-secret-backend-max-ttl]' }),
  defaultTTLVal: fillable('[data-test-ttl-value]', { scope: '[data-test-secret-backend-default-ttl]' }),
  defaultTTLUnit: fillable('[data-test-ttl-unit]', { scope: '[data-test-secret-backend-default-ttl]' }),
});
