/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  settled,
  visit,
  fillIn,
  currentURL,
  waitFor,
  findAll,
  triggerKeyEvent,
  find,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { runCmd, createNS } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from '../helpers/general-selectors';
import { NAMESPACE_PICKER_SELECTORS } from '../helpers/namespace-picker';

import sinon from 'sinon';

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  let fetchSpy;

  hooks.beforeEach(() => {
    fetchSpy = sinon.spy(window, 'fetch');
    return login();
  });

  hooks.afterEach(() => {
    fetchSpy.restore();
  });

  test('it focuses the search input field when the component is loaded', async function (assert) {
    assert.expect(1);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the search input field is focused
    const searchInput = find(NAMESPACE_PICKER_SELECTORS.searchInput);
    assert.strictEqual(
      document.activeElement,
      searchInput,
      'The search input field is focused on component load'
    );
  });

  test('it navigates to the matching namespace when Enter is pressed', async function (assert) {
    assert.expect(2);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Simulate typing into the search input
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep/boop');
    assert
      .dom(NAMESPACE_PICKER_SELECTORS.searchInput)
      .hasValue('beep/boop', 'The search input field has the correct value');

    // Simulate pressing Enter
    await triggerKeyEvent(NAMESPACE_PICKER_SELECTORS.searchInput, 'keydown', 'Enter');

    // Verify navigation to the matching namespace
    assert.strictEqual(
      this.owner.lookup('service:router').currentURL,
      '/vault/dashboard?namespace=beep%2Fboop',
      'Navigates to the correct namespace when Enter is pressed'
    );
  });

  test('it filters namespaces based on search input', async function (assert) {
    assert.expect(7);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify all namespaces are displayed initially
    assert.dom(NAMESPACE_PICKER_SELECTORS.link()).exists('Namespace link(s) exist');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      5,
      'All namespaces are displayed initially'
    );

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 3 namespaces are displayed after searching for "beep"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      3,
      'Display 3 namespaces matching "beep" after searching'
    );

    // Verify 1 namespace is displayed after searching for "bop"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'bop');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      1,
      'Display 1 namespace matching "bop" after searching'
    );

    // Verify no namespaces are displayed after searching for "other"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'other');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      'No namespaces are displayed after searching for "other"'
    );

    // Clear the search input & verify all namespaces are displayed again
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, '');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      5,
      'All namespaces are displayed after clearing search input'
    );
  });

  test('it updates the namespace list after clicking "Refresh list"', async function (assert) {
    assert.expect(3);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the namespace list was fetched on load
    let listNamespaceRequests = fetchSpy
      .getCalls()
      .filter((call) => call.args[0].includes('/v1/sys/internal/ui/namespaces'));
    assert.strictEqual(
      listNamespaceRequests.length,
      1,
      'The network call to the specific endpoint was made twice (once on load, once on refresh)'
    );

    // Refresh the list of namespaces
    assert.dom(NAMESPACE_PICKER_SELECTORS.refreshList).exists('Refresh list button exists');
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Verify that the namespace list was fetched on refresh
    listNamespaceRequests = fetchSpy
      .getCalls()
      .filter((call) => call.args[0].includes('/v1/sys/internal/ui/namespaces'));
    assert.strictEqual(
      listNamespaceRequests.length,
      2,
      'The network call to the specific endpoint was made twice (once on load, once on refresh)'
    );
  });

  test('it displays the "Manage" button with the correct URL', async function (assert) {
    assert.expect(1);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify the "Manage" button is rendered and has the correct URL
    assert
      .dom('[href="/ui/vault/access/namespaces"]')
      .exists('The "Manage" button is displayed with the correct URL');
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
