/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import hbs from 'htmlbars-inline-precompile';
import { click, fillIn, render, findAll } from '@ember/test-helpers';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import PolicyForm from 'vault/forms/policy';

module('Integration | Component | policy-form-modal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.form = new PolicyForm({ name: 'modal-test', enforcement_level: 'hard-mandatory' });
    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();

    this.renderComponent = () =>
      render(hbs`
      <PolicyFormModal
        @form={{this.form}}
        @onSave={{this.onSave}}
        @onCancel={{this.onCancel}}
      />
    `);
  });

  test('it should render empty state when type is not selected', async function (assert) {
    await this.renderComponent();
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No policy type selected', 'shows empty state when no type is selected');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Select a policy type to continue creating.',
        'shows empty state message when no type is selected'
      );
  });

  test('it should render correct policy types in select', async function (assert) {
    await this.renderComponent();
    const options = findAll(`${GENERAL.selectByAttr('policyType')} option`).map((option) => option.value);
    assert.deepEqual(options, ['', 'acl', 'rgp'], 'renders correct policy types in select dropdown');
  });

  test('it should disable rgp without sentinel feature', async function (assert) {
    sinon.stub(this.owner.lookup('service:version'), 'hasSentinel').value(false);
    await this.renderComponent();
    const options = findAll(`${GENERAL.selectByAttr('policyType')} option`);
    assert.true(options[2].disabled, 'RGP option is disabled without sentinel feature');
  });

  test('it should render tabs when type is selected', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.selectByAttr('policyType'), 'acl');
    assert.dom('[data-test-tab-your-policy]').exists('renders policy tab when type is selected');
    assert.dom('[data-test-policy-form]').exists('renders policy form when tab is active');

    await click('[data-test-tab-example-policy]');
    assert.dom('[data-test-policy-example]').exists('renders policy example when tab is selected');
  });

  test('it should render form correctly based on policy type', async function (assert) {
    await this.renderComponent();
    await fillIn(GENERAL.selectByAttr('policyType'), 'acl');
    assert.dom('[data-test-policy-form]').exists('renders policy form when type is selected');
    assert
      .dom(GENERAL.inputByAttr('enforcement_level'))
      .doesNotExist('enforcement level input is not rendered for ACL policy');

    this.form.policyType = null;
    await this.renderComponent();
    await fillIn(GENERAL.selectByAttr('policyType'), 'rgp');
    assert
      .dom(GENERAL.inputByAttr('enforcement_level'))
      .exists('enforcement level input is rendered for RGP policy');
  });
});
