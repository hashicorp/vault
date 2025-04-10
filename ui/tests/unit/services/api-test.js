/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

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
    this.deleteControlGroupToken = sinon.spy(controlGroupService, 'deleteControlGroupToken');

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
    const warnings = ['warning1', 'warning2'];
    const response = new Response(JSON.stringify({ warnings }));

    await this.apiService.showWarnings({ response });

    assert.true(this.info.firstCall.calledWith(warnings[0]), 'First warning message is shown');
    assert.true(this.info.secondCall.calledWith(warnings[1]), 'Second warning message is shown');
  });

  test('it should delete control group token', async function (assert) {
    await this.apiService.deleteControlGroupToken({ url: this.url });

    assert.true(this.tokenForUrl.calledWith(this.url), 'Url is passed to tokenForUrl method');
    assert.true(
      this.deleteControlGroupToken.calledWith(this.wrapInfo.accessor),
      'Control group token is deleted'
    );
  });

  test('it should format error response', async function (assert) {
    const e = { data: { error: 'Something went wrong' } };
    const response = new Response(JSON.stringify(e), { status: 400 });

    const errorResponse = await this.apiService.formatErrorResponse({ response, url: this.url });
    const error = await errorResponse.json();
    const expectedError = {
      ...e,
      httpStatus: 400,
      path: this.url,
      errors: ['Something went wrong'],
    };

    assert.deepEqual(error, expectedError, 'Error is reformated and returned');
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
});
