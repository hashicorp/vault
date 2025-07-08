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
import { AUTH_METHOD_LOGIN_DATA } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { ENTERPRISE_LOGIN_METHODS, supportedTypes } from 'vault/utils/supported-login-methods';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';

module('Integration | Component | auth | form template', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    window.localStorage.clear();
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };

    this.alternateView = null;
    this.defaultView = { view: 'dropdown', tabData: null };
    this.initialFormState = { initialAuthType: 'token', showAlternate: false };
    this.onSuccess = sinon.spy();
    this.visibleMountTypes = null;

    this.renderComponent = () => {
      return render(hbs`
         <Auth::FormTemplate
          @alternateView={{this.alternateView}}
          @cluster={{this.cluster}}
          @defaultView={{this.defaultView}}
          @initialFormState={{this.initialFormState}}
          @onSuccess={{this.onSuccess}}
          @visibleMountTypes={{this.visibleMountTypes}}
        />`);
    };
  });

  // test to select each method is in "ent" module to include enterprise methods
  test('it selects token by default', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('token');
  });

  test('it does not show toggle buttons if @alternateView does not exist', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render');
    assert
      .dom(GENERAL.button('Sign in with other methods'))
      .doesNotExist('"Sign in with other methods" does not render');
  });

  test('it initializes with preset auth type', async function (assert) {
    this.initialFormState = { initialAuthType: 'userpass' };
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('userpass');
  });

  test('it displays errors', async function (assert) {
    const authenticateStub = sinon.stub(this.owner.lookup('service:auth'), 'authenticate');
    authenticateStub.throws('permission denied');
    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error Authentication failed: permission denied: Sinon-provided permission denied');
    authenticateStub.restore();
  });

  test('dropdown does not include enterprise methods on community versions', async function (assert) {
    this.version.type = 'community';
    const supported = supportedTypes(false);
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

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      const defaultTabs = {
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
            path: 'my_oidc/',
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
      // all computed by the parent, in this case the initial tabs are the same as visible mount types
      // but that isn't always the case
      this.visibleMountTypes = Object.keys(defaultTabs);
      this.defaultView = { type: 'tabs', tabData: defaultTabs };
      this.alternateView = { type: 'dropdown', tabData: null };
      this.initialFormState = { initialAuthType: 'userpass', showAlternate: false };
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

    test('it clicks "Sign in with other methods" and toggles to other view', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'tabs render by default');
      assert.dom(GENERAL.backButton).doesNotExist();
      await click(GENERAL.button('Sign in with other methods'));
      assert
        .dom(GENERAL.button('Sign in with other methods'))
        .doesNotExist('"Sign in with other methods" does not render after it is clicked');
      assert
        .dom(GENERAL.selectByAttr('auth type'))
        .exists('clicking "Sign in with other methods" renders dropdown instead of tabs');
      await click(GENERAL.backButton);
      assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render after it is clicked');
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'clicking "Back" renders tabs again');
      assert
        .dom(GENERAL.button('Sign in with other methods'))
        .exists('"Sign in with other methods" renders again');
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
      await click(GENERAL.button('Sign in with other methods'));
      assert.dom(GENERAL.selectByAttr('auth type')).exists('it renders dropdown instead of tabs');
      await click(GENERAL.backButton);
      // assert tab selection is reset
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'false');
      assert.dom(AUTH_FORM.tabBtn('token')).hasAttribute('aria-selected', 'false');
    });

    test('it preselects tab from initialFormState', async function (assert) {
      this.initialFormState = { initialAuthType: 'oidc', showAlternate: false };
      await this.renderComponent();
      assert.dom(AUTH_FORM.authForm('oidc')).exists('oidc form renders');
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'true');
    });

    test('it renders dropdown and preselects type if initialFormState is not a tab', async function (assert) {
      this.initialFormState = { initialAuthType: 'ldap', showAlternate: true };
      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('ldap');
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();

      assert.dom(GENERAL.backButton).exists('"Back" button renders');
      assert
        .dom(GENERAL.button('Sign in with other methods'))
        .doesNotExist('"Sign in with other methods" does not render');
    });
  });

  // tests with "enterprise" in the title are filtered out from CE test runs
  // naming the module 'ent' so these tests still run on the CE repo
  module('ent', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
      this.namespaceQueryParam = '';
    });

    // in the ent module to test ALL supported login methods
    // iterating in tests should generally be avoided, but purposefully wanted to test the component
    // renders as expected as auth types change
    test('it selects each supported auth type and renders its form and relevant fields', async function (assert) {
      const authMethodTypes = supportedTypes(true);
      const totalFields = Object.values(AUTH_METHOD_LOGIN_DATA).reduce(
        (sum, obj) => sum + Object.keys(obj).length,
        0
      );
      // 3 assertions per method, plus an assertion for each expected field
      assert.expect(3 * authMethodTypes.length + totalFields); // count at time of writing is 40

      await this.renderComponent();
      for (const authType of authMethodTypes) {
        const loginData = AUTH_METHOD_LOGIN_DATA[authType];

        const fields = Object.keys(loginData);
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

    test('dropdown includes enterprise methods', async function (assert) {
      const supported = supportedTypes(true);
      assert.expect(supported.length);
      await this.renderComponent();

      const dropdownOptions = findAll(`${GENERAL.selectByAttr('auth type')} option`).map((o) => o.value);
      supported.forEach((m) => {
        assert.true(dropdownOptions.includes(m), `dropdown includes supported method: ${m}`);
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
      await click(GENERAL.submitButton);
    });
  });
});
