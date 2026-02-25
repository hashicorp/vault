/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from '../helpers/general-selectors';
import { overrideResponse } from '../helpers/stubs';
import { setupMirage } from 'ember-cli-mirage/test-support';

// The route "vault.cluster.not-found" catches any unmatched routes within the cluster (e.g., /vault/fake-route)
// Otherwise the most closely related error sub-state should render.
module('Acceptance | router error handling', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it handles when a route does not exist at the app route level', async function (assert) {
    await visit('/not-real-route');
    assert.strictEqual(currentRouteName(), 'vault.error', 'it redirects to error route');
    assert.dom('[data-test-sidebar-nav-panel="Cluster"]').doesNotExist('sidebar nav panel does not render');
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found ');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('Sorry, we were unable to find any content at not-real-route.');
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Double check the URL or return to the dashboard. Go to dashboard');
    assert.dom('.hds-application-state').hasClass('align-self-center');
    assert.dom('[data-test-app-footer]').exists('app footer still renders');
  });

  test('it handles when a route does not exist at the cluster route level', async function (assert) {
    const route = 'some-fake-route';
    await visit(`/vault/${route}`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.not-found', 'it redirects to not-found route');
    assert.dom('[data-test-sidebar-nav-panel="Cluster"]').exists('sidebar nav panel still renders');
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found');
    assert.dom(GENERAL.pageError.message).hasText(`Sorry, we were unable to find any content at ${route}.`);
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Double check the URL or return to the dashboard. Go to dashboard');
    assert.dom('[data-test-app-footer]').exists('app footer still renders');
  });

  // Since there is no `secrets/error.hbs` template errors bubble up to the cluster route
  test('it handles when a secrets engine is does not exist', async function (assert) {
    const path = 'notarealengine';
    await visit(`/vault/secrets-engines/${path}/configuration/general-settings`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.error', 'it redirects to cluster error route');
    assert.dom('[data-test-sidebar-nav-panel="Cluster"]').exists('sidebar nav panel still renders');
    assert.dom('[data-test-app-footer]').exists('app footer still renders');
    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        `preflight capability check returned 403, please ensure client's policies grant access to path "${path}/"`,
        'it renders api message'
      );
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Double check the URL or return to the dashboard. Go to dashboard');
  });

  // There IS a `secrets/backends/error.hbs` template which is what renders here
  test('it handles when a secret path does not exist', async function (assert) {
    await visit('/vault/secrets-engines/cubbyhole/show/not-real');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.error',
      'it redirects to secret error substate'
    );
    assert.dom(GENERAL.pageError.title(404)).hasText('ERROR 404 Not found');
    assert
      .dom(GENERAL.pageError.message)
      .hasText('Sorry, we were unable to find any content at /v1/cubbyhole/not-real.');
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Try going back to the root and navigating from there. Go back');
  });

  test('it handles when the API returns a permission denied error at the cluster level', async function (assert) {
    // Mock an endpoint to return 403
    this.server.get('/sys/internal/ui/mounts', () => overrideResponse(403));

    await visit('/vault/secrets-engines');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.error',
      'it bubbles up the error to the cluster error substate'
    );
    assert.dom(GENERAL.pageError.title(403)).hasText('ERROR 403 Not authorized');
    assert
      .dom(GENERAL.pageError.message)
      .hasText(
        'You are not authorized to access content at /v1/sys/internal/ui/mounts.',
        'message includes the API request url NOT the browser URL'
      );
    assert
      .dom(GENERAL.pageError.error)
      .hasTextContaining('Double check the URL or return to the dashboard.', 'it renders footer message');
    assert.dom(`${GENERAL.pageError.error} a`).hasText('Go to dashboard', 'Dashboard link renders');
    assert.dom(`${GENERAL.pageError.error} a`).hasAttribute('href', '/');
  });
});
