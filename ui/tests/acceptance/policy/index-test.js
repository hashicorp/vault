/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  currentURL,
  currentRouteName,
  settled,
  fillIn,
  visit,
  click,
  waitFor,
  waitUntil,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import localStorage from 'vault/lib/local-storage';

const SELECT = {
  policyByName: (name) => `[data-test-policy-link="${name}"]`,
  filterBar: '[data-test-component="navigate-input"]',
  createPolicy: '[data-test-policy-create-link]',
  policyTitle: '[data-test-policy-name]',
  listBreadcrumb: '[data-test-policy-list-link] a',
};

module('Acceptance | policies/acl', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    await login();
    // dismiss wizard
    localStorage.setItem('dismissed-wizards', ['acl-policy']);
  });

  test('it lists default and root acls', async function (assert) {
    await visit('/vault/policies/acl');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
    await fillIn(SELECT.filterBar, 'default');
    await waitFor(SELECT.policyByName('default'));
    assert.dom(SELECT.policyByName('default')).exists('default policy shown in the list');
    await fillIn(SELECT.filterBar, 'root');
    // root isn't clickable so it has a different selector
    assert.dom('[data-test-policy-name]').hasText('root', 'root policy shown in the list');
  });

  test('it navigates to show when clicking on the link', async function (assert) {
    await visit('/vault/policies/acl');
    await fillIn(SELECT.filterBar, 'default');
    await waitFor(SELECT.policyByName('default'));
    await click(SELECT.policyByName('default'));
    assert.strictEqual(currentRouteName(), 'vault.cluster.policy.show');
    assert.strictEqual(currentURL(), '/vault/policy/acl/default');
  });

  test('it allows deletion of policies with dots in names', async function (assert) {
    const POLICY = 'path "*" { capabilities = ["list"]}';
    const policyName = 'list.policy';
    await runCmd(`write sys/policies/acl/${policyName} policy=${window.btoa(POLICY)}`);
    await settled();
    await visit('/vault/policies/acl');
    await fillIn(SELECT.filterBar, policyName);
    await waitFor(SELECT.policyByName(policyName));
    assert.dom(SELECT.policyByName(policyName)).exists('policy is shown in list');
    await click(`${SELECT.policyByName(policyName)} [data-test-popup-menu-trigger]`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.dom(SELECT.policyByName(policyName)).doesNotExist('policy is deleted successfully');
  });

  // https://github.com/hashicorp/vault/issues/4395
  test('it properly fetches policies when the name ends in a ,', async function (assert) {
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `${this.uid}-policy-symbol,.`;

    await visit('/vault/policies/acl');
    // new policy creation
    await click(SELECT.createPolicy);

    await fillIn(GENERAL.inputByAttr('name'), policyName);
    await click(GENERAL.radioByAttr('code'));
    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, policyString);

    await click(GENERAL.submitButton);
    assert.strictEqual(
      currentURL(),
      `/vault/policy/acl/${policyName}`,
      'navigates to policy show on successful save'
    );
    assert.dom(GENERAL.button('Edit policy')).exists({ count: 1 }, 'shows the edit toggle');
  });

  test('it can create and delete correctly', async function (assert) {
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `Policy test ${this.uid}`;
    const policyLower = policyName.toLowerCase();

    await visit('/vault/policies/acl');
    // new policy creation
    await click(SELECT.createPolicy);

    await fillIn(GENERAL.inputByAttr('name'), policyName);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText(`Error 'policy' parameter not supplied or empty`, 'renders error message on save');
    await click(GENERAL.radioByAttr('code'));
    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, policyString);

    await click(GENERAL.submitButton);

    await waitUntil(() => currentURL() === `/vault/policy/acl/${encodeURIComponent(policyLower)}`);
    assert.strictEqual(
      currentURL(),
      `/vault/policy/acl/${encodeURIComponent(policyLower)}`,
      'navigates to policy show on successful save'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(policyLower, 'displays the policy name on the show page');
    assert.dom(GENERAL.latestFlashContent).hasText(`ACL policy "${policyLower}" was successfully created.`);
    await click(GENERAL.breadcrumbAtIdx(1));

    assert.strictEqual(currentURL(), `/vault/policies/acl`, 'navigates to policy list from breadcrumb');
    // List of policies can get long quickly -- filter for the policy to make the test more robust
    await fillIn(SELECT.filterBar, policyLower);
    assert
      .dom(`[data-test-policy-link="${policyLower}"]`)
      .exists({ count: 1 }, 'new policy shown in the list');

    // policy deletion
    await click(SELECT.policyByName(policyLower));
    await click(GENERAL.button('Edit policy'));
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    await waitUntil(() => currentURL() === `/vault/policies/acl`);
    assert.strictEqual(
      currentURL(),
      `/vault/policies/acl`,
      'navigates to policy list on successful deletion'
    );
    assert
      .dom(`[data-test-policy-item="${policyLower}"]`)
      .doesNotExist('deleted policy is not shown in the list');
  });
});
