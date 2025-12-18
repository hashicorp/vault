/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kv | kv-header', function (hooks) {
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
          <KvHeader
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
            @pageTitle={{this.pageTitle}}
            @mountName={{this.mountName}}
            @secretPath={{this.secretPath}}
          >
            <:tabs>
              <li><LinkTo @route="list" data-test-secrets-tab="Secrets">Secrets</LinkTo></li>
            </:tabs>

            <:badges>
              <Hds::Badge @text="version 2" data-test-badge />
            </:badges>

            <:toolbarActions>
              <ToolbarLink @route="secrets.create" @type="add">Create secret</ToolbarLink>
            </:toolbarActions>

            <:toolbarFilters>
              <p>stuff here</p>
            </:toolbarFilters>
          </KvHeader>
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
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create new version', 'displays custom title.');
    assert.dom(GENERAL.icon('key-values')).doesNotExist('Does not show icon if not at engine level.');
  });

  test('it renders a title and copy button for @secretPath', async function (assert) {
    assert.expect(3);
    this.secretPath = 'my/secret/path';
    this.pageTitle = this.secretPath;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('my/secret/path', 'displays path');
    assert.dom(GENERAL.copyButton).exists('renders copy button for path');
    assert.dom('[data-test-icon="clipboard-copy"]').exists('renders copy icon');
  });

  test('it renders a title, icon and tag if engine view', async function (assert) {
    assert.expect(3);
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'secrets' },
    ];
    this.backend = { id: this.backend, icon: 'key-values' };
    this.pageTitle = this.backend.id;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(`${this.pageTitle}`, 'Mount path renders for title.');
    assert.dom(GENERAL.badge()).hasText('version 2', 'version badge renders in header.');
    assert.dom(GENERAL.icon('key-values')).exists('An icon renders next to title.');
  });

  test('it renders tabs', async function (assert) {
    assert.expect(1);
    await this.renderComponent();
    assert.dom('[data-test-secrets-tab="Secrets"]').hasText('Secrets', 'Secrets tab renders');
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
