/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | tabs', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.tabData = {
      userpass: [
        {
          path: 'userpass/',
          description: 'platform team only',
          options: {},
          type: 'userpass',
        },
        {
          path: 'userpass2/',
          description: 'backup login method',
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
    this.selectedAuthMethod = '';
    this.handleTabClick = sinon.spy();
    this.renderComponent = () => {
      return render(hbs`
      <Auth::Tabs
        @authTabData={{this.tabData}}
        @handleTabClick={{this.handleTabClick}}
        @selectedAuthMethod={{this.selectedAuthMethod}}
      />`);
    };
  });

  test('it renders tabs', async function (assert) {
    const expectedTabs = [
      { type: 'userpass', display: 'Userpass' },
      { type: 'oidc', display: 'OIDC' },
      { type: 'token', display: 'Token' },
    ];

    await this.renderComponent();
    expectedTabs.forEach((m) => {
      assert.dom(AUTH_FORM.tabBtn(m.type)).exists(`${m.type} renders as a tab`);
      assert.dom(AUTH_FORM.tabBtn(m.type)).hasText(m.display, `${m.type} renders expected display name`);
    });
  });

  test('it selects first tab if no @selectedAuthMethod exists', async function (assert) {
    await this.renderComponent();
    assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
  });

  test('it renders a dropdown if multiple mount paths are returned', async function (assert) {
    this.selectedAuthMethod = 'userpass';
    await this.renderComponent();
    const dropdownOptions = findAll(`${GENERAL.selectByAttr('path')} option`).map((o) => o.value);
    const expectedPaths = ['userpass/', 'userpass2/'];
    expectedPaths.forEach((p) => {
      assert.true(dropdownOptions.includes(p), `dropdown includes path: ${p}`);
    });
  });

  test('it renders a description for the selected mount path', async function (assert) {
    // if a method type does not have any mounts with listing_visibility="unauth" its value is `null`
    // add a tab that has no mount data
    this.tabData.ldap = null;
    // update selected method and re-render component when a tab is clicked
    this.handleTabClick = async (selection) => {
      this.selectedAuthMethod = selection;
      await this.renderComponent();
    };
    this.selectedAuthMethod = 'userpass';

    await this.renderComponent();
    assert
      .dom(GENERAL.selectByAttr('path'))
      .hasValue('userpass/', 'it preselects first mount for "userpass"');
    assert.dom(AUTH_FORM.description).hasText('platform team only', 'it renders the mount description');
    await fillIn(GENERAL.selectByAttr('path'), 'userpass2/');
    assert.dom(GENERAL.selectByAttr('path')).hasValue('userpass2/', 'it selects the second mount');
    assert
      .dom(AUTH_FORM.description)
      .hasText('backup login method', 'it renders relevant description when mount dropdown updates');
    // select a different tab
    await click(AUTH_FORM.tabBtn('token'));
    assert
      .dom(AUTH_FORM.description)
      .hasText('token based credentials', 'it updates the mount description when a new tab is selected');
    // select method without a description
    await click(AUTH_FORM.tabBtn('oidc'));
    assert.dom(AUTH_FORM.tabBtn('oidc')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('my_oidc/');
    assert.dom(AUTH_FORM.description).doesNotExist();
    // select method with no mount data
    await click(AUTH_FORM.tabBtn('ldap'));
    assert.dom(AUTH_FORM.tabBtn('ldap')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.inputByAttr('path')).doesNotExist();
    assert.dom(AUTH_FORM.description).doesNotExist();
  });

  test('it does not render a description if only one tab exists and it does not have mount data', async function (assert) {
    // if a method type does not have any mounts with listing_visibility="unauth" its value is `null`
    this.tabData = { userpass: null };
    this.selectedAuthMethod = 'userpass';
    await this.renderComponent();
    assert.dom(AUTH_FORM.tabBtn('userpass')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.inputByAttr('path')).doesNotExist();
    assert.dom(AUTH_FORM.description).doesNotExist();
  });

  test('it renders hidden input if only one mount path is returned', async function (assert) {
    this.selectedAuthMethod = 'oidc';
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('my_oidc/');
  });

  test('it calls handleTabClick with tab method type', async function (assert) {
    await this.renderComponent();
    await click(AUTH_FORM.tabBtn('oidc'));
    const [actual] = this.handleTabClick.lastCall.args;
    assert.strictEqual(actual, 'oidc');
  });
});
