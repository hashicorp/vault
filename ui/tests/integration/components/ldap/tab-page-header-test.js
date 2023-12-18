/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | ldap | TabPageHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'ldap_64e858b1',
        path: 'ldap-test/',
        type: 'ldap',
      },
    });
    this.model = this.store.peekRecord('secret-engine', 'ldap-test');
    this.mount = this.model.path.slice(0, -1);
    this.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: this.mount }];
  });

  test('it should render breadcrumbs', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-breadcrumbs] li:nth-child(1) a').hasText('secrets', 'Secrets breadcrumb renders');

    assert
      .dom('[data-test-breadcrumbs] li:nth-child(2)')
      .containsText(this.mount, 'Mount path breadcrumb renders');
  });

  test('it should render title', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-header-title] svg')
      .hasClass('flight-icon-folder-users', 'Correct icon renders in title');
    assert.dom('[data-test-header-title]').hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-tab="overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-tab="roles"]').hasText('Roles', 'Roles tab renders');
    assert.dom('[data-test-tab="libraries"]').hasText('Libraries', 'Libraries tab renders');
    assert.dom('[data-test-tab="config"]').hasText('Configuration', 'Configuration tab renders');
  });

  test('it should yield toolbar blocks', async function (assert) {
    await render(
      hbs`
      <TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}}>
        <:toolbarFilters>
          <span data-test-filters>Toolbar filters</span>
        </:toolbarFilters>
        <:toolbarActions>
          <span data-test-actions>Toolbar actions</span>
        </:toolbarActions>
      </TabPageHeader>
    `,
      { owner: this.engine }
    );

    assert
      .dom('.toolbar-filters [data-test-filters]')
      .hasText('Toolbar filters', 'Block is yielded for toolbar filters');
    assert
      .dom('.toolbar-actions [data-test-actions]')
      .hasText('Toolbar actions', 'Block is yielded for toolbar actions');
  });
});
