/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { click } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { SYS_INTERNAL_UI_MOUNTS } from 'vault/tests/helpers/auth/auth-helpers';
import setupTestContext from './setup-test-context';

/* 
  Login settings are an enterprise only feature but the component is version agnostic (and subsequently so are these tests) 
  because fetching login settings happens in the route only for enterprise versions.
  Each combination must be tested with and without visible mounts (i.e. tuned with listing_visibility="unauth")
  1. default+backups: default type set, backup types set 
  2. default only: no backup types
  3. backup only: backup types set without a default
   */
module('Integration | Component | auth | page | ent login settings', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    setupTestContext(this);
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
    await click(GENERAL.button('Sign in with other methods'));
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
    await click(GENERAL.button('Sign in with other methods'));
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
    assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
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
    assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
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
      await click(GENERAL.button('Sign in with other methods'));
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
      assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
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
      assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
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
      await click(GENERAL.button('Sign in with other methods'));
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
      await click(GENERAL.button('Sign in with other methods'));
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
      await click(GENERAL.button('Sign in with other methods'));
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
        assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
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
        await click(GENERAL.button('Sign in with other methods'));
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
        assert.dom(GENERAL.button('Sign in with other methods')).doesNotExist();
        await click(GENERAL.backButton);
        assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
        await click(GENERAL.button('Sign in with other methods'));
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
        await click(GENERAL.button('Sign in with other methods'));
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
