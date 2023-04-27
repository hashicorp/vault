/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';

const SELECTORS = {
  nameInput: '[data-test-policy-input="name"]',
  uploadFileToggle: '[data-test-policy-edit-toggle]',
  policyEditor: '[data-test-policy-editor]',
  policyUpload: '[data-test-text-file-input]',
  saveButton: '[data-test-policy-save]',
  cancelButton: '[data-test-policy-cancel]',
  error: '[data-test-error]',
};

module('Integration | Component | policy-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('policy/acl');
    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();
    this.server = new Pretender(function () {
      this.put('/v1/sys/policies/acl/bad-policy', () => {
        return [
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({ errors: ['An error occurred'] }),
        ];
      });
      this.put('/v1/sys/policies/acl/**', () => {
        return [204, { 'Content-Type': 'application/json' }];
      });
      this.put('/v1/sys/policies/rgp/**', () => {
        return [204, { 'Content-Type': 'application/json' }];
      });
    });
  });
  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it renders the form for new ACL policy', async function (assert) {
    const policy = `
    path "secret/*" {
      capabilities = [ "create", "read", "update", "list" ]
    }
    `;
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.nameInput).exists({ count: 1 }, 'Name input exists');
    assert.dom(SELECTORS.nameInput).hasNoText('Name field is not filled');
    assert.dom(SELECTORS.uploadFileToggle).exists({ count: 1 }, 'Upload file toggle exists');
    await fillIn(SELECTORS.nameInput, 'Foo');
    assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
    await fillIn(`${SELECTORS.policyEditor} textarea`, policy);
    assert.strictEqual(this.model.policy, policy, 'Policy editor sets policy on model');
    assert.ok(this.onSave.notCalled);
    assert.dom(SELECTORS.saveButton).hasText('Create policy');
    await click(SELECTORS.saveButton);
    assert.ok(this.onSave.calledOnceWith(this.model));
  });

  test('it renders the form for new RGP policy', async function (assert) {
    const model = this.store.createRecord('policy/rgp');
    const policy = `
    path "secret/*" {
      capabilities = [ "create", "read", "update", "list" ]
    }
    `;
    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.nameInput).exists({ count: 1 }, 'Name input exists');
    assert.dom(SELECTORS.nameInput).hasNoText('Name field is not filled');
    assert.dom(SELECTORS.uploadFileToggle).exists({ count: 1 }, 'Upload file toggle exists');
    await fillIn(SELECTORS.nameInput, 'Foo');
    assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
    await fillIn(`${SELECTORS.policyEditor} textarea`, policy);
    assert.strictEqual(this.model.policy, policy, 'Policy editor sets policy on model');
    assert.ok(this.onSave.notCalled);
    assert.dom(SELECTORS.saveButton).hasText('Create policy');
    await click(SELECTORS.saveButton);
    assert.ok(this.onSave.calledOnceWith(this.model));
  });

  test('it toggles to upload a new policy and uploads file', async function (assert) {
    const policy = `
    path "auth/token/lookup-self" {
      capabilities = ["read"]
    }`;
    this.file = new File([policy], 'test-policy.hcl');
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.uploadFileToggle).exists({ count: 1 }, 'Upload file toggle exists');
    assert.dom(SELECTORS.policyEditor).exists({ count: 1 }, 'Policy editor is shown');
    assert.dom(SELECTORS.policyUpload).doesNotExist('Policy upload is not shown');
    await click(SELECTORS.uploadFileToggle);
    assert.dom(SELECTORS.policyUpload).exists({ count: 1 }, 'Policy upload is shown after toggle');
    assert.dom(SELECTORS.policyEditor).doesNotExist('Policy editor is not shown');
    await triggerEvent(SELECTORS.policyUpload, 'change', { files: [this.file] });
    assert.dom(SELECTORS.nameInput).hasValue('test-policy', 'it fills in policy name');
    await click(SELECTORS.saveButton);
    assert.propEqual(this.onSave.lastCall.args[0].policy, policy, 'policy content saves in correct format');
  });

  test('it renders the form to edit existing ACL policy', async function (assert) {
    const model = this.store.createRecord('policy/acl', {
      name: 'bar',
      policy: 'some policy content',
    });
    model.save();

    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.nameInput).doesNotExist('Name input is not rendered');
    assert.dom(SELECTORS.uploadFileToggle).doesNotExist('Upload file toggle does not exist');

    await fillIn(`${SELECTORS.policyEditor} textarea`, 'updated-');
    assert.strictEqual(
      this.model.policy,
      'updated-some policy content',
      'Policy editor updates policy value on model'
    );
    assert.ok(this.onSave.notCalled);
    assert.dom(SELECTORS.saveButton).hasText('Save', 'Save button text is correct');
    await click(SELECTORS.saveButton);
    assert.ok(this.onSave.calledOnceWith(this.model));
  });
  test('it renders the form to edit existing RGP policy', async function (assert) {
    const model = this.store.createRecord('policy/rgp', {
      name: 'bar',
      policy: 'some policy content',
    });
    model.save();

    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.nameInput).doesNotExist('Name input is not rendered');
    assert.dom(SELECTORS.uploadFileToggle).doesNotExist('Upload file toggle does not exist');

    await fillIn(`${SELECTORS.policyEditor} textarea`, 'updated-');
    assert.strictEqual(
      this.model.policy,
      'updated-some policy content',
      'Policy editor updates policy value on model'
    );
    assert.ok(this.onSave.notCalled);
    assert.dom(SELECTORS.saveButton).hasText('Save', 'Save button text is correct');
    await click(SELECTORS.saveButton);
    assert.ok(this.onSave.calledOnceWith(this.model));
  });
  test('it shows the error message on form when save fails', async function (assert) {
    const model = this.store.createRecord('policy/acl', {
      name: 'bad-policy',
      policy: 'some policy content',
    });

    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    await click(SELECTORS.saveButton);
    assert.ok(this.onSave.notCalled);
    assert.dom(SELECTORS.error).includesText('An error occurred');
  });
});
