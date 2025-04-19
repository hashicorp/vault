/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, settled, visit, fillIn, currentURL, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { runCmd, createNS } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from '../helpers/general-selectors';
import { NAMESPACE_PICKER_SELECTORS } from '../helpers/namespace-picker';

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it clears namespaces when you log out', async function (assert) {
    const ns = 'foo';
    await runCmd(createNS(ns), false);
    const token = await runCmd(`write -field=client_token auth/token/create policies=default`);
    await login(token);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    assert.dom(NAMESPACE_PICKER_SELECTORS.link()).hasText('root', 'root renders as current namespace');
    assert
      .dom(`${NAMESPACE_PICKER_SELECTORS.link()} svg${GENERAL.icon('check')}`)
      .exists('The root namespace is selected');
  });

  // TODO: revisit test name/description, is this still relevant? A '/' prefix is stripped from namespace on login form
  test('it shows nested namespaces if you log in with a namespace starting with a /', async function (assert) {
    assert.expect(6);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    const nses = ['beep', 'boop', 'bop'];
    for (const [i, ns] of nses.entries()) {
      await runCmd(createNS(ns), false);
      await settled();

      // the namespace path will include all of the namespaces up to this point
      const targetNamespace = nses.slice(0, i + 1).join('/');
      const url = `/vault/secrets?namespace=${targetNamespace}`;

      // this is usually triggered when creating a ns in the form -- trigger a reload of the namespaces manually
      await click(NAMESPACE_PICKER_SELECTORS.toggle);

      // refresh the list of namespaces
      await waitFor(NAMESPACE_PICKER_SELECTORS.refreshList);
      await click(NAMESPACE_PICKER_SELECTORS.refreshList);

      // check that the full namespace path, like "beep/boop", shows in the toggle display
      await waitFor(NAMESPACE_PICKER_SELECTORS.link(targetNamespace));
      assert
        .dom(NAMESPACE_PICKER_SELECTORS.link(targetNamespace))
        .hasText(targetNamespace, `shows the namespace ${targetNamespace} in the toggle component`);

      // because quint does not like page reloads, visiting url directly instead of clicking on namespace in toggle
      await visit(url);
    }

    await loginNs('/beep/boop');
    await settled();

    // Open the namespace picker & wait for it to render
    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await waitFor(`svg${GENERAL.icon('check')}`);

    // Find the selected element with the check icon & ensure it exists
    const checkIcon = document.querySelector(
      `${NAMESPACE_PICKER_SELECTORS.link()} svg${GENERAL.icon('check')}`
    );
    assert.ok(checkIcon, 'A selected namespace link with the check icon exists');

    // Get the selected namespace with the data-test-namespace-link attribute & ensure it exists
    const selectedNamespace = checkIcon.closest(NAMESPACE_PICKER_SELECTORS.link());
    assert.ok(selectedNamespace, 'The selected namespace link exists');

    // Verify that the selected namespace has the correct data-test-namespace-link attribute and path value
    assert.strictEqual(
      selectedNamespace.getAttribute('data-test-namespace-link'),
      'beep/boop',
      'The current namespace does not begin or end with /'
    );
  });

  test('it shows the regular namespace toolbar when not managed', async function (assert) {
    // This test is the opposite of the test in managed-namespace-test
    await logout();
    assert.strictEqual(currentURL(), '/vault/auth?with=token', 'Does not redirect');
    assert.dom('[data-test-namespace-toolbar]').exists('Normal namespace toolbar exists');
    assert.dom(AUTH_FORM.managedNsRoot).doesNotExist('Managed namespace indicator does not exist');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Root)');
    await fillIn('input#namespace', '/foo/bar ');
    const encodedNamespace = encodeURIComponent('foo/bar');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}&with=token`,
      'correctly sanitizes namespace'
    );
  });
});
