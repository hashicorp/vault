/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | kubernetes | Page::Roles', function (hooks) {
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
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-config'),
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-role'),
    });
    this.backend = this.store.peekRecord('secret-engine', 'kubernetes-test');
    this.config = this.store.peekRecord('kubernetes/config', 'kubernetes-test');
    this.roles = this.store.peekAll('kubernetes/role');
    this.filterValue = '';
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.backend.id },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Roles @config={{this.config}} @backend={{this.backend}} @roles={{this.roles}} @filterValue={{this.filterValue}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render tab page header and config cta', async function (assert) {
    this.config = null;
    await this.renderComponent();
    assert.dom('.title svg').hasClass('flight-icon-kubernetes', 'Kubernetes icon renders in title');
    assert.dom('.title').hasText('kubernetes-test', 'Mount path renders in title');
    assert.dom('[data-test-toolbar-roles-action]').hasText('Create role', 'Toolbar action has correct text');
    assert
      .dom('[data-test-toolbar-roles-action] svg')
      .hasClass('flight-icon-plus', 'Toolbar action has correct icon');
    assert.dom('[data-test-nav-input]').exists('Roles filter input renders');
    assert.dom('[data-test-config-cta]').exists('Config cta renders');
  });

  test('it should render create roles cta', async function (assert) {
    this.roles = null;
    await this.renderComponent();
    assert.dom('[data-test-empty-state-title]').hasText('No roles yet', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'When created, roles will be listed here. Create a role to start generating service account tokens.',
        'Message renders'
      );
    assert.dom('[data-test-empty-state-actions] a').hasText('Create role', 'Action renders');
  });

  test('it should render no matches filter message', async function (assert) {
    this.roles = [];
    this.filterValue = 'test';
    await this.renderComponent();
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('There are no roles matching "test"', 'Filter message renders');
  });

  test('it should render roles list', async function (assert) {
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        'kubernetes/role': ['root'],
      },
    }));
    await this.renderComponent();
    assert.dom('[data-test-list-item-content] svg').hasClass('flight-icon-user', 'List item icon renders');
    assert
      .dom('[data-test-list-item-content]')
      .hasText(this.roles.firstObject.name, 'List item name renders');
    await click('[data-test-popup-menu-trigger]');
    assert.dom('[data-test-details]').hasText('Details', 'Details link renders in menu');
    assert.dom('[data-test-edit]').hasText('Edit', 'Edit link renders in menu');
    assert.dom('[data-test-delete]').hasText('Delete', 'Details link renders in menu');
  });
});
