import { currentURL, visit as _visit, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import auth from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const visit = async url => {
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
  await consoleComponent.runCommands(`write -field=token auth/token/create policies=default -wrap-ttl=3m`);
  return consoleComponent.lastLogOutput;
};

const setupWrapping = async () => {
  await auth.logout();
  await auth.visit();
  await auth.tokenInput('root').submit();
  let wrappedToken = await wrappedAuth();
  return wrappedToken;
};
module('Acceptance | redirect_to query param functionality', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    // normally we'd use the auth.logout helper to visit the route and reset the app, but in this case that
    // also routes us to the auth page, and then all of the transitions from the auth page get redirected back
    // to the auth page resulting in no redirect_to query param being set
    localStorage.clear();
  });
  test('redirect to a route after authentication', async function(assert) {
    let url = '/vault/secrets/secret/create';
    await visit(url);
    assert.equal(
      currentURL(),
      `/vault/auth?redirect_to=${encodeURIComponent(url)}&with=token`,
      'encodes url for the query param'
    );
    // the login method on this page does another visit call that we don't want here
    await auth.tokenInput('root').submit();
    assert.equal(currentURL(), url, 'navigates to the redirect_to url after auth');
  });

  test('redirect from root does not include redirect_to', async function(assert) {
    let url = '/';
    await visit(url);
    assert.equal(currentURL(), `/vault/auth`, 'there is no redirect_to query param');
  });

  test('redirect to a route after authentication with a query param', async function(assert) {
    let url = '/vault/secrets/secret/create?initialKey=hello';
    await visit(url);
    assert.equal(
      currentURL(),
      `/vault/auth?redirect_to=${encodeURIComponent(url)}`,
      'encodes url for the query param'
    );
    await auth.tokenInput('root').submit();
    assert.equal(currentURL(), url, 'navigates to the redirect_to with the query param after auth');
  });

  test('redirect to logout with wrapped token authenticates you', async function(assert) {
    let wrappedToken = await setupWrapping();
    let url = '/vault/secrets/cubbyhole/create';

    await auth.logout({
      redirect_to: url,
      wrapped_token: wrappedToken,
    });

    assert.equal(currentURL(), url, 'authenticates then navigates to the redirect_to url after auth');
  });
});
