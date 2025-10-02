/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | kv | kv-page-header', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';

    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'secrets' },
      { label: this.path, route: 'secrets.secret.details', model: this.path },
      { label: 'Edit' },
    ];

    this.renderComponent = () =>
      render(
        hbs`
          <KvPageHeader
            @breadcrumbs={{this.breadcrumbs}}
            @pageTitle={{this.pageTitle}}
            @mountName={{this.mountName}}
            @secretPath={{this.secretPath}}
          >
            <:tabLinks>
              <li><LinkTo @route="list" data-test-secrets-tab="Secrets">Secrets</LinkTo></li>
              <li><LinkTo @route="configuration" data-test-secrets-tab="Configuration">Configuration</LinkTo></li>
            </:tabLinks>

            <:toolbarActions>
              <ToolbarLink @route="secrets.create" @type="add">Create secret</ToolbarLink>
            </:toolbarActions>

            <:toolbarFilters>
              <p>stuff here</p>
            </:toolbarFilters>
          </KvPageHeader>
        `,
        { owner: this.engine }
      );
  });

  test('it renders breadcrumbs', async function (assert) {
    assert.expect(4);
    await this.renderComponent();
    assert.dom('[data-test-breadcrumbs] li:nth-child(1) a').hasText('Secrets', 'Secrets breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').hasText(this.backend, 'engine name renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(3) a').hasText(this.path, 'secret path renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(4) .hds-breadcrumb__current')
      .exists('final breadcrumb renders and it is not a link.');
  });

  test('it renders a custom title for @pageTitle', async function (assert) {
    assert.expect(2);
    this.pageTitle = 'Create new version';
    await this.renderComponent();
    assert.dom('[data-test-header-title]').hasText('Create new version', 'displays custom title.');
    assert.dom('[data-test-header-title] svg').doesNotExist('Does not show icon if not at engine level.');
  });

  test('it renders a title and copy button for @secretPath', async function (assert) {
    assert.expect(3);
    this.secretPath = 'my/secret/path';
    await this.renderComponent();
    assert.dom('[data-test-header-title]').hasText('my/secret/path', 'displays path');
    assert.dom('[data-test-header-title] button').exists('renders copy button for path');
    assert.dom('[data-test-icon="clipboard-copy"]').exists('renders copy icon');
  });

  test('it renders a title, icon and tag if engine view', async function (assert) {
    assert.expect(2);
    this.mountName = this.backend;
    await this.renderComponent();
    assert
      .dom('[data-test-header-title]')
      .hasText(`${this.backend} version 2`, 'Mount path and version tag render for title.');
    assert
      .dom('[data-test-header-title] [data-test-icon="key-values"]')
      .exists('An icon renders next to title.');
  });

  test('it renders tabs', async function (assert) {
    assert.expect(2);
    await this.renderComponent();
    assert.dom('[data-test-secrets-tab="Secrets"]').hasText('Secrets', 'Secrets tab renders');
    assert
      .dom('[data-test-secrets-tab="Configuration"]')
      .hasText('Configuration', 'Configuration tab renders');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    assert.expect(1);
    await this.renderComponent();
    assert.dom('.toolbar-actions').exists('Block is yielded for toolbar actions');
  });

  test('it should yield block for toolbar filters', async function (assert) {
    assert.expect(1);
    await this.renderComponent();
    assert.dom('.toolbar-filters').exists('Block is yielded for toolbar filters');
  });
});
