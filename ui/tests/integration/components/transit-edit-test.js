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
import TransitKeyResource from 'vault/resources/transit/key';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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

  test('edit submits only configuration with numeric versions and displays API errors', async function (assert) {
    const api = this.owner.lookup('service:api');
    const configure = sinon.stub(api.secrets, 'transitConfigureKey').rejects(new Error('Update denied'));
    sinon.stub(api, 'parseError').resolves({ message: 'Update denied' });
    const transition = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.form = new TransitKeyForm({
      backend: 'transit',
      id: 'key',
      name: 'key',
      keys: { 1: 123, 2: 456, 3: 789 },
      latest_version: 3,
      min_decryption_version: 1,
      min_encryption_version: 0,
      auto_rotate_period: 0,
      deletion_allowed: false,
    });
    await render(hbs`<TransitFormEdit @form={{this.form}} @capabilities={{hash canUpdate=true}} />`);
    await fillIn('#min_decryption_version', '2');
    assert
      .dom('#min_encryption_version option')
      .exists({ count: 3 }, 'Latest plus the two allowed encryption versions are available');
    await fillIn('#min_encryption_version', '3');
    await click(GENERAL.submitButton);
    assert.deepEqual(
      configure.firstCall.args,
      [
        'key',
        'transit',
        {
          auto_rotate_period: '0',
          min_decryption_version: 2,
          min_encryption_version: 3,
          deletion_allowed: false,
        },
      ],
      'configuration excludes identity and read-only key metadata from the body'
    );
    assert.dom(GENERAL.messageError).containsText('Update denied', 'a failed update is visible');
    assert.false(transition.called, 'failure leaves the user on the editable form');
    configure.resolves();
    await click(GENERAL.submitButton);
    assert.true(
      transition.calledWith('vault.cluster.secrets.backend.show', 'transit', 'key', {
        queryParams: { tab: 'details' },
      }),
      'a successful update returns to details'
    );
  });

  test('failed deletion keeps the form and displays the API error', async function (assert) {
    const api = this.owner.lookup('service:api');
    const deletion = sinon.stub(api.secrets, 'transitDeleteKey').rejects(new Error('Delete denied'));
    sinon.stub(api, 'parseError').resolves({ message: 'Delete denied' });
    const transition = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    await render(hbs`<TransitFormEdit @form={{this.form}} @capabilities={{hash canDelete=true}} />`);
    await click('[data-test-confirm-action-trigger]');
    await click(GENERAL.confirmButton);
    assert.deepEqual(deletion.firstCall.args, ['some-key', 'transit'], 'deletion uses the key identity');
    assert.dom(GENERAL.messageError).containsText('Delete denied', 'a failed deletion is visible');
    assert.false(transition.called, 'failed deletion does not leave the form');
  });

  test('show reads a Resource and reports rotation failures', async function (assert) {
    const api = this.owner.lookup('service:api');
    const rotate = sinon.stub(api.secrets, 'transitRotateKey').rejects(new Error('Rotate denied'));
    sinon.stub(api, 'parseError').resolves({ message: 'Rotate denied' });
    const danger = sinon.spy(this.owner.lookup('service:flash-messages'), 'danger');
    this.key = new TransitKeyResource({
      backend: 'transit',
      id: 'some-key',
      type: 'aes256-gcm96',
      keys: { 1: 1633024800 },
      min_decryption_version: 1,
    });
    await render(hbs`<TransitEdit @key={{this.key}} @mode="show" @tab="versions"
      @capabilities={{hash canRotate=true}} />`);
    assert
      .dom('[data-test-transit-version="1"]')
      .containsText('Version 1', 'Resource versions are displayed');
    await click('[data-test-transit-key-rotate]');
    await click(GENERAL.confirmButton);
    assert.deepEqual(rotate.firstCall.args, ['some-key', 'transit', {}], 'rotation uses Resource identity');
    assert.true(danger.calledWith('Rotate denied'), 'the parsed rotation error is flashed');
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
