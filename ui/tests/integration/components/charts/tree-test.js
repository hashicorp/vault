/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, settled, triggerKeyEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const TREE_DATA = [
  {
    name: 'Vault',
    children: [
      {
        name: 'Secrets',
        children: [
          { name: 'KV Store', value: 100 },
          { name: 'Database', value: 75 },
          { name: 'PKI', value: 50 },
        ],
      },
      {
        name: 'Auth',
        children: [
          { name: 'LDAP', value: 25 },
          { name: 'OIDC', value: 30 },
          { name: 'Userpass', value: 15 },
        ],
      },
      {
        name: 'Policies',
        children: [
          { name: 'ACL', value: 45 },
          { name: 'RGP', value: 20 },
          { name: 'EGP', value: 35 },
        ],
      },
    ],
  },
];

const TREE_OPTIONS = {
  title: 'Vault Hierarchy Tree',
  height: '400px',
  width: '600px',
  tree: {
    rootTitle: 'root',
  },
  theme: 'g10',
};
const SELECTORS = {
  tree: (title) => (title ? `[data-test-tree-chart="${title}"]` : '[data-test-tree-chart]'),
};

module('Integration | Component | tree-chart', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = TREE_DATA;
    this.options = TREE_OPTIONS;
    this.title = 'Test';

    this.renderComponent = async () => {
      await render(hbs`
          <TreeChart @data={{this.data}} @options={{this.options}} @title={{this.title}}/>
      `);
    };
  });

  test('it renders carbon tree chart', async function (assert) {
    await this.renderComponent();
    assert
      .dom(SELECTORS.tree(this.title))
      .hasAttribute('role', 'group')
      .hasAttribute('aria-label', this.title);
    assert.dom(SELECTORS.tree(this.title)).containsText(TREE_OPTIONS.title, 'title is rendered');
    assert.dom(SELECTORS.tree()).containsText(TREE_OPTIONS.tree.rootTitle, 'root node title is rendered');
    assert.dom(SELECTORS.tree()).hasTextContaining('Secrets', 'nodes are rendered');
  });

  test('it preserves toolbar controls and keyboard access to the overflow menu', async function (assert) {
    await this.renderComponent();
    assert
      .dom(`${SELECTORS.tree()} .toolbar-control`)
      .doesNotHaveAttribute('role', 'toolbar wrappers are not interactive');
    assert
      .dom(`${SELECTORS.tree()} button[aria-label="Show as table"]`)
      .hasAttribute('aria-disabled', 'false', 'table view remains available');
    assert.dom(`${SELECTORS.tree()} button[aria-label="Make fullscreen"]`).exists();

    await click(`${SELECTORS.tree()} button[aria-label="More options"]`);
    assert
      .dom(`${SELECTORS.tree()} button[aria-label="More options"]`)
      .hasAttribute('aria-expanded', 'true', 'overflow menu opens');
    assert
      .dom(`${SELECTORS.tree()} button[role="menuitem"]`)
      .isFocused('export action receives keyboard focus');
    await triggerKeyEvent(`${SELECTORS.tree()} button[role="menuitem"]`, 'keydown', 'Escape');
    assert
      .dom(`${SELECTORS.tree()} button[aria-label="More options"]`)
      .hasAttribute('aria-expanded', 'false', 'Escape closes the overflow menu');
  });

  test('it handles data updates', async function (assert) {
    await this.renderComponent();
    assert.dom(SELECTORS.tree()).exists('initial chart renders');

    this.set('data', [
      {
        name: 'Updated',
        children: [
          { name: 'Child 1', value: 10 },
          { name: 'Child 2', value: 20 },
        ],
      },
    ]);
    await settled();

    assert.dom(SELECTORS.tree()).hasTextContaining('Updated', 'new node is rendered');
    assert.dom(SELECTORS.tree()).doesNotHaveTextContaining('Secrets', 'old nodes are not rendered');
    assert
      .dom(`${SELECTORS.tree()} .toolbar-control`)
      .doesNotHaveAttribute('role', 'toolbar semantics remain valid after updating data');
    await click(`${SELECTORS.tree()} button[aria-label="More options"]`);
    assert
      .dom(`${SELECTORS.tree()} button[role="menuitem"]`)
      .isFocused('overflow menu remains accessible after updating data');
  });
});
