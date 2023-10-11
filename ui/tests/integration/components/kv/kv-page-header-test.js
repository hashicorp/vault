import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';

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
    assert.dom('[data-test-breadcrumbs] li:nth-child(1) a').hasText('secrets', 'Secrets breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').hasText(this.backend, 'engine name renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(3) a').hasText(this.path, 'secret path renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(4)')
      .hasText('/ edit', 'final breadcrumb renders and it is not a link.');
  });

  test('it renders a custom title for a non engine view', async function (assert) {
    assert.expect(2);
    await render(hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @pageTitle="Create new version"/>`, {
      owner: this.engine,
    });
    assert.dom('[data-test-header-title]').hasText('Create new version', 'displays custom title.');
    assert.dom('[data-test-header-title] svg').doesNotExist('Does not show icon if not at engine level.');
  });

  test('it renders a title, icon and tag if engine view', async function (assert) {
    assert.expect(2);
    await render(hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName={{this.backend}} />`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-header-title]')
      .hasText(`${this.backend} Version 2`, 'Mount path and Version tag render for title.');
    assert.dom('[data-test-header-title] span').hasClass('hs-icon', 'An icon renders next to title.');
  });

  test('it renders tabs', async function (assert) {
    assert.expect(2);
    await render(
      hbs`<KvPageHeader @breadcrumbs={{this.breadcrumbs}} @mountName="my-engine">
  <:tabLinks>
    <LinkTo @route="list" data-test-secrets-tab="Secrets">Secrets</LinkTo>
    <LinkTo @route="configuration" data-test-secrets-tab="Configuration">Configuration</LinkTo>
  </:tabLinks>
  </KvPageHeader>
    `,
      {
        owner: this.engine,
      }
    );
    assert.dom('[data-test-secrets-tab="Secrets"]').hasText('Secrets', 'Secrets tab renders');
    assert
      .dom('[data-test-secrets-tab="Configuration"]')
      .hasText('Configuration', 'Configuration tab renders');
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
