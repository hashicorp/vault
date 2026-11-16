/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn, find, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import RoleAwsForm from 'vault/forms/aws-role';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';

module('Integration | Component | role aws edit', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.form = new RoleAwsForm({ backend: 'aws', credential_type: 'iam_user' }, { isNew: true });
    this.root = {
      text: 'aws',
      label: 'aws',
      path: 'vault.cluster.secrets.backend.list-root',
      model: 'aws',
    };
    this.router = this.owner.lookup('service:router');
    this.transition = sinon.stub(this.router, 'transitionTo');
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('rejects missing names and surfaces API write errors without navigating', async function (assert) {
    this.server.post('/aws/roles/example', () => overrideResponse(400, { errors: ['Invalid AWS role'] }));
    await render(hbs`<RoleAwsEdit @form={{this.form}} @mode="create" @root={{this.root}} />`);
    await click(GENERAL.submitButton);
    assert.false(this.transition.called, 'Invalid forms do not navigate');
    assert.dom('[data-test-message-error]').exists('Validation errors are shown');
    await fillIn(GENERAL.inputByAttr('name'), 'example');
    await click(GENERAL.submitButton);
    assert.dom('[data-test-message-error]').includesText('Invalid AWS role');
    assert.false(this.transition.called, 'Failed writes remain on the form');
  });

  test('shows the credential type and only authorized actions', async function (assert) {
    this.form = new RoleAwsForm({
      backend: 'aws',
      name: 'example',
      credential_type: 'iam_user',
    });
    this.capabilities = { canEdit: true, canDelete: false, canGenerate: false };
    await render(hbs`
      <RoleAwsEdit @form={{this.form}} @mode="show" @root={{this.root}} @capabilities={{this.capabilities}} />
    `);
    assert
      .dom(GENERAL.infoRowValue('Credential type'))
      .hasText('iam_user', 'Displays the role credential type');
    assert.dom('[data-test-backend-credentials]').doesNotExist('Generation is hidden without update access');
    assert.dom('[data-test-confirm-action-trigger]').doesNotExist('Delete is hidden without delete access');
    assert.dom().includesText('Edit role');
  });

  for (const [credentialType, arnField] of [
    ['iam_user', 'policy_arns'],
    ['assumed_role', 'role_arns'],
  ]) {
    test(`submits ${credentialType} ARN and policy widgets through the form`, async function (assert) {
      assert.expect(3);
      const policy = '{"Version":"2012-10-17","Statement":[]}';
      this.server.post('/aws/roles/example', (_, request) => {
        assert.deepEqual(
          JSON.parse(request.requestBody),
          {
            credential_type: credentialType,
            policy_arns: arnField === 'policy_arns' ? ['arn:example'] : [],
            role_arns: arnField === 'role_arns' ? ['arn:example'] : [],
            policy_document: policy,
          },
          'Rendered widgets feed the API payload'
        );
        return {};
      });
      await render(hbs`<RoleAwsEdit @form={{this.form}} @mode="create" @root={{this.root}} />`);
      await fillIn(GENERAL.inputByAttr('name'), 'example');
      await fillIn(GENERAL.inputByAttr('credential_type'), credentialType);
      await fillIn(GENERAL.stringListByIdx(0), 'arn:example');
      setCodeEditorValue(codemirror(), policy);
      await click(GENERAL.submitButton);
      assert.true(
        this.transition.calledOnceWithExactly('vault.cluster.secrets.backend.show', 'example'),
        'Navigates after saving'
      );
      assert.dom('[data-test-message-error]').doesNotExist('Successful writes have no error');
    });
  }

  test('details deletion reports failures and supports a successful retry', async function (assert) {
    this.form = new RoleAwsForm({ backend: 'aws', name: 'example', credential_type: 'iam_user' });
    this.capabilities = { canDelete: true };
    const remove = sinon.stub(this.owner.lookup('service:api').secrets, 'awsDeleteRole');
    remove.onFirstCall().rejects(new Error('Deletion denied'));
    remove.onSecondCall().resolves();
    await render(
      hbs`<RoleAwsEdit @form={{this.form}} @mode="show" @root={{this.root}} @capabilities={{this.capabilities}} />`
    );
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.dom('[data-test-message-error]').includesText('Deletion denied', 'Shows the delete failure');
    assert.false(this.transition.called, 'Remains on details after failure');
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.true(remove.calledWithExactly('example', 'aws'), 'Uses the selected mount and role');
    assert.true(
      this.transition.calledOnceWithExactly('vault.cluster.secrets.backend.list-root'),
      'Returns to the list after successful deletion'
    );
  });

  test('pending saves disable submission and ignore duplicate submit events', async function (assert) {
    this.form.data.name = 'example';
    let resolve;
    const request = new Promise((done) => {
      resolve = done;
    });
    const write = sinon.stub(this.owner.lookup('service:api').secrets, 'awsWriteRole').returns(request);
    await render(hbs`<RoleAwsEdit @form={{this.form}} @mode="create" @root={{this.root}} />`);
    const submission = click(GENERAL.submitButton);
    try {
      await waitFor(`${GENERAL.submitButton}:disabled`);
      assert.dom(GENERAL.submitButton).isDisabled('Cannot submit while saving');
      find('form').dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
      assert.strictEqual(write.callCount, 1, 'A second submit does not duplicate the API write');
      assert.false(this.transition.called, 'Does not navigate before the response');
    } finally {
      resolve();
      await submission;
    }
    assert.dom(GENERAL.submitButton).isNotDisabled('Restores submission after saving');
    assert.true(this.transition.calledOnce, 'Navigates once after success');
  });
});
