import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { settled, currentURL } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import auth from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

const wrappedAuth = async () => {
  await consoleComponent.runCommands(`write -field=token auth/token/create policies=default -wrap-ttl=3m`);
  await settled();
  return consoleComponent.lastLogOutput;
};

const setupWrapping = async () => {
  await auth.logout();
  await settled();
  await auth.visit();
  await settled();
  await auth.tokenInput('root').submit();
  await settled();
  const token = await wrappedAuth();
  await auth.logout();
  await settled();
  return token;
};
module('Acceptance | wrapped_token query param functionality', function (hooks) {
  setupApplicationTest(hooks);

  test('it authenticates you if the query param is present', async function (assert) {
    const token = await setupWrapping();
    await auth.visit({ wrapped_token: token });
    await settled();
    assert.strictEqual(currentURL(), '/vault/secrets', 'authenticates and redirects to home');
  });

  test('it authenticates when used with the with=token query param', async function (assert) {
    const token = await setupWrapping();
    await auth.visit({ wrapped_token: token, with: 'token' });
    await settled();
    assert.strictEqual(currentURL(), '/vault/secrets', 'authenticates and redirects to home');
  });
});
