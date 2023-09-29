/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, settled, visit, fillIn, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

const shell = create(consoleClass);

const createNS = async (name) => shell.runCommands(`write sys/namespaces/${name} -force`);

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it clears namespaces when you log out', async function (assert) {
    const ns = 'foo';
    await createNS(ns);
    await shell.runCommands(`write -field=client_token auth/token/create policies=default`);
    const token = shell.lastLogOutput;
    await logout.visit();
    await authPage.login(token);
    await click('[data-test-namespace-toggle]');
    assert.dom('[data-test-current-namespace]').hasText('root', 'root renders as current namespace');
    assert.dom('[data-test-namespace-link]').doesNotExist('Additional namespace have been cleared');
  });

  test('it shows nested namespaces if you log in with a namespace starting with a /', async function (assert) {
    assert.expect(5);

    await click('[data-test-namespace-toggle]');

    const nses = ['beep', 'boop', 'bop'];
    for (const [i, ns] of nses.entries()) {
      await createNS(ns);
      await settled();
      // the namespace path will include all of the namespaces up to this point
      const targetNamespace = nses.slice(0, i + 1).join('/');
      const url = `/vault/secrets?namespace=${targetNamespace}`;
      // this is usually triggered when creating a ns in the form -- trigger a reload of the namespaces manually
      await click('[data-test-refresh-namespaces]');
      // check that the single namespace "beep" or "boop" not "beep/boop" shows in the toggle display
      assert
        .dom(`[data-test-namespace-link="${targetNamespace}"]`)
        .hasText(ns, `shows the namespace ${ns} in the toggle component`);
      // because quint does not like page reloads, visiting url directing instead of clicking on namespace in toggle
      await visit(url);
    }

    await logout.visit();
    await settled();
    await authPage.visit({ namespace: '/beep/boop' });
    await settled();
    await authPage.tokenInput('root').submit();
    await settled();
    await click('[data-test-namespace-toggle]');

    assert.dom('[data-test-current-namespace]').hasText('/beep/boop/', 'current namespace begins with a /');
    assert
      .dom('[data-test-namespace-link="beep/boop/bop"]')
      .exists('renders the link to the nested namespace');
  });

  test('it shows the regular namespace toolbar when not managed', async function (assert) {
    // This test is the opposite of the test in managed-namespace-test
    await logout.visit();
    assert.strictEqual(currentURL(), '/vault/auth?with=token', 'Does not redirect');
    assert.dom('[data-test-namespace-toolbar]').exists('Normal namespace toolbar exists');
    assert
      .dom('[data-test-managed-namespace-toolbar]')
      .doesNotExist('Managed namespace toolbar does not exist');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Root)');
    await fillIn('input#namespace', '/foo');
    const encodedNamespace = encodeURIComponent('/foo');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}&with=token`,
      'Does not prepend root to namespace'
    );
  });
});
