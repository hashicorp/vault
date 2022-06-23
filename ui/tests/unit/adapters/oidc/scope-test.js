import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/key', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/scope';
    this.data = {
      name: 'foo-scope',
      template: '{ "groups": {{identity.entity.groups.names}} }',
      description: 'A simple scope example.',
    };
    this.path = '/identity/oidc/scope/foo-scope';
  });

  testHelper(test);
});
