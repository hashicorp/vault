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
    assert.dom('[data-test-error]').exists({ count: 1 }, 'renders error messages on save');
    find('.CodeMirror').get(0).CodeMirror.setValue(policyString);
  });
  click('[data-test-policy-save]');
  andThen(function() {
    assert.equal(
      currentURL(),
      `/vault/policy/acl/${encodeURIComponent(policyLower)}`,
      'navigates to policy show on successful save'
    );
    assert.dom('[data-test-policy-name]').hasText(policyLower, 'displays the policy name on the show page');
    assert.dom('[data-test-flash-message].is-info').doesNotExist('no flash message is displayed on save');
  });
  click('[data-test-policy-list-link]');
  andThen(function() {
    assert
      .dom(`[data-test-policy-link="${policyLower}"]`)
      .exists({ count: 1 }, 'new policy shown in the list');
  });

  // policy deletion
  click(`[data-test-policy-link="${policyLower}"]`);
  click('[data-test-policy-edit-toggle]');
  click('[data-test-policy-delete] button');
  click('[data-test-confirm-button]');
  andThen(function() {
    assert.equal(currentURL(), `/vault/policies/acl`, 'navigates to policy list on successful deletion');
    assert
      .dom(`[data-test-policy-item="${policyLower}"]`)
      .doesNotExist('deleted policy is not shown in the list');
  });
});

// https://github.com/hashicorp/vault/issues/4395
test('it properly fetches policies when the name ends in a ,', function(assert) {
  const now = new Date().getTime();
  const policyString = 'path "*" { capabilities = ["update"]}';
  const policyName = `${now}-symbol,.`;

  page.visit({ type: 'acl' });
  // new policy creation
  click('[data-test-policy-create-link]');
  fillIn('[data-test-policy-input="name"]', policyName);
  andThen(function() {
    find('.CodeMirror').get(0).CodeMirror.setValue(policyString);
  });
  click('[data-test-policy-save]');
  andThen(function() {
    assert.equal(
      currentURL(),
      `/vault/policy/acl/${policyName}`,
      'navigates to policy show on successful save'
    );
    assert.dom('[data-test-policy-edit-toggle]').exists({ count: 1 }, 'shows the edit toggle');
  });
});
