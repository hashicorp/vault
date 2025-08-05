/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { fillInLoginFields, SYS_INTERNAL_UI_MOUNTS } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CSP_ERROR } from 'vault/components/auth/page';
import { setupTotpMfaResponse } from 'vault/tests/helpers/mfa/mfa-helpers';

module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };
    this.directLinkData = null;
    this.loginSettings = null;
    this.namespaceQueryParam = '';
    this.oidcProviderQueryParam = '';
    this.onAuthSuccess = sinon.spy();
    this.onNamespaceUpdate = sinon.spy();
    this.visibleAuthMounts = false;

    this.renderComponent = () => {
      return render(hbs`
        <Auth::Page
          @cluster={{this.cluster}}
          @directLinkData={{this.directLinkData}}
          @loginSettings={{this.loginSettings}}
          @namespaceQueryParam={{this.namespaceQueryParam}}
          @oidcProviderQueryParam={{this.oidcProviderQueryParam}}
          @onAuthSuccess={{this.onAuthSuccess}}
          @onNamespaceUpdate={{this.onNamespaceUpdate}}
          @visibleAuthMounts={{this.visibleAuthMounts}}
        />
        `);
    };

    // in the real world more info is returned by auth requests
    // only including pertinent data for testing
    this.authRequest = (url) => this.server.post(url, () => ({ auth: { policies: ['default'] } }));
  });

  test('it renders error on CSP violation', async function (assert) {
    assert.expect(2);
    this.cluster.standby = true;
    await this.renderComponent();
    assert.dom(GENERAL.pageError.error).doesNotExist();
    this.owner.lookup('service:csp-event').handleEvent({ violatedDirective: 'connect-src' });
    await waitFor(GENERAL.pageError.error);
    assert.dom(GENERAL.pageError.error).hasText(CSP_ERROR);
  });

  test('it renders splash logo and disables namespace input when oidc provider query param is present', async function (assert) {
    this.oidcProviderQueryParam = 'myprovider';
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
    assert
      .dom(AUTH_FORM.helpText)
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.'
      );
  });

  test('it calls onNamespaceUpdate', async function (assert) {
    assert.expect(1);
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  test('it passes query param to namespace input', async function (assert) {
    this.version.features = ['Namespaces'];
    this.namespaceQueryParam = 'ns-1';
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).hasValue(this.namespaceQueryParam);
  });

  test('it does not render the namespace input on community', async function (assert) {
    this.version.type = 'community';
    this.version.features = [];
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
  });

  test('it does not render the namespace input on enterprise without the "Namespaces" feature', async function (assert) {
    this.version.type = 'enterprise';
    this.version.features = [];
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('namespace')).doesNotExist();
  });

  test('it selects type in the dropdown if direct link just has type', async function (assert) {
    this.directLinkData = { type: 'oidc' };
    await this.renderComponent();
    assert.dom(AUTH_FORM.tabBtn('oidc')).doesNotExist('tab does not render');
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('oidc', 'dropdown has type selected');
    assert.dom(AUTH_FORM.authForm('oidc')).exists();
    assert.dom(GENERAL.inputByAttr('role')).exists();
    await click(AUTH_FORM.advancedSettings);
    assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
    assert.dom(GENERAL.backButton).doesNotExist();
    assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.visibleAuthMounts = SYS_INTERNAL_UI_MOUNTS;
      window.localStorage.clear();
    });

    test('it formats and renders tabs if visible auth mounts exist', async function (assert) {
      await this.renderComponent();
      const expectedTabs = [
        { type: 'userpass', display: 'Userpass' },
        { type: 'oidc', display: 'OIDC' },
        { type: 'ldap', display: 'LDAP' },
      ];

      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist('dropdown does not render');
      // there are 4 mount paths returned in visibleAuthMounts above,
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

    test('it renders dropdown as alternate view', async function (assert) {
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'tabs render by default');
      assert.dom(GENERAL.backButton).doesNotExist();
      await click(AUTH_FORM.otherMethodsBtn);
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('button disappears after it is clicked');
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('userpass', 'dropdown has userpass selected');
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

    module('with a direct link', function (hooks) {
      hooks.beforeEach(function () {
        // if path exists, the mount has listing_visibility="unauth"
        this.directLinkIsVisibleMount = { path: 'my_oidc/', type: 'oidc' };
        this.directLinkIsJustType = { type: 'okta' };
      });

      test('it selects type in the dropdown if direct link is just type', async function (assert) {
        this.directLinkData = this.directLinkIsJustType;
        await this.renderComponent();
        assert.dom(AUTH_FORM.tabBtn('okta')).doesNotExist('tab does not render');
        assert.dom(GENERAL.selectByAttr('auth type')).hasValue('okta', 'dropdown has type selected');
        assert.dom(AUTH_FORM.authForm('okta')).exists();
        assert.dom(GENERAL.inputByAttr('username')).exists();
        assert.dom(GENERAL.inputByAttr('password')).exists();
        await click(AUTH_FORM.advancedSettings);
        assert.dom(GENERAL.inputByAttr('path')).exists({ count: 1 });
        assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
        assert.dom(GENERAL.backButton).exists('back button renders because tabs exist for other methods');
        await click(GENERAL.backButton);
        assert
          .dom(AUTH_FORM.tabBtn('userpass'))
          .hasAttribute('aria-selected', 'true', 'first tab is selected on back');
      });

      test('it renders single method view instead of tabs if direct link includes path', async function (assert) {
        this.directLinkData = this.directLinkIsVisibleMount;
        await this.renderComponent();
        assert.dom(AUTH_FORM.authForm('oidc')).exists;
        assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders tab for type');
        assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
        assert.dom(GENERAL.inputByAttr('role')).exists();
        assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
        assert.dom(GENERAL.inputByAttr('path')).hasValue('my_oidc/');
        assert.dom(AUTH_FORM.otherMethodsBtn).exists('"Sign in with other methods" renders');
        assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist();
        assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
        assert.dom(GENERAL.backButton).doesNotExist();
      });

      test('it prioritizes auth type from canceled mfa instead of direct link for path', async function (assert) {
        assert.expect(1);
        this.directLinkData = this.directLinkIsVisibleMount;
        const authType = 'okta';
        const { loginData, url } = REQUEST_DATA.username;
        const requestUrl = url({ path: authType, username: loginData?.username });
        this.server.post(requestUrl, () => setupTotpMfaResponse(authType));

        await this.renderComponent();
        await click(AUTH_FORM.otherMethodsBtn);
        await fillIn(AUTH_FORM.selectMethod, authType);
        await fillInLoginFields(loginData);
        await click(AUTH_FORM.login);
        await waitFor('[data-test-mfa-description]'); // wait until MFA validation renders
        await click(GENERAL.backButton);
        assert.dom(AUTH_FORM.selectMethod).hasValue(authType, 'Okta is selected in dropdown');
      });

      test('it prioritizes auth type from canceled mfa instead of direct link with type', async function (assert) {
        assert.expect(1);
        this.directLinkData = this.directLinkIsJustType;
        const authType = 'userpass';
        const { loginData, url } = REQUEST_DATA.username;
        const requestUrl = url({ path: authType, username: loginData?.username });
        this.server.post(requestUrl, () => setupTotpMfaResponse(authType));
        await this.renderComponent();
        await fillIn(AUTH_FORM.selectMethod, authType);
        await fillInLoginFields(loginData);
        await click(AUTH_FORM.login);
        await click(GENERAL.backButton);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
      });
    });
  });

  const REQUEST_DATA = {
    username: {
      loginData: { username: 'matilda', password: 'password' },
      url: ({ path, username }) => `/auth/${path}/login/${username}`,
    },
    github: {
      loginData: { token: 'mysupersecuretoken' },
      url: ({ path }) => `/auth/${path}/login`,
    },
  };

  // only testing methods that submit via AuthForm (and not separate, child component)
  const AUTH_METHOD_TEST_CASES = [
    { authType: 'github', options: REQUEST_DATA.github },
    { authType: 'userpass', options: REQUEST_DATA.username },
    { authType: 'ldap', options: REQUEST_DATA.username },
    { authType: 'okta', options: REQUEST_DATA.username },
    { authType: 'radius', options: REQUEST_DATA.username },
  ];

  for (const { authType, options } of AUTH_METHOD_TEST_CASES) {
    test(`${authType}: it calls onAuthSuccess on submit for default path`, async function (assert) {
      assert.expect(1);
      const { loginData, url } = options;
      const requestUrl = url({ path: authType, username: loginData?.username });
      this.authRequest(requestUrl);

      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields(loginData);
      await click(AUTH_FORM.login);
      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = {
        namespace: '',
        token: `vault-${authType}☃1`,
        isRoot: false,
      };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
    });

    test(`${authType}: it calls onAuthSuccess on submit for custom path`, async function (assert) {
      assert.expect(1);
      const customPath = `${authType}-custom`;
      const { loginData, url } = options;
      const loginDataWithPath = { ...loginData, path: customPath };
      // pass custom path to request URL
      const requestUrl = url({ path: customPath, username: loginData?.username });
      this.authRequest(requestUrl);

      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      // await fillIn(AUTH_FORM.selectMethod, authType);
      // toggle mount path input to specify custom path
      await fillInLoginFields(loginDataWithPath, { toggleOptions: true });
      await click(AUTH_FORM.login);

      const [actual] = this.onAuthSuccess.lastCall.args;
      const expected = {
        namespace: '',
        token: `vault-${authType}☃1`,
        isRoot: false,
      };
      assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
    });

    test('it preselects auth type from canceled mfa', async function (assert) {
      assert.expect(1);
      const { loginData, url } = options;
      const requestUrl = url({ path: authType, username: loginData?.username });
      this.server.post(requestUrl, () => setupTotpMfaResponse(authType));

      await this.renderComponent();
      await fillIn(AUTH_FORM.selectMethod, authType);
      await fillInLoginFields(loginData);
      await click(AUTH_FORM.login);
      await click(GENERAL.backButton);
      assert.dom(AUTH_FORM.selectMethod).hasValue(authType, `${authType} is selected in dropdown`);
    });
  }

  // token makes a GET request so test separately
  test('token: it calls onAuthSuccess on submit', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return { data: { policies: ['default'] } };
    });

    await this.renderComponent();
    await fillIn(AUTH_FORM.selectMethod, 'token');
    // await fillIn(AUTH_FORM.selectMethod, 'token');
    await fillInLoginFields({ token: 'mysupersecuretoken' });
    await click(AUTH_FORM.login);
    const [actual] = this.onAuthSuccess.lastCall.args;
    const expected = {
      namespace: '',
      token: `vault-token☃1`,
      isRoot: false,
    };
    assert.propEqual(actual, expected, `onAuthSuccess called with: ${JSON.stringify(actual)}`);
  });

  /* 
  Login settings are an enterprise only feature but the component is version agnostic (and subsequently so are these tests) 
  because fetching login settings happens in the route only for enterprise versions.
  Each combination must be tested with and without visible mounts (i.e. tuned with listing_visibility="unauth")
  1. default+backups: default type set, backup types set 
  2. default only: no backup types
  3. backup only: backup types set without a default
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
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
      assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
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

    test('(default+backups): it initially renders default type if backup types include the default method', async function (assert) {
      this.loginSettings = {
        defaultType: 'userpass',
        backupTypes: ['ldap', 'userpass', 'oidc'],
      };
      await this.renderComponent();
      assert.dom(GENERAL.backButton).doesNotExist('it renders default view');
      assert.dom(AUTH_FORM.tabBtn('userpass')).hasText('Userpass', 'it renders default method');
      assert
        .dom(AUTH_FORM.tabs)
        .exists({ count: 1 }, 'it is rendering the default view because only one tab renders');

      await click(AUTH_FORM.otherMethodsBtn);
      assert.dom(GENERAL.backButton).exists('it toggles to backup method view');
      assert.dom(AUTH_FORM.tabs).exists({ count: 3 }, 'it renders 3 backup type tabs');
      assert
        .dom(AUTH_FORM.tabBtn('ldap'))
        .hasAttribute('aria-selected', 'true', 'it selects the first backup type');
    });

    test('(default only): it renders default type without backup methods', async function (assert) {
      this.loginSettings.backupTypes = null;
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
      assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
      assert.dom(GENERAL.backButton).doesNotExist();
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
    });

    test('(backups only): it initially renders backup types if no default is set', async function (assert) {
      this.loginSettings.defaultType = '';
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('oidc')).doesNotExist();
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
        this.visibleAuthMounts = SYS_INTERNAL_UI_MOUNTS;
      });

      test('(default+backups): it hides advanced settings for both views', async function (assert) {
        await this.renderComponent();
        assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
        assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
        this.assertPathInput(assert, { isHidden: true, value: 'my_oidc/' });
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
        assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
        assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
        assert.dom(AUTH_FORM.authForm('oidc')).exists();
        this.assertPathInput(assert, { isHidden: true, value: 'my_oidc/' });
        assert.dom(GENERAL.backButton).doesNotExist();
        assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist();
      });

      test('(backups only): it hides advanced settings and renders hidden input', async function (assert) {
        this.loginSettings.defaultType = '';
        await this.renderComponent();
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
        this.mountData = (path) => ({ [path]: SYS_INTERNAL_UI_MOUNTS[path] });
      });

      test('(default+backups): it hides advanced settings for default with visible mount but it renders for backups', async function (assert) {
        this.visibleAuthMounts = { ...this.mountData('my_oidc/') };
        await this.renderComponent();
        this.assertPathInput(assert, { isHidden: true, value: 'my_oidc/' });
        await click(AUTH_FORM.otherMethodsBtn);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        await this.assertPathInput(assert);
        await click(AUTH_FORM.tabBtn('ldap'));
        await this.assertPathInput(assert);
      });

      test('(default+backups): it only renders advanced settings for method without mounts', async function (assert) {
        // default and only one backup method have visible mounts
        this.visibleAuthMounts = {
          ...this.mountData('my_oidc/'),
          ...this.mountData('userpass/'),
          ...this.mountData('userpass2/'),
        };
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
        this.visibleAuthMounts = { ...this.mountData('ldap/') };
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
        this.visibleAuthMounts = { ...this.mountData('ldap/') };
        await this.renderComponent();
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        assert.dom(AUTH_FORM.advancedSettings).exists();
        await click(AUTH_FORM.tabBtn('ldap'));
        this.assertPathInput(assert, { isHidden: true, value: 'ldap/' });
      });
    });

    module('@directLinkData overrides login settings', function (hooks) {
      hooks.beforeEach(function () {
        this.mountData = SYS_INTERNAL_UI_MOUNTS;
      });

      module('when there are no visible mounts at all', function (hooks) {
        hooks.beforeEach(function () {
          this.visibleAuthMounts = null;
          this.directLinkData = { type: 'okta' };
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
          this.visibleAuthMounts = {
            ...this.mountData,
            'my-okta/': {
              description: '',
              options: null,
              type: 'okta',
            },
          };
          this.directLinkData = { path: 'my-okta/', type: 'okta' };
        });

        const testHelper = async (assert) => {
          assert.dom(AUTH_FORM.tabBtn('okta')).hasText('Okta', 'it renders preferred method');
          assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
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
          // only type is present in directLinkData because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'okta' };
          this.visibleAuthMounts = this.mountData;
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
          // only type exists because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'oidc' };
          this.visibleAuthMounts = this.mountData;
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
});
