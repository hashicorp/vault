/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | kubernetes | KubernetesHeader', function (hooks) {
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
      hbs`<KubernetesHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
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
      hbs`<KubernetesHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      {
        owner: this.engine,
      }
    );
    assert
      .dom(GENERAL.icon('kubernetes-color'))
      .hasClass('hds-icon-kubernetes-color', 'Correct icon renders in title');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(
      hbs`<KubernetesHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}}/>`,
      {
        owner: this.engine,
      }
    );
    assert.dom('[data-test-tab="overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-tab="roles"]').hasText('Roles', 'Roles tab renders');
  });

  test('it should render filter for roles', async function (assert) {
    await render(
      hbs`<KubernetesHeader @model={{this.model}} @filterRoles={{true}} @query="test" @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      { owner: this.engine }
    );
    assert.dom(GENERAL.filterInputExplicit).hasValue('test', 'Filter renders with provided value');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    await render(
      hbs`
      <KubernetesHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}}>
        <span data-test-yield>It yields!</span>
      </KubernetesHeader>
    `,
      { owner: this.engine }
    );

    assert
      .dom('.toolbar-actions [data-test-yield]')
      .hasText('It yields!', 'Block is yielded for toolbar actions');
  });

  test('it should render a dropdown when configRoute is omitted', async function (assert) {
    await render(
      hbs`<KubernetesHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.dropdownToggle('Manage')).hasText('Manage', 'Manage dropdown renders');
    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Configure')).exists('Configure dropdown item renders');
    assert.dom(GENERAL.menuItem('Delete')).exists('Configure dropdown item renders');
  });

  test('it should render exit configuration button when configRoute is provided', async function (assert) {
    await render(
      hbs`<KubernetesHeader @configRoute="configuration" @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} @handleSearch={{this.handleSearch}} @handleInput={{this.handleInput}} @handleKeyDown={{this.handleKeyDown}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(GENERAL.button('Exit configuration')).exists('Exit configuration button renders');
  });
});
