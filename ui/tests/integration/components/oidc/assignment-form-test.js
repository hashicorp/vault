/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import OidcAssingmentForm from 'vault/forms/oidc/assignment';
import sinon from 'sinon';

module('Integration | Component | oidc/assignment-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', () => {});

    this.assignment = {
      name: 'test',
      entity_ids: ['1234-12345'],
      group_ids: ['abcdef-123'],
    };
    this.entities = [{ id: '1234-12345', name: 'test-entity' }];
    this.groups = [{ id: 'abcdef-123', name: 'test-group' }];
    this.onCancel = sinon.spy();
    this.onSave = sinon.spy();

    this.writeStub = sinon.stub(this.owner.lookup('service:api').identity, 'oidcWriteAssignment').resolves();

    this.renderComponent = (assignment) => {
      this.form = new OidcAssingmentForm(assignment, { isNew: !assignment });
      return render(hbs`
        <Oidc::AssignmentForm
          @form={{this.form}}
          @entities={{this.entities}}
          @groups={{this.groups}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `);
    };

    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it should save new assignment', async function (assert) {
    assert.expect(6);

    await this.renderComponent();

    assert.dom('[data-test-oidc-assignment-save]').hasText('Create', 'Save button has correct label');
    await click('[data-test-oidc-assignment-save]');

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required.', 'Validation message is shown for name');
    assert
      .dom(GENERAL.validationErrorByAttr('target'))
      .hasText('At least one entity or group is required.', 'Validation message is shown for target');
    assert
      .dom('[data-test-invalid-form-alert]')
      .hasText('There are 2 errors with this form.', 'Renders form error count');

    await fillIn('[data-test-input="name"]', 'test');
    await click('[data-test-component="search-select"]#entities .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await click('[data-test-oidc-assignment-save]');

    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
    assert.true(
      this.writeStub.calledWith('test', { entity_ids: ['1234-12345'] }),
      'API is called with correct payload'
    );
  });

  test('it should populate fields with model data on edit view and update an assignment', async function (assert) {
    assert.expect(5);

    await this.renderComponent(this.assignment);

    assert.dom('[data-test-oidc-assignment-save]').hasText('Update', 'Save button has correct label');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test', 'Name input is populated with model value');
    assert
      .dom('[data-test-search-select="entities"] [data-test-smaller-id]')
      .hasText('1234-12345', 'entity id renders in selected option');
    assert
      .dom('[data-test-search-select="groups"] [data-test-smaller-id]')
      .hasText('abcdef-123', 'group id renders in selected option');
  });

  test('it should use fallback component on create if no permissions for entities or groups', async function (assert) {
    assert.expect(2);

    this.entities = [];
    this.groups = [];

    await this.renderComponent();

    assert
      .dom('[data-test-component="search-select"]#entities [data-test-component="string-list"]')
      .exists('entities string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#groups [data-test-component="string-list"]')
      .exists('groups string list fallback component exists');
  });

  test('it should use fallback component on edit if no permissions for entities or groups', async function (assert) {
    assert.expect(8);

    this.entities = [];
    this.groups = [];

    await this.renderComponent(this.assignment);

    assert
      .dom('[data-test-component="search-select"]#entities [data-test-component="string-list"]')
      .exists('entities string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#entities [data-test-string-list-input="0"]')
      .hasValue('1234-12345', 'first row pre-populated with model entity');
    assert
      .dom(
        '[data-test-component="search-select"]#entities [data-test-string-list-row="0"] [data-test-string-list-button="delete"]'
      )
      .exists('first row renders delete icon');
    assert
      .dom(
        '[data-test-component="search-select"]#entities [data-test-string-list-row="1"] [data-test-string-list-button="add"]'
      )
      .exists('second row renders add icon');

    assert
      .dom('[data-test-component="search-select"]#groups [data-test-component="string-list"]')
      .exists('groups string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#groups [data-test-string-list-input="0"]')
      .hasValue('abcdef-123', 'first row pre-populated with model group');
    assert
      .dom(
        '[data-test-component="search-select"]#groups [data-test-string-list-row="0"] [data-test-string-list-button="delete"]'
      )
      .exists('first row renders delete icon');
    assert
      .dom(
        '[data-test-component="search-select"]#groups [data-test-string-list-row="1"] [data-test-string-list-button="add"]'
      )
      .exists('second row renders add icon');
  });
});
