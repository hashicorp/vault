/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import KvForm from 'vault/forms/secrets/kv';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.path = 'my-secret';
    this.form = new KvForm({
      path: this.path,
      custom_metadata: { foo: 'bar' },
      max_versions: 15,
      cas_required: true,
      delete_version_after: '4h30m',
    });
    this.backend = 'my-kv';
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path, route: 'secret.details', model: this.path },
      { label: 'Metadata' },
    ];
    this.capabilities = { canUpdateMetadata: true };
    this.onCancel = sinon.spy();
    this.onSave = sinon.spy();

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Metadata::Edit
            @form={{this.form}}
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
            @capabilities={{this.capabilities}}
            @onCancel={{this.onCancel}}
            @onSave={{this.onSave}}
          />
        `,
        { owner: this.engine }
      );
  });

  test('it renders all inputs for a model that has all default values', async function (assert) {
    assert.expect(5);

    this.form.data.custom_metadata = null;
    this.form.data.delete_version_after = null;

    await this.renderComponent();

    assert.dom(FORM.kvRow).exists({ count: 1 }, 'renders one kv row for custom metadata');
    assert.dom(FORM.inputByAttr('max_versions')).exists('renders Max versions.');
    assert.dom(FORM.inputByAttr('cas_required')).exists('renders Required Check and Set.');
    assert
      .dom('[data-test-toggle-label="Automate secret deletion"]')
      .exists('the label for automate secret deletion renders.');
    assert
      .dom(FORM.ttlValue('Automate secret deletion'))
      .doesNotExist('the toggle for secret deletion is not triggered.');
  });

  test('it displays previous inputs from metadata record and saves new values', async function (assert) {
    assert.expect(7);

    await this.renderComponent();

    assert.dom(FORM.keyInput()).hasValue('foo', 'renders custom metadata key');
    assert.dom(FORM.valueInput()).hasValue('bar', 'renders custom metadata value');
    assert
      .dom(FORM.inputByAttr('max_versions'))
      .hasValue('15', 'renders Max versions that was on the record.');
    assert
      .dom(FORM.inputByAttr('cas_required'))
      .hasValue('on', 'renders Required Check and Set that was on the record.');
    assert
      .dom(FORM.ttlValue('Automate secret deletion'))
      .hasValue('270', 'renders Automate secret deletion that was on the record.'); // 4h30m = 270m

    // update values
    await fillIn(FORM.keyInput(1), 'last');
    await fillIn(FORM.valueInput(1), 'value');
    await fillIn(FORM.inputByAttr('max_versions'), '8');
    await click(FORM.inputByAttr('cas_required'));
    await fillIn(FORM.ttlValue('Automate secret deletion'), '60'); // 60m = 3600s

    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'kvV2WriteMetadata').resolves();
    const metadata = {
      max_versions: '8',
      cas_required: false,
      delete_version_after: '3600s',
      custom_metadata: {
        foo: 'bar',
        last: 'value',
      },
    };
    await click(FORM.saveBtn);
    assert.true(this.writeStub.calledWith(this.path, this.backend, metadata), 'updated metadata is saved');
    assert.true(this.onSave.called, 'onSave action is called');
  });

  test('it displays validation errors and does not save inputs on cancel', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    // trigger validation error
    await fillIn(FORM.inputByAttr('max_versions'), 'a');
    await click(FORM.saveBtn);
    assert
      .dom(FORM.validationError('max_versions'))
      .hasText('Maximum versions must be a number.', 'Validation message is shown for max_versions');

    await click(FORM.cancelBtn);
    assert.true(this.onCancel.called, 'onCancel action is called');
  });

  test('it shows an empty state if user does not have metadata update permissions', async function (assert) {
    assert.expect(1);

    this.capabilities.canUpdateMetadata = false;
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('You do not have permissions to edit metadata');
  });
});
