/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, find, findAll, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { AUTH_METHOD_MAP } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import {
  ALL_LOGIN_METHODS,
  BASE_LOGIN_METHODS,
  ENTERPRISE_LOGIN_METHODS,
} from 'vault/utils/supported-login-methods';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';

module('Integration | Component | auth | form template', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    window.localStorage.clear();
    this.version = this.owner.lookup('service:version');
    this.visibleMountsByType = null;
    this.canceledMfaAuth = '';
    this.cluster = { id: '1' };
    this.directLinkData = null;
    this.handleNamespaceUpdate = sinon.spy();
    this.loginSettings = null;
    this.namespaceQueryParam = '';
    this.oidcProviderQueryParam = '';
    this.onSuccess = sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
         <Auth::FormTemplate
          @canceledMfaAuth={{this.canceledMfaAuth}}
          @cluster={{this.cluster}}
          @directLinkData={{this.directLinkData}}
          @handleNamespaceUpdate={{this.handleNamespaceUpdate}}
          @loginSettings={{this.loginSettings}}
          @namespaceQueryParam={{this.namespaceQueryParam}}
          @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
          @onSuccess={{this.onSuccess}}
          @visibleMountsByType={{this.visibleMountsByType}}
        />`);
    };
  });

  // test to select each method is in "ent" module to include enterprise methods
  test('it selects token by default', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('token');
  });

  test('it selects @canceledMfaAuth by default', async function (assert) {
    this.canceledMfaAuth = 'ldap';
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('ldap');
    assert.dom(GENERAL.inputByAttr('username')).exists();
    assert.dom(GENERAL.inputByAttr('password')).exists();
  });

  test('it selects type in the dropdown if @directLinkData data just contains type', async function (assert) {
    this.directLinkData = { type: 'oidc', isVisibleMount: false };
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('oidc');
    assert.dom(GENERAL.inputByAttr('role')).exists();
    await click(AUTH_FORM.advancedSettings);
    assert.dom(GENERAL.inputByAttr('path')).exists();
    assert.dom(GENERAL.backButton).doesNotExist();
    assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
  });

  test('it does not show toggle buttons when listing visibility is not set', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render');
    assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
  });

  test('it displays errors', async function (assert) {
    const authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    authenticateStub.throws('permission denied');
    await this.renderComponent();
    await click(AUTH_FORM.login);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error Authentication failed: permission denied: Sinon-provided permission denied');
    authenticateStub.restore();
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.visibleMountsByType = {
        userpass: [
          {
            path: 'userpass/',
            description: '',
            options: {},
            type: 'userpass',
          },
          {
            path: 'userpass2/',
            description: '',
            options: {},
            type: 'userpass',
          },
        ],
        oidc: [
          {
            path: 'my-oidc/',
            description: '',
            options: {},
            type: 'oidc',
          },
        ],
        token: [
          {
            path: 'token/',
            description: 'token based credentials',
            options: null,
            type: 'token',
          },
        ],
      };
    });

    test('it renders mounts configured with listing_visibility="unuath"', async function (assert) {
      const expectedTabs = [
        { type: 'userpass', display: 'Userpass' },
        { type: 'oidc', display: 'OIDC' },
        { type: 'token', display: 'Token' },
      ];

      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      // there are 4 mount paths returned in the stubbed sys/internal/ui/mounts response above,
      // but two are of the same type so only expect 3 tabs
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'it groups mount paths by type and renders 3 tabs');
      expectedTabs.forEach((m) => {
        assert.dom(AUTH_FORM.tabBtn(m.type)).exists(`${m.type} renders as a tab`);
        assert.dom(AUTH_FORM.tabBtn(m.type)).hasText(m.display, `${m.type} renders expected display name`);
      });
      assert
        .dom(AUTH_FORM.tabBtn('userpass'))
        .hasAttribute('aria-selected', 'true', 'it selects the first type by default');
    });

    test('it selects each auth tab and renders form for that type', async function (assert) {
      await this.renderComponent();
      const assertSelected = (type) => {
        assert.dom(AUTH_FORM.authForm(type)).exists(`${type}: form renders when tab is selected`);
        assert.dom(AUTH_FORM.tabBtn(type)).hasAttribute('aria-selected', 'true');
      };
      const assertUnselected = (type) => {
        assert.dom(AUTH_FORM.authForm(type)).doesNotExist(`${type}: form does NOT render`);
        assert.dom(AUTH_FORM.tabBtn(type)).hasAttribute('aria-selected', 'false');
      };
      // click through each tab
      await click(AUTH_FORM.tabBtn('userpass'));
      assertSelected('userpass');
      assertUnselected('oidc');
      assertUnselected('token');
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();

      await click(AUTH_FORM.tabBtn('oidc'));
      assertSelected('oidc');
      assertUnselected('token');
      assertUnselected('userpass');
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();

      await click(AUTH_FORM.tabBtn('token'));
      assertSelected('token');
      assertUnselected('oidc');
      assertUnselected('userpass');
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
    });

    test('it renders the mount description', async function (assert) {
      await this.renderComponent();
      await click(AUTH_FORM.tabBtn('token'));
      assert.dom('section p').hasText('token based credentials');
    });

    test('it renders a dropdown if multiple mount paths are returned', async function (assert) {
      await this.renderComponent();
      await click(AUTH_FORM.tabBtn('userpass'));
      const dropdownOptions = findAll(`${GENERAL.selectByAttr('path')} option`).map((o) => o.value);
      const expectedPaths = ['userpass/', 'userpass2/'];
      expectedPaths.forEach((p) => {
        assert.true(dropdownOptions.includes(p), `dropdown includes path: ${p}`);
      });
    });

    test('it renders hidden input if only one mount path is returned', async function (assert) {
      await this.renderComponent();
      await click(AUTH_FORM.tabBtn('oidc'));
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
    });

    test('it clicks "Sign in with other methods" and renders standard form', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'tabs render by default');
      assert.dom(GENERAL.backButton).doesNotExist();
      await click(AUTH_FORM.otherMethodsBtn);
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('button disappears after it is clicked');
      assert
        .dom(GENERAL.selectByAttr('auth type'))
        .hasValue('token', 'dropdown renders and resets to select "Token"');
      await fillIn(GENERAL.selectByAttr('auth type'), 'userpass');
      assert.dom(AUTH_FORM.advancedSettings).exists('toggle renders even though userpass has visible mounts');
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
      assert.dom(GENERAL.inputByAttr('path')).hasValue('', 'it renders empty custom path input');
      await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
      assert.dom(AUTH_FORM.advancedSettings).exists('toggle renders even though oidc has a visible mount');
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
      assert.dom(GENERAL.inputByAttr('path')).hasValue('', 'it renders empty custom path input');
      await click(GENERAL.backButton);
      assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render after it is clicked');
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'clicking "Back" renders tabs again');
      assert.dom(AUTH_FORM.otherMethodsBtn).exists('"Sign in with other methods" renders again');
    });

    test('it resets selected tab after clicking "Sign in with other methods" and then "Back"', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'false');
      assert.dom(AUTH_FORM.tabBtn('token')).hasAttribute('aria-selected', 'false');

      // select a different tab before clicking "Sign in with other methods"
      await click(AUTH_FORM.tabBtn('oidc'));
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'true');
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'false');
      await click(AUTH_FORM.otherMethodsBtn);
      assert.dom(GENERAL.selectByAttr('auth type')).exists('it renders dropdown instead of tabs');
      await click(GENERAL.backButton);
      // assert tab selection is reset
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'false');
      assert.dom(AUTH_FORM.tabBtn('token')).hasAttribute('aria-selected', 'false');
    });

    test('it preselects tab if @canceledMfaAuth is a tab', async function (assert) {
      this.canceledMfaAuth = 'oidc';
      await this.renderComponent();
      assert.dom(AUTH_FORM.authForm('oidc')).exists('oidc form renders');
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'true');
    });

    test('if @canceledMfaAuth is NOT a tab, dropdown renders with type selected instead of tabs', async function (assert) {
      this.canceledMfaAuth = 'ldap';
      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('ldap');
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();

      assert.dom(GENERAL.backButton).exists('"Back" button renders');
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
    });

    // if mount data exists, the mount has listing_visibility="unauth"
    test('it renders single mount view instead of tabs if @directLinkData data exists and includes mount data', async function (assert) {
      this.directLinkData = { path: 'my-oidc/', type: 'oidc', isVisibleMount: true };
      await this.renderComponent();
      assert.dom(AUTH_FORM.preferredMethod('oidc')).hasText('OIDC', 'it renders mount type');
      assert.dom(GENERAL.inputByAttr('role')).exists();
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
      assert.dom(AUTH_FORM.otherMethodsBtn).exists('"Sign in with other methods" renders');

      assert.dom(AUTH_FORM.tabBtn('oidc')).doesNotExist('tab does not render');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist();
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
      assert.dom(GENERAL.backButton).doesNotExist();
    });

    test('it does not render tabs if @directLinkData data exists and just includes type', async function (assert) {
      // set a type that is NOT in a visible mount because mount data exists otherwise
      this.directLinkData = { type: 'ldap', isVisibleMount: false };
      await this.renderComponent();

      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('ldap', 'dropdown has type selected');
      assert.dom(AUTH_FORM.authForm('ldap')).exists();
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists();

      assert.dom(AUTH_FORM.preferredMethod('ldap')).doesNotExist('single mount view does not render');
      assert.dom(AUTH_FORM.tabBtn('ldap')).doesNotExist('tab does not render');
      assert
        .dom(GENERAL.backButton)
        .exists('back button renders because listing_visibility="unauth" for other mounts');
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
    });
  });

  module('community', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'community';
    });

    test('it does not render the namespace input on community', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
    });

    test('dropdown does not include enterprise methods', async function (assert) {
      const supported = BASE_LOGIN_METHODS.map((m) => m.type);
      const unsupported = ENTERPRISE_LOGIN_METHODS.map((m) => m.type);
      assert.expect(supported.length + unsupported.length);
      await this.renderComponent();
      const dropdownOptions = findAll(`${GENERAL.selectByAttr('auth type')} option`).map((o) => o.value);

      supported.forEach((m) => {
        assert.true(dropdownOptions.includes(m), `dropdown includes supported method: ${m}`);
      });
      unsupported.forEach((m) => {
        assert.false(dropdownOptions.includes(m), `dropdown does NOT include unsupported method: ${m}`);
      });
    });
  });

  // tests with "enterprise" in the title are filtered out from CE test runs
  // naming the module 'ent' so these tests still run on the CE repo
  module('ent', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
      this.version.features = ['Namespaces'];
      this.namespaceQueryParam = '';
    });

    test('it does not render the namespace input if version does not include feature', async function (assert) {
      this.version.features = [];
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
    });

    // in th ent module to test ALL supported login methods
    // iterating in tests should generally be avoided, but purposefully wanted to test the component
    // renders as expected as auth types change
    test('it selects each supported auth type and renders its form and relevant fields', async function (assert) {
      const fieldCount = AUTH_METHOD_MAP.map((m) => Object.keys(m.options.loginData).length);
      const sum = fieldCount.reduce((a, b) => a + b, 0);
      const methodCount = AUTH_METHOD_MAP.length;
      // 3 assertions per method, plus an assertion for each expected field
      assert.expect(3 * methodCount + sum); // count at time of writing is 40

      await this.renderComponent();
      for (const method of AUTH_METHOD_MAP) {
        const { authType, options } = method;

        const fields = Object.keys(options.loginData);
        await fillIn(GENERAL.selectByAttr('auth type'), authType);

        assert.dom(GENERAL.selectByAttr('auth type')).hasValue(authType), `${authType}: it selects type`;
        assert.dom(AUTH_FORM.authForm(authType)).exists(`${authType}: it renders form component`);

        // token is the only method that does not support a custom mount path
        if (authType !== 'token') {
          // jwt and oidc render the same component so the toggle remains open switching between those types
          const element = find(AUTH_FORM.advancedSettings);
          if (element.ariaExpanded === 'false') {
            await click(AUTH_FORM.advancedSettings);
          }
        }

        const assertion = authType === 'token' ? 'doesNotExist' : 'exists';
        assert.dom(GENERAL.inputByAttr('path'))[assertion](`${authType}: mount path input ${assertion}`);

        fields.forEach((field) => {
          assert.dom(GENERAL.inputByAttr(field)).exists(`${authType}: ${field} input renders`);
        });
      }
    });

    test('it disables namespace input when an oidc provider query param exists', async function (assert) {
      this.oidcProviderQueryParam = 'myprovider';
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
    });

    test('dropdown includes enterprise methods', async function (assert) {
      const supported = ALL_LOGIN_METHODS.map((m) => m.type);
      assert.expect(supported.length);
      await this.renderComponent();

      const dropdownOptions = findAll(`${GENERAL.selectByAttr('auth type')} option`).map((o) => o.value);
      supported.forEach((m) => {
        assert.true(dropdownOptions.includes(m), `dropdown includes supported method: ${m}`);
      });
    });

    test('it sets namespace for hvd managed clusters', async function (assert) {
      this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespaceQueryParam = 'admin/west-coast';
      await this.renderComponent();
      assert.dom(AUTH_FORM.managedNsRoot).hasValue('/admin');
      assert.dom(AUTH_FORM.managedNsRoot).hasAttribute('readonly');
      assert.dom(GENERAL.inputByAttr('namespace')).hasValue('/west-coast');
    });
  });

  // Login settings are an enterprise only feature but the component is version agnostic
  // because fetching login customizations happens on enterprise only in the route.
  /* 
  TEST CASES 
  All need to be tested with and without visible mounts (i.e. tuned with listing_visibility="unauth")
  1. default type set, backup types set 
  2. default type set, no backup types
  3. no default type, backup types set 
   */
  module('ent login settings', function (hooks) {
    hooks.beforeEach(function () {
      this.loginSettings = {
        defaultType: 'oidc',
        backupTypes: ['userpass', 'ldap'],
      };

      this.assertPathInput = async (assert, { isHidden = false, value = '' } = {}) => {
        // the path input can render behind the "Advanced settings" toggle or as a hidden input.
        // Assert it only renders once and is the expected input
        if (!isHidden) {
          await click(AUTH_FORM.advancedSettings);
          assert.dom(GENERAL.inputByAttr('path')).exists('it renders mount path input');
        }
        if (isHidden) {
          assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
          assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
          assert.dom(GENERAL.inputByAttr('path')).hasValue(value);
        }
        assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
      };
    });

    test('(default+backups): it initially renders default type and toggles to view backup methods', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.preferredMethod('oidc')).hasText('OIDC', 'it renders default method');
      assert.dom(AUTH_FORM.authForm('oidc')).exists();
      assert.dom(GENERAL.backButton).doesNotExist();
      await this.assertPathInput(assert);
      await click(AUTH_FORM.otherMethodsBtn);
      assert.dom(GENERAL.backButton).exists();
      assert.dom(AUTH_FORM.tabs).exists({ count: 2 }, 'it renders 2 backup type tabs');
      assert
        .dom(AUTH_FORM.tabBtn('userpass'))
        .hasAttribute('aria-selected', 'true', 'it selects the first backup type');
      await this.assertPathInput(assert);
      await click(AUTH_FORM.tabBtn('ldap'));
      assert.dom(AUTH_FORM.tabBtn('ldap')).hasAttribute('aria-selected', 'true', 'it selects ldap tab');
      await this.assertPathInput(assert);
    });

    test('(default only): it renders default type without backup methods', async function (assert) {
      this.loginSettings.backupTypes = null;
      await this.renderComponent();
      assert.dom(AUTH_FORM.preferredMethod('oidc')).hasText('OIDC', 'it renders default method');
      assert.dom(GENERAL.backButton).doesNotExist();
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
    });

    test('(backups only): it initially renders backup types if no default is set', async function (assert) {
      this.loginSettings.defaultType = '';
      await this.renderComponent();
      assert.dom(AUTH_FORM.preferredMethod('oidc')).doesNotExist();
      assert.dom(AUTH_FORM.tabs).exists({ count: 2 }, 'it renders 2 backup type tabs');
      assert
        .dom(AUTH_FORM.tabBtn('userpass'))
        .hasAttribute('aria-selected', 'true', 'it selects the first backup type');
      await this.assertPathInput(assert);
      assert.dom(GENERAL.backButton).doesNotExist();
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
    });

    module('all methods have visible mounts', function (hooks) {
      hooks.beforeEach(function () {
        this.loginSettings = {
          defaultType: 'oidc',
          backupTypes: ['userpass', 'ldap'],
        };
        this.visibleMountsByType = {
          userpass: [
            {
              path: 'userpass/',
              description: '',
              options: {},
              type: 'userpass',
            },
            {
              path: 'userpass2/',
              description: '',
              options: {},
              type: 'userpass',
            },
          ],
          oidc: [
            {
              path: 'my-oidc/',
              description: '',
              options: {},
              type: 'oidc',
            },
          ],
          ldap: [
            {
              path: 'ldap/',
              description: '',
              options: null,
              type: 'ldap',
            },
          ],
        };
      });

      test('(default+backups): it hides advanced settings for both views', async function (assert) {
        await this.renderComponent();
        assert.dom(AUTH_FORM.preferredMethod('oidc')).hasText('OIDC', 'it renders default method');
        this.assertPathInput(assert, { isHidden: true, value: 'my-oidc/' });
        await click(AUTH_FORM.otherMethodsBtn);
        assert.dom(AUTH_FORM.tabs).exists({ count: 2 }, 'it renders 2 backup type tabs');
        assert
          .dom(AUTH_FORM.tabBtn('userpass'))
          .hasAttribute('aria-selected', 'true', 'it selects the first backup type');
        assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
        assert.dom(GENERAL.inputByAttr('path')).doesNotExist();
        assert.dom(GENERAL.selectByAttr('path')).exists(); // dropdown renders because userpass has 2 mount paths
        await click(AUTH_FORM.tabBtn('ldap'));
        this.assertPathInput(assert, { isHidden: true, value: 'ldap/' });
      });

      test('(default only): it hides advanced settings and renders hidden input', async function (assert) {
        this.loginSettings.backupTypes = null;
        await this.renderComponent();
        assert.dom(AUTH_FORM.preferredMethod('oidc')).hasText('OIDC', 'it renders default method');
        assert.dom(AUTH_FORM.authForm('oidc')).exists();
        this.assertPathInput(assert, { isHidden: true, value: 'my-oidc/' });
        assert.dom(GENERAL.backButton).doesNotExist();
        assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
      });

      test('(backups only): it hides advanced settings and renders hidden input', async function (assert) {
        this.loginSettings.defaultType = '';
        await this.renderComponent();
        assert.dom(AUTH_FORM.preferredMethod('oidc')).doesNotExist();
        assert.dom(AUTH_FORM.tabs).exists({ count: 2 }, 'it renders 2 backup type tabs');
        assert
          .dom(AUTH_FORM.tabBtn('userpass'))
          .hasAttribute('aria-selected', 'true', 'it selects the first backup type');
        assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
        assert.dom(GENERAL.backButton).doesNotExist();
        assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
      });
    });

    module('only some methods have visible mounts', function (hooks) {
      hooks.beforeEach(function () {
        this.loginSettings = {
          defaultType: 'oidc',
          backupTypes: ['userpass', 'ldap'],
        };
        this.mountData = {
          userpass: [
            {
              path: 'userpass/',
              description: '',
              options: {},
              type: 'userpass',
            },
            {
              path: 'userpass2/',
              description: '',
              options: {},
              type: 'userpass',
            },
          ],
          oidc: [
            {
              path: 'my-oidc/',
              description: '',
              options: {},
              type: 'oidc',
            },
          ],
          ldap: [
            {
              path: 'ldap/',
              description: '',
              options: null,
              type: 'ldap',
            },
          ],
        };
      });

      test('(default+backups): it hides advanced settings for default with visible mount but it renders for backups', async function (assert) {
        this.visibleMountsByType = { oidc: this.mountData.oidc };
        await this.renderComponent();
        this.assertPathInput(assert, { isHidden: true, value: 'my-oidc/' });
        await click(AUTH_FORM.otherMethodsBtn);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        await this.assertPathInput(assert);
        await click(AUTH_FORM.tabBtn('ldap'));
        await this.assertPathInput(assert);
      });

      test('(default+backups): it only renders advanced settings for method without mounts', async function (assert) {
        // default and only one backup method have visible mounts
        this.visibleMountsByType = { oidc: this.mountData.oidc, userpass: this.mountData.userpass };
        await this.renderComponent();
        assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
        await click(AUTH_FORM.otherMethodsBtn);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        assert.dom(GENERAL.selectByAttr('path')).exists();
        assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
        await click(AUTH_FORM.tabBtn('ldap'));
        assert.dom(AUTH_FORM.advancedSettings).exists();
      });

      test('(default+backups): it hides advanced settings for single backup method with mounts', async function (assert) {
        this.visibleMountsByType = { ldap: this.mountData.ldap };
        await this.renderComponent();
        assert.dom(AUTH_FORM.authForm('oidc')).exists();
        assert.dom(AUTH_FORM.advancedSettings).exists();
        await click(AUTH_FORM.otherMethodsBtn);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        assert.dom(AUTH_FORM.advancedSettings).exists();
        await click(AUTH_FORM.tabBtn('ldap'));
        this.assertPathInput(assert, { isHidden: true, value: 'ldap/' });
      });

      test('(backups only): it hides advanced settings for single method with mounts', async function (assert) {
        this.loginSettings.defaultType = '';
        this.visibleMountsByType = { ldap: this.mountData.ldap };
        await this.renderComponent();
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        assert.dom(AUTH_FORM.advancedSettings).exists();
        await click(AUTH_FORM.tabBtn('ldap'));
        this.assertPathInput(assert, { isHidden: true, value: 'ldap/' });
      });
    });

    module('@directLinkData overrides login settings', function (hooks) {
      hooks.beforeEach(function () {
        this.mountData = {
          userpass: [
            {
              path: 'userpass/',
              description: '',
              options: {},
              type: 'userpass',
            },
            {
              path: 'userpass2/',
              description: '',
              options: {},
              type: 'userpass',
            },
          ],
          oidc: [
            {
              path: 'my-oidc/',
              description: '',
              options: {},
              type: 'oidc',
            },
          ],
          ldap: [
            {
              path: 'ldap/',
              description: '',
              options: null,
              type: 'ldap',
            },
          ],
        };
      });

      module('when there are no visible mounts at all', function (hooks) {
        hooks.beforeEach(function () {
          this.visibleMountsByType = null;
          this.directLinkData = { type: 'okta', isVisibleMount: false };
        });

        const testHelper = (assert) => {
          assert.dom(AUTH_FORM.selectMethod).hasValue('okta');
          assert.dom(AUTH_FORM.authForm('okta')).exists();
          assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
          assert.dom(GENERAL.backButton).doesNotExist();
        };

        test('(default+backups): it renders standard view and selects @directLinkData type from dropdown', async function (assert) {
          await this.renderComponent();
          testHelper(assert);
        });

        test('(default only): it renders standard view and selects @directLinkData type from dropdown', async function (assert) {
          this.loginSettings.backupTypes = null;
          await this.renderComponent();
          testHelper(assert);
        });

        test('(backups only): it renders standard view and selects @directLinkData type from dropdown', async function (assert) {
          this.loginSettings.defaultType = '';
          await this.renderComponent();
          testHelper(assert);
        });
      });

      module('when param matches a visible mount path and other visible mounts exist', function (hooks) {
        hooks.beforeEach(function () {
          this.visibleMountsByType = {
            ...this.mountData,
            okta: [
              {
                path: 'my-okta/',
                description: '',
                options: null,
                type: 'okta',
              },
            ],
          };
          this.directLinkData = { path: 'my-okta/', type: 'okta', isVisibleMount: true };
        });

        const testHelper = async (assert) => {
          assert.dom(AUTH_FORM.preferredMethod('okta')).hasText('Okta');
          assert.dom(AUTH_FORM.authForm('okta'));
          assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
          assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
          assert.dom(GENERAL.inputByAttr('path')).hasValue('my-okta/');
          assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
          await click(AUTH_FORM.otherMethodsBtn);
          assert
            .dom(GENERAL.selectByAttr('auth type'))
            .exists('it renders dropdown after clicking "Sign in with other"');
        };

        test('(default+backups): it renders single mount view for @directLinkData', async function (assert) {
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(default only): it renders single mount view for @directLinkData', async function (assert) {
          this.loginSettings.backupTypes = null;
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(backups only): it renders single mount view for @directLinkData', async function (assert) {
          this.loginSettings.defaultType = '';
          await this.renderComponent();
          await testHelper(assert);
        });
      });

      module('when param matches a type and other visible mounts exist', function (hooks) {
        hooks.beforeEach(function () {
          // isVisibleMount is false because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'okta', isVisibleMount: false };
          this.visibleMountsByType = this.mountData;
        });

        const testHelper = async (assert) => {
          assert.dom(GENERAL.backButton).exists('back button renders because other methods have tabs');
          assert.dom(AUTH_FORM.selectMethod).hasValue('okta');
          assert.dom(AUTH_FORM.authForm('okta')).exists();
          assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
          await click(GENERAL.backButton);
          assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
          await click(AUTH_FORM.otherMethodsBtn);
          assert.dom(AUTH_FORM.selectMethod).exists('it toggles back to dropdown');
        };

        test('(default+backups): it selects @directLinkData type from dropdown and toggles to tab view', async function (assert) {
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(default only): it selects @directLinkData type from dropdown and toggles to tab view', async function (assert) {
          this.loginSettings.backupTypes = null;
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(backups only): it selects @directLinkData type from dropdown and toggles to tab view', async function (assert) {
          this.loginSettings.defaultType = '';
          await this.renderComponent();
          await testHelper(assert);
        });
      });

      module('when param matches a type that matches other visible mounts', function (hooks) {
        hooks.beforeEach(function () {
          // isVisibleMount is false because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'oidc', isVisibleMount: false };
          this.visibleMountsByType = this.mountData;
        });

        const testHelper = async (assert) => {
          assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'true');
          assert.dom(AUTH_FORM.authForm('oidc')).exists();
          assert.dom(GENERAL.backButton).doesNotExist();
          await click(AUTH_FORM.otherMethodsBtn);
          assert.dom(AUTH_FORM.selectMethod).exists('it toggles to view dropdown');
          await click(GENERAL.backButton);
          assert.dom(AUTH_FORM.tabs).exists('it toggles back to tabs');
        };

        test('(default+backups): it selects @directLinkData type tab and toggles to dropdown view', async function (assert) {
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(default only): it selects @directLinkData type tab and toggles to dropdown view', async function (assert) {
          this.loginSettings.backupTypes = null;
          await this.renderComponent();
          await testHelper(assert);
        });

        test('(backups only): it selects @directLinkData type tab and toggles to dropdown view', async function (assert) {
          this.loginSettings.defaultType = '';
          await this.renderComponent();
          await testHelper(assert);
        });
      });
    });
  });

  // AUTH METHOD SPECIFIC TESTS
  // since the template yields each auth <form> some assertions are best done here instead of
  // in the corresponding the Auth::Form::<Type> integration tests
  module('oidc-jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.store = this.owner.lookup('service:store');
      this.routerStub = (path) =>
        sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns(`/auth/${path}/oidc/callback`);
    });

    test('it re-requests the auth_url when authType changes', async function (assert) {
      this.routerStub('oidc');
      assert.expect(2); // auth_url should be hit twice, one for each type selection
      let expectedType = 'oidc';
      this.server.post(`/auth/:path/oidc/auth_url`, (_, req) => {
        assert.strictEqual(
          req.params.path,
          expectedType,
          `it makes request to auth_url for selected type: ${expectedType}`
        );
        return { data: { auth_url: '123-example.com' } };
      });
      await this.renderComponent();
      // auth_url should be requested once when "oidc" is selected
      await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
      // auth_url should be requested again when "jwt" is selected
      expectedType = 'jwt';
      await fillIn(GENERAL.selectByAttr('auth type'), 'jwt');
    });

    // for simplicity the auth types are configured as their namesake but type isn't relevant.
    // these tests assert that CONFIG changes from OIDC -> JWT render correctly and vice versa
    // so the order the requests are hit is what matters.
    test('"OIDC" to "JWT" configuration: it updates the form when the auth_url response changes', async function (assert) {
      this.routerStub('oidc');
      this.server.post(`/auth/oidc/oidc/auth_url`, () => ({ data: { auth_url: '123-example.com' } })); // this return means mount is configured as oidc
      this.server.post(`/auth/jwt/oidc/auth_url`, () => overrideResponse(400, { errors: [ERROR_JWT_LOGIN] })); // this return means the mount is configured as jwt
      await this.renderComponent();

      // select mount configured for OIDC first
      await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
      assert.dom(GENERAL.inputByAttr('jwt')).doesNotExist();
      // then select mount configured for JWT
      await fillIn(GENERAL.selectByAttr('auth type'), 'jwt');
      assert.dom(GENERAL.inputByAttr('jwt')).exists();
    });

    test('"JWT" to "OIDC" configuration: it updates the form when the auth_url response changes', async function (assert) {
      this.routerStub('oidc');
      this.server.post(`/auth/jwt/oidc/auth_url`, () => overrideResponse(400, { errors: [ERROR_JWT_LOGIN] })); // this return means the mount is configured as jwt
      this.server.post(`/auth/oidc/oidc/auth_url`, () => ({ data: { auth_url: '123-example.com' } })); // this return means mount is configured as oidc
      await this.renderComponent();

      // select mount configured for JWT first
      await fillIn(GENERAL.selectByAttr('auth type'), 'jwt');
      assert.dom(GENERAL.inputByAttr('jwt')).exists();

      // then select mount configured for OIDC
      await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
      assert.dom(GENERAL.inputByAttr('jwt')).doesNotExist();
    });

    test('it should retain role input value when mount path changes', async function (assert) {
      assert.expect(2);
      this.routerStub('foo-oidc');
      const auth_url = 'http://dev-foo-bar.com';
      this.server.post('/auth/:path/oidc/auth_url', (_, req) => {
        const { role, redirect_uri } = JSON.parse(req.requestBody);
        const goodRequest =
          req.params.path === 'foo-oidc' &&
          role === 'foo' &&
          redirect_uri.includes('/auth/foo-oidc/oidc/callback');
        if (goodRequest) {
          return { data: { auth_url } };
        } else {
          return overrideResponse(400, { errors: [ERROR_JWT_LOGIN] });
        }
      });

      window.open = (url) => {
        assert.strictEqual(url, auth_url, 'auth_url is returned when required params are passed');
      };

      await this.renderComponent();

      await fillIn(GENERAL.selectByAttr('auth type'), 'oidc');
      await fillIn(GENERAL.inputByAttr('role'), 'foo');
      await click(AUTH_FORM.advancedSettings);
      await fillIn(GENERAL.inputByAttr('role'), 'foo');
      await fillIn(GENERAL.inputByAttr('path'), 'foo-oidc');
      assert.dom(GENERAL.inputByAttr('role')).hasValue('foo', 'role is retained when mount path is changed');
      await click(AUTH_FORM.login);
    });
  });
});
