/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kv | kv-page-header', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.version = 2;
    this.id = kvDataPath(this.backend, this.path, this.version);
    this.payload = {
      backend: this.backend,
      path: this.path,
      version: 2,
    };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.id,
      ...this.payload,
    });

    this.model = this.store.peekRecord('kv/data', this.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'secrets' },
      { label: this.model.path, route: 'secrets.secret.details', model: this.model.path },
      { label: 'edit' },
    ];
  });

  test('it renders breadcrumbs', async function (assert) {
    assert.expect(4);
    await render(hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @pageTitle="Create new version"/>`, {
      owner: this.engine,
    });
    assert.dom(GENERAL.breadcrumbAtIdx(1)).hasText('secrets', 'Secrets breadcrumb renders');
    assert.dom(GENERAL.breadcrumbAtIdx(2)).hasText(this.backend, 'engine name renders');
    assert.dom(GENERAL.breadcrumbAtIdx(3)).hasText(this.path, 'secret path renders');
    assert.dom(GENERAL.breadcrumbAtIdx(4)).hasClass('hds-breadcrumb__current');
  });

  test('it renders a custom title for a non engine view', async function (assert) {
    assert.expect(2);
    await render(hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @pageTitle="Create new version"/>`, {
      owner: this.engine,
    });
    assert.dom(GENERAL.title).hasText('Create new version', 'displays custom title.');
    assert.dom(`${GENERAL.title} svg`).doesNotExist('Does not show icon if not at engine level.');
  });

  test('it renders a title, icon and tag if engine view', async function (assert) {
    assert.expect(2);
    await render(hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName={{this.backend}} />`, {
      owner: this.engine,
    });
    assert
      .dom(GENERAL.title)
      .hasText(`${this.backend} version 2`, 'Mount path and version tag render for title.');
    assert.dom(`${GENERAL.title} [data-test-icon="key-values"]`).exists('An icon renders next to title.');
  });

  test('it renders tabs', async function (assert) {
    assert.expect(2);
    await render(
      hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName="my-engine">
  <:tabLinks>
    <li><LinkTo @route="list" data-test-tab="Secrets">Secrets</LinkTo></li>
    <li><LinkTo @route="configuration" data-test-tab="Configuration">Configuration</LinkTo></li>
  </:tabLinks>
  </KvPageHeader>
    `,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.tab('Secrets')).hasText('Secrets', 'Secrets tab renders');
    assert.dom(GENERAL.tab('Configuration')).hasText('Configuration', 'Configuration tab renders');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    assert.expect(1);
    await render(
      hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName="my-engine">
      <:toolbarActions>
      <ToolbarLink @route="secrets.create" @type="add">Create secret</ToolbarLink>
    </:toolbarActions>
  </KvPageHeader>
    `,
      { owner: this.engine }
    ),
      assert.dom('.toolbar-actions').exists('Block is yielded for toolbar actions');
  });

  test('it should yield block for toolbar filters', async function (assert) {
    assert.expect(1);
    await render(
      hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName="my-engine">
      <:toolbarFilters>
      <p>stuff here</p>
    </:toolbarFilters>
  </KvPageHeader>
    `,
      { owner: this.engine }
    ),
      assert.dom('.toolbar-filters').exists('Block is yielded for toolbar filters');
  });
});
