/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
module('Integration | Component | kv-v2 | Page::Secret::Paths', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.canReadMetadata = true;
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];

    this.renderComponent = async () => {
      await render(
        hbs`
      <Page::Secret::Paths
        @path={{this.path}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
        @canReadMetadata={{this.canReadMetadata}}
      />
      `,
        { owner: this.engine }
      );
    };
  });

  test('it renders tabs', async function (assert) {
    await this.renderComponent();
    const tabs = ['Secret', 'Metadata', 'Paths', 'Version History'];
    for (const tab of tabs) {
      assert.dom(PAGE.secretTab(tab)).hasText(tab);
    }
  });

  test('it hides version history when cannot READ metadata', async function (assert) {
    this.canReadMetadata = false;
    await this.renderComponent();
    const tabs = ['Secret', 'Metadata', 'Paths'];
    for (const tab of tabs) {
      assert.dom(PAGE.secretTab(tab)).hasText(tab);
    }
    assert.dom(PAGE.secretTab('Version History')).doesNotExist();
  });

  test('it renders header', async function (assert) {
    await this.renderComponent();
    assert.dom(PAGE.breadcrumbs).hasText(`Secrets ${this.backend} ${this.path}`);
    assert.dom(PAGE.title).hasText(this.path);
  });

  test('it renders commands which is the uncondensed version of KvPathsCard', async function (assert) {
    await this.renderComponent();
    assert.dom(PAGE.paths.codeSnippet('cli')).exists();
    assert.dom(PAGE.paths.codeSnippet('api')).exists();
  });
});
