/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';
import TransitKeyForm from 'vault/forms/transit/key';

const SELECTORS = {
  createForm: '[data-test-transit-create-form]',
  editForm: '[data-test-transit-edit-form]',
  ttlToggle: '[data-test-ttl-toggle="Auto-rotation period"]',
  ttlValue: '[data-test-ttl-value="Auto-rotation period"]',
};
module('Integration | Component | transit-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub('transit-backend/keys/some-key', ['sudo'])
    );
    this.model = { backend: 'transit', id: 'some-key' };
    this.backendCrumb = {
      label: 'transit',
      text: 'transit',
      path: 'vault.cluster.secrets.backend.list-root',
      model: 'transit',
    };
    this.form = new TransitKeyForm(this.model, { isNew: false });
  });

  test('it renders in create mode and updates model', async function (assert) {
    await render(hbs`
    <TransitEdit
      @form={{this.form}}
      @mode="create"
      @root={{this.backendCrumb}}
      @preferAdvancedEdit={{false}}
    />`);

    assert.dom(SELECTORS.createForm).exists();
    assert.dom(SELECTORS.ttlToggle).isNotChecked();

    // confirm model params update when ttl changes
    assert.strictEqual(this.form.data.auto_rotate_period, undefined);
    await click(SELECTORS.ttlToggle);

    assert.dom(SELECTORS.ttlValue).hasValue('30'); // 30 days
    assert.strictEqual(this.form.data.auto_rotate_period, '720h');

    await fillIn(SELECTORS.ttlValue, '10'); // 10 days
    assert.strictEqual(this.form.data.auto_rotate_period, '240h');
  });

  test('it renders edit form correctly when key has autoRotatePeriod=0', async function (assert) {
    this.form.data.auto_rotate_period = 0;
    this.form.data.keys = {
      1: 1684882652000,
    };
    await render(hbs`
      <TransitEdit
        @form={{this.form}}
        @mode="edit"
        @root={{this.backendCrumb}}
        @preferAdvancedEdit={{false}}
      />`);
    assert.dom(SELECTORS.editForm).exists();
    assert.dom(SELECTORS.ttlToggle).isNotChecked();

    assert.strictEqual(this.form.data.auto_rotate_period, 0);

    await click(SELECTORS.ttlToggle);
    assert.dom(SELECTORS.ttlToggle).isChecked();
    assert.dom(SELECTORS.ttlValue).hasValue('30');
    assert.strictEqual(this.form.data.auto_rotate_period, '720h', 'model value changes with toggle');

    await click(SELECTORS.ttlToggle);
    assert.strictEqual(this.form.data.auto_rotate_period, '0s'); // reverts to original value when toggled back on
  });

  test('it renders edit form correctly when key has non-zero rotation period', async function (assert) {
    this.form.data.auto_rotate_period = '5h';
    this.form.data.keys = {
      1: 1684882652000,
    };
    await render(hbs`
      <TransitEdit
        @form={{this.form}}
        @mode="edit"
        @root={{this.backendCrumb}}
        @preferAdvancedEdit={{false}}
      />`);

    assert.dom(SELECTORS.editForm).exists();
    assert.dom(SELECTORS.ttlToggle).isChecked();

    await click(SELECTORS.ttlToggle);
    assert.dom(SELECTORS.ttlToggle).isNotChecked();
    assert.strictEqual(
      this.form.data.auto_rotate_period,
      '0s',
      'model value changes back to 0 when toggled off'
    );

    await click(SELECTORS.ttlToggle);
    assert.strictEqual(
      this.form.data.auto_rotate_period,
      '5h',
      'model value changes to original value when toggled back on'
    );

    await fillIn(SELECTORS.ttlValue, '10');
    assert.strictEqual(this.form.data.auto_rotate_period, '10h');
  });
});
