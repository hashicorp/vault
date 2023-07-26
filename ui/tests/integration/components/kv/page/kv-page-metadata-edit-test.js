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
import { kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';
import { SELECTORS as PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';

module('Integration | Component | kv | Page::Secret::MetadataEdit', function (hooks) {
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
    // Used to specifically test a model with custom_metadata and non-default inputs.
    this.metadataModelEdit = store.peekRecord('kv/metadata', data.id);
    // Used to specifically test a model with no custom_metadata and default values.
    this.metadataModelCreate = store.createRecord('kv/metadata', {
      backend: 'kv-engine',
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
      <Page::Secret::MetadataEdit
        @metadata={{this.metadataModelCreate}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(PAGE.edit.kvRow).exists({ count: 1 }, 'renders one kv row when model is new.');
    assert.dom(PAGE.edit.inputByAttr('maxVersions')).exists('renders Max versions.');
    assert.dom(PAGE.edit.inputByAttr('casRequired')).exists('renders Required Check and Set.');
    assert
      .dom('[data-test-toggle-label="Automate secret deletion"]')
      .exists('the label for automate secret deletion renders.');
    assert
      .dom(PAGE.edit.automateSecretDeletion)
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
      <Page::Secret::MetadataEdit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert
      .dom(PAGE.edit.kvRow)
      .exists({ count: 4 }, 'renders all kv rows including previous data and one extra to fill out.');
    assert
      .dom(PAGE.edit.inputByAttr('maxVersions'))
      .hasValue('15', 'renders Max versions that was on the record.');
    assert
      .dom(PAGE.edit.inputByAttr('casRequired'))
      .hasValue('on', 'renders Required Check and Set that was on the record.');
    assert
      .dom(PAGE.edit.automateSecretDeletion)
      .hasValue('12319', 'renders Automate secret deletion that was on the record.');

    // change the "Additional option" values
    await click('[data-test-kv-delete-row="0"]'); // delete the first kv row
    const keys = document.querySelectorAll('[data-test-kv-key]');
    const values = document.querySelectorAll('[data-test-kv-value]');
    await fillIn(keys[2], 'last');
    await fillIn(values[2], 'value');
    await fillIn(PAGE.edit.inputByAttr('maxVersions'), '8');
    await click(PAGE.edit.inputByAttr('casRequired'));
    await fillIn(PAGE.edit.automateSecretDeletion, '1000');
    // save test and check record
    this.server.post('/kv-engine/metadata/my-secret', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        backend: 'kv-engine',
        path: 'my-secret',
        max_versions: 8,
        cas_required: false,
        delete_version_after: '1000s',
        custom_metadata: {
          baz: '5c07d823-3810-48f6-a147-4c06b5219e84',
          foo: 'abc',
          last: 'value',
        },
      };
      assert.deepEqual(expected, data, 'POST request made to save metadata with correct properties.');
    });
    await click('[data-test-kv-metadata-save]');
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
      <Page::Secret::MetadataEdit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    // trigger validation error
    await fillIn(PAGE.edit.inputByAttr('maxVersions'), 'a');
    await click('[data-test-kv-metadata-save]');
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('Maximum versions must be a number.', 'Validation message is shown for max_versions');

    await click(PAGE.edit.metadataCancel);
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
      <Page::Secret::MetadataEdit
        @metadata={{this.metadataModelEdit}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(SELECTORS.emptyStateTitle).hasText('You do not have permissions to edit metadata');
  });
});
