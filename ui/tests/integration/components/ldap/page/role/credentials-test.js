/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { duration } from 'core/helpers/format-duration';
import { dateFormat } from 'core/helpers/date-format';

module('Integration | Component | ldap | Page::Role::Credentials', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'ldap-test', route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: 'test-role', route: 'roles.role' },
      { label: 'Credentials' },
    ];
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
  });

  test('it should render page title and breadcrumbs', async function (assert) {
    this.creds = {};
    await render(
      hbs`<Page::Role::Credentials @credentials={{this.creds}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );

    assert.dom('[data-test-header-title]').hasText('Credentials', 'Page title renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(1)')
      .containsText('ldap-test', 'Overview breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').containsText('Roles', 'Roles breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText('test-role', 'Role breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(4)')
      .containsText('Credentials', 'Credentials breadcrumb renders');
  });

  test('it should render error', async function (assert) {
    this.error = { errors: ['Failed to fetch credentials for role'] };

    await render(hbs`<Page::Role::Credentials @error={{this.error}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    assert.dom('[data-test-page-error-details]').hasText(this.error.errors[0], 'Error renders');
  });

  test('it should render fields for static role', async function (assert) {
    const fields = [
      {
        label: 'Last Vault rotation',
        value: () => dateFormat([this.creds.last_vault_rotation, 'MMM d yyyy, h:mm:ss aaa'], {}),
      },
      { label: 'Password', key: 'password', isMasked: true },
      { label: 'Username', key: 'username' },
      { label: 'Rotation period', value: () => duration([this.creds.rotation_period]) },
      { label: 'Time remaining', value: () => duration([this.creds.ttl]) },
    ];
    this.creds = this.server.create('ldap-credential', 'static');

    await render(
      hbs`<Page::Role::Credentials @credentials={{this.creds}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );

    for (const field of fields) {
      assert
        .dom(`[data-test-row-label="${field.label}"]`)
        .hasText(field.label, `${field.label} label renders`);

      if (field.isMasked) {
        await click(`[data-test-value-div="${field.label}"] [data-test-button="toggle-masked"]`);
      }

      const value = field.value ? field.value() : this.creds[field.key];
      assert.dom(`[data-test-value-div="${field.label}"]`).hasText(value, `${field.label} value renders`);
    }

    await click('[data-test-done]');
    assert.true(
      this.transitionStub.calledOnceWith('vault.cluster.secrets.backend.ldap.roles.role.details'),
      'Transitions to correct route on done'
    );
  });

  test('it should render fields for dynamic role', async function (assert) {
    const fields = [
      { label: 'Distinguished Name', value: () => this.creds.distinguished_names.join(', ') },
      { label: 'Username', key: 'username', isMasked: true },
      { label: 'Password', key: 'password', isMasked: true },
      { label: 'Lease ID', key: 'lease_id' },
      { label: 'Lease duration', value: () => duration([this.creds.lease_duration]) },
      { label: 'Lease renewable', value: () => (this.creds.renewable ? 'True' : 'False') },
    ];
    this.creds = this.server.create('ldap-credential', 'dynamic');

    await render(
      hbs`<Page::Role::Credentials @credentials={{this.creds}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );

    assert
      .dom('[data-test-alert-description]')
      .hasText(
        'You won’t be able to access these credentials later, so please copy them now.',
        'Alert renders for dynamic roles'
      );

    for (const field of fields) {
      assert
        .dom(`[data-test-row-label="${field.label}"]`)
        .hasText(field.label, `${field.label} label renders`);

      if (field.isMasked) {
        await click(`[data-test-value-div="${field.label}"] [data-test-button="toggle-masked"]`);
      }

      const value = field.value ? field.value() : this.creds[field.key];
      assert.dom(`[data-test-value-div="${field.label}"]`).hasText(value, `${field.label} value renders`);
    }

    await click('[data-test-done]');
    assert.true(
      this.transitionStub.calledOnceWith('vault.cluster.secrets.backend.ldap.roles.role.details'),
      'Transitions to correct route on done'
    );
  });
});
