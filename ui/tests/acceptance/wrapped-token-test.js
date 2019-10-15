import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import auth from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

const wrappedAuth = async () => {
  await consoleComponent.runCommands(`write -field=token auth/token/create policies=default -wrap-ttl=3m`);
  return consoleComponent.lastLogOutput;
};

const setupWrapping = async () => {
  await auth.logout();
  await auth.visit();
  await auth.tokenInput('root').submit();
  let token = await wrappedAuth();
  await auth.logout();
  return token;
};
module('Acceptance | wrapped_token query param functionality', function(hooks) {
  setupApplicationTest(hooks);

  test('it authenticates you if the query param is present', async function(assert) {
    let token = await setupWrapping();
    await auth.visit({ wrapped_token: token });
    assert.equal(currentURL(), '/vault/secrets', 'authenticates and redirects to home');
  });

  test('it authenticates when used with the with=token query param', async function(assert) {
    let token = await setupWrapping();
    await auth.visit({ wrapped_token: token, with: 'token' });
    assert.equal(currentURL(), '/vault/secrets', 'authenticates and redirects to home');
  });
});
