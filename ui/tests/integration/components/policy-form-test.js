/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, settled, triggerEvent, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { FORM } from 'vault/tests/helpers/form-selectors';

async function setEditorValue(value) {
  await waitFor('.cm-editor');
  const editor = codemirror(GENERAL.codemirror);
  setCodeEditorValue(editor, value);
  return settled();
}

module('Integration | Component | policy-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    // Set model here with "ACL" policy type for so PolicyForm component consistently has a @model arg
    this.model = this.store.createRecord('policy/acl');
    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();
    this.isCompact = undefined;
    this.server.put('/sys/policies/acl/:name', (_, req) => {
      if (req.params.name === 'bad-policy') {
        return overrideResponse(400, { errors: ['An error occurred'] });
      }
      return overrideResponse(204);
    });
    this.server.put('/sys/policies/rgp/:name', () => overrideResponse(204));
    this.server.put('/sys/policies/egp/:name', () => overrideResponse(204));

    this.assertNoVisualEditor = (assert, msg = 'it does not render visual policy builder') => {
      assert.dom(GENERAL.radioByAttr()).doesNotExist('it does not render radio options');
      assert.dom(GENERAL.codemirror).exists('JSON editor renders');
      assert.dom(GENERAL.fieldByAttr('visual editor')).doesNotExist(msg);
    };

    this.renderComponent = () => {
      return render(
        hbs`<PolicyForm
          @model={{this.model}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
          @isCompact={{this.isCompact}}
        />`
      );
    };
  });

  // Tests that are policy type agnostic are below, otherwise tests are organized by module for each type: ACL, RGP, EGP
  test('it toggles to upload a new policy and uploads file', async function (assert) {
    const policy = `path "auth/token/lookup-self" { capabilities = ["read"] }`;
    const file = new File([policy], 'test-policy.hcl');
    await this.renderComponent();
    assert.dom(GENERAL.toggleInput('Upload file')).exists({ count: 1 }, 'Upload file toggle exists');
    assert.dom(GENERAL.fieldByAttr('visual editor')).exists('Visual editor renders');
    assert.dom(GENERAL.fileInput).doesNotExist('Policy upload is not shown');
    // Click upload file toggle
    await click(GENERAL.toggleInput('Upload file'));
    assert.dom(GENERAL.fileInput).exists({ count: 1 }, 'Policy upload is shown after toggle');
    assert.dom(FORM.header('Policy rules')).doesNotExist();
    assert.dom(FORM.description('Policy rules')).doesNotExist();
    assert
      .dom(GENERAL.fieldByAttr('visual editor'))
      .doesNotExist('Visual editor is hidden when "Upload file" is selected');
    assert
      .dom(GENERAL.radioByAttr())
      .doesNotExist('it does not render radio buttons "Upload file" is selected');
    // Upload file
    await triggerEvent(GENERAL.fileInput, 'change', { files: [file] });
    await waitFor('.cm-editor');
    assert.dom(GENERAL.toggleInput('Upload file')).isNotChecked('Upload file is unchecked after upload');
    assert.dom(GENERAL.codemirror).exists().hasTextContaining(policy, 'code editor renders policy from file');
    assert.dom(GENERAL.inputByAttr('name')).hasValue('test-policy', 'it fills in policy name');
    assert.dom(GENERAL.radioByAttr('visual')).exists().isNotChecked();
    assert.dom(GENERAL.radioByAttr('code')).exists().isChecked();
    await click(GENERAL.submitButton);
    assert.propEqual(this.onSave.lastCall.args[0].policy, policy, 'policy content saves in correct format');
  });

  test('it renders all elements by default (when not compact)', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('visual editor')).exists();
    assert
      .dom(GENERAL.accordionButton('Automation snippets'))
      .exists('Automation snippets render')
      .hasAttribute('aria-expanded', 'false');
    assert.dom(GENERAL.radioByAttr()).exists({ count: 2 }, 'radio buttons render for each editor');
    assert.dom(GENERAL.radioByAttr('visual')).exists();
    assert.dom(GENERAL.radioByAttr('code')).exists();
    await click(GENERAL.radioByAttr('code'));
    assert.dom(GENERAL.button('How to write a policy')).exists();
  });

  test('it hides elements when isCompact', async function (assert) {
    this.isCompact = true;
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('visual editor')).doesNotExist('Visual editor does not render');
    assert.dom(GENERAL.accordionButton('Automation snippets')).doesNotExist();
    assert.dom(GENERAL.radioByAttr('visual')).doesNotExist();
    assert.dom(GENERAL.radioByAttr('code')).doesNotExist();
    assert.dom(GENERAL.button('How to write a policy')).doesNotExist();
  });

  test('it does not create a new policy when the cancel button is clicked', async function (assert) {
    const policy = `
    path "secret/*" {
      capabilities = [ "create", "read", "update", "list" ]
    }
    `;
    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('name'), 'Foo');
    assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
    await click(GENERAL.radioByAttr('code'));
    await setEditorValue(policy);
    assert.strictEqual(this.model.policy, policy, 'Policy editor sets policy on model');

    await click(GENERAL.cancelButton);
    assert.true(this.onSave.notCalled, 'onSave is not called yet');
    assert.true(this.onCancel.calledOnce, 'Form calls onCancel');
  });

  test('it does not save edits when the cancel button is clicked', async function (assert) {
    this.model.name = 'foo';
    this.model.policy = 'some policy content';
    this.model.save();
    await this.renderComponent();
    await setEditorValue('updated');
    assert.strictEqual(this.model.policy, 'updated', 'Policy editor updates policy value on model');
    await click(GENERAL.cancelButton);
    assert.true(this.onSave.notCalled, 'onSave is not called yet');
    assert.true(this.onCancel.calledOnce, 'Form calls onCancel');

    await this.renderComponent();
    assert.strictEqual(
      this.model.policy,
      'some policy content',
      'Policy editor shows original policy content, meaning that onCancel worked successfully'
    );
  });

  test('it shows the error message on form when save fails', async function (assert) {
    this.model.name = 'bad-policy';
    this.model.policy = 'some policy content';
    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert.true(this.onSave.notCalled, 'onSave is not called yet');
    assert.dom(GENERAL.messageError).includesText('An error occurred');
  });
  // End shared functionality tests

  module('ACL', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.store.createRecord('policy/acl');
      this.policy = `path "secret/*" {
      capabilities = [ "create", "read", "update", "list" ]
    }
    `;
    });

    test('it renders the form for new ACL policy', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('name')).exists({ count: 1 }, 'Name input exists');
      assert.dom(GENERAL.inputByAttr('name')).hasNoText('Name field is not filled');
      // Assert visual policy editor default state
      assert.dom(GENERAL.radioByAttr('visual')).exists().isChecked();
      assert.dom(GENERAL.radioByAttr('code')).exists().isNotChecked();
      assert.dom(GENERAL.codemirror).doesNotExist('JSON editor does not render by default');
      assert.dom(GENERAL.fieldByAttr('visual editor')).exists('it renders visual policy editor by default');
      assert.dom(GENERAL.toggleInput('Upload file')).exists({ count: 1 }, 'Upload file toggle exists');
      assert.dom(GENERAL.submitButton).hasText('Create policy');
    });

    test('it saves a new ACL policy using the code editor', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'Foo');
      assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
      await click(GENERAL.radioByAttr('code'));
      await setEditorValue(this.policy);
      assert.strictEqual(this.model.policy, this.policy, 'Policy editor sets policy on model');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
      const [actual] = this.onSave.lastCall.args;
      assert.strictEqual(actual.policy, this.policy, 'onSave is called with expected policy');
    });

    test('it renders the form to edit existing ACL policy', async function (assert) {
      this.model.name = 'bar';
      this.model.policy = this.policy;
      this.model.save();
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('name')).doesNotExist('Name input is not rendered');
      assert.dom(GENERAL.toggleInput('Upload file')).doesNotExist('Upload file toggle does not exist');
      this.assertNoVisualEditor(assert, 'it does not render visual editor when editing an ACL policy');

      await setEditorValue('updated');
      assert.strictEqual(this.model.policy, 'updated', 'Policy editor updates policy value on model');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      assert.dom(GENERAL.submitButton).hasText('Save', 'Save button text is correct');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
    });

    test('it renders the correct title for ACL example for the policy example modal', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.radioByAttr('code'));
      await click(GENERAL.button('How to write a policy'));
      assert.dom(GENERAL.modal.container('Example policy')).exists('Modal renders');
      assert.dom(GENERAL.modal.header('Example policy')).hasText('Example ACL Policy');
    });

    // Only ACL policy types support the visual editor
    test('it toggles between visual and code editors', async function (assert) {
      await this.renderComponent();
      // Assert default state
      assert.dom(GENERAL.radioByAttr('visual')).exists().isChecked();
      assert.dom(GENERAL.radioByAttr('code')).exists().isNotChecked();
      assert.dom(GENERAL.codemirror).doesNotExist('JSON editor does not render by default');
      assert
        .dom(GENERAL.fieldByAttr('visual editor'))
        .hasTextContaining('Rule Show preview')
        .exists('it renders visual policy editor by default');
      // Select Code editor
      await click(GENERAL.radioByAttr('code'));
      assert.dom(GENERAL.radioByAttr('visual')).exists().isNotChecked();
      assert.dom(GENERAL.radioByAttr('code')).exists().isChecked();
      assert.dom(GENERAL.codemirror).exists('code editor renders after selecting "Code editor"');
      assert.dom(GENERAL.fieldByAttr('visual editor')).doesNotExist('visual editor no longer renders');
      // Go back to Visual editor
      await click(GENERAL.radioByAttr('visual'));
      assert.dom(GENERAL.radioByAttr('visual')).exists().isChecked();
      assert.dom(GENERAL.radioByAttr('code')).exists().isNotChecked();
      assert.dom(GENERAL.codemirror).doesNotExist();
      assert
        .dom(GENERAL.fieldByAttr('visual editor'))
        .hasTextContaining('Rule Show preview')
        .exists('Visual editor renders after selecting radio');
    });

    test('it saves a new ACL policy using the visual editor', async function (assert) {
      const expectedPolicy = `path "first/path" {
    capabilities = ["read"]
}
path "second/path" {
    capabilities = ["update"]
}`;
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'Foo');
      assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
      // Set up first rule
      await fillIn(GENERAL.inputByAttr('path'), 'first/path');
      await click(GENERAL.checkboxByAttr('read'));
      // Add second rule
      await click(GENERAL.button('Add rule'));
      await fillIn(`${GENERAL.cardContainer('1')} ${GENERAL.inputByAttr('path')}`, 'second/path');
      await click(`${GENERAL.cardContainer('1')} ${GENERAL.checkboxByAttr('update')}`);
      // Save policy
      assert.strictEqual(this.model.policy, expectedPolicy, 'Policy editor sets policy on model');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
      const [actual] = this.onSave.lastCall.args;
      assert.strictEqual(actual.policy, expectedPolicy, 'save is called with expected policy');
    });

    // Automation snippets are only supported for "ACL" policy types at this time
    test('it renders empty snippets by default', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.accordionButton('Automation snippets')).exists();
      await click(GENERAL.accordionButton('Automation snippets'));
      assert.dom(GENERAL.hdsTab('terraform')).exists().hasAttribute('aria-selected', 'true');
      assert.dom(GENERAL.hdsTab('cli')).exists().hasAttribute('aria-selected', 'false');
      const expectedTfvp = `resource "vault_policy" "<local identifier>" {
      name = "<policy name>"
    
      policy = <<EOT
      path "" {
        capabilities = []
    }
    EOT
    }`;
      assert
        .dom(GENERAL.fieldByAttr('terraform'))
        .hasText(expectedTfvp, 'it renders empty terraform snippet');
      const expectedCli = `vault policy write <policy name> - <<EOT
      path "" {
        capabilities = []
    }
    EOT`;
      assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli, 'it renders empty cli snippet');
    });

    test('it updates snippets as policy changes', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputByAttr('name'), 'my-secure-policy');
      await fillIn(GENERAL.inputByAttr('path'), 'my/super/secret/*');
      await click(GENERAL.checkboxByAttr('patch'));
      await click(GENERAL.accordionButton('Automation snippets'));
      const expectedTfvp = `resource "vault_policy" "<local identifier>" {
  name = "my-secure-policy"

  policy = <<EOT
  path "my/super/secret/*" {
    capabilities = ["patch"]
}
EOT
}`;
      assert.dom(GENERAL.fieldByAttr('terraform')).hasText(expectedTfvp);
      const expectedCli = `vault policy write my-secure-policy - <<EOT
  path "my/super/secret/*" {
    capabilities = ["patch"]
}
EOT`;
      assert.dom(GENERAL.fieldByAttr('cli')).hasText(expectedCli);
    });

    module('switch editors modal', function (hooks) {
      hooks.beforeEach(function () {
        this.originalPolicy = `path "secret/data/*" {
    capabilities = ["read"]
}`;
        this.editedPolicy = `path "secret/data/*" {
    capabilities = ["read", "list"]
}`;
        this.setupSwitch = async () => {
          await this.renderComponent();
          // Use visual editor
          await fillIn(GENERAL.inputByAttr('path'), 'secret/data/*');
          await click(GENERAL.checkboxByAttr('read'));
          // Switch to code editor and set with edited policy
          await click(GENERAL.radioByAttr('code'));
          await setEditorValue(this.editedPolicy);
        };
      });

      test('it renders modal when switching back to visual editor after using code editor', async function (assert) {
        await this.setupSwitch();
        // Switch back to visual editor
        await click(GENERAL.radioByAttr('visual'));
        assert
          .dom(GENERAL.modal.container('warning'))
          .exists()
          .hasText(
            'Switching editors? Changes made in the Code editor will not be carried over to the Visual editor. Do you want to switch and discard changes? Switch and discard changes Cancel'
          );
        assert.dom(GENERAL.icon('alert-triangle')).exists();
        assert.dom(GENERAL.confirmButton).exists();
        assert.dom(GENERAL.cancelButton).exists();
      });

      test('it renders modal when only formatting has changed', async function (assert) {
        // Same as original policy but remove indents before "capabilities"
        this.editedPolicy = `path "secret/data/*" {
capabilities = ["read"]
}`;
        await this.setupSwitch();
        // Switch back to visual editor
        await click(GENERAL.radioByAttr('visual'));
        assert.dom(GENERAL.modal.container('warning')).exists();
      });

      test('it does NOT render modal when code editor is updated with original policy', async function (assert) {
        // Clear editor
        this.editedPolicy = '';
        await this.setupSwitch();
        // Set with original policy
        await setEditorValue(this.originalPolicy);
        // Switch back to visual editor
        await click(GENERAL.radioByAttr('visual'));
        assert.dom(GENERAL.modal.container('warning')).doesNotExist('modal does not render');
        assert.dom(GENERAL.radioByAttr('visual')).isChecked('visual editor is checked');
        assert.dom(GENERAL.radioByAttr('code')).isNotChecked();
      });

      test('it does not switch editors after clicking "Cancel"', async function (assert) {
        await this.setupSwitch();
        // Switch back to visual editor
        await click(GENERAL.radioByAttr('visual'));
        await click(GENERAL.cancelButton);
        assert.dom(GENERAL.modal.container('warning')).doesNotExist('Clicking cancel closes the modal');
        assert
          .dom(GENERAL.radioByAttr('visual'))
          .isNotChecked('visual editor is not checked after clicking "Cancel"');
        assert.dom(GENERAL.fieldByAttr('visual editor')).doesNotExist('visual editor does not render');
        assert.dom(GENERAL.radioByAttr('code')).isChecked('code editor is still checked');
        assert
          .dom(GENERAL.codemirror)
          .exists()
          .hasTextContaining('list', 'code editor renders edited policy');
      });

      test('it switches editors and reverts changes after confirming', async function (assert) {
        await this.setupSwitch();
        // Switch back to visual editor
        await click(GENERAL.radioByAttr('visual'));
        await click(GENERAL.confirmButton);
        assert.dom(GENERAL.modal.container('warning')).doesNotExist('confirming switch closes the modal');
        assert.dom(GENERAL.radioByAttr('visual')).isChecked('visual editor is checked after "Cancel"');
        assert.dom(GENERAL.fieldByAttr('visual editor')).exists('visual editor renders');
        assert.dom(GENERAL.inputByAttr('path')).hasValue('secret/data/*');
        assert.dom(GENERAL.checkboxByAttr('read')).isChecked('"read" is still checked');
        assert.dom(GENERAL.radioByAttr('code')).isNotChecked('code editor is no longer checked');
        // Confirm code editor reverts to original policy and not edited one
        await click(GENERAL.radioByAttr('code'));
        await waitFor('.cm-editor');
        assert.dom(GENERAL.codemirror).doesNotHaveTextContaining('list');
      });
    });
  });

  module('RGP', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.store.createRecord('policy/rgp');
      this.policy = `import "strings"
precond = rule {
    strings.has_prefix(request.path, "sys/policies/admin")
}
main = rule when precond {
    identity.entity.metadata.role is "Team Lead" or
      identity.entity.name is "James Thomas"
}`;
    });

    test('it renders the form for new RGP policy', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('name')).exists({ count: 1 }, 'Name input exists');
      assert.dom(GENERAL.inputByAttr('name')).hasNoText('Name field is not filled');
      assert.dom(GENERAL.toggleInput('Upload file')).exists({ count: 1 }, 'Upload file toggle exists');
      this.assertNoVisualEditor(assert, 'it hides visual editor for RGP policy types');

      await fillIn(GENERAL.inputByAttr('name'), 'Foo');
      assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
      await setEditorValue(this.policy);
      assert.strictEqual(this.model.policy, this.policy, 'Policy editor sets policy on model');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      assert.dom(GENERAL.submitButton).hasText('Create policy');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
    });

    test('it renders the form to edit existing RGP policy', async function (assert) {
      this.model.name = 'bar';
      this.model.policy = this.policy;
      this.model.save();
      await this.renderComponent();
      assert.dom(GENERAL.inputByAttr('name')).doesNotExist('Name input is not rendered');
      assert.dom(GENERAL.toggleInput('Upload file')).doesNotExist('Upload file toggle does not exist');

      await setEditorValue('updated');
      assert.strictEqual(this.model.policy, 'updated', 'Policy editor updates policy value on model');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      assert.dom(GENERAL.submitButton).hasText('Save', 'Save button text is correct');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
    });

    test('it renders the correct title for RGP example for the policy example modal', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.button('How to write a policy'));
      assert.dom(GENERAL.modal.header('Example policy')).hasText('Example RGP Policy');
    });

    test('it does not render automation snippets', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.accordionButton('Automation snippets')).doesNotExist();
    });
  });

  module('EGP', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.store.createRecord('policy/egp');
      this.policy = `import "time"
workdays = rule {
    time.now.weekday > 0 and time.now.weekday < 6
}
workhours = rule {
    time.now.hour > 7 and time.now.hour < 18
}
main = rule {
    workdays and workhours
}
`;
    });

    test('it renders the form for new EGP policy', async function (assert) {
      await this.renderComponent();

      assert.dom(GENERAL.inputByAttr('name')).exists({ count: 1 }, 'Name input exists');
      assert.dom(GENERAL.inputByAttr('name')).hasNoText('Name field is not filled');
      assert.dom(GENERAL.toggleInput('Upload file')).exists({ count: 1 }, 'Upload file toggle exists');
      this.assertNoVisualEditor(assert, 'it hides visual editor for EGP policy types');

      await fillIn(GENERAL.inputByAttr('name'), 'Foo');
      assert.strictEqual(this.model.name, 'foo', 'Input sets name on model to lowercase input');
      await setEditorValue(this.policy);
      assert.strictEqual(this.model.policy, this.policy, 'Policy editor sets policy on model');
      assert.dom(GENERAL.fieldByAttr('paths')).exists('Paths field exists');
      assert.dom(GENERAL.stringListByIdx(0)).exists('0 field exists');
      await fillIn(GENERAL.stringListByIdx(0), 'my path');
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      assert.dom(GENERAL.submitButton).hasText('Create policy');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
    });

    test('it renders the form to edit existing EGP policy', async function (assert) {
      this.model.name = 'bar';
      this.model.policy = this.policy;
      this.model.paths = ['first path'];
      this.model.save();
      await this.renderComponent();

      assert.dom(GENERAL.inputByAttr('name')).doesNotExist('Name input is not rendered');
      assert.dom(GENERAL.toggleInput('Upload file')).doesNotExist('Upload file toggle does not exist');
      await setEditorValue('updated');
      assert.strictEqual(this.model.policy, 'updated', 'Policy editor updates policy value on model');
      await fillIn(GENERAL.stringListByIdx(1), 'second path');
      assert.strictEqual(
        JSON.stringify(this.model.paths),
        '["first path","second path"]',
        'Second path field is updated on model'
      );
      assert.true(this.onSave.notCalled, 'onSave is not called yet');
      assert.dom(GENERAL.submitButton).hasText('Save', 'Save button text is correct');
      await click(GENERAL.submitButton);
      assert.true(this.onSave.calledOnceWith(this.model), 'onSave is called with model');
    });

    test('it renders the correct title for EGP example for the policy example modal', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.button('How to write a policy'));
      assert.dom(GENERAL.modal.header('Example policy')).hasText('Example EGP Policy');
    });

    test('it does not render automation snippets', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.accordionButton('Automation snippets')).doesNotExist();
    });
  });
});
