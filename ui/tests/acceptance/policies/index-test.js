/**
 * Copyright (c) HashiCorp, Inc.
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

import authPage from 'vault/tests/pages/auth';
import { runCmd } from 'vault/tests/helpers/commands';
import codemirror from 'vault/tests/helpers/codemirror';

const SELECT = {
  policyByName: (name) => `[data-test-policy-link="${name}"]`,
  filterBar: '[data-test-component="navigate-input"]',
  delete: '[data-test-confirm-action-trigger]',
  confirmDelete: '[data-test-confirm-button]',
  createLink: '[data-test-policy-create-link]',
  nameInput: '[data-test-policy-input="name"]',
  save: '[data-test-policy-save]',
  createError: '[data-test-message-error]',
  policyTitle: '[data-test-policy-name]',
  listBreadcrumb: '[data-test-policy-list-link] a',
};
module('Acceptance | policies/acl', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
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
    await click(SELECT.delete);
    await click(SELECT.confirmDelete);
    assert.dom(SELECT.policyByName(policyName)).doesNotExist('policy is deleted successfully');
  });

  // https://github.com/hashicorp/vault/issues/4395
  test('it properly fetches policies when the name ends in a ,', async function (assert) {
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `${this.uid}-policy-symbol,.`;

    await visit('/vault/policies/acl');
    // new policy creation
    await click(SELECT.createLink);

    await fillIn(SELECT.nameInput, policyName);
    codemirror().setValue(policyString);
    await click(SELECT.save);
    assert.strictEqual(
      currentURL(),
      `/vault/policy/acl/${policyName}`,
      'navigates to policy show on successful save'
    );
    assert.dom('[data-test-policy-edit-toggle]').exists({ count: 1 }, 'shows the edit toggle');
  });

  test('it can create and delete correctly', async function (assert) {
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `Policy test ${this.uid}`;
    const policyLower = policyName.toLowerCase();

    await visit('/vault/policies/acl');
    // new policy creation
    await click(SELECT.createLink);

    await fillIn(SELECT.nameInput, policyName);
    await click(SELECT.save);
    assert
      .dom(SELECT.createError)
      .hasText(`Error 'policy' parameter not supplied or empty`, 'renders error message on save');
    codemirror().setValue(policyString);
    await click(SELECT.save);

    await waitUntil(() => currentURL() === `/vault/policy/acl/${encodeURIComponent(policyLower)}`);
    assert.strictEqual(
      currentURL(),
      `/vault/policy/acl/${encodeURIComponent(policyLower)}`,
      'navigates to policy show on successful save'
    );
    assert.dom(SELECT.policyTitle).hasText(policyLower, 'displays the policy name on the show page');
    // will fail if you have a license about to expire.
    assert.dom('[data-test-flash-message].is-info').doesNotExist('no flash message is displayed on save');
    await click(SELECT.listBreadcrumb);

    assert.strictEqual(currentURL(), `/vault/policies/acl`, 'navigates to policy list from breadcrumb');
    // List of policies can get long quickly -- filter for the policy to make the test more robust
    await fillIn(SELECT.filterBar, policyLower);
    assert
      .dom(`[data-test-policy-link="${policyLower}"]`)
      .exists({ count: 1 }, 'new policy shown in the list');

    // policy deletion
    await click(SELECT.policyByName(policyLower));

    await click('[data-test-policy-edit-toggle]');

    await click('[data-test-confirm-action-trigger]');

    await click('[data-test-confirm-button]');
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
