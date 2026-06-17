/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | usage-reporting/namespace-picker', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the dropdown with the first namespace selected by default', async function (assert) {
    this.set('namespaces', ['root', 'ns1', 'ns2']);
    this.set('onNamespaceChange', sinon.stub());
    await render(hbs`
      <UsageReporting::NamespacePicker
        @namespaces={{this.namespaces}}
        @onNamespaceChange={{this.onNamespaceChange}}
      />
    `);
    assert.dom('[data-test-vault-reporting-namespace-picker]').exists('namespace picker dropdown renders');
    // The toggle button shows the selected namespace
    assert
      .dom('[data-test-vault-reporting-namespace-picker] button')
      .hasText('root', 'first namespace is selected by default');
  });

  test('it lists all namespaces in the dropdown', async function (assert) {
    this.set('namespaces', ['root', 'ns1', 'ns2']);
    this.set('onNamespaceChange', sinon.stub());
    await render(hbs`
      <UsageReporting::NamespacePicker
        @namespaces={{this.namespaces}}
        @onNamespaceChange={{this.onNamespaceChange}}
      />
    `);
    await click('[data-test-vault-reporting-namespace-picker] button');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="root"]')
      .exists('root namespace item renders');
    assert.dom('[data-test-vault-reporting-namespace-menu-item="ns1"]').exists('ns1 namespace item renders');
    assert.dom('[data-test-vault-reporting-namespace-menu-item="ns2"]').exists('ns2 namespace item renders');
  });

  test('it filters namespaces on search input', async function (assert) {
    this.set('namespaces', ['root', 'ns1', 'ns2']);
    this.set('onNamespaceChange', sinon.stub());
    await render(hbs`
      <UsageReporting::NamespacePicker
        @namespaces={{this.namespaces}}
        @onNamespaceChange={{this.onNamespaceChange}}
      />
    `);
    await click('[data-test-vault-reporting-namespace-picker] button');
    await fillIn('[data-test-vault-reporting-namespace-search]', 'ns1');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="ns1"]')
      .exists('matching namespace is visible');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="root"]')
      .doesNotExist('non-matching namespace is hidden');
    assert
      .dom('[data-test-vault-reporting-namespace-menu-item="ns2"]')
      .doesNotExist('non-matching namespace ns2 is hidden');
  });

  test('it calls @onNamespaceChange when a namespace is selected', async function (assert) {
    this.set('namespaces', ['root', 'ns1', 'ns2']);
    const stub = sinon.stub();
    this.set('onNamespaceChange', stub);
    await render(hbs`
      <UsageReporting::NamespacePicker
        @namespaces={{this.namespaces}}
        @onNamespaceChange={{this.onNamespaceChange}}
      />
    `);
    await click('[data-test-vault-reporting-namespace-picker] button');
    await click('[data-test-vault-reporting-namespace-menu-item="ns1"]');
    assert.ok(stub.calledOnceWith('ns1'), 'onNamespaceChange is called with selected namespace');
  });
});
