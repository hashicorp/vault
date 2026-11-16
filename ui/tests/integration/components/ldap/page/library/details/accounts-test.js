/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';
import { setupDetailsTest } from 'vault/tests/helpers/ldap/ldap-helpers';

module('Integration | Component | ldap | Page::Library::Details::Accounts', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);
  setupDetailsTest(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = () =>
      render(
        hbs`<Page::Library::Details::Accounts @library={{this.library}} @capabilities={{this.capabilities}} @statuses={{this.statuses}} />`,
        { owner: this.engine }
      );
  });

  test('it should render account cards', async function (assert) {
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();

    assert.dom('[data-test-account-name="foo.bar"]').hasText('foo.bar', 'Account name renders');
    assert
      .dom('[data-test-account-status="foo.bar"]')
      .hasText('Unavailable', 'Correct badge renders for checked out account');
    assert
      .dom('[data-test-account-status="bar.baz"]')
      .hasText('Available', 'Correct badge renders for available account');

    await click('[data-test-check-out]');
    await fillIn('[data-test-ttl-value="TTL"]', 4);
    await click('[data-test-check-out="save"]');

    const didTransition = transitionStub.calledWith(
      'vault.cluster.secrets.backend.ldap.libraries.library.check-out',
      { queryParams: { ttl: '4h' } }
    );
    assert.true(didTransition, 'Transitions to check out route on action click');

    assert.dom('[data-test-checked-out-card]').exists('Accounts checked out card renders');

    assert
      .dom(`${GENERAL.codeBlock('accounts')} code`)
      .hasText(
        'vault lease renew ldap-test/library/test-library/check-out/:lease_id',
        'Renew cli command renders with backend path'
      );
  });

  test('it should block a check-out TTL longer than Vault supports', async function (assert) {
    assert.expect(5);
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();
    await click('[data-test-check-out]');
    await fillIn(GENERAL.ttl.input('TTL'), '9999999999');
    await fillIn(TTL_PICKER.unitInputByLabel('TTL'), 'd');
    await click('[data-test-check-out="save"]');

    assert.dom(TTL_PICKER.error).hasText('TTL is too large. Maximum allowed value is 106,751 days.');
    assert.true(transitionStub.notCalled, 'does not transition to the check out route');
    assert.dom('[data-test-check-out="save"]').exists('the prompt stays open');

    await fillIn(GENERAL.ttl.input('TTL'), '106751');
    await click('[data-test-check-out="save"]');
    assert.true(transitionStub.calledOnce, 'transitions once the TTL is valid');
    assert.strictEqual(transitionStub.firstCall.args[1].queryParams.ttl, '106751d');
  });

  test('it should block a check-out TTL too long to be stored as a number', async function (assert) {
    assert.expect(2);
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();
    await click('[data-test-check-out]');
    await fillIn(GENERAL.ttl.input('TTL'), '9'.repeat(400));
    await fillIn(TTL_PICKER.unitInputByLabel('TTL'), 'd');
    await click('[data-test-check-out="save"]');

    assert.dom(TTL_PICKER.error).exists('the picker shows its error');
    assert.true(transitionStub.notCalled, 'does not transition to the check out route');
  });

  test('it should not carry an over-limit check-out TTL into the next prompt', async function (assert) {
    assert.expect(3);
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();
    await click('[data-test-check-out]');
    await fillIn(GENERAL.ttl.input('TTL'), '9999999999');
    await fillIn(TTL_PICKER.unitInputByLabel('TTL'), 'd');
    await click('[data-test-check-out="cancel"]');
    await click('[data-test-check-out]');

    assert.dom(TTL_PICKER.error).doesNotExist('the reopened prompt shows no error');
    await click('[data-test-check-out="save"]');
    assert.true(transitionStub.calledOnce, 'the library default TTL checks out');
    assert.strictEqual(transitionStub.firstCall.args[1].queryParams.ttl, null, 'no stale TTL is sent');
  });
});
