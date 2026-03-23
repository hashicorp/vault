/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const SELECTORS = {
  intro: '[data-test-intro]',
  policyByName: (name) => `[data-test-policy-link="${name}"]`,
};

const policiesMockModel = [
  {
    id: 'default',
    name: 'default',
    policy: undefined,
    policyType: 'acl',
    canEdit: true,
    canRead: true,
  },
  {
    id: 'root',
    name: 'root',
    policy: undefined,
    policyType: 'acl',
  },
];
policiesMockModel.meta = {
  currentPage: 1,
  lastPage: 1,
  nextPage: 1,
  prevPage: 0,
  total: 2,
  filteredTotal: 2,
  pageSize: 15,
};

const customPoliciesMockModel = [
  {
    id: 'default',
    name: 'default',
    policy: undefined,
    policyType: 'acl',
    canEdit: true,
    canRead: true,
  },
  {
    id: 'root',
    name: 'root',
    policy: undefined,
    policyType: 'acl',
  },
  {
    id: 'custom',
    name: 'custom',
    policy: undefined,
    policyType: 'acl',
    canDelete: true,
    canEdit: true,
    canRead: true,
  },
];
customPoliciesMockModel.meta = {
  currentPage: 1,
  lastPage: 1,
  nextPage: 1,
  prevPage: 0,
  total: 3,
  filteredTotal: 3,
  pageSize: 15,
};

module('Integration | Component | page/policies', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = policiesMockModel;
    this.policyType = 'acl';
    this.store = this.owner.lookup('service:store');
    this.wizardService = this.owner.lookup('service:wizard');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.refreshSpy = sinon.spy();
    this.version = this.owner.lookup('service:version');

    // Stub flash message methods
    sinon.stub(this.flashMessages, 'success');
    sinon.stub(this.flashMessages, 'danger');

    this.renderComponent = () => {
      return render(hbs`
        <Page::Policies
          @filter={{null}}
          @model={{this.model}}
          @policyType={{this.policyType}}
          @onRefresh={{this.refreshSpy}}
        />
      `);
    };
  });

  hooks.afterEach(async function () {
    // Ensure clean state
    localStorage.clear();
  });

  test('it shows wizard intro page initially when only default policies exist', async function (assert) {
    await this.renderComponent();
    assert.false(this.wizardService.isDismissed('acl-policy'), 'Wizard is not dismissed initially');
    assert.dom(SELECTORS.intro).exists('ACL Policies intro page is rendered');
  });

  test('it does not show intro page when dismissed', async function (assert) {
    // Dismiss the wizard before rendering
    this.wizardService.dismiss('acl-policy');
    await this.renderComponent();

    assert.true(this.wizardService.isDismissed('acl-policy'), 'Wizard is marked as dismissed');
    assert.dom(SELECTORS.intro).doesNotExist('ACL Policies intro page is not rendered');
  });

  test('it does not show intro page when more than default policies exist', async function (assert) {
    this.model = customPoliciesMockModel;
    await this.renderComponent();
    assert.dom(SELECTORS.intro).doesNotExist('Intro page is not rendered when custom policies exist');
    assert.dom(GENERAL.button('intro')).doesNotExist('"New to ACL Policies?" button is not rendered');
  });

  test('it dismisses the intro page when clicking skip', async function (assert) {
    await this.renderComponent();

    assert.false(this.wizardService.isDismissed('acl-policy'), 'Intro page is not dismissed initially');
    assert.dom(SELECTORS.intro).exists('intro page is rendered');

    await click(GENERAL.button('Skip'));

    assert.true(
      this.wizardService.isDismissed('acl-policy'),
      'Wizard was marked as dismissed after clicking Skip'
    );
    assert.true(this.refreshSpy.calledOnce, 'onRefresh callback was called');
    assert.dom(SELECTORS.introButton).exists('"New to ACL Policies?" button is visible');
  });

  test('it re-renders the intro page as modal when clicking "New to ACL Policies?" button', async function (assert) {
    // Dismiss the wizard first
    this.wizardService.dismiss('acl-policy');
    await this.renderComponent();

    assert.dom(SELECTORS.intro).doesNotExist('ACL Policies intro page is not initially rendered');
    await click(GENERAL.button('intro'));

    assert.false(this.wizardService.isDismissed('acl-policy'), 'Wizard dismissal state was reset');
    assert.dom(SELECTORS.intro).exists('ACL Policies intro page is now rendered as modal');
  });

  test('it does not show the intro page or intro button for non-acl policy types', async function (assert) {
    this.policyType = 'rgp';

    await this.renderComponent();

    assert.dom(SELECTORS.intro).doesNotExist('ACL Policies intro page is not rendered for non-acl policies');
    assert
      .dom(GENERAL.button('intro'))
      .doesNotExist('"New to ACL Policies?" button is not shown for non-acl policies');
  });

  test('it successfully deletes an ACL policy', async function (assert) {
    this.model = customPoliciesMockModel;

    // Mock the API service to simulate successful deletion
    this.server.delete('/sys/policies/acl/:name', () => {
      return { data: null };
    });

    await this.renderComponent();
    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify success message was shown
    assert.ok(
      this.flashMessages.success.calledWith(`Successfully deleted policy: custom`),
      'Success flash message was displayed'
    );

    // Verify refresh was called
    assert.ok(this.refreshSpy.calledOnce, 'onRefresh callback was called after successful delete');
  });

  test('it handles error when deleting an ACL policy fails', async function (assert) {
    this.model = customPoliciesMockModel;

    // Mock the API service to simulate failed deletion
    this.server.delete('/v1/sys/policies/acl/:name', () => {
      return overrideResponse(500);
    });

    await this.renderComponent();
    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify error message was shown
    assert.ok(this.flashMessages.danger.calledOnce, 'Error flash message was displayed');

    // Verify refresh was NOT called on error
    assert.ok(this.refreshSpy.notCalled, 'onRefresh callback was not called after failed delete');
  });

  test('enterprise: it successfully deletes an RGP policy', async function (assert) {
    this.model = customPoliciesMockModel;
    this.policyType = 'rgp';
    this.version.type = 'enterprise';
    this.version.features = ['Sentinel'];

    // Mock the API service to simulate successful deletion
    this.server.delete('/sys/policies/rgp/:name', () => {
      return { data: null };
    });

    await this.renderComponent();

    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify success message was shown
    assert.ok(
      this.flashMessages.success.calledWith(`Successfully deleted policy: custom`),
      'Success flash message was displayed'
    );

    // Verify refresh was called
    assert.ok(this.refreshSpy.calledOnce, 'onRefresh callback was called after successful delete');
  });

  test('enterprise: it handles error when deleting an RGP policy fails', async function (assert) {
    this.model = customPoliciesMockModel;
    this.policyType = 'rgp';
    this.version.type = 'enterprise';
    this.version.features = ['Sentinel'];

    // Mock the API service to simulate failed deletion
    this.server.delete('/sys/policies/rgp/:name', () => {
      return overrideResponse(500);
    });

    await this.renderComponent();

    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify error message was shown
    assert.ok(this.flashMessages.danger.calledOnce, 'Error flash message was displayed');

    // Verify refresh was NOT called on error
    assert.ok(this.refreshSpy.notCalled, 'onRefresh callback was not called after failed delete');
  });

  test('enterprise: it successfully deletes an EGP policy', async function (assert) {
    this.model = customPoliciesMockModel;
    this.policyType = 'egp';
    this.version.type = 'enterprise';
    this.version.features = ['Sentinel'];

    // Mock the API service to simulate successful deletion
    this.server.delete('/sys/policies/egp/:name', () => {
      return { data: null };
    });

    await this.renderComponent();

    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify success message was shown
    assert.ok(
      this.flashMessages.success.calledWith(`Successfully deleted policy: custom`),
      'Success flash message was displayed'
    );

    // Verify refresh was called
    assert.ok(this.refreshSpy.calledOnce, 'onRefresh callback was called after successful delete');
  });

  test('enterprise: it handles error when deleting an EGP policy fails', async function (assert) {
    this.model = customPoliciesMockModel;
    this.policyType = 'egp';
    this.version.type = 'enterprise';
    this.version.features = ['Sentinel'];

    // Mock the API service to simulate failed deletion
    this.server.delete('/sys/policies/egp/:name', () => {
      return overrideResponse(500);
    });

    await this.renderComponent();

    await click(`${SELECTORS.policyByName('custom')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    // Verify error message was shown
    assert.ok(this.flashMessages.danger.calledOnce, 'Error flash message was displayed');

    // Verify refresh was NOT called on error
    assert.ok(this.refreshSpy.notCalled, 'onRefresh callback was not called after failed delete');
  });
});
