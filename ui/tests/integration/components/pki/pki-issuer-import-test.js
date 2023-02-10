import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { issuerPemBundle } from 'vault/tests/helpers/pki/values';

module('Integration | Component | pki issuer import', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/issuer');
    this.backend = 'pki-test';
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = this.backend;
    this.pemBundle = issuerPemBundle;
  });

  test('it renders import and updates model', async function (assert) {
    assert.expect(3);
    await render(
      hbs`
      <PkiCaCertificateImport
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
      `,
      { owner: this.engine }
    );

    assert.dom('[data-test-pki-ca-cert-import-form]').exists('renders form');
    assert.dom('[data-test-component="text-file"]').exists('renders text file input');
    await click('[data-test-text-toggle]');
    await typeIn('[data-test-text-file-textarea]', this.pemBundle);
    assert.strictEqual(this.model.pemBundle, this.pemBundle);
  });

  test('it sends correct payload to import endpoint', async function (assert) {
    assert.expect(3);
    this.server.post(`/${this.backend}/issuers/import/bundle`, (schema, req) => {
      assert.ok(true, 'Request made to the correct endpoint to import issuer');
      const request = JSON.parse(req.requestBody);
      assert.propEqual(
        request,
        {
          pem_bundle: `${this.pemBundle}`,
        },
        'sends params in correct type'
      );
      return {};
    });

    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(
      hbs`
      <PkiCaCertificateImport
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
         @adapterOptions={{hash import=true}}
       />
      `,
      { owner: this.engine }
    );

    await click('[data-test-text-toggle]');
    await typeIn('[data-test-text-file-textarea]', this.pemBundle);
    assert.strictEqual(this.model.pemBundle, this.pemBundle);
    await click('[data-test-pki-ca-cert-import]');
  });

  test('it should unload record on cancel', async function (assert) {
    assert.expect(2);
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');
    await render(
      hbs`
        <PkiCaCertificateImport
          @model={{this.model}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `,
      { owner: this.engine }
    );

    await click('[data-test-pki-ca-cert-cancel]');
    assert.true(this.model.isDestroyed, 'new model is unloaded on cancel');
  });
});
