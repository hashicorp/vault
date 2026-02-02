/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { generateBreadcrumbs } from 'vault/tests/helpers/ldap/ldap-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import LdapConfigForm from 'vault/forms/secrets/ldap/config';

const selectors = {
  radioCard: '[data-test-radio-card="OpenLDAP"]',
  save: '[data-test-config-save]',
  binddn: '[data-test-field="binddn"] input',
  bindpass: '[data-test-input="bindpass"]',
};

module('Integration | Component | ldap | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');

  const fillAndSubmit = async (rotate) => {
    await click(selectors.radioCard);
    await fillIn(selectors.binddn, 'foo');
    await fillIn(selectors.bindpass, 'bar');
    await click(selectors.save);
    const buttonLabel = rotate === 'without' ? 'Save without rotating' : 'Save and rotate';
    await click(GENERAL.button(buttonLabel));
    return { binddn: 'foo', bindpass: 'bar', schema: 'openldap', groupattr: 'cn', userattr: 'cn' };
  };

  hooks.beforeEach(function () {
    this.newForm = new LdapConfigForm({}, { isNew: true });
    this.existingConfig = {
      schema: 'openldap',
      binddn: 'cn=vault,ou=Users,dc=hashicorp,dc=com',
      bindpass: 'foobar',
    };
    this.editForm = new LdapConfigForm(this.existingConfig);
    this.breadcrumbs = generateBreadcrumbs('ldap', 'configure');
    this.model = { promptConfig: true, form: this.newForm }; // most of the tests use newForm but set this to editForm when needed

    this.owner.lookup('service:secret-mount-path').update('ldap-new');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    const { secrets } = this.owner.lookup('service:api');
    this.configStub = sinon.stub(secrets, 'ldapConfigure').resolves();
    this.rotateStub = sinon.stub(secrets, 'ldapRotateRootCredentials').resolves();

    this.renderComponent = () =>
      render(hbs`<Page::Configure @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
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
      .dom(GENERAL.validationErrorByAttr('binddn'))
      .hasText('Administrator distinguished name is required.', 'Validation message renders for binddn');
    assert
      .dom(GENERAL.validationErrorByAttr('bindpass'))
      .hasText('Administrator password is required.', 'Validation message renders for bindpass');
    assert
      .dom('[data-test-invalid-form-message]')
      .hasText('There are 2 errors with this form.', 'Invalid form message renders');
  });

  test('it should save new configuration without rotating root password', async function (assert) {
    assert.expect(2);

    await this.renderComponent();
    const payload = await fillAndSubmit('without');

    assert.true(
      this.configStub.calledWith('ldap-new', payload),
      'Config save called with correct mount path'
    );
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should save new configuration and rotate root password', async function (assert) {
    assert.expect(3);

    await this.renderComponent();
    const payload = await fillAndSubmit('with');
    assert.true(
      this.configStub.calledWith('ldap-new', payload),
      'Config save called with correct mount path'
    );
    assert.true(this.rotateStub.calledWith('ldap-new'), 'Rotate root called with correct mount path');
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });

  test('it should populate fields when editing form', async function (assert) {
    this.model = { promptConfig: true, form: this.editForm };

    await this.renderComponent();

    assert.dom(selectors.radioCard).isChecked('Correct radio card is checked for schema value');
    assert.dom(selectors.binddn).hasValue(this.existingConfig.binddn, 'binddn value renders');

    await fillIn(selectors.binddn, 'foobar');
    await click('[data-test-config-cancel]');

    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.ldap.configuration'),
      'Transitions to configuration route on save success'
    );
  });
});
