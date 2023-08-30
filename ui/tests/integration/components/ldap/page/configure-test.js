/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { generateBreadcrumbs } from 'vault/tests/helpers/ldap';

const selectors = {
  radioCard: '[data-test-radio-card="OpenLDAP"]',
  save: '[data-test-config-save]',
  binddn: '[data-test-field="binddn"] input',
  bindpass: '[data-test-field="bindpass"] input',
};

module('Integration | Component | ldap | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  const fillAndSubmit = async (rotate) => {
    await click(selectors.radioCard);
    await fillIn(selectors.binddn, 'foo');
    await fillIn(selectors.bindpass, 'bar');
    await click(selectors.save);
    await click(`[data-test-save-${rotate}-rotate]`);
  };

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.newModel = this.store.createRecord('ldap/config', { backend: 'ldap-new' });
    this.existingConfig = {
      schema: 'openldap',
      binddn: 'cn=vault,ou=Users,dc=hashicorp,dc=com',
      bindpass: 'foobar',
    };
    this.store.pushPayload('ldap/config', {
      modelName: 'ldap/config',
      backend: 'ldap-edit',
      ...this.existingConfig,
    });
    this.editModel = this.store.peekRecord('ldap/config', 'ldap-edit');
    this.breadcrumbs = generateBreadcrumbs('ldap', 'configure');
    this.model = this.newModel; // most of the tests use newModel but set this to editModel when needed
    this.renderComponent = () => {
      return render(
        hbs`<div id="modal-wormhole"></div><Page::Configure @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
        {
          owner: this.engine,
        }
      );
    };
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
  });

  test('it should render empty state when schema is not selected', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-empty-state-title]').hasText('Choose an option', 'Empty state title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText('Pick an option above to see available configuration options', 'Empty state title renders');
    assert.dom(selectors.save).isDisabled('Save button is disabled when schema is not selected');

    await click(selectors.radioCard);
    assert
      .dom('[data-test-component="empty-state"]')
      .doesNotExist('Empty state is hidden when schema is selected');
  });

  test('it should render validation messages for invalid form', async function (assert) {
    await this.renderComponent();

    await click(selectors.radioCard);
    await click(selectors.save);

    assert
      .dom('[data-test-field="binddn"] [data-test-inline-error-message]')
      .hasText('Administrator distinguished name is required.', 'Validation message renders for binddn');
    assert
      .dom('[data-test-field="bindpass"] [data-test-inline-error-message]')
      .hasText('Administrator password is required.', 'Validation message renders for bindpass');
    assert
      .dom('[data-test-invalid-form-message] p')
      .hasText('There are 2 errors with this form.', 'Invalid form message renders');
  });

  test('it should save new configuration without rotating root password', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-new/config', () => {
      assert.ok(true, 'POST request made to save config');
      return new Response(204, {});
    });

    await this.renderComponent();
    await fillAndSubmit('without');

    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should save new configuration and rotate root password', async function (assert) {
    assert.expect(3);

    this.server.post('/ldap-new/config', () => {
      assert.ok(true, 'POST request made to save config');
      return new Response(204, {});
    });
    this.server.post('/ldap-new/rotate-root', () => {
      assert.ok(true, 'POST request made to rotate root password');
      return new Response(204, {});
    });

    await this.renderComponent();
    await fillAndSubmit('with');

    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should populate fields when editing form', async function (assert) {
    this.model = this.editModel;

    await this.renderComponent();

    assert.dom(selectors.radioCard).isChecked('Correct radio card is checked for schema value');
    assert.dom(selectors.binddn).hasValue(this.existingConfig.binddn, 'binddn value renders');

    await fillIn(selectors.binddn, 'foobar');
    await click('[data-test-config-cancel]');

    assert.strictEqual(this.model.binddn, this.existingConfig.binddn, 'Model is rolled back on cancel');
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });
});
