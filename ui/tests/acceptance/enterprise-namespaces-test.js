import { click } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

const shell = create(consoleClass);

const createNS = async name => {
  await shell.runCommands(`write sys/namespaces/${name} -force`);
};

const switchToNS = async name => {
  await click('[data-test-namespace-toggle]');
  await click(`[data-test-namespace-link="${name}"]`);
  await click('[data-test-namespace-toggle]');
};

module('Acceptance | Enterprise | namespaces', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it clears namespaces when you log out', async function(assert) {
    let ns = 'foo';
    await createNS(ns);
    await shell.runCommands(`write -field=client_token auth/token/create policies=default`);
    let token = shell.lastLogOutput;
    await logout.visit();
    await authPage.login(token);
    assert.dom('[data-test-namespace-toggle]').doesNotExist('does not show the namespace picker');
    await logout.visit();
  });

  test('it shows nested namespaces if you log in with a namspace starting with a /', async function(assert) {
    let nses = ['beep', 'boop', 'bop'];
    for (let [i, ns] of nses.entries()) {
      await createNS(ns);
      // this is usually triggered when creating a ns in the form, here we'll trigger a reload of the
      // namespaces manually
      await this.owner.lookup('service:namespace').findNamespacesForUser.perform();
      if (i === nses.length - 1) {
        break;
      }
      // the namespace path will include all of the namespaces up to this point
      let targetNamespace = nses.slice(0, i + 1).join('/');
      await switchToNS(targetNamespace);
    }
    await logout.visit();
    await authPage.visit({ namespace: '/beep/boop' });
    await authPage.tokenInput('root').submit();
    await click('[data-test-namespace-toggle]');
    assert.dom('[data-test-current-namespace]').hasText('/beep/boop/', 'current namespace begins with a /');
    assert
      .dom('[data-test-namespace-link="beep/boop/bop"]')
      .exists('renders the link to the nested namespace');
  });
});
