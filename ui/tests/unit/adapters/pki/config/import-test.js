import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/config/import', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
  });

  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/config/import');
    assert.ok(adapter);
  });

  test('it should make request to correct endpoint on create', async function (assert) {
    assert.expect(4);
    this.server.post(`${this.backend}/config/ca`, (url, { requestBody }) => {
      assert.ok(true, `request made to correct endpoint ${url}`);
      assert.deepEqual(
        JSON.parse(requestBody),
        { pem_bundle: 'abcdefg' },
        'has correct payload for config/ca'
      );
      return {
        data: {},
      };
    });
    this.server.post(`${this.backend}/intermediate/set-signed`, (url, { requestBody }) => {
      assert.ok(true, `request made to correct endpoint ${url}`);
      assert.deepEqual(
        JSON.parse(requestBody),
        { certificate: 'zyxwvut' },
        'has correct payload for intermediate'
      );
      return {
        data: {},
      };
    });

    await this.store
      .createRecord('pki/config/import', {
        pemBundle: 'abcdefg',
      })
      .save();

    await this.store
      .createRecord('pki/config/import', {
        certificate: 'zyxwvut',
      })
      .save();
  });
});
