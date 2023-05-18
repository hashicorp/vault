/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-tidy-form';

module('Integration | Component | pki tidy form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.server.post('/sys/capabilities-self', () => {});

    this.manualTidy = this.store.createRecord('pki/tidy', { backend: 'pki-manual-tidy' });
    this.store.pushPayload('pki/tidy', {
      modelName: 'pki/tidy',
      id: 'pki-auto-tidy',
    });
    this.autoTidy = this.store.peekRecord('pki/tidy', 'pki-auto-tidy');
  });

  test('it should render tidy fields', async function (assert) {
    this.version.version = '1.14.1+ent';
    await render(hbs`<PkiTidyForm @tidy={{this.tidy}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom(SELECTORS.tidyCertStoreLabel).hasText('Tidy the certificate store');
    assert.dom(SELECTORS.tidyRevocationList).hasText('Tidy the revocation list (CRL)');
    assert.dom(SELECTORS.safetyBufferTTL).exists();
    assert.dom(SELECTORS.safetyBufferInput).hasValue('3');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('d');
  });

  test('it should change the attributes on the model', async function (assert) {
    await render(hbs`<PkiTidyForm @tidy={{this.tidy}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    await click(SELECTORS.tidyCertStoreCheckbox);
    await click(SELECTORS.tidyRevocationCheckbox);
    await fillIn(SELECTORS.safetyBufferInput, '5');
    assert.true(this.tidy.tidyCertStore);
    assert.true(this.tidy.tidyRevocationQueue);
    assert.dom(SELECTORS.safetyBufferInput).hasValue('5');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('d');
    assert.strictEqual(this.tidy.safetyBuffer, '120h');
  });

  test('it updates auto-tidy config', async function (assert) {
    assert.expect(4);
    this.server.post('/pki-auto-tidy/config/auto-tidy', (schema, req) => {
      assert.ok(true, 'Request made to update auto-tidy');
      assert.propEqual(JSON.parse(req.requestBody), { enabled: false }, 'response contains auto-tidy params');
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
      <PkiTidyForm
        @tidy={{this.autoTidy}}
        @tidyType="auto"
        @onSave={{this.onSave}}
        @onCancel={{this.onCancel}}
      />
    `,
      { owner: this.engine }
    );

    await click(SELECTORS.tidySave);
    await click(SELECTORS.tidyCancel);
  });

  test('it saves and performs manual tidy', async function (assert) {
    assert.expect(4);

    this.server.post('/pki-manual-tidy/tidy', (schema, req) => {
      assert.ok(true, 'Request made to perform manual tidy');
      assert.propEqual(JSON.parse(req.requestBody), {}, 'response contains manual tidy params');
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
      <PkiTidyForm
        @tidy={{this.manualTidy}}
        @tidyType="manual"
        @onSave={{this.onSave}}
        @onCancel={{this.onCancel}}
      />
    `,
      { owner: this.engine }
    );

    await click(SELECTORS.tidySave);
    await click(SELECTORS.tidyCancel);
  });
});
