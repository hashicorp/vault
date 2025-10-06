/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, find, render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

import { setupMirage } from 'ember-cli-mirage/test-support';
import recoveryHandler from 'vault/mirage/handlers/recovery';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Integration | Component | recovery/snapshots/snapshot-manage', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    recoveryHandler(this.server);

    const snapshot = this.server.create('snapshot');

    const api = this.owner.lookup('service:api');
    let namespaces = [];
    try {
      const response = await api.sys.internalUiListNamespaces();
      namespaces = response.keys ?? [];
    } catch {
      namespaces = [];
    }

    this.model = {
      snapshot,
      namespaces,
    };

    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';

    this.renderComponent = () =>
      render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);
  });
  const scenarios = [
    {
      engine: 'Cubbyhole',
      namespace: 'root',
      mount: 'cubbyhole',
      path: 'nested/my-path',
    },
    {
      engine: 'KV',
      namespace: 'root',
      mount: 'kv',
      path: 'my-path',
    },
    {
      engine: 'Cubbyhole',
      namespace: 'root',
      mount: 'cubbyhole',
      path: 'nested/my-path',
    },
    {
      engine: 'Database',
      namespace: 'root',
      mount: 'database',
      path: 'test-static-path',
    },
  ];
  scenarios.forEach((scenario) => {
    test(`it performs a read operation in the root namespace context for ${scenario.namespace} namespace item - Engine: ${scenario.engine}`, async function (assert) {
      await this.renderComponent();
      await click(GENERAL.selectByAttr('namespace'));
      await click(`[data-test-option="${scenario.namespace}"]`);
      await click(GENERAL.selectByAttr('mount'));
      await click(`[data-test-option="${scenario.mount}"]`);
      await fillIn(GENERAL.inputByAttr('resourcePath'), scenario.path);

      await click(GENERAL.button('read'));
      await waitFor('[data-test-read-secrets]');

      // Open modal
      assert.dom('[data-test-read-secrets]').exists('renders read modal');
      assert.dom(GENERAL.maskedInput).exists('renders secret data');

      // Close modal
      await click(GENERAL.button('close'));
      assert.dom('[data-test-read-secrets]').doesNotExist('read modal closed');
    });

    test(`it performs a recover operation successfully in root namespace for ${scenario.namespace} namespace item - Engine: ${scenario.engine}`, async function (assert) {
      await this.renderComponent();

      await click(GENERAL.selectByAttr('namespace'));
      await click(`[data-test-option="${scenario.namespace}"]`);
      await click(GENERAL.selectByAttr('mount'));
      await click(`[data-test-option="${scenario.mount}"]`);
      await fillIn(GENERAL.inputByAttr('resourcePath'), scenario.path);

      await click(GENERAL.button('recover'));

      assert.dom(GENERAL.inlineAlert).containsText('Success', 'shows success message');
      assert.dom(GENERAL.inlineAlert).containsText(scenario.path, 'shows the recovered path');

      await click(GENERAL.inputByAttr('copy'));
      await fillIn(GENERAL.inputByAttr('copyPath'), scenario.path + '-copy');
      await click(GENERAL.button('recover'));

      assert.dom(GENERAL.inlineAlert).containsText('Success', 'shows success message');
      assert.dom(GENERAL.inlineAlert).containsText(scenario.path + '-copy', 'shows the recovered copy path');
    });
  });

  test('it renders manual input if mounts request does not have data', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/ui/mounts', () => ({
      // Cubbyhole is enabled by default and cannot be disabled so supportedMounts should always
      // have values if the request was successful, but handling just in case.
      data: { secret: {} },
    }));
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);
    assert.dom(GENERAL.inputByAttr('manual-mount-path')).exists();
    assert.dom(GENERAL.selectByAttr('mount')).doesNotExist();
  });

  test('it renders manual input if mounts request errors', async function (assert) {
    assert.expect(14);
    this.server.get('/sys/internal/ui/mounts', () => overrideResponse(403));
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);
    assert.dom(GENERAL.button('Type')).exists().hasText('Type');
    assert.dom(GENERAL.inputByAttr('manual-mount-path')).exists().isDisabled();
    assert.dom(GENERAL.selectByAttr('mount')).doesNotExist();
    await click(GENERAL.button('Type'));
    assert.dom(GENERAL.button('Type')).hasAttribute('aria-expanded', 'true');
    // Select kv
    await click(GENERAL.radioByAttr('kv'));
    assert
      .dom(GENERAL.button('Type'))
      .hasAttribute('aria-expanded', 'false', 'toggle closes after selecting type');
    assert
      .dom(GENERAL.inputByAttr('manual-mount-path'))
      .isNotDisabled('input is enabled after selecting a type');
    assert.dom(GENERAL.button('Type')).hasText('KV v1');
    await fillIn(GENERAL.inputByAttr('manual-mount-path'), 'kv-mount');
    // Select cubbyhole
    await click(GENERAL.button('Type'));
    await click(GENERAL.radioByAttr('cubbyhole'));
    assert
      .dom(GENERAL.inputByAttr('manual-mount-path'))
      .hasAttribute('readonly', '', 'cubbyhole input is readonly')
      .hasValue('cubbyhole', 'input has default cubbyhole path');
    assert.dom(GENERAL.button('Type')).hasText('Cubbyhole');
    // Select database
    await click(GENERAL.button('Type'));
    await click(GENERAL.radioByAttr('database'));
    assert
      .dom(GENERAL.inputByAttr('manual-mount-path'))
      .hasValue('', 'input clears when selecting a new type');
    assert.dom(GENERAL.button('Type')).hasText('Database');
  });

  test('it displays loaded snapshot card', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);
    assert.dom(GENERAL.badge('status')).hasText('Ready', 'status badge renders');
  });

  test('it displays namespace selector for root namespace', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);

    assert.dom(GENERAL.selectByAttr('namespace')).exists('namespace selector is visible in root namespace');
  });

  test('it validates form fields before read/recover operations', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);
    // Try to read without selecting mount or resource path
    await click(GENERAL.button('read'));

    assert.dom(GENERAL.validationErrorByAttr('mount')).hasText('Please select a secret mount');
    assert.dom(GENERAL.validationErrorByAttr('resourcePath')).hasText('Please enter a resource path');
  });

  test('it clears form selections', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="1"]');
    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index]');
    await fillIn(GENERAL.inputByAttr('resourcePath'), 'test-path');

    await click(GENERAL.button('clear'));

    const nsSelect = find(GENERAL.selectByAttr('namespace'));
    assert.strictEqual(nsSelect.textContent.trim(), 'root', 'namespace was reset');

    const mountSelect = find(GENERAL.selectByAttr('mount'));
    assert.strictEqual(mountSelect.textContent.trim(), 'Select a mount here', 'mount is cleared');

    assert.dom(GENERAL.inputByAttr('resourcePath')).hasValue('', 'resource path is cleared');
  });

  test('it displays error alert when read operation fails', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);

    await fillIn(GENERAL.inputByAttr('resourcePath'), 'nonexistent-secret');
    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index="1.0"]');
    await click(GENERAL.button('read'));

    assert.dom(GENERAL.inlineAlert).containsText('Error', 'shows error alert');
  });

  test('it toggles JSON view in read modal', async function (assert) {
    await render(hbs`<Recovery::Page::Snapshots::SnapshotManage @model={{this.model}}/>`);

    await fillIn(GENERAL.inputByAttr('resourcePath'), 'test-secret');
    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index]');
    await click(GENERAL.button('read'));
    await waitFor('[data-test-read-secrets]');

    assert.dom('[data-test-read-secrets]').exists('read modal opens');

    await click(GENERAL.toggleInput('snapshot-read-secrets'));
    assert.dom('.hds-code-block').exists('renders JSON view');
  });
});
