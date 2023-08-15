/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { createSecretsEngine, generateBreadcrumbs } from 'vault/tests/helpers/ldap';

module('Integration | Component | ldap | Page::Roles', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.store = this.owner.lookup('service:store');
    this.backend = createSecretsEngine(this.store);
    this.breadcrumbs = generateBreadcrumbs(this.backend.id);

    for (const type of ['static', 'dynamic']) {
      this.store.pushPayload('ldap/role', {
        modelName: 'ldap/role',
        backend: 'ldap-test',
        type,
        ...this.server.create('ldap-role', type, { name: `${type}-test` }),
      });
    }
    this.backend = this.store.peekRecord('secret-engine', 'ldap-test');
    this.roles = this.store.peekAll('ldap/role');
    this.promptConfig = false;

    this.renderComponent = () => {
      return render(
        hbs`<Page::Roles
          @promptConfig={{this.promptConfig}}
          @backendModel={{this.backend}}
          @roles={{this.roles}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render tab page header and config cta', async function (assert) {
    this.promptConfig = true;

    await this.renderComponent();

    assert.dom('.title svg').hasClass('flight-icon-folder-users', 'LDAP icon renders in title');
    assert.dom('.title').hasText('ldap-test', 'Mount path renders in title');
    assert
      .dom('[data-test-toolbar-action="config"]')
      .hasText('Configure LDAP', 'Correct toolbar action renders');
    assert.dom('[data-test-config-cta]').exists('Config cta renders');
  });

  test('it should render create roles cta', async function (assert) {
    this.roles = null;

    await this.renderComponent();

    assert.dom('[data-test-toolbar-action="role"]').hasText('Create role', 'Toolbar action has correct text');
    assert
      .dom('[data-test-toolbar-action="role"] svg')
      .hasClass('flight-icon-plus', 'Toolbar action has correct icon');
    assert
      .dom('[data-test-filter-input]')
      .doesNotExist('Roles filter input is hidden when roles have not been created');
    assert.dom('[data-test-empty-state-title]').hasText('No roles created yet', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Roles in Vault will allow you to manage LDAP credentials. Create a role to get started.',
        'Message renders'
      );
    assert.dom('[data-test-empty-state-actions] a').hasText('Create role', 'Action renders');
  });

  test('it should render roles list', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-list-item-content] svg').hasClass('flight-icon-user', 'List item icon renders');
    assert
      .dom('[data-test-role="static-test"]')
      .hasText(this.roles.firstObject.name, 'List item name renders');
    assert
      .dom('[data-test-role-type-badge="static-test"]')
      .hasText(this.roles.firstObject.type, 'List item type badge renders');

    await click('[data-test-popup-menu-trigger]');
    assert.dom('[data-test-edit]').hasText('Edit', 'Edit link renders in menu');
    assert.dom('[data-test-get-creds]').hasText('Get credentials', 'Get credentials link renders in menu');
    assert
      .dom('[data-test-rotate-creds]')
      .hasText('Rotate credentials', 'Rotate credentials link renders in menu');
    assert.dom('[data-test-details]').hasText('Details', 'Details link renders in menu');
    assert.dom('[data-test-delete]').hasText('Delete', 'Details link renders in menu');

    await click('[data-test-popup-menu-trigger]:last-of-type');
    assert.dom('[data-test-rotate-creds]').doesNotExist('Rotate credentials link is hidden for dynamic type');
  });

  test('it should filter roles', async function (assert) {
    await this.renderComponent();

    await fillIn('[data-test-filter-input]', 'foo');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('There are no roles matching "foo"', 'Filter message renders');

    await fillIn('[data-test-filter-input]', 'static');
    assert.dom('[data-test-list-item-content]').exists({ count: 1 }, 'List is filtered with correct results');

    await fillIn('[data-test-filter-input]', '');
    assert
      .dom('[data-test-list-item-content]')
      .exists({ count: 2 }, 'All roles are displayed when filter is cleared');
  });
});
