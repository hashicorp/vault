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

const SELECTORS = {
  intro: '[data-test-intro]',
};

const defaultMethodModel = {
  methods: [
    {
      id: 'token/',
      type: 'token',
      path: 'token/',
      methodType: 'token',
      icon: 'token',
      accessor: 'auth_token_12345',
    },
  ],
  capabilities: {},
};

const multipleMethodsModel = {
  methods: [
    {
      id: 'token/',
      type: 'token',
      path: 'token/',
      methodType: 'token',
      icon: 'token',
      accessor: 'auth_token_12345',
    },
    {
      id: 'userpass/',
      type: 'userpass',
      path: 'userpass/',
      methodType: 'userpass',
      icon: 'user',
      accessor: 'auth_userpass_67890',
    },
  ],
  capabilities: {},
};

module('Integration | Component | page/methods', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Auth Methods' },
    ];
    this.model = defaultMethodModel;
    this.router = this.owner.lookup('service:router');
    this.refreshSpy = sinon.stub(this.router, 'refresh');
    this.wizardService = this.owner.lookup('service:wizard');

    this.renderComponent = () => {
      return render(hbs`
        <Page::Methods
          @model={{this.model}}
          @breadcrumbs={{this.breadcrumbs}}
        />
      `);
    };
  });

  hooks.afterEach(async function () {
    // Ensure clean state
    localStorage.clear();
  });

  test('it shows wizard intro page initially when only default method exists', async function (assert) {
    await this.renderComponent();
    assert.false(this.wizardService.isDismissed('auth-methods'), 'Wizard is not dismissed initially');
    assert.dom(SELECTORS.intro).exists('Auth Methods intro page is rendered');
  });

  test('it does not show intro page when dismissed', async function (assert) {
    // Dismiss the wizard before rendering
    this.wizardService.dismiss('auth-methods');
    await this.renderComponent();

    assert.true(this.wizardService.isDismissed('auth-methods'), 'Wizard is marked as dismissed');
    assert.dom(SELECTORS.intro).doesNotExist('Auth Methods intro page is not rendered');
  });

  test('it does not show intro page when more than default method exists', async function (assert) {
    this.model = multipleMethodsModel;
    await this.renderComponent();

    assert.dom(SELECTORS.intro).doesNotExist('Intro page is not rendered when custom methods exist');
    assert.dom(GENERAL.button('intro')).doesNotExist('"New to Auth Methods?" button is not rendered');
  });

  test('it dismisses the intro page when clicking skip', async function (assert) {
    await this.renderComponent();

    assert.false(this.wizardService.isDismissed('auth-methods'), 'Intro page is not dismissed initially');
    assert.dom(SELECTORS.intro).exists('intro page is rendered');

    await click(GENERAL.button('Skip'));

    assert.true(
      this.wizardService.isDismissed('auth-methods'),
      'Wizard was marked as dismissed after clicking Skip'
    );
    assert.true(this.refreshSpy.calledOnce, 'onRefresh callback was called');
    assert.dom(SELECTORS.introButton).exists('"New to ACL Policies?" button is visible');
  });

  test('it re-renders the intro page as modal when clicking "New to Auth Methods?" button', async function (assert) {
    // Dismiss the wizard first
    this.wizardService.dismiss('auth-methods');
    await this.renderComponent();

    assert.dom(SELECTORS.intro).doesNotExist('Auth Methods intro page is not initially rendered');
    await click(GENERAL.button('intro'));

    assert.false(this.wizardService.isDismissed('auth-methods'), 'Wizard dismissal state was reset');
    assert.dom(SELECTORS.intro).exists('Auth Methods intro page is now rendered as modal');
  });
});
