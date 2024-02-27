/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const SELECTORS = {
  nameInput: '[data-test-policy-input="name"]',
  uploadFileToggle: '[data-test-policy-edit-toggle]',
  policyEditor: '[data-test-policy-editor]',
  policyUpload: '[data-test-text-file-input]',
  saveButton: '[data-test-policy-save]',
  cancelButton: '[data-test-policy-cancel]',
  error: '[data-test-message-error]',
  // For example modal:
  exampleButton: '[data-test-policy-example-button]',
  exampleModal: '[data-test-policy-example-modal]',
  exampleModalTitle: '[data-test-modal-title]',
  exampleModalClose: '[data-test-modal-close-button]',
  // For additional fields for EGP policy:
  fields: (name) => `[data-test-field=${name}]`,
  pathsInput: (index) => `[data-test-string-list-input="${index}"]`,
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
      this.put('/v1/sys/policies/egp/**', () => {
        return [204, { 'Content-Type': 'application/json' }];
      });
    });
    setRunOptions({
      rules: {
        // TODO: fix JSONEditor/CodeMirror
        label: { enabled: false },
        'label-title-only': { enabled: false },
      },
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

  test('it renders the form for new EGP policy', async function (assert) {
    const model = this.store.createRecord('policy/egp');
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
    assert.dom(SELECTORS.fields('paths')).exists('Paths field exists');
    assert.dom(SELECTORS.pathsInput('0')).exists('0 field exists');
    await fillIn(SELECTORS.pathsInput('0'), 'my path');
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

  test('it renders the form to edit existing EGP policy', async function (assert) {
    const model = this.store.createRecord('policy/egp', {
      name: 'bar',
      policy: 'some policy content',
      paths: ['first path'],
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
    await fillIn(SELECTORS.pathsInput('1'), 'second path');
    assert.strictEqual(
      JSON.stringify(this.model.paths),
      '["first path","second path"]',
      'Second path field is updated on model'
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

  test('it does not create a new policy when the cancel button is clicked', async function (assert) {
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
    await fillIn(SELECTORS.nameInput, 'Foo');
    assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
    await fillIn(`${SELECTORS.policyEditor} textarea`, policy);
    assert.strictEqual(this.model.policy, policy, 'Policy editor sets policy on model');

    await click(SELECTORS.cancelButton);
    assert.ok(this.onSave.notCalled);
    assert.ok(this.onCancel.calledOnce, 'Form calls onCancel');
  });

  test('it does not save edits when the cancel button is clicked', async function (assert) {
    const model = this.store.createRecord('policy/acl', {
      name: 'foo',
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
    await fillIn(`${SELECTORS.policyEditor} textarea`, 'updated-');
    assert.strictEqual(
      this.model.policy,
      'updated-some policy content',
      'Policy editor updates policy value on model'
    );
    await click(SELECTORS.cancelButton);
    assert.ok(this.onSave.notCalled);
    assert.ok(this.onCancel.calledOnce, 'Form calls onCancel');

    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.strictEqual(
      this.model.policy,
      'some policy content',
      'Policy editor shows original policy content, meaning that onCancel worked successfully'
    );
  });

  test('it does not render the button and modal for the policy example if not specified to', async function (assert) {
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    `);
    assert.dom(SELECTORS.exampleModal).doesNotExist('Modal for the policy example does not exist');
    assert.dom(SELECTORS.exampleButton).doesNotExist('Button for the policy example modal does not exist');
  });

  test('it renders the button and modal for the policy example when specified to', async function (assert) {
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
      @renderPolicyExampleModal={{true}}
    />
        `);
    assert.dom(SELECTORS.exampleButton).exists({ count: 1 }, 'Modal for the policy example exists');
    assert.dom(SELECTORS.exampleButton).exists({ count: 1 }, 'Button for the policy example modal exists');
  });

  test('it renders the correct title for ACL example for the policy example modal', async function (assert) {
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
      @renderPolicyExampleModal={{true}}
    />
        `);
    await click(SELECTORS.exampleButton);
    assert.dom(SELECTORS.exampleModalTitle).hasText('Example ACL Policy');
  });

  test('it renders the correct title for RGP example for the policy example modal', async function (assert) {
    const model = this.store.createRecord('policy/rgp');
    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
      @renderPolicyExampleModal={{true}}
    />
        `);
    await click(SELECTORS.exampleButton);
    assert.dom(SELECTORS.exampleModalTitle).hasText('Example RGP Policy');
  });

  test('it renders the correct title for EGP example for the policy example modal', async function (assert) {
    const model = this.store.createRecord('policy/egp');
    this.set('model', model);
    await render(hbs`
    <PolicyForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
      @renderPolicyExampleModal={{true}}
    />
        `);
    await click(SELECTORS.exampleButton);
    assert.dom(SELECTORS.exampleModalTitle).hasText('Example EGP Policy');
  });
});
