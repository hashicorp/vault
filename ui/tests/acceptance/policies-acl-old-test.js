import { click, fillIn, findAll, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/policies/index';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | policies (old)', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('policies', async function(assert) {
    const now = new Date().getTime();
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `Policy test ${now}`;
    const policyLower = policyName.toLowerCase();

    await page.visit({ type: 'acl' });
    // new policy creation
    await click('[data-test-policy-create-link]');
    await fillIn('[data-test-policy-input="name"]', policyName);
    await click('[data-test-policy-save]');
    assert.dom('[data-test-error]').exists({ count: 1 }, 'renders error messages on save');
    findAll('.CodeMirror')[0].CodeMirror.setValue(policyString);
    await click('[data-test-policy-save]');
    assert.equal(
      currentURL(),
      `/vault/policy/acl/${encodeURIComponent(policyLower)}`,
      'navigates to policy show on successful save'
    );
    assert.dom('[data-test-policy-name]').hasText(policyLower, 'displays the policy name on the show page');
    assert.dom('[data-test-flash-message].is-info').doesNotExist('no flash message is displayed on save');
    await click('[data-test-policy-list-link]');
    assert
      .dom(`[data-test-policy-link="${policyLower}"]`)
      .exists({ count: 1 }, 'new policy shown in the list');

    // policy deletion
    await click(`[data-test-policy-link="${policyLower}"]`);
    await click('[data-test-policy-edit-toggle]');
    await click('[data-test-policy-delete] button');
    await click('[data-test-confirm-button]');
    assert.equal(currentURL(), `/vault/policies/acl`, 'navigates to policy list on successful deletion');
    assert
      .dom(`[data-test-policy-item="${policyLower}"]`)
      .doesNotExist('deleted policy is not shown in the list');
  });

  // https://github.com/hashicorp/vault/issues/4395
  test('it properly fetches policies when the name ends in a ,', async function(assert) {
    const now = new Date().getTime();
    const policyString = 'path "*" { capabilities = ["update"]}';
    const policyName = `${now}-symbol,.`;

    await page.visit({ type: 'acl' });
    // new policy creation
    await click('[data-test-policy-create-link]');
    await fillIn('[data-test-policy-input="name"]', policyName);
    findAll('.CodeMirror')[0].CodeMirror.setValue(policyString);
    await click('[data-test-policy-save]');
    assert.equal(
      currentURL(),
      `/vault/policy/acl/${policyName}`,
      'navigates to policy show on successful save'
    );
    assert.dom('[data-test-policy-edit-toggle]').exists({ count: 1 }, 'shows the edit toggle');
  });
});
