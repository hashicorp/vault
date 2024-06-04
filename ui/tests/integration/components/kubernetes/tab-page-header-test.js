/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | kubernetes | TabPageHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      },
    });
    this.model = this.store.peekRecord('secret-engine', 'kubernetes-test');
    this.mount = this.model.path.slice(0, -1);
    this.breadcrumbs = [{ label: 'Secrets', route: 'secrets', linkExternal: true }, { label: this.mount }];
    this.handleSearch = sinon.spy();
    this.handleInput = sinon.spy();
    this.handleKeyDown = sinon.spy();
  });

  test('it should render breadcrumbs', async function (assert) {
    await render(
      hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom('[data-test-breadcrumbs] li:nth-child(1) a').hasText('Secrets', 'Secrets breadcrumb renders');

    assert
      .dom('[data-test-breadcrumbs] li:nth-child(2)')
      .containsText(this.mount, 'Mount path breadcrumb renders');
  });

  test('it should render title', async function (assert) {
    await render(
      hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      {
        owner: this.engine,
      }
    );
    assert
      .dom('[data-test-header-title] svg')
      .hasClass('flight-icon-kubernetes-color', 'Correct icon renders in title');
    assert.dom('[data-test-header-title]').hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(
      hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}}/>`,
      {
        owner: this.engine,
      }
    );
    assert.dom('[data-test-tab="overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-tab="roles"]').hasText('Roles', 'Roles tab renders');
    assert.dom('[data-test-tab="config"]').hasText('Configuration', 'Configuration tab renders');
  });

  test('it should render filter for roles', async function (assert) {
    await render(
      hbs`<TabPageHeader @model={{this.model}} @filterRoles={{true}} @query="test" @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      { owner: this.engine }
    );
    assert.dom(GENERAL.filterInputExplicit).hasValue('test', 'Filter renders with provided value');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    await render(
      hbs`
      <TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}}>
        <span data-test-yield>It yields!</span>
      </TabPageHeader>
    `,
      { owner: this.engine }
    );

    assert
      .dom('.toolbar-actions [data-test-yield]')
      .hasText('It yields!', 'Block is yielded for toolbar actions');
  });
});
