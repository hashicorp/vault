/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';

module('Integration | Component | kubernetes | Page::Roles', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kubernetes-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.secretsEngine = new SecretsEngineResource({
      accessor: 'kubernetes_f3400dee',
      path: 'kubernetes-test/',
      type: 'kubernetes',
    });
    this.roleName = 'role-0';
    this.roles = [this.roleName];
    this.filterValue = '';
    this.breadcrumbs = [{ label: 'Secrets', route: 'secrets', linkExternal: true }, { label: this.backend }];
    this.promptConfig = false;

    const path = this.owner
      .lookup('service:capabilities')
      .pathFor('kubernetesRole', { backend: this.backend, name: this.roleName });
    this.capabilities = { [path]: { canRead: true, canUpdate: true, canDelete: true } };

    this.renderComponent = () => {
      return render(
        hbs`<Page::Roles
          @promptConfig={{this.promptConfig}}
          @secretsEngine={{this.secretsEngine}}
          @roles={{this.roles}}
          @capabilities={{this.capabilities}}
          @filterValue={{this.filterValue}}
          @breadcrumbs={{this.breadcrumbs}}
        />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render tab page header and config cta', async function (assert) {
    this.promptConfig = true;
    await this.renderComponent();
    assert
      .dom(GENERAL.icon('kubernetes-color'))
      .hasClass('hds-icon-kubernetes-color', 'Kubernetes icon renders in title');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('kubernetes-test', 'Mount path renders in title');
    assert
      .dom('[data-test-toolbar-roles-action]')
      .doesNotExist('Create role', 'Toolbar action does not render when not configured');
    assert
      .dom(GENERAL.filterInputExplicit)
      .doesNotExist('Roles filter input does not render when not configured');
    assert.dom(GENERAL.emptyStateTitle).hasText('Kubernetes not configured');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure Kubernetes');
  });

  test('it should render create roles cta', async function (assert) {
    this.roles = null;
    await this.renderComponent();
    assert.dom('[data-test-toolbar-roles-action]').hasText('Create role', 'Toolbar action has correct text');
    assert
      .dom('[data-test-toolbar-roles-action] svg')
      .hasClass('hds-icon-plus', 'Toolbar action has correct icon');
    assert.dom(GENERAL.filterInputExplicit).exists('Roles filter input renders');
    assert.dom(GENERAL.emptyStateTitle).hasText('No roles yet', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'When created, roles will be listed here. Create a role to start generating service account tokens.',
        'Message renders'
      );
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
    await this.renderComponent();
    assert.dom('[data-test-list-item-content] svg').hasClass('hds-icon-user', 'List item icon renders');
    assert.dom('[data-test-list-item-content]').hasText(this.roles[0], 'List item name renders');
    await click(GENERAL.menuTrigger);
    assert.dom('[data-test-details]').hasText('Details', 'Details link renders in menu');
    assert.dom('[data-test-edit]').hasText('Edit', 'Edit link renders in menu');
    assert.dom('[data-test-delete]').hasText('Delete', 'Details link renders in menu');
  });
});
