/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response as MirageResponse } from 'miragejs';
import sinon from 'sinon';
import AwsCredentialForm from 'vault/forms/aws-credential';

module('Unit | Component | generate credentials', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const manager = this.owner.lookup('component-manager:glimmer');
    const componentClass = this.owner.factoryFor('component:generate-credentials').class;
    this.component = manager.createComponent(componentClass, {
      named: { backendType: 'aws', backendPath: 'aws-test', roleName: 'example', awsRoleType: 'iam_user' },
    });
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  const cases = [
    ['IAM ignores TTL and ARN', 'iam_user', '3h', 'arn', 'GET', 'creds', null],
    ['federation TTL ignores ARN', 'federation_token', '3h', 'arn', 'POST', 'creds', { ttl: '3h' }],
    ['federation without TTL ignores ARN', 'federation_token', undefined, 'arn', 'POST', 'creds', {}],
    ['session with TTL', 'session_token', '3h', undefined, 'POST', 'creds', { ttl: '3h' }],
    ['session without TTL', 'session_token', undefined, undefined, 'POST', 'creds', {}],
    ['assumed role without ARN or TTL', 'assumed_role', undefined, undefined, 'POST', 'creds', {}],
    ['assumed role with empty TTL', 'assumed_role', '', undefined, 'POST', 'creds', {}],
    ['assumed role with TTL only', 'assumed_role', '3h', undefined, 'POST', 'creds', { ttl: '3h' }],
    [
      'assumed role with ARN and TTL',
      'assumed_role',
      '3h',
      'arn',
      'POST',
      'creds',
      { ttl: '3h', role_arn: 'arn' },
    ],
    ['assumed role with ARN only', 'assumed_role', '', 'arn', 'POST', 'creds', { role_arn: 'arn' }],
  ];

  cases.forEach(([description, credential_type, ttl, role_arn, method, endpoint, body]) => {
    test(description, async function (assert) {
      assert.expect(6);
      this.server[method.toLowerCase()](`/aws-test/${endpoint}/example`, (_, req) => {
        assert.strictEqual(req.method, method, 'Uses the credential-type-specific HTTP method');
        assert.deepEqual(JSON.parse(req.requestBody || 'null'), body, 'Sends only allowed parameters');
        assert.strictEqual(req.queryParams.ttl, undefined, 'TTL is never sent as a query parameter');
        return {
          data: { access_key: 'access', secret_key: 'secret', security_token: 'token' },
          lease_id: 'aws-test/creds/example/lease',
          lease_duration: 3600,
          renewable: false,
        };
      });
      this.component.credentialForm = new AwsCredentialForm({ credential_type, ttl, role_arn });
      await this.component.create(new Event('submit'));
      assert.true(this.component.hasGenerated, 'Successful generation switches to results');
      assert.false(this.component.loading, 'Loading is cleared');
      assert.strictEqual(
        this.component.model.lease_id,
        'aws-test/creds/example/lease',
        'Preserves the lease'
      );
    });
  });

  test('role-read-denied allows selecting a type and updates inputs', function (assert) {
    const manager = this.owner.lookup('component-manager:glimmer');
    const componentClass = this.owner.factoryFor('component:generate-credentials').class;
    const component = manager.createComponent(componentClass, {
      named: { backendType: 'aws', backendPath: 'aws', roleName: 'role' },
    });
    assert.deepEqual(
      component.formFields.map((field) => field.name),
      ['credential_type'],
      'Defaults to IAM'
    );
    assert.true(
      component.helpText.includes('do not have permissions'),
      'Explains the missing role read access'
    );
    component.credentialForm.set('credential_type', 'assumed_role');
    assert.deepEqual(
      component.formFields.map((field) => field.name),
      ['credential_type', 'ttl', 'role_arn'],
      'Selecting assumed role exposes TTL and ARN'
    );
    component.credentialForm.set('credential_type', 'session_token');
    assert.deepEqual(
      component.formFields.map((field) => field.name),
      ['credential_type', 'ttl'],
      'Selecting session token removes the ARN input'
    );
  });

  test('API errors keep the inputs available for retry', async function (assert) {
    const error = new Error('generation failed');
    sinon.stub(this.owner.lookup('service:api').secrets, 'awsGenerateCredentialsRaw').rejects(error);
    this.component.credentialForm.set('ttl', '3h');
    await this.component.create(new Event('submit'));
    assert.strictEqual(this.component.errorMessage, 'generation failed', 'Stores the parsed error message');
    assert.false(this.component.hasGenerated, 'Does not switch to results on failure');
    assert.false(this.component.loading, 'Clears loading on failure');
    assert.strictEqual(this.component.credentialForm.data.ttl, '3h', 'Preserves inputs for retry');
  });

  test('empty success responses are reported as errors, not generated credentials', async function (assert) {
    this.server.get('/aws-test/creds/example', () => ({ data: {} }));
    await this.component.create(new Event('submit'));
    assert.strictEqual(this.component.errorMessage, 'AWS did not return credentials. Please try again.');
    assert.false(this.component.hasGenerated, 'Does not display an empty credentials result');
    assert.false(this.component.loading, 'Allows retrying an incomplete response');
  });

  test('Vault permission errors are parsed without switching to results', async function (assert) {
    this.server.get('/aws-test/creds/example', () => {
      return new MirageResponse(403, {}, { errors: ['permission denied'] });
    });
    await this.component.create(new Event('submit'));
    assert.strictEqual(this.component.errorMessage, 'permission denied', 'Shows the Vault API error');
    assert.false(this.component.hasGenerated, 'Does not show empty credentials for a denied request');
    assert.false(this.component.loading, 'Clears loading after a denied request');
  });

  test('control groups retain the token and show the authorization instructions', async function (assert) {
    const api = this.owner.lookup('service:api');
    const controlGroup = this.owner.lookup('service:control-group');
    const response = { isControlGroupError: true, token: 'wrapped-token', accessor: 'accessor' };
    sinon.stub(api.secrets, 'awsGenerateCredentialsRaw').rejects(new Error('Control Group encountered'));
    sinon.stub(api, 'parseError').resolves({ message: 'Control Group encountered', response });
    const saveToken = sinon.stub(controlGroup, 'saveTokenFromError');
    sinon.stub(controlGroup, 'logFromError').returns({ content: 'Authorize this request before retrying' });
    await this.component.create(new Event('submit'));
    assert.true(saveToken.calledOnceWithExactly(response), 'Saves the token for the next request');
    assert.strictEqual(
      this.component.errorMessage,
      'Authorize this request before retrying',
      'Shows instructions'
    );
    assert.false(this.component.hasGenerated, 'Authorization must finish before showing results');
    assert.false(this.component.loading, 'Retry is available');
  });

  test('reset clears one-time results and restores form defaults', function (assert) {
    this.component.hasGenerated = true;
    this.component.model = { access_key: 'one-time-secret' };
    this.component.errorMessage = 'old error';
    this.component.credentialForm.set('ttl', '');
    this.component.reset();
    assert.false(this.component.hasGenerated, 'Returns to generation');
    assert.strictEqual(this.component.model, undefined, 'Drops previously generated credentials');
    assert.strictEqual(this.component.errorMessage, null, 'Drops previous errors');
    assert.strictEqual(this.component.credentialForm.data.ttl, '3600s', 'Restores the one-hour TTL');
  });

  test('pending generation ignores duplicate requests until credentials arrive', async function (assert) {
    let resolve;
    const request = new Promise((done) => {
      resolve = done;
    });
    const generate = sinon
      .stub(this.owner.lookup('service:api').secrets, 'awsGenerateCredentialsRaw')
      .returns(request);
    const submission = this.component.create(new Event('submit'));
    try {
      assert.true(this.component.loading, 'Marks the request as pending');
      await this.component.create(new Event('submit'));
      assert.strictEqual(generate.callCount, 1, 'Only one credential request is sent');
      assert.false(this.component.hasGenerated, 'Waits for actual credentials before showing results');
    } finally {
      resolve({ raw: new Response(JSON.stringify({ data: { access_key: 'access' } })) });
      await submission;
    }
    assert.true(this.component.hasGenerated, 'Shows credentials after the response');
    assert.false(this.component.loading, 'Clears pending state');
  });
});
