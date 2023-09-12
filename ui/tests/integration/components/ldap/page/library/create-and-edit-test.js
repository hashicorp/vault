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
import sinon from 'sinon';

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

    this.store = this.owner.lookup('service:store');
    this.newModel = this.store.createRecord('ldap/library', { backend: 'ldap-test' });

    this.libraryData = this.server.create('ldap-library', { name: 'test-library' });
    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      ...this.libraryData,
    });

    this.breadcrumbs = [
      { label: 'ldap', route: 'overview' },
      { label: 'libraries', route: 'libraries' },
      { label: 'create' },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Library::CreateAndEdit @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should populate form when editing', async function (assert) {
    this.model = this.store.peekRecord('ldap/library', this.libraryData.name);

    await this.renderComponent();

    assert.dom('[data-test-input="name"]').hasValue(this.libraryData.name, 'Name renders');
    [0, 1].forEach((index) => {
      assert
        .dom(`[data-test-string-list-input="${index}"]`)
        .hasValue(this.libraryData.service_account_names[index], 'Service account renders');
    });
    assert.dom('[data-test-ttl-value="Default lease TTL"]').hasAnyValue('Default lease ttl renders');
    assert.dom('[data-test-ttl-value="Max lease TTL"]').hasAnyValue('Max lease ttl renders');
    const checkInValue = this.libraryData.disable_check_in_enforcement ? 'Disabled' : 'Enabled';
    assert
      .dom(`[data-test-input="disable_check_in_enforcement"] input#${checkInValue}`)
      .isChecked('Correct radio is checked for check-in enforcement');
  });

  test('it should go back to list route and clean up model on cancel', async function (assert) {
    this.model = this.store.peekRecord('ldap/library', this.libraryData.name);
    const spy = sinon.spy(this.model, 'rollbackAttributes');

    await this.renderComponent();
    await click('[data-test-cancel]');

    assert.ok(spy.calledOnce, 'Model is rolled back on cancel');
    assert.ok(this.transitionCalledWith('libraries'), 'Transitions to libraries list route on cancel');
  });

  test('it should validate form fields', async function (assert) {
    this.model = this.newModel;

    await this.renderComponent();
    await click('[data-test-save]');

    assert
      .dom('[data-test-field-validation="name"] p')
      .hasText('Library name is required.', 'Name validation error renders');
    assert
      .dom('[data-test-field-validation="service_account_names"] p')
      .hasText('At least one service account is required.', 'Service account name validation error renders');
    assert
      .dom('[data-test-invalid-form-message] p')
      .hasText('There are 2 errors with this form.', 'Invalid form message renders');
  });

  test('it should create new library', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-test/library/new-library', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        service_account_names: 'foo@bar.com,bar@baz.com',
        ttl: '24h',
        max_ttl: '24h',
        disable_check_in_enforcement: true,
      };
      assert.deepEqual(data, expected, 'POST request made with correct properties when creating library');
    });

    this.model = this.newModel;

    await this.renderComponent();

    await fillIn('[data-test-input="name"]', 'new-library');
    await fillIn('[data-test-string-list-input="0"]', 'foo@bar.com');
    await click('[data-test-string-list-button="add"]');
    await fillIn('[data-test-string-list-input="1"]', 'bar@baz.com');
    await click('[data-test-string-list-button="add"]');
    await click('[data-test-input="disable_check_in_enforcement"] input#Disabled');
    await click('[data-test-save]');

    assert.ok(
      this.transitionCalledWith('libraries.library.details', 'new-library'),
      'Transitions to library details route on save success'
    );
  });

  test('it should save edited library with correct properties', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-test/library/test-library', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        service_account_names: this.libraryData.service_account_names[1],
        ttl: this.libraryData.ttl,
        max_ttl: this.libraryData.max_ttl,
        disable_check_in_enforcement: true,
      };
      assert.deepEqual(expected, data, 'POST request made to save library with correct properties');
    });

    this.model = this.store.peekRecord('ldap/library', this.libraryData.name);

    await this.renderComponent();

    await click('[data-test-string-list-button="delete"]');
    await click('[data-test-input="disable_check_in_enforcement"] input#Disabled');
    await click('[data-test-save]');

    assert.ok(
      this.transitionCalledWith('libraries.library.details', 'test-library'),
      'Transitions to library details route on save success'
    );
  });
});
