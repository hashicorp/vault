/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, find, findAll, render, typeIn } from '@ember/test-helpers';
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
import { Response } from 'miragejs';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';

module('Integration | Component | auth | form template', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };
    this.wrappedToken = '';
    this.namespaceQueryParam = '';
    this.oidcProviderQueryParam = '';
    this.onAuthResponse = sinon.spy();
    this.onNamespaceChange = sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
         <Auth::FormTemplate
          @wrappedToken={{this.wrappedToken}}
          @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
          @cluster={{this.cluster}}
          @handleNamespaceUpdate={{this.onNamespaceChange}}
          @namespace={{this.namespaceQueryParam}}
          @onSuccess={{this.onAuthResponse}}
        />`);
    };
  });

  // test to select each method is in "ent" module to include enterprise methods
  test('it selects token by default', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('token');
  });

  test('it does not show toggle buttons when listing visibility is not set', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render');
    assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render ');
  });

  test('it calls sys/internal/ui/mounts on initial render', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/ui/mounts', (_, req) => {
      assert.true(true, 'request is made to /sys/internal/ui/mounts');
      assert.strictEqual(
        req.requestHeaders['X-Vault-Namespace'],
        undefined,
        'it does not pass a namespace header'
      );
      return {};
    });

    await this.renderComponent();
  });

  test('it fails gracefully if sys/internal/ui/mounts request errors', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/ui/mounts', () => {
      assert.true(true, 'request is made to /sys/internal/ui/mounts');
      return new Response(500, {}, { errors: ['something wrong with urls'] });
    });
    await this.renderComponent();
    assert.dom(GENERAL.selectByAttr('auth type')).exists();
  });

  test('it displays errors', async function (assert) {
    await this.renderComponent();
    await click(AUTH_FORM.login);
    // this error message text is because the auth service is not stubbed in this test
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: permission denied');
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.server.get('/sys/internal/ui/mounts', () => {
        return {
          data: {
            auth: {
              'userpass/': {
                description: '',
                options: {},
                type: 'userpass',
              },
              'userpass2/': {
                description: '',
                options: {},
                type: 'userpass',
              },
              'my-oidc/': {
                description: '',
                options: {},
                type: 'oidc',
              },
              'token/': {
                description: 'token based credentials',
                options: null,
                type: 'token',
              },
            },
          },
        };
      });
    });

    test('it renders mounts configured with listing_visibility="unuath"', async function (assert) {
      const expectedTabs = [
        { type: 'userpass', display: 'Username' },
        { type: 'oidc', display: 'OIDC' },
        { type: 'token', display: 'Token' },
      ];
      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      // there are 4 mount paths returned in the stubbed sys/internal/ui/mounts response above,
      // but two are of the same type so only expect 3 tabs
      assert.dom(AUTH_FORM.tabs()).exists({ count: 3 }, 'it groups mount paths by type and renders 3 tabs');
      expectedTabs.forEach((m) => {
        assert.dom(AUTH_FORM.tabs(m.type)).exists(`${m.type} renders as a tab`);
        assert.dom(AUTH_FORM.tabs(m.type)).hasText(m.display, `${m.type} renders expected display name`);
      });
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
      assert.dom('section p').hasText('token based credentials data-test-description');
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

    test('it renders a readonly input if only one mount path is returned', async function (assert) {
      await this.renderComponent();
      await click(AUTH_FORM.tabBtn('oidc'));
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('readonly');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
    });

    test('it clicks "Sign in with other methods"', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs()).exists({ count: 3 }, 'tabs render by default');
      assert.dom(GENERAL.backButton).doesNotExist();
      await click(AUTH_FORM.otherMethodsBtn);
      assert
        .dom(AUTH_FORM.otherMethodsBtn)
        .doesNotExist('"Sign in with other methods" does not render after it is clicked');
      assert
        .dom(GENERAL.selectByAttr('auth type'))
        .exists('clicking "Sign in with other methods" renders dropdown instead of tabs');
      await click(GENERAL.backButton);
      assert.dom(GENERAL.backButton).doesNotExist('"Back" button does not render after it is clicked');
      assert.dom(AUTH_FORM.tabs()).exists({ count: 3 }, 'clicking "Back" renders tabs again');
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

    test('it re-requests mount data when a namespace is inputted', async function (assert) {
      assert.expect(3);
      const expectedNs = 'test-ns1';

      let count = 0;
      this.server.get('/sys/internal/ui/mounts', () => {
        count++;
        const msg = count === 1 ? 'on initial render' : 'when namespace is inputted';
        assert.true(true, `/sys/internal/ui/mounts is called ${msg}`);
        return {};
      });

      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('namespace'), expectedNs);
      const [actual] = this.onNamespaceChange.lastCall.args;
      assert.strictEqual(actual, expectedNs, 'callback has expected args');
    });

    test('it re-requests mount data when namespace input is prefilled and then updated', async function (assert) {
      assert.expect(3);
      this.namespaceQueryParam = 'admin';
      const childNs = '/test-ns1';

      let count = 0;
      this.server.get('/sys/internal/ui/mounts', () => {
        count++;
        const msg = count === 1 ? 'on initial render' : 'when namespace updates';
        assert.true(true, `/sys/internal/ui/mounts is called ${msg}`);
        return {};
      });

      await this.renderComponent();
      await typeIn(GENERAL.inputByAttr('namespace'), childNs);
      const [actual] = this.onNamespaceChange.lastCall.args;
      assert.strictEqual(actual, `${this.namespaceQueryParam}${childNs}`, 'callback has expected args');
    });

    test('it sets namespace for hvd managed clusters', async function (assert) {
      this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespaceQueryParam = 'admin/west-coast';
      await this.renderComponent();
      assert.dom(AUTH_FORM.managedNsRoot).hasValue('/admin');
      assert.dom(AUTH_FORM.managedNsRoot).hasAttribute('readonly');
      assert.dom(GENERAL.inputByAttr('namespace')).hasValue('/west-coast');
    });

    test('it does NOT display tabs when updated namespace has no visible mounts', async function (assert) {
      assert.expect(4);
      let count = 0;
      this.server.get('/sys/internal/ui/mounts', () => {
        count++;
        const mounts = {
          data: {
            auth: {
              'userpass2/': {
                description: '',
                options: {},
                type: 'userpass',
              },
            },
          },
        };
        // mocks re-requesting the endpoint when namespace changes by returning
        // mounts on initial request, then when a namespace is inputted a second request is made which return NO mounts
        const response = count === 1 ? mounts : {};
        return response;
      });

      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs('userpass')).exists('userpass renders as a tab');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      await fillIn(GENERAL.inputByAttr('namespace'), 'admin');
      assert.dom(AUTH_FORM.tabs()).doesNotExist('tabs do not render');
      assert.dom(GENERAL.selectByAttr('auth type')).exists('dropdown renders');
    });

    test('it DOES display tabs when updated namespace has visible mounts', async function (assert) {
      assert.expect(4);
      let count = 0;
      this.server.get('/sys/internal/ui/mounts', () => {
        count++;
        const mounts = {
          data: {
            auth: {
              'userpass2/': {
                description: '',
                options: {},
                type: 'userpass',
              },
            },
          },
        };
        // mocks re-requesting the endpoint when namespace changes by returning
        // no mounts on initial request, then when a namespace is inputted a second request is made which return mounts
        const response = count === 1 ? {} : mounts;
        return response;
      });

      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs()).doesNotExist('tabs do not render');
      assert.dom(GENERAL.selectByAttr('auth type')).exists('dropdown renders');
      // fire off second request to sys/internal/mounts
      await fillIn(GENERAL.inputByAttr('namespace'), 'admin');
      assert.dom(AUTH_FORM.tabs('userpass')).exists('userpass renders as a tab');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
    });
  });

  // AUTH METHOD SPECIFIC TESTS
  // since the template yields each auth <form> some assertions are best done here instead of
  // in the corresponding the Auth::Form::<Type> integration tests
  module('oidc-jwt', function (hooks) {
    hooks.beforeEach(async function () {
      this.store = this.owner.lookup('service:store');
      this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'urlFor').returns('123-example.com');
    });

    test('it re-requests the auth_url when authType changes', async function (assert) {
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
  });
});
