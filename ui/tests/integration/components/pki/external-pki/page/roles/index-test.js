/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import sinon from 'sinon';
import { CreationMethod } from 'vault/utils/constants/snippet';

module('Integration | Component | pki | external-pki | ExternalPki::Page::Roles', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.setupModel = (roles = []) => {
      return {
        engine: new SecretsEngineResource({
          accessor: 'pki-external-ca_e158c567',
          type: 'pki-external-ca',
          path: 'my-pki-external-ca/',
        }),
        roles,
      };
    };

    this.renderComponent = () =>
      render(hbs`<ExternalPki::Page::Roles @model={{this.model}} />`, { owner: this.engine });
  });

  test('it renders empty state when no roles exist', async function (assert) {
    this.model = this.setupModel([]);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No roles exist yet');
    // Implementation select should be visible
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM)).exists();
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).exists();
    assert.dom(GENERAL.textDisplay('0')).hasText('Create a role');
    // Search and refresh should not be visible
    assert.dom(GENERAL.inputSearch('Filter by name')).doesNotExist();
    assert.dom(GENERAL.button('Refresh')).doesNotExist();
  });

  module('with roles', function (hooks) {
    hooks.beforeEach(function () {
      this.roles = [
        'web-server',
        'api-gateway',
        'database-server',
        'load-balancer',
        'prod-server',
        'potato',
        'API-gateway',
      ];
      this.model = this.setupModel(this.roles);
    });

    test('it renders list of roles', async function (assert) {
      await this.renderComponent();

      // Empty state should not be visible
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      // Search and refresh should be visible
      assert.dom(GENERAL.inputSearch('Filter by name')).exists();
      assert.dom(GENERAL.button('Refresh')).exists();
      // Table should render with all roles
      assert.dom(GENERAL.listItem()).exists({ count: 7 });
      this.roles.forEach((r) => {
        assert.dom(GENERAL.linkTo(r)).hasText(r, `table renders link for role: ${r}`);
      });
    });

    test('it filters roles by search input', async function (assert) {
      await this.renderComponent();

      // Initially all roles are visible
      assert.dom(GENERAL.listItem()).exists({ count: 7 });
      assert.dom(GENERAL.pagination).hasTextContaining('1–7 of 7');
      // Filter by partial match "a"
      await fillIn(GENERAL.inputSearch('Filter by name'), 'a');
      assert.dom(GENERAL.listItem()).exists({ count: 5 }, 'it filters matching roles');
      assert.dom(GENERAL.pagination).hasTextContaining('1–5 of 5');
      assert.dom(GENERAL.linkTo('api-gateway')).exists();
      assert.dom(GENERAL.linkTo('API-gateway')).exists();
      assert.dom(GENERAL.linkTo('database-server')).exists();
      assert.dom(GENERAL.linkTo('load-balancer')).exists();
      assert.dom(GENERAL.linkTo('potato')).exists();
      assert.dom(GENERAL.linkTo('web-server')).doesNotExist();
      assert.dom(GENERAL.linkTo('prod-server')).doesNotExist();
      // Filter by "api"
      await fillIn(GENERAL.inputSearch('Filter by name'), 'server');
      assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'shows 3 matching roles');
      assert.dom(GENERAL.pagination).hasTextContaining('1–3 of 3');
      assert.dom(GENERAL.linkTo('web-server')).exists();
      assert.dom(GENERAL.linkTo('prod-server')).exists();
      assert.dom(GENERAL.linkTo('database-server')).exists();
      assert.dom(GENERAL.linkTo('api-gateway')).doesNotExist();
      assert.dom(GENERAL.linkTo('load-balancer')).doesNotExist();
      assert.dom(GENERAL.linkTo('potato')).doesNotExist();
      // Clear filter
      await fillIn(GENERAL.inputSearch('Filter by name'), '');
      assert.dom(GENERAL.listItem()).exists({ count: 7 }, 'shows all roles again');
      assert.dom(GENERAL.pagination).hasTextContaining('1–7 of 7');
    });

    test('it paginates filtered roles', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.paginationSizeSelector, '5');
      // Initially all roles are visible
      assert.dom(GENERAL.listItem()).exists({ count: 5 });
      assert.dom(GENERAL.pagination).hasTextContaining('1–5 of 7');
      await click(GENERAL.nextPage);
      assert.dom(GENERAL.listItem()).exists({ count: 2 });
      assert.dom(GENERAL.pagination).hasTextContaining('6–7 of 7');

      // Filter by partial match "e"
      await fillIn(GENERAL.inputSearch('Filter by name'), 'e');
      assert.dom(GENERAL.listItem()).exists({ count: 5 }, 'it filters matching roles');
      assert.dom(GENERAL.pagination).hasTextContaining('1–5 of 6');
      await click(GENERAL.nextPage);
      assert.dom(GENERAL.listItem()).exists({ count: 1 });
      assert.dom(GENERAL.pagination).hasTextContaining('6–6 of 6');
    });

    test('it shows empty state when search has no matches', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.inputSearch('Filter by name'), 'nonexistent-role');
      assert.dom(GENERAL.listItem()).doesNotExist();
      assert.dom(GENERAL.emptyStateTitle).hasText('No roles matching: nonexistent-role');

      // Implementation select should not be visible during search
      assert.dom(GENERAL.radioCardByAttr()).doesNotExist();
    });

    test('search is case sensitive', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.inputSearch('Filter by name'), 'API');
      assert.dom(GENERAL.listItem()).exists({ count: 1 });
      assert.dom(GENERAL.linkTo('API-gateway')).exists();
      assert.dom(GENERAL.linkTo('api-gateway')).doesNotExist();
    });

    test('it calls refresh when refresh button is clicked', async function (assert) {
      const router = this.owner.lookup('service:router');
      const refreshStub = sinon.stub(router, 'refresh');

      // Mock currentRoute
      const currentRouteName = 'vault.cluster.secrets.backend.pki.external.roles';
      sinon.stub(router, 'currentRoute').value({
        parent: { name: currentRouteName },
      });
      await this.renderComponent();
      await click(GENERAL.button('Refresh'));
      assert.true(refreshStub.calledOnce, 'refresh was called once');
      assert.true(
        refreshStub.calledWith('vault.cluster.secrets.backend.pki.external.roles'),
        'refresh was called with correct route'
      );
    });
  });
});
