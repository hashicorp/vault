/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { module, test } from 'qunit';
import sinon from 'sinon';
import V2Form from 'vault/forms/v2/v2-form';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | form/v2/apply', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    // Create a minimal FormConfig for testing
    const formConfig = {
      name: 'test-resource',
      path: '/v1/test/resource',
      title: 'Test Resource',
      payload: {
        name: 'test-name',
        description: 'test-description',
      },
      submit: sinon.stub().resolves({ id: 'test-123' }),
      sections: [
        {
          name: 'basic',
          fields: [
            {
              name: 'name',
              label: 'Name',
              type: 'TextInput',
            },
          ],
        },
      ],
    };

    this.form = new V2Form(formConfig);
    this.onBack = sinon.spy();
    this.onApply = sinon.spy();
    this.onDone = sinon.spy();
  });

  test('it renders the apply component', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    assert.dom(GENERAL.textDisplay('Step title')).hasText('Choose your implementation method');
    assert.dom('.hds-form-radio-card').exists({ count: 3 }, 'renders three creation method options');
  });

  test('it renders Terraform as the default selected option', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const radios = this.element.querySelectorAll('input[type="radio"]');
    assert.true(radios[0].checked, 'Terraform option (first radio) is checked by default');
  });

  test('it shows Terraform code snippet by default', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    assert.dom(GENERAL.fieldByAttr('terraform')).exists('Terraform code block is rendered');
    assert.dom('.hds-code-block').exists('Code block component is rendered');
  });

  test('it allows changing to API/CLI creation method', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[1]); // API/CLI is second option

    assert.true(radios[1].checked, 'API/CLI option is now checked');
    // API/CLI shows code snippets too, just not the Terraform-specific download button
    const buttons = this.element.querySelectorAll('button');
    const downloadButton = Array.from(buttons).find((btn) => btn.textContent.includes('Export as tf file'));
    assert.notOk(downloadButton, 'Terraform download button is hidden for API/CLI');
  });

  test('it allows changing to UI workflow method', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[2]); // UI workflow is third option

    assert.true(radios[2].checked, 'UI workflow option is now checked');
    assert.dom('.hds-code-block').doesNotExist('Code snippets are hidden for UI workflow');
  });

  test('it shows Apply changes button only for UI workflow', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    // Initially on Terraform - no Apply button
    let buttons = this.element.querySelectorAll('button');
    let applyButton = Array.from(buttons).find((btn) => btn.textContent.includes('Apply changes'));
    assert.notOk(applyButton, 'Apply changes button not shown for Terraform');

    // Switch to UI workflow
    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[2]); // UI workflow is third option

    buttons = this.element.querySelectorAll('button');
    applyButton = Array.from(buttons).find((btn) => btn.textContent.includes('Apply changes'));
    assert.ok(applyButton, 'Apply changes button shown for UI workflow');
  });

  test('it always shows Back and Done buttons', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const buttons = this.element.querySelectorAll('button');
    const backButton = Array.from(buttons).find((btn) => btn.textContent.includes('Back'));
    const doneButton = Array.from(buttons).find((btn) => btn.textContent.includes('Done & exit'));

    assert.ok(backButton, 'Back button is rendered');
    assert.ok(doneButton, 'Done & exit button is rendered');
  });

  test('it calls onBack when Back button is clicked', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const buttons = this.element.querySelectorAll('button');
    const backButton = Array.from(buttons).find((btn) => btn.textContent.includes('Back'));
    await click(backButton);

    assert.ok(this.onBack.calledOnce, 'onBack callback was called');
  });

  test('it calls onDone when Done button is clicked', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const buttons = this.element.querySelectorAll('button');
    const doneButton = Array.from(buttons).find((btn) => btn.textContent.includes('Done & exit'));
    await click(doneButton);

    assert.ok(this.onDone.calledOnce, 'onDone callback was called');
  });

  test('it calls onApply when Apply changes button is clicked', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    // Switch to UI workflow to show Apply button
    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[2]); // UI workflow is third option

    const buttons = this.element.querySelectorAll('button');
    const applyButton = Array.from(buttons).find((btn) => btn.textContent.includes('Apply changes'));
    await click(applyButton);

    assert.ok(this.onApply.calledOnce, 'onApply callback was called');
  });

  test('it shows download button for Terraform snippet', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const buttons = this.element.querySelectorAll('button');
    const downloadButton = Array.from(buttons).find((btn) => btn.textContent.includes('Export as tf file'));
    assert.ok(downloadButton, 'Download button is rendered for Terraform');
  });

  test('it hides download button for API/CLI method', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[1]); // API/CLI is second option

    const buttons = this.element.querySelectorAll('button');
    const downloadButton = Array.from(buttons).find((btn) => btn.textContent.includes('Export as tf file'));
    assert.notOk(downloadButton, 'Download button is not shown for API/CLI');
  });

  test('it shows appropriate descriptions for each creation method', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    // Check for description text in the rendered cards
    assert.dom('.hds-form-radio-card').exists({ count: 3 }, 'Three radio cards rendered');
    assert.dom(this.element).includesText('Infrastructure as Code', 'Terraform description present');
    assert.dom(this.element).includesText('Vault CLI or REST API', 'API/CLI description present');
    assert.dom(this.element).includesText('Apply changes immediately', 'UI workflow description present');
  });

  test('it shows edit configuration section for non-UI methods', async function (assert) {
    await render(hbs`
      <Form::V2::Apply
        @form={{this.form}}
        @onBack={{this.onBack}}
        @onApply={{this.onApply}}
        @onDone={{this.onDone}}
      />
    `);

    // Terraform (default) should show edit configuration
    assert.dom('h2').includesText('Edit configuration', 'Edit configuration section shown for Terraform');

    // Switch to UI workflow
    const radios = this.element.querySelectorAll('input[type="radio"]');
    await click(radios[2]); // UI workflow is third option

    const headings = this.element.querySelectorAll('h2');
    const editConfigHeading = Array.from(headings).find((h) => h.textContent.includes('Edit configuration'));
    assert.notOk(editConfigHeading, 'Edit configuration section hidden for UI workflow');
  });
});
