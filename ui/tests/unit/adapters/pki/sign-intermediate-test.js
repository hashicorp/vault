import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/sign-intermediate', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.payload = {
      issuerRef: 'my-issuer-id',
    };
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/sign-intermediate');
    assert.ok(adapter);
  });

  test('it calls the correct endpoint on save', async function (assert) {
    assert.expect(2);

    this.server.post(`${this.backend}/issuer/my-issuer-id/sign-intermediate`, () => {
      assert.ok(true, 'correct endpoint called');
      return {
        request_id: 'unique-request-id',
        data: {
          foo: 'bar',
        },
      };
    });

    const result = await this.store.createRecord('pki/sign-intermediate', this.payload).save();
    assert.strictEqual(result.id, 'unique-request-id', 'Resulting model has ID matching request ID');
  });
});
