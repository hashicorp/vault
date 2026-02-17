/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import LdapLibraryForm from 'vault/forms/secrets/ldap/library';

module('Integration | Component | ldap | Page::Library::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    const routerStub = sinon.stub(router, 'transitionTo');
    this.transitionCalledWith = (routeName, name) => {
      const route = `vault.cluster.secrets.backend.ldap.${routeName}`;
      const args = name ? [route, name] : [route];
      return routerStub.calledWith(...args);
    };

    this.apiStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapLibraryConfigure').resolves();

    this.backend = 'ldap-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.createForm = new LdapLibraryForm(
      {
        ttl: '24h',
        max_ttl: '24h',
        disable_check_in_enforcement: 'Enabled',
      },
      { isNew: true }
    );
    this.libraryData = this.server.create('ldap-library', { name: 'test-library' });
    delete this.libraryData.id;
    this.editForm = new LdapLibraryForm({ ...this.libraryData, disable_check_in_enforcement: 'Disabled' });
    this.form = this.editForm;

    this.breadcrumbs = [
      { label: 'ldap', route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: 'Create' },
    ];

    this.renderComponent = () =>
      render(hbs`<Page::Library::CreateAndEdit @form={{this.form}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
  });

  test('it should populate form when editing', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-input="name"]').hasValue(this.libraryData.name, 'Name renders');
    assert.dom('[data-test-input="name"]').isDisabled('Name field is disabled when editing');
    [(0, 1)].forEach((index) => {
      assert
        .dom(`[data-test-string-list-input="${index}"]`)
        .hasValue(this.libraryData.service_account_names[index], 'Service account renders');
    });
    assert.dom('[data-test-ttl-value="Default lease TTL"]').hasAnyValue('Default lease ttl renders');
    assert.dom('[data-test-ttl-value="Max lease TTL"]').hasAnyValue('Max lease ttl renders');
    assert
      .dom('[data-test-input-group="disable_check_in_enforcement"] input#Disabled')
      .isChecked('Correct radio is checked for check-in enforcement');
  });

  test('it should go back to list route on cancel', async function (assert) {
    await this.renderComponent();
    await click('[data-test-cancel]');

    assert.ok(this.transitionCalledWith('libraries'), 'Transitions to libraries list route on cancel');
  });

  test('it should validate form fields', async function (assert) {
    this.form = this.createForm;

    await this.renderComponent();
    await click('[data-test-submit]');

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Library name is required.', 'Name validation error renders');
    assert
      .dom(GENERAL.validationErrorByAttr('service_account_names'))
      .hasText('At least one service account is required.', 'Service account name validation error renders');
    assert
      .dom('[data-test-invalid-form-message]')
      .hasText('There are 2 errors with this form.', 'Invalid form message renders');
  });

  test('it should create new library', async function (assert) {
    assert.expect(2);

    this.form = this.createForm;
    await this.renderComponent();

    const service_account_names = ['foo@bar.com', 'bar@baz.com'];
    await fillIn('[data-test-input="name"]', 'new-library');
    await fillIn('[data-test-string-list-input="0"]', service_account_names[0]);
    await click('[data-test-string-list-button="add"]');
    await fillIn('[data-test-string-list-input="1"]', service_account_names[1]);
    await click('[data-test-string-list-button="add"]');
    await click('[data-test-input-group="disable_check_in_enforcement"] input#Disabled');
    await click('[data-test-submit]');

    const payload = {
      service_account_names,
      disable_check_in_enforcement: true,
      ttl: '24h',
      max_ttl: '24h',
    };
    assert.true(
      this.apiStub.calledWith('new-library', this.backend, payload),
      'API called to configure new library'
    );
    assert.ok(
      this.transitionCalledWith('libraries.library.details', 'new-library'),
      'Transitions to library details route on save success'
    );
  });

  test('it should save edited library with correct properties', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    await click('[data-test-string-list-row="0"] [data-test-string-list-button="delete"]');
    await click('[data-test-input-group="disable_check_in_enforcement"] input#Disabled');
    await click('[data-test-submit]');

    const payload = {
      service_account_names: [this.libraryData.service_account_names[1]],
      ttl: this.libraryData.ttl,
      max_ttl: this.libraryData.max_ttl,
      disable_check_in_enforcement: true,
    };
    assert.true(
      this.apiStub.calledWith(this.libraryData.name, this.backend, payload),
      'API called to configure existing library'
    );
    assert.ok(
      this.transitionCalledWith('libraries.library.details', 'test-library'),
      'Transitions to library details route on save success'
    );
  });
});
