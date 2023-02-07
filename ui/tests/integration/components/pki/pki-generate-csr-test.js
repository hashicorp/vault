/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | PkiGenerateCsr', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.owner.lookup('service:secretMountPath').update('pki-test');
    this.model = this.owner
      .lookup('service:store')
      .createRecord('pki/action', { actionType: 'generate-csr' });

    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
        'pki-test/issuers/generate/intermediate/exported': ['root'],
      },
    }));
  });

  test('it should render fields and save', async function (assert) {
    assert.expect(9);

    this.server.post('/pki-test/issuers/generate/intermediate/exported', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.common_name, 'foo', 'Request made to correct endpoint on save');
    });

    this.onSave = () => assert.ok(true, 'onSave action fires');

    await render(hbs`<PkiGenerateCsr @model={{this.model}} @onSave={{this.onSave}} />`, {
      owner: this.engine,
    });

    const fields = [
      'type',
      'commonName',
      'excludeCnFromSans',
      'format',
      'serialNumber',
      'addBasicConstraints',
    ];
    fields.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} form field renders`);
    });

    assert.dom('[data-test-toggle-group]').exists({ count: 3 }, 'Toggle groups render');

    await fillIn('[data-test-input="type"]', 'exported');
    await fillIn('[data-test-input="commonName"]', 'foo');
    await click('[data-test-save]');
  });

  test('it should display validation errors', async function (assert) {
    assert.expect(4);

    this.onCancel = () => assert.ok(true, 'onCancel action fires');

    await render(hbs`<PkiGenerateCsr @model={{this.model}} @onCancel={{this.onCancel}} />`, {
      owner: this.engine,
    });

    await click('[data-test-save]');

    assert
      .dom('[data-test-field-validation="type"]')
      .hasText('Type is required.', 'Type validation error renders');
    assert
      .dom('[data-test-field="commonName"] [data-test-inline-alert]')
      .hasText('Common name is required.', 'Common name validation error renders');
    assert.dom('[data-test-alert]').hasText('There are 2 errors with this form.', 'Alert renders');

    await click('[data-test-cancel]');
  });
});
