/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import kubernetesHandlers from 'vault/mirage/handlers/kubernetes';
import authPage from 'vault/tests/pages/auth';
import { visit, click, currentRouteName } from '@ember/test-helpers';
import { selectChoose } from 'ember-power-select/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { KUBERNETES_OVERVIEW } from 'vault/tests/helpers/kubernetes/kubernetes-selectors';

module('Acceptance | kubernetes | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    kubernetesHandlers(this.server);
    this.createScenario = (shouldConfigureRoles = true) =>
      shouldConfigureRoles ? kubernetesScenario(this.server) : kubernetesScenario(this.server, false);

    this.visitOverview = () => {
      return visit('/vault/secrets/kubernetes/kubernetes/overview');
    };
    this.validateRoute = (assert, route, message) => {
      assert.strictEqual(currentRouteName(), `vault.cluster.secrets.backend.kubernetes.${route}`, message);
    };
    return authPage.login();
  });

  test('it should transition to configuration page during empty state', async function (assert) {
    assert.expect(1);
    await this.visitOverview();
    await click(`${GENERAL.emptyStateActions} a`);
    this.validateRoute(assert, 'configure', 'Transitions to Configure route on click');
  });

  test('it should transition to view roles', async function (assert) {
    assert.expect(1);
    this.createScenario();
    await this.visitOverview();
    await click(KUBERNETES_OVERVIEW.rolesCardLink);
    this.validateRoute(assert, 'roles.index', 'Transitions to roles route on View Roles click');
  });

  test('it should transition to create roles', async function (assert) {
    assert.expect(1);
    this.createScenario(false);
    await this.visitOverview();
    await click(KUBERNETES_OVERVIEW.rolesCardLink);
    this.validateRoute(assert, 'roles.create', 'Transitions to roles route on Create Roles click');
  });

  test('it should transition to generate credentials', async function (assert) {
    assert.expect(1);
    await this.createScenario();
    await this.visitOverview();
    await selectChoose('.search-select', 'role-0');
    await click(KUBERNETES_OVERVIEW.generateCredentialsCardButton);
    this.validateRoute(assert, 'roles.role.credentials', 'Transitions to roles route on Generate click');
  });
});
