import { currentURL, currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import page from 'vault/tests/pages/policies/index';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

module('Acceptance | policies/acl', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it lists default and root acls', async function (assert) {
    await page.visit({ type: 'acl' });
    await settled();
    assert.strictEqual(currentURL(), '/vault/policies/acl');
    assert.ok(page.findPolicyByName('root'), 'root policy shown in the list');
    assert.ok(page.findPolicyByName('default'), 'default policy shown in the list');
  });

  test('it navigates to show when clicking on the link', async function (assert) {
    await page.visit({ type: 'acl' });
    await settled();
    await page.findPolicyByName('default').click();
    await settled();
    assert.strictEqual(currentRouteName(), 'vault.cluster.policy.show');
    assert.strictEqual(currentURL(), '/vault/policy/acl/default');
  });

  test('it allows deletion of policies with dots in names', async function (assert) {
    const POLICY = 'path "*" { capabilities = ["list"]}';
    const policyName = 'list.policy';
    await consoleComponent.runCommands([`write sys/policies/acl/${policyName} policy=${btoa(POLICY)}`]);
    await settled();
    await page.visit({ type: 'acl' });
    await settled();
    const policy = page.row.filterBy('name', policyName)[0];
    assert.ok(policy, 'policy is shown in the list');
    await policy.menu();
    await settled();
    await page.delete().confirmDelete();
    await settled();
    assert.notOk(page.findPolicyByName(policyName), 'policy is deleted successfully');
  });
});
