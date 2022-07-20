import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/assignment', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/assignment';
    this.data = {
      name: 'foo-assignment',
      // ARG TODO when Jordan/Claire are back there are some issues here with the names of the params and how we have the serializer setup. Instead of redoing the serializer in this PR, I'm going to comment this out and revisit to push through this pr.
      // entity_ids: ['my-entity'],
      // group_ids: ['my-group'],
    };
    this.path = '/identity/oidc/assignment/foo-assignment';
  });

  testHelper(test);
});
