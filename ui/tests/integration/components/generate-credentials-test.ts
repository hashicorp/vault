/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import type ApiService from 'vault/services/api';

module('Integration | Component | generate credentials', function (hooks) {
  setupRenderingTest(hooks);

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('binds TTL and role ARN through the form and renders masked credentials with lease metadata', async function (assert) {
    const api = this.owner.lookup('service:api') as ApiService;
    const generate = sinon.stub(api.secrets, 'awsGenerateCredentialsWithParametersRaw').resolves({
      raw: new Response(
        JSON.stringify({
          data: { access_key: 'access-value', secret_key: 'secret-value', security_token: 'token-value' },
          lease_id: 'aws/creds/example/lease',
          lease_duration: 3600,
          renewable: false,
        })
      ),
    });
    const copy = sinon.stub(navigator.clipboard, 'writeText').resolves();
    await render(hbs`
      <GenerateCredentials @backendType="aws" @backendPath="aws" @roleName="example" @awsRoleType="assumed_role" />
    `);
    assert.dom(GENERAL.toggleInput('TTL')).isChecked('Uses the one-hour default TTL');
    await click(GENERAL.toggleInput('TTL'));
    await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:example');
    await click(GENERAL.submitButton);
    assert.deepEqual(
      generate.firstCall.args[0],
      {
        name: 'example',
        aws_mount_path: 'aws',
        AwsGenerateCredentialsWithParametersRequest: { role_arn: 'arn:example' },
      },
      'FormField updates reach the request and disabled TTL is omitted'
    );
    assert.dom('[data-test-warning]').exists('Warns that credentials are only available once');
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('access-value', 'Renders the access key');
    assert
      .dom(GENERAL.infoRowValue('Lease ID'))
      .hasText('aws/creds/example/lease', 'Renders top-level lease ID');
    assert.dom(GENERAL.infoRowValue('Lease duration')).hasText('3600', 'Renders lease duration');
    assert.dom(GENERAL.infoRowValue('Renewable')).hasText('No', 'Renders false renewability');
    assert.dom(GENERAL.infoRowValue('Secret key')).includesText('***********', 'Masks the secret key');
    assert
      .dom(GENERAL.infoRowValue('Security token'))
      .includesText('***********', 'Masks the security token');
    await click(`${GENERAL.infoRowValue('Secret key')} [data-test-button="toggle-masked"]`);
    assert.dom(GENERAL.infoRowValue('Secret key')).includesText('secret-value', 'Can reveal the secret');
    await click(`${GENERAL.infoRowValue('Security token')} [data-test-copy-button]`);
    assert.true(copy.calledWith('token-value'), 'Copies the unmasked token');
    await click('[data-test-copy-credentials]');
    assert.deepEqual(
      JSON.parse(copy.lastCall.args[0]),
      {
        accessKey: 'access-value',
        secretKey: 'secret-value',
        securityToken: 'token-value',
        leaseId: 'aws/creds/example/lease',
      },
      'Preserves the existing clipboard JSON keys while copying API values and the lease ID'
    );
    await click('[data-test-copy-lease]');
    assert.true(copy.calledWith('aws/creds/example/lease'), 'Copies the lease separately');
  });

  test('role-read-denied allows selecting an STS credential type and updating TTL', async function (assert) {
    const api = this.owner.lookup('service:api') as ApiService;
    const generate = sinon.stub(api.secrets, 'awsGenerateCredentialsWithParametersRaw').resolves({
      raw: new Response(JSON.stringify({ data: { access_key: 'access', session_token: 'session-token' } })),
    });
    await render(hbs`<GenerateCredentials @backendType="aws" @backendPath="aws" @roleName="example" />`);
    assert
      .dom('[data-test-help-text]')
      .includesText('do not have permissions to read this role', 'Explains the manual type selection');
    await fillIn(GENERAL.inputByAttr('credential_type'), 'session_token');
    await fillIn(GENERAL.ttl.input('TTL'), '3');
    await click(GENERAL.submitButton);
    assert.deepEqual(
      generate.firstCall.args[0],
      {
        name: 'example',
        aws_mount_path: 'aws',
        AwsGenerateCredentialsWithParametersRequest: { ttl: '10800s' },
      },
      'Sends edited session token TTL in the credentials request body'
    );
    assert
      .dom(GENERAL.infoRowValue('Security token'))
      .includesText('***********', 'Masks modern session tokens');
  });

  test('displays generation errors without losing inputs and clears them on retry', async function (assert) {
    const api = this.owner.lookup('service:api') as ApiService;
    const generate = sinon.stub(api.secrets, 'awsGenerateCredentialsWithParametersRaw');
    generate.onFirstCall().rejects(new Error('AWS credentials could not be generated'));
    generate.onSecondCall().resolves({
      raw: new Response(JSON.stringify({ data: { access_key: 'access' } })),
    });
    await render(hbs`
      <GenerateCredentials @backendType="aws" @backendPath="aws" @roleName="example" @awsRoleType="assumed_role" />
    `);
    await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:retry');
    await click(GENERAL.submitButton);
    assert.dom('[data-test-message-error]').includesText('AWS credentials could not be generated');
    assert
      .dom(GENERAL.inputByAttr('role_arn'))
      .hasValue('arn:retry', 'Preserves inputs after a failed request');
    assert.dom(GENERAL.submitButton).isNotDisabled('Allows retry after failure');
    await click(GENERAL.submitButton);
    assert.dom('[data-test-warning]').exists('Shows results after a successful retry');
    assert.dom('[data-test-message-error]').doesNotExist('Clears the previous failure');
    assert.dom(GENERAL.infoRowLabel('Secret key')).doesNotExist('Omits absent masked credentials');
  });

  test('generation stays disabled until a pending request resolves', async function (assert) {
    const api = this.owner.lookup('service:api') as ApiService;
    let resolve!: (response: { raw: Response }) => void;
    const request = new Promise<{ raw: Response }>((done) => {
      resolve = done;
    });
    const generate = sinon.stub(api.secrets, 'awsGenerateCredentialsRaw').returns(request);
    await render(hbs`
      <GenerateCredentials @backendType="aws" @backendPath="aws" @roleName="example" @awsRoleType="iam_user" />
    `);
    const submission = click(GENERAL.submitButton);
    try {
      await waitFor(`${GENERAL.submitButton}:disabled`);
      assert.dom(GENERAL.submitButton).isDisabled('Disables the button during generation');
      assert.strictEqual(generate.callCount, 1, 'Sends one request');
      assert.dom('[data-test-secret-generate-form]').exists('Retains the form until credentials arrive');
    } finally {
      resolve({ raw: new Response(JSON.stringify({ data: { access_key: 'access' } })) });
      await submission;
    }
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('access', 'Shows the returned credentials');
  });
});
