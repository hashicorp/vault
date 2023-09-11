/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv | Page::Secret::Metadata::Edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.onCancel = sinon.spy();
    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const data = this.server.create('kv-metadatum', 'withCustomMetadata');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    store.pushPayload('kv/metadatum', {
      modelName: 'kv/metadata',
      ...data,
    });
    // Used to test a model with custom_metadata and non-default inputs.
    this.metadataModelEdit = store.peekRecord('kv/metadata', data.id);
    // Used to test a model with no custom_metadata and default values.
    this.metadataModelCreate = store.createRecord('kv/metadata', {
      backend: 'kv-engine-new',
      path: 'my-secret-new',
    });
  });

  test('it renders all inputs for a model that has all default values', async function (assert) {
    assert.expect(5);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadataModelCreate.backend, route: 'list' },
      { label: this.metadataModelCreate.path, route: 'secret.details', model: this.metadataModelCreate.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`
      <Page::Secret::Metadata::Edit
        @metadata={{this.metadataModelCreate}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(FORM.kvRow).exists({ count: 1 }, 'renders one kv row when model is new.');
    assert.dom(FORM.inputByAttr('maxVersions')).exists('renders Max versions.');
    assert.dom(FORM.inputByAttr('casRequired')).exists('renders Required Check and Set.');
    assert
      .dom('[data-test-toggle-label="Automate secret deletion"]')
      .exists('the label for automate secret deletion renders.');
    assert
      .dom(FORM.ttlValue('Automate secret deletion'))
      .doesNotExist('the toggle for secret deletion is not triggered.');
  });

  test('it displays previous inputs from metadata record and saves new values', async function (assert) {
    assert.expect(5);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadataModelEdit.backend, route: 'list' },
      { label: this.metadataModelEdit.path, route: 'secret.details', model: this.metadataModelEdit.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`
      <Page::Secret::Metadata::Edit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert
      .dom(FORM.kvRow)
      .exists({ count: 4 }, 'renders all kv rows including previous data and one extra to fill out.');
    assert
      .dom(FORM.inputByAttr('maxVersions'))
      .hasValue('15', 'renders Max versions that was on the record.');
    assert
      .dom(FORM.inputByAttr('casRequired'))
      .hasValue('on', 'renders Required Check and Set that was on the record.');
    assert
      .dom(FORM.ttlValue('Automate secret deletion'))
      .hasValue('12319', 'renders Automate secret deletion that was on the record.');

    // change the "Additional option" values
    await click(FORM.deleteRow()); // delete the first kv row
    await fillIn(FORM.keyInput(2), 'last');
    await fillIn(FORM.valueInput(2), 'value');
    await fillIn(FORM.inputByAttr('maxVersions'), '8');
    await click(FORM.inputByAttr('casRequired'));
    await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');
    // save test and check record
    this.server.post('/kv-engine/metadata/my-secret', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        max_versions: 8,
        cas_required: false,
        delete_version_after: '1000s',
        custom_metadata: {
          baz: '5c07d823-3810-48f6-a147-4c06b5219e84',
          foo: 'abc',
          last: 'value',
        },
      };
      assert.propEqual(data, expected, 'POST request made to save metadata with correct properties.');
    });
    await click(FORM.saveBtn);
  });

  test('it displays validation errors and does not save inputs on cancel', async function (assert) {
    assert.expect(2);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadataModelEdit.backend, route: 'list' },
      { label: this.metadataModelEdit.path, route: 'secret.details', model: this.metadataModelEdit.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`
      <Page::Secret::Metadata::Edit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    // trigger validation error
    await fillIn(FORM.inputByAttr('maxVersions'), 'a');
    await click(FORM.saveBtn);
    assert
      .dom(FORM.inlineAlert)
      .hasText('Maximum versions must be a number.', 'Validation message is shown for max_versions');

    await click(FORM.cancelBtn);
    assert.strictEqual(this.metadataModelEdit.maxVersions, 15, 'Model is rolled back on cancel.');
  });

  test('it shows an empty state if user does not have metadata update permissions', async function (assert) {
    assert.expect(1);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub('list'));
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadataModelEdit.backend, route: 'list' },
      { label: this.metadataModelEdit.path, route: 'secret.details', model: this.metadataModelEdit.path },
      { label: 'metadata' },
    ];
    await render(
      hbs`
      <Page::Secret::Metadata::Edit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(PAGE.emptyStateTitle).hasText('You do not have permissions to edit metadata');
  });
});
