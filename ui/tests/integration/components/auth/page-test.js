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
import { fillInLoginFields, VISIBLE_MOUNTS } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CSP_ERROR } from 'vault/components/auth/page';

module('Integration | Component | auth | page', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.cluster = { id: '1' };
    this.onAuthSuccess = sinon.spy();
    this.onNamespaceUpdate = sinon.spy();
    this.visibleAuthMounts = null;
    this.directLinkData = null;
    this.loginSettings = null;

    this.renderComponent = () => {
      return render(hbs`
        <Auth::Page
          @cluster={{this.cluster}}
          @directLinkData={{this.directLinkData}}
          @loginSettings={{this.loginSettings}}
          @namespaceQueryParam={{this.nsQp}}
          @oidcProviderQueryParam={{this.providerQp}}
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

  test('it renders splash logo when oidc provider query param is present', async function (assert) {
    this.providerQp = 'myprovider';
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert
      .dom(AUTH_FORM.helpText)
      .hasText(
        'Once you log in, you will be redirected back to your application. If you require login credentials, contact your administrator.'
      );
  });

  test('it disables namespace input when oidc provider query param is present', async function (assert) {
    this.providerQp = 'myprovider';
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    assert.dom(AUTH_FORM.logo).exists();
    assert.dom(GENERAL.inputByAttr('namespace')).isDisabled();
  });

  test('it calls onNamespaceUpdate', async function (assert) {
    assert.expect(1);
    this.version.features = ['Namespaces'];
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  test('it calls onNamespaceUpdate for HVD managed clusters', async function (assert) {
    assert.expect(2);
    this.version.features = ['Namespaces'];
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.nsQp = 'admin';
    await this.renderComponent();

    assert.dom(GENERAL.inputByAttr('namespace')).hasValue('');
    await fillIn(GENERAL.inputByAttr('namespace'), 'mynamespace');
    const [actual] = this.onNamespaceUpdate.lastCall.args;
    assert.strictEqual(actual, 'mynamespace', `onNamespaceUpdate called with: ${actual}`);
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.visibleAuthMounts = VISIBLE_MOUNTS;
      window.localStorage.clear();
    });

    test('it formats tab data if visible auth mounts exist', async function (assert) {
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

    test('it selects type in the dropdown if @directLinkData references NON visible type', async function (assert) {
      this.directLinkData = { type: 'okta', isVisibleMount: false };
      await this.renderComponent();
      assert.dom(GENERAL.selectByAttr('auth type')).hasValue('okta', 'dropdown has type selected');
      assert.dom(AUTH_FORM.authForm('okta')).exists();
      assert.dom(GENERAL.inputByAttr('username')).exists();
      assert.dom(GENERAL.inputByAttr('password')).exists();
      await click(AUTH_FORM.advancedSettings);
      assert.dom(GENERAL.inputByAttr('path')).exists();
      assert.dom(AUTH_FORM.tabBtn('okta')).doesNotExist('tab does not render');
      assert
        .dom(GENERAL.backButton)
        .exists('back button renders because listing_visibility="unauth" for other mounts');
      assert.dom(AUTH_FORM.otherMethodsBtn).doesNotExist('"Sign in with other methods" does not render');
    });

    test('it renders single mount view instead of tabs if @directLinkData data references a visible type', async function (assert) {
      this.directLinkData = { path: 'my-oidc/', type: 'oidc', isVisibleMount: true };
      await this.renderComponent();
      assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders tab for type');
      assert.dom(GENERAL.inputByAttr('role')).exists();
      assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
      assert.dom(AUTH_FORM.otherMethodsBtn).exists('"Sign in with other methods" renders');
      assert.dom(GENERAL.selectByAttr('auth type')).doesNotExist();
      assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
      assert.dom(GENERAL.backButton).doesNotExist();
    });
  });

  /* 
  Login settings are an enterprise only feature but the component is version agnostic
  because fetching login customizations happens in the route for enterprise clusters only.
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
        this.visibleAuthMounts = VISIBLE_MOUNTS;
      });

      test('(default+backups): it hides advanced settings for both views', async function (assert) {
        await this.renderComponent();
        assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
        assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
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
        assert.dom(AUTH_FORM.tabBtn('oidc')).hasText('OIDC', 'it renders default method');
        assert.dom(AUTH_FORM.tabs).exists({ count: 1 }, 'only one tab renders');
        assert.dom(AUTH_FORM.authForm('oidc')).exists();
        this.assertPathInput(assert, { isHidden: true, value: 'my-oidc/' });
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
        this.mountData = (path) => ({ [path]: VISIBLE_MOUNTS[path] });
      });

      test('(default+backups): it hides advanced settings for default with visible mount but it renders for backups', async function (assert) {
        this.visibleAuthMounts = { ...this.mountData('my-oidc/') };
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
        this.visibleAuthMounts = {
          ...this.mountData('my-oidc/'),
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
        this.mountData = VISIBLE_MOUNTS;
      });

      module('when there are no visible mounts at all', function (hooks) {
        hooks.beforeEach(function () {
          this.visibleAuthMounts = null;
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
          this.visibleAuthMounts = {
            ...this.mountData,
            'my-okta/': {
              description: '',
              options: null,
              type: 'okta',
            },
          };
          this.directLinkData = { path: 'my-okta/', type: 'okta', isVisibleMount: true };
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
          // isVisibleMount is false because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'okta', isVisibleMount: false };
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
          // isVisibleMount is false because the query param does not match a path with listing_visibility="unauth"
          this.directLinkData = { type: 'oidc', isVisibleMount: false };
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
      // await fillIn(AUTH_FORM.selectMethod, authType);
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
});
