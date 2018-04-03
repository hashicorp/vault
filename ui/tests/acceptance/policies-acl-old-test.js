import Ember from 'ember';
import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/policies/index';

let adapterException;
let loggerError;
moduleForAcceptance('Acceptance | policies (old)', {
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
  },
});

test('policies', function(assert) {
  const now = new Date().getTime();
  const policyString = 'path "*" { capabilities = ["update"]}';
  const policyName = `Policy test ${now}`;
  const policyLower = policyName.toLowerCase();

  page.visit({ type: 'acl' });
  // new policy creation
  click('[data-test-policy-create-link]');
  fillIn('[data-test-policy-input="name"]', policyName);
  click('[data-test-policy-save]');
  andThen(function() {
    assert.equal(find('[data-test-error]').length, 1, 'renders error messages on save');
    find('.CodeMirror').get(0).CodeMirror.setValue(policyString);
  });
  click('[data-test-policy-save]');
  andThen(function() {
    assert.equal(
      currentURL(),
      `/vault/policy/acl/${encodeURIComponent(policyLower)}`,
      'navigates to policy show on successful save'
    );
    assert.equal(
      find('[data-test-policy-name]').text().trim(),
      policyLower,
      'displays the policy name on the show page'
    );
    assert.equal(
      find('[data-test-flash-message].is-info').length,
      0,
      'no flash message is displayed on save'
    );
  });
  click('[data-test-policy-list-link]');
  andThen(function() {
    assert.equal(find(`[data-test-policy-link="${policyLower}"]`).length, 1, 'new policy shown in the list');
  });

  // policy deletion
  click(`[data-test-policy-link="${policyLower}"]`);
  click('[data-test-policy-edit-toggle]');
  click('[data-test-policy-delete] button');
  click('[data-test-confirm-button]');
  andThen(function() {
    assert.equal(currentURL(), `/vault/policies/acl`, 'navigates to policy list on successful deletion');
    assert.equal(
      find(`[data-test-policy-item="${policyLower}"]`).length,
      0,
      'deleted policy is not shown in the list'
    );
  });
});
