/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, visit as _visit, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import auth from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const visit = async (url) => {
  try {
    await _visit(url);
  } catch (e) {
    if (e.message !== 'TransitionAborted') {
      throw e;
    }
  }

  await settled();
};

const consoleComponent = create(consoleClass);

const wrappedAuth = async () => {
  await consoleComponent.toggle();
  await settled();
  await consoleComponent.runCommands(
    `write -field=token auth/token/create policies=default -wrap-ttl=5m`,
    false
  );
  await settled();
  // because of flaky test, trying to capture the token using a dom selector instead of the page object
  const token = document.querySelector('[data-test-component="console/log-text"] pre').textContent;
  if (token.includes('Error')) {
    throw new Error(`Error mounting secrets engine: ${token}`);
  }
  return token;
};

const setupWrapping = async () => {
  await auth.logout();
  await settled();
  await auth.visit();
  await settled();
  await auth.authType('token');
  await auth.tokenInput('root').submit();
  await settled();
  const wrappedToken = await wrappedAuth();
  return wrappedToken;
};
module('Acceptance | redirect_to query param functionality', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    // normally we'd use the auth.logout helper to visit the route and reset the app, but in this case that
    // also routes us to the auth page, and then all of the transitions from the auth page get redirected back
    // to the auth page resulting in no redirect_to query param being set
    localStorage.clear();
  });
  test('redirect to a route after authentication', async function (assert) {
    const url = '/vault/secrets/secret/kv/create';
    await visit(url);

    assert.ok(
      currentURL().includes(`redirect_to=${encodeURIComponent(url)}`),
      `encodes url for the query param in ${currentURL()}`
    );
    await auth.authType('token');
    // the login method on this page does another visit call that we don't want here
    await auth.tokenInput('root').submit();
    await settled();
    assert.strictEqual(currentURL(), url, 'navigates to the redirect_to url after auth');
  });

  test('redirect from root does not include redirect_to', async function (assert) {
    const url = '/';
    await visit(url);
    assert.ok(currentURL().indexOf('redirect_to') < 0, 'there is no redirect_to query param');
  });

  test('redirect to a route after authentication with a query param', async function (assert) {
    const url = '/vault/secrets/secret/kv/create?initialKey=hello';
    await visit(url);
    assert.ok(
      currentURL().includes(`?redirect_to=${encodeURIComponent(url)}`),
      'encodes url for the query param'
    );
    await auth.authType('token');
    await auth.tokenInput('root').submit();
    await settled();
    assert.strictEqual(currentURL(), url, 'navigates to the redirect_to with the query param after auth');
  });

  test('redirect to logout with wrapped token authenticates you', async function (assert) {
    const wrappedToken = await setupWrapping();
    const url = '/vault/secrets/cubbyhole/create';

    await auth.logout({
      redirect_to: url,
      wrapped_token: wrappedToken,
    });
    await settled();

    assert.strictEqual(currentURL(), url, 'authenticates then navigates to the redirect_to url after auth');
  });
});
