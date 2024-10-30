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
import { validate } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

module('Integration | Component | auth | page ', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.auth = this.owner.lookup('service:auth');
    this.cluster = { id: '1' };
    this.selectedAuth = 'token';
    this.onSuccess = sinon.spy();

    this.renderComponent = async () => {
      return render(hbs`
        <Auth::LoginForm
          @wrappedToken={{this.wrappedToken}}
          @cluster={{this.cluster}}
          @namespace={{this.namespaceQueryParam}}
          @selectedAuth={{this.authMethod}}
          @onSuccess={{this.onSuccess}}
        />
        `);
    };
  });
  const CSP_ERR_TEXT = `Error This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
  test('it renders error on CSP violation', async function (assert) {
    assert.expect(2);
    this.cluster.standby = true;
    await this.renderComponent();
    assert.dom(GENERAL.messageError).doesNotExist();
    this.owner.lookup('service:csp-event').handleEvent({ violatedDirective: 'connect-src' });
    await settled();
    assert.dom(GENERAL.messageError).hasText(CSP_ERR_TEXT);
  });

  test('it renders with vault style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['Not allowed'] });
    });

    await this.renderComponent();
    await click(AUTH_FORM.login);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: Not allowed');
  });

  test('it renders AdapterError style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['API Error here'] });
    });

    await this.renderComponent();
    await click(AUTH_FORM.login);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error Authentication failed: API Error here', 'shows the error from the API');
  });

  test('it calls auth service authenticate method with expected args', async function (assert) {
    assert.expect(1);
    const authenticateStub = sinon.stub(this.auth, 'authenticate');
    this.selectedAuth = 'foo/'; // set to a non-default path
    this.server.get('/sys/internal/ui/mounts', () => {
      return {
        data: {
          auth: {
            'foo/': {
              type: 'userpass',
            },
          },
        },
      };
    });

    await this.renderComponent();
    await fillIn(AUTH_FORM.input('username'), 'sandy');
    await fillIn(AUTH_FORM.input('password'), '1234');
    await click(AUTH_FORM.login);
    const [actual] = authenticateStub.lastCall.args;
    const expectedArgs = {
      backend: 'userpass',
      clusterId: '1',
      data: {
        username: 'sandy',
        password: '1234',
        path: 'foo',
      },
      selectedAuth: 'foo/',
    };
    assert.propEqual(
      actual,
      expectedArgs,
      `it calls auth service authenticate method with expected args: ${JSON.stringify(actual)} `
    );
  });

  test('it calls onSuccess with expected args', async function (assert) {
    assert.expect(3);
    this.server.get(`auth/token/lookup-self`, () => {
      return {
        data: {
          policies: ['default'],
        },
      };
    });

    await this.renderComponent();
    await fillIn(AUTH_FORM.input('token'), 'mytoken');
    await click(AUTH_FORM.login);
    const [authResponse, backendType, data] = this.onSuccess.lastCall.args;
    const expected = { isRoot: false, namespace: '', token: 'vault-token☃1' };

    assert.propEqual(
      authResponse,
      expected,
      `it calls onSuccess with response: ${JSON.stringify(authResponse)} `
    );
    assert.strictEqual(backendType, 'token', `it calls onSuccess with backend type: ${backendType}`);
    assert.propEqual(data, { token: 'mytoken' }, `it calls onSuccess with data: ${JSON.stringify(data)}`);
  });

  test('it makes a request to unwrap if passed a wrappedToken and logs in', async function (assert) {
    assert.expect(3);
    const authenticateStub = sinon.stub(this.auth, 'authenticate');
    this.wrappedToken = '54321';

    this.server.post('/sys/wrapping/unwrap', (_, req) => {
      assert.strictEqual(req.url, '/v1/sys/wrapping/unwrap', 'makes call to unwrap the token');
      assert.strictEqual(
        req.requestHeaders['X-Vault-Token'],
        this.wrappedToken,
        'uses passed wrapped token for the unwrap'
      );
      return {
        auth: {
          client_token: '12345',
        },
      };
    });

    await this.renderComponent();
    later(() => cancelTimers(), 50);
    await settled();
    const [actual] = authenticateStub.lastCall.args;
    assert.propEqual(
      actual,
      {
        backend: 'token',
        clusterId: '1',
        data: {
          token: '12345',
        },
        selectedAuth: 'token',
      },
      `it calls auth service authenticate method with correct args: ${JSON.stringify(actual)} `
    );
  });

  test('it should set nonce value as uuid for okta method type', async function (assert) {
    assert.expect(4);
    this.server.post('/auth/okta/login/foo', (_, req) => {
      const { nonce } = JSON.parse(req.requestBody);
      assert.true(validate(nonce), 'Nonce value passed as uuid for okta login');
      return {
        auth: {
          client_token: '12345',
          policies: ['default'],
        },
      };
    });

    await this.renderComponent();

    await fillIn(GENERAL.selectByAttr('auth-method'), 'okta');
    await fillIn(AUTH_FORM.input('username'), 'foo');
    await fillIn(AUTH_FORM.input('password'), 'bar');
    await click(AUTH_FORM.login);
    assert
      .dom('[data-test-okta-number-challenge]')
      .hasText(
        'To finish signing in, you will need to complete an additional MFA step. Please wait... Back to login',
        'renders okta number challenge on submit'
      );
    await click(GENERAL.backButton);
    assert.dom(AUTH_FORM.form).exists('renders auth form on return to login');
    assert.dom(GENERAL.selectByAttr('auth-method')).hasValue('okta', 'preserves method type on back');
  });
});
