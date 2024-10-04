/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, _cancelTimers as cancelTimers } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

module('Integration | Component | auth form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.selectedAuth = 'token';
    this.performAuth = sinon.spy();
    this.renderComponent = async () => {
      return render(hbs`
        <AuthForm
          @wrappedToken={{this.wrappedToken}}
          @cluster={{this.cluster}}
          @selectedAuth={{this.selectedAuth}}
          @performAuth={{this.performAuth}}
          @authIsRunning={{this.authIsRunning}}
          @delayIsIdle={{this.delayIsIdle}}
        />`);
    };
  });

  test('it calls performAuth on submit', async function (assert) {
    await this.renderComponent();
    await fillIn(AUTH_FORM.input('token'), '123token');
    await click(AUTH_FORM.login);
    const [type, data] = this.performAuth.lastCall.args;
    assert.strictEqual(type, 'token', 'performAuth is called with type');
    assert.propEqual(data, { token: '123token' }, 'performAuth is called with data');
  });

  test('it disables sign in button when authIsRunning', async function (assert) {
    this.authIsRunning = true;
    await this.renderComponent();
    assert.dom(AUTH_FORM.login).isDisabled('sign in button is disabled');
    assert.dom(`${AUTH_FORM.login} [data-test-icon="loading"]`).exists('sign in button renders loading icon');
  });

  test('it renders alert info message when delayIsIdle', async function (assert) {
    this.delayIsIdle = true;
    this.authIsRunning = true;
    await this.renderComponent();
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'If login takes longer than usual, you may need to check your device for an MFA notification, or contact your administrator if login times out.'
      );
  });

  test('it renders no tabs when no methods are passed', async function (assert) {
    this.server.get('/sys/internal/ui/mounts', () => {
      return {
        data: {
          auth: {
            'approle/': {
              type: 'approle',
            },
          },
        },
      };
    });
    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).doesNotExist();
  });

  test('it renders all the supported methods and Other tab when methods are present', async function (assert) {
    this.server.get('/sys/internal/ui/mounts', () => {
      return {
        data: {
          auth: {
            'foo/': {
              type: 'userpass',
            },
            'approle/': {
              type: 'approle',
            },
          },
        },
      };
    });

    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).exists({ count: 2 });
    assert.dom(AUTH_FORM.tabs('foo')).exists('tab uses the path in the label');
    assert.dom(AUTH_FORM.tabs('other')).exists('second tab is the Other tab');
  });

  test('it renders the description', async function (assert) {
    this.selectedAuth = null;
    this.server.get('/sys/internal/ui/mounts', () => {
      return {
        data: {
          auth: {
            'approle/': {
              type: 'userpass',
              description: 'app description',
            },
          },
        },
      };
    });
    await this.renderComponent();
    assert.dom(AUTH_FORM.description).hasText('app description');
  });

  test('it renders no tabs when no supported methods are present in passed methods', async function (assert) {
    const methods = {
      'approle/': {
        type: 'approle',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });
    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).doesNotExist();
  });

  test('it shows an error if unwrap errors', async function (assert) {
    assert.expect(1);
    this.wrappedToken = '54321';
    this.server.post('/sys/wrapping/unwrap', () => {
      return new Response(
        400,
        { 'Content-Type': 'application/json' },
        { errors: ['There was an error unwrapping!'] }
      );
    });

    await this.renderComponent();
    later(() => cancelTimers(), 50);
    await settled();
    assert.dom(GENERAL.messageError).hasText('Error Token unwrap failed: There was an error unwrapping!');
  });

  test('it should retain oidc role when mount path is changed', async function (assert) {
    assert.expect(2);

    const auth_url = 'http://dev-foo-bar.com';
    this.server.post('/auth/:path/oidc/auth_url', (_, req) => {
      const { role, redirect_uri } = JSON.parse(req.requestBody);
      const goodRequest =
        req.params.path === 'foo-oidc' &&
        role === 'foo' &&
        redirect_uri.includes('/auth/foo-oidc/oidc/callback');

      return new Response(
        goodRequest ? 200 : 400,
        { 'Content-Type': 'application/json' },
        JSON.stringify(
          goodRequest ? { data: { auth_url } } : { errors: [`role "${role}" could not be found`] }
        )
      );
    });
    window.open = (url) => {
      assert.strictEqual(url, auth_url, 'auth_url is returned when required params are passed');
    };

    this.owner.lookup('service:router').reopen({
      urlFor(route, { auth_path }) {
        return `/auth/${auth_path}/oidc/callback`;
      },
    });

    await this.renderComponent();

    await fillIn(GENERAL.selectByAttr('auth-method'), 'oidc');
    await fillIn(AUTH_FORM.input('role'), 'foo');
    await click(AUTH_FORM.moreOptions);
    await fillIn(AUTH_FORM.input('role'), 'foo');
    await fillIn(AUTH_FORM.mountPathInput, 'foo-oidc');
    assert.dom(AUTH_FORM.input('role')).hasValue('foo', 'role is retained when mount path is changed');
    await click(AUTH_FORM.login);
  });
});
