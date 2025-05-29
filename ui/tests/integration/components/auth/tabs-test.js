/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, findAll, render } from '@ember/test-helpers';
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
      { type: 'userpass', display: 'Username & Password' },
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

  test('it renders the mount description', async function (assert) {
    this.selectedAuthMethod = 'token';
    await this.renderComponent();
    assert.dom(AUTH_FORM.description).hasText('token based credentials');
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

  test('it renders hidden input if only one mount path is returned', async function (assert) {
    this.selectedAuthMethod = 'oidc';
    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('path')).hasAttribute('type', 'hidden');
    assert.dom(GENERAL.inputByAttr('path')).hasValue('my-oidc/');
  });

  test('it calls handleTabClick with tab method type', async function (assert) {
    await this.renderComponent();
    await click(AUTH_FORM.tabBtn('oidc'));
    const [actual] = this.handleTabClick.lastCall.args;
    assert.strictEqual(actual, 'oidc');
  });
});
