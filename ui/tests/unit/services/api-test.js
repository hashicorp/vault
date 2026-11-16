/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import config from 'vault/config/environment';
import { RequiredError } from '@hashicorp/vault-client-typescript';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Unit | Service | api', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.apiService = this.owner.lookup('service:api');

    const authService = this.owner.lookup('service:auth');
    this.setLastFetch = sinon.spy(authService, 'setLastFetch');
    this.currentToken = sinon.stub(authService, 'currentToken').value('foobar');

    const namespaceService = this.owner.lookup('service:namespace');
    this.namespace = sinon.stub(namespaceService, 'path').value('another-ns');

    const controlGroupService = this.owner.lookup('service:control-group');
    this.wrapInfo = { token: 'ctrl-group', accessor: '84tfdfd5pQ5vOOEMxC2o3Ymt' };
    this.tokenForUrl = sinon.stub(controlGroupService, 'tokenForUrl').returns(this.wrapInfo);
    this.tokenToUnwrap = sinon.stub(controlGroupService, 'tokenToUnwrap').value(this.wrapInfo);
    this.deleteControlGroupToken = sinon.spy(controlGroupService, 'deleteControlGroupToken');
    this.isRequestedPathLocked = sinon.stub(controlGroupService, 'isRequestedPathLocked').returns(true);

    const flashMessageService = this.owner.lookup('service:flash-messages');
    this.info = sinon.spy(flashMessageService, 'info');

    this.url = '/v1/sys/capabilities-self';
  });

  test('it should set last fetch time', async function (assert) {
    await this.apiService.setLastFetch({ url: '/v1/sys/health' });
    assert.true(this.setLastFetch.notCalled, 'Last fetch is not set for polling url');

    await this.apiService.setLastFetch({ url: '/v1/auth/token/lookup-self' });
    assert.true(this.setLastFetch.calledOnce, 'Last fetch is set for non polling url');
  });

  test('it should get control group token', async function (assert) {
    const context = {
      url: this.url,
      init: {
        method: 'GET',
        headers: { 'X-Vault-Token': 'root' },
      },
    };

    this.tokenForUrl.returns(undefined);
    const noTokenContext = await this.apiService.getControlGroupToken(context);

    assert.true(this.tokenForUrl.calledWith(context.url), 'Url is passed to tokenForUrl method');
    assert.deepEqual(context, noTokenContext, 'Original context is returned when no token is present');

    this.tokenForUrl.returns(this.wrapInfo);

    const { token } = this.wrapInfo;
    const tokenContext = await this.apiService.getControlGroupToken(context);
    const newContext = {
      url: '/v1/sys/wrapping/unwrap',
      init: {
        method: 'POST',
        headers: { 'X-Vault-Token': token },
        body: JSON.stringify({ token }),
      },
    };

    assert.deepEqual(tokenContext, newContext, 'New context is returned when token is present');
  });

  test('it should set default headers', async function (assert) {
    const {
      init: { headers },
    } = await this.apiService.setHeaders({ init: { method: 'PATCH' } });

    assert.strictEqual(
      headers.get('X-Vault-Token'),
      'foobar',
      'Token header is set with value from auth service'
    );
    assert.strictEqual(
      headers.get('X-Vault-Namespace'),
      'another-ns',
      'Namespace header is set with value from namespace service'
    );
    assert.strictEqual(
      headers.get('Content-Type'),
      'application/merge-patch+json',
      'Content type header is set for PATCH method'
    );
  });

  test('it should override default headers when set on request init', async function (assert) {
    const initHeaders = {
      'X-Vault-Token': 'root',
      'X-Vault-Namespace': 'ns1',
    };

    const {
      init: { headers },
    } = await this.apiService.setHeaders({ init: { headers: initHeaders } });

    assert.strictEqual(headers.get('X-Vault-Token'), 'root', 'Token header set on request init is preserved');
    assert.strictEqual(
      headers.get('X-Vault-Namespace'),
      'ns1',
      'Namespace header set on request init is preserved'
    );
  });

  test('it should show warnings', async function (assert) {
    const warnings = JSON.stringify({ warnings: ['warning1', 'warning2'] });
    const response = new Response(warnings, { headers: { 'Content-Type': 'application/json' } });

    await this.apiService.showWarnings({ response });

    assert.true(this.info.firstCall.calledWith('warning1'), 'First warning message is shown');
    assert.true(this.info.secondCall.calledWith('warning2'), 'Second warning message is shown');
  });

  test('it should not attempt to set warnings for empty response', async function (assert) {
    const response = new Response();
    await this.apiService.showWarnings({ response });
    assert.true(this.info.notCalled, 'No warning messages are shown');
  });

  test('it should check for control group', async function (assert) {
    const headers = new Headers({ 'Content-Type': 'application/json', 'X-Vault-Wrap-TTL': 1800 });
    const body = { data: null, wrap_info: this.wrapInfo };
    const init = { headers: new Headers({ 'X-Vault-Token': this.wrapInfo.token }) };
    const apiResponse = new Response(JSON.stringify(body), { headers });

    const response = await this.apiService.checkControlGroup({ url: this.url, response: apiResponse, init });

    assert.true(
      this.deleteControlGroupToken.calledWith(this.wrapInfo.accessor),
      'Control group token is deleted'
    );
    assert.true(this.isRequestedPathLocked.calledWith(body, '1800'), 'isRequestedPathLocked called');

    assert.strictEqual(response.status, 403, 'Response status is updated to 403 for control group error');
    const ctrlError = await response.json();
    const expectedError = {
      message: 'Control Group encountered',
      isControlGroupError: true,
      ...this.wrapInfo,
    };
    assert.deepEqual(ctrlError, expectedError, 'Control group error is returned in response body');
  });

  test('it should build headers', async function (assert) {
    const headerMap = {
      token: 'foobar',
      namespace: 'ns1',
      wrap: '10s',
    };

    const token = await this.apiService.buildHeaders({ token: headerMap.token });
    assert.deepEqual(token.headers, { 'X-Vault-Token': headerMap.token }, 'Token header is set');

    const namespace = await this.apiService.buildHeaders({ namespace: headerMap.namespace });
    assert.deepEqual(
      namespace.headers,
      { 'X-Vault-Namespace': headerMap.namespace },
      'Namespace header is set'
    );

    const wrapTTL = await this.apiService.buildHeaders({ wrap: headerMap.wrap });
    assert.deepEqual(wrapTTL.headers, { 'X-Vault-Wrap-TTL': '10s' }, 'Wrap TTL header is set');

    const multi = await this.apiService.buildHeaders(headerMap);
    assert.deepEqual(
      multi.headers,
      {
        'X-Vault-Token': 'foobar',
        'X-Vault-Namespace': 'ns1',
        'X-Vault-Wrap-TTL': '10s',
      },
      'All supported headers are set'
    );
  });

  test('it should map list response to array', async function (assert) {
    const response = {
      key_info: {
        key1: { title: 'Title 1' },
        key2: { title: 'Title 2' },
      },
      keys: ['key1', 'key2'],
    };
    const expectedResponse = [
      { id: 'key1', title: 'Title 1' },
      { id: 'key2', title: 'Title 2' },
    ];
    const listData = this.apiService.keyInfoToArray(response);
    assert.deepEqual(listData, expectedResponse, 'List response is mapped to array');
  });

  module('Raw requests', function (hooks) {
    hooks.beforeEach(function () {
      this.tokenForUrl.returns(null);
      this.tokenToUnwrap.value(null);
      this.isRequestedPathLocked.returns(false);
      this.fetch = sinon.stub(window, 'fetch');
      this.respond = (body, status = 200) =>
        this.fetch.resolves(
          new Response(body === undefined ? null : JSON.stringify(body), {
            status,
            headers: body === undefined ? {} : { 'Content-Type': 'application/json' },
          })
        );
      this.lastRequest = () => {
        const [url, init] = this.fetch.lastCall.args;
        return { url, method: init.method, body: init.body, headers: new Headers(init.headers) };
      };
    });

    hooks.afterEach(function () {
      this.fetch.restore();
    });

    test('it sends the path, method, body and default headers', async function (assert) {
      this.respond({ data: { ok: true } });
      const response = await this.apiService.request.post('/auth/userpass/users/bob', { password: 'p' });

      const { url, method, body, headers } = this.lastRequest();
      assert.strictEqual(url, '/v1/auth/userpass/users/bob', 'path is prefixed with the API base path');
      assert.strictEqual(method, 'POST', 'method is sent');
      assert.strictEqual(body, JSON.stringify({ password: 'p' }), 'body is serialized');
      assert.strictEqual(headers.get('X-Vault-Token'), 'foobar', 'token header is set');
      assert.strictEqual(headers.get('X-Vault-Namespace'), 'another-ns', 'namespace header is set');
      assert.deepEqual(await response.json(), { data: { ok: true } }, 'the raw response is returned');
    });

    test('it keeps headers passed by the caller', async function (assert) {
      this.respond({ data: {} });
      await this.apiService.request.put('/sys/unseal', { key: 'k' }, { 'X-Vault-Namespace': '' });
      assert.strictEqual(this.lastRequest().headers.get('X-Vault-Namespace'), '', 'explicit header wins');
    });

    test('it throws a ResponseError for unsuccessful responses', async function (assert) {
      this.respond({ errors: ['permission denied'] }, 403);
      try {
        await this.apiService.request.delete('/auth/userpass/users/bob');
        assert.ok(false, 'request should reject');
      } catch (error) {
        const { status, message } = await this.apiService.parseError(error);
        assert.strictEqual(status, 403, 'status is parsed');
        assert.strictEqual(message, 'permission denied', 'message is parsed');
      }
    });

    test('it sets the last fetch time', async function (assert) {
      this.respond({ data: {} });
      await this.apiService.request.get('/auth/userpass/users/bob');
      assert.true(this.setLastFetch.calledOnce, 'last fetch is set for raw requests');
    });

    test('it shows response warnings', async function (assert) {
      this.respond({ data: {}, warnings: ['token_period exceeds max TTL'] });
      await this.apiService.request.post('/auth/cert/certs/web', { token_period: '100h' });
      assert.true(this.info.calledOnceWith('token_period exceeds max TTL'), 'warning is flashed');
    });

    test('it tolerates empty successful responses', async function (assert) {
      this.respond(undefined, 204);
      const response = await this.apiService.request.delete('/auth/userpass/users/bob');
      assert.strictEqual(response.status, 204, 'empty response resolves');
      assert.true(this.info.notCalled, 'no warnings are flashed');
    });

    test('it turns a Control Group response into a control group error', async function (assert) {
      this.isRequestedPathLocked.returns(true);
      this.respond({ data: null, wrap_info: { ...this.wrapInfo, creation_path: 'auth/userpass/users/bob' } });
      try {
        await this.apiService.request.post('/auth/userpass/users/bob', { password: 'p' });
        assert.ok(false, 'request should reject');
      } catch (error) {
        const { status, response } = await this.apiService.parseError(error);
        assert.strictEqual(status, 403, 'locked request rejects with 403');
        assert.true(response.isControlGroupError, 'response is flagged as a control group error');
        assert.strictEqual(response.accessor, this.wrapInfo.accessor, 'wrap info is included');
      }
    });

    test('it unwraps an approved Control Group token instead of repeating the request', async function (assert) {
      this.tokenForUrl.returns(this.wrapInfo);
      this.tokenToUnwrap.value(this.wrapInfo);
      this.respond({ data: { token_policies: ['default'] } });
      const response = await this.apiService.request.get('/auth/userpass/users/bob');

      const { url, method, body, headers } = this.lastRequest();
      assert.strictEqual(url, '/v1/sys/wrapping/unwrap', 'request is sent to the unwrap endpoint');
      assert.strictEqual(method, 'POST', 'unwrap uses POST');
      assert.strictEqual(body, JSON.stringify({ token: this.wrapInfo.token }), 'wrapped token is sent');
      assert.strictEqual(headers.get('X-Vault-Token'), this.wrapInfo.token, 'wrapped token authenticates');
      assert.true(this.deleteControlGroupToken.calledWith(this.wrapInfo.accessor), 'used token is deleted');
      assert.deepEqual(
        await response.json(),
        { data: { token_policies: ['default'] } },
        'unwrapped data is returned'
      );
    });
  });

  module('Error parsing', function () {
    test('it should correctly parse message from error', async function (assert) {
      let e = await this.apiService.parseError(getErrorResponse(undefined, 400));
      assert.strictEqual(e.message, 'first error, second error', 'Builds message from errors');

      e = await this.apiService.parseError(
        getErrorResponse({ errors: [], message: 'there were some errors' }, 400)
      );
      assert.strictEqual(e.message, 'there were some errors', 'Returns message when errors are empty');

      const error = new Error('some js type error');
      e = await this.apiService.parseError(error);
      assert.strictEqual(e.message, error.message, 'Returns message from generic Error');

      e = await this.apiService.parseError('some random error');
      assert.strictEqual(e.message, 'An error occurred, please try again', 'Returns default fallback');

      const fallback = 'Everything is broken, sorry';
      e = await this.apiService.parseError('some random error', fallback);
      assert.strictEqual(e.message, fallback, 'Returns custom fallback');
    });

    test('it should return status', async function (assert) {
      const { status } = await this.apiService.parseError(getErrorResponse());
      assert.strictEqual(status, 404, 'Returns the status code from the response');
    });

    test('it should return path', async function (assert) {
      const { path } = await this.apiService.parseError(getErrorResponse());
      assert.strictEqual(path, '/v1/test/error/parsing', 'Returns the path from the request url');
    });

    test('it should return error response', async function (assert) {
      const error = {
        errors: ['something bad happened', 'something else bad too'],
        message: 'all bad things occurred',
      };
      const { response } = await this.apiService.parseError(getErrorResponse(error, 400));
      assert.deepEqual(response, error, 'Returns the original error response');
    });

    test('it should log out error in development environment', async function (assert) {
      const consoleStub = sinon.stub(console, 'error');
      sinon.stub(config, 'environment').value('development');
      const error = new Error('some js type error');
      await this.apiService.parseError(error);
      assert.true(consoleStub.calledWith('API Error:', error));
      sinon.restore();
    });

    test('it should return the generic fallback message for RequiredError instead of exposing internals, while logging the real error', async function (assert) {
      const consoleStub = sinon.stub(console, 'error');
      const error = new RequiredError(
        'EntityUpdateAliasByIdRequest',
        'Required parameter "EntityUpdateAliasByIdRequest" was null or undefined when calling entityUpdateAliasById().'
      );
      const e = await this.apiService.parseError(error);
      assert.strictEqual(
        e.message,
        'An error occurred, please try again',
        'Returns the generic fallback instead of the raw internal error'
      );
      assert.true(consoleStub.calledWith(error), 'Logs the real error for debugging visibility');

      const fallback = 'Something went wrong saving this item';
      const eWithCustomFallback = await this.apiService.parseError(error, fallback);
      assert.strictEqual(eWithCustomFallback.message, fallback, 'Returns the custom fallback when provided');
      sinon.restore();
    });
  });
});
