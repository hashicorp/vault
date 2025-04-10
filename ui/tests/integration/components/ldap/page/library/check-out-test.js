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

module('Integration | Component | ldap | Page::Library::CheckOut', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.creds = {
      account: 'foo.bar',
      password: 'password',
      lease_id: 'ldap/library/test/check-out/123',
      lease_duration: 86400,
      renewable: true,
    };
    this.breadcrumbs = [
      { label: 'ldap-test', route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: 'test-library', route: 'libraries.library' },
      { label: 'Check-Out' },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Library::CheckOut @credentials={{this.creds}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render page title and breadcrumbs', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-header-title]').hasText('Check-out', 'Page title renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(1)')
      .containsText('ldap-test', 'Overview breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(2) a')
      .containsText('Libraries', 'Libraries breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText('test-library', 'Library breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(4)')
      .containsText('Check-Out', 'Check-out breadcrumb renders');
  });

  test('it should render check out information and credentials', async function (assert) {
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();

    assert
      .dom('[data-test-alert-description]')
      .hasText(
        'You won’t be able to access these credentials later, so please copy them now.',
        'Warning alert renders'
      );
    assert.dom('[data-test-row-value="Account name"]').hasText('foo.bar', 'Account name renders');
    await click('[data-test-button="toggle-masked"]');
    assert.dom('[data-test-value-div="Password"] .masked-value').hasText('password', 'Password renders');
    assert
      .dom('[data-test-row-value="Lease ID"]')
      .hasText('ldap/library/test/check-out/123', 'Lease ID renders');
    assert
      .dom('[data-test-value-div="Lease renewable"] svg')
      .hasClass('hds-icon-check-circle', 'Lease renewable true icon renders');
    assert
      .dom('[data-test-value-div="Lease renewable"] svg')
      .hasClass('has-text-success', 'Lease renewable true icon color renders');
    assert.dom('[data-test-value-div="Lease renewable"] span').hasText('True', 'Lease renewable renders');

    this.creds.renewable = false;
    await this.renderComponent();
    assert
      .dom('[data-test-value-div="Lease renewable"] svg')
      .hasClass('hds-icon-x-circle', 'Lease renewable false icon renders');
    assert
      .dom('[data-test-value-div="Lease renewable"] svg')
      .hasClass('has-text-danger', 'Lease renewable false icon color renders');
    assert.dom('[data-test-value-div="Lease renewable"] span').hasText('False', 'Lease renewable renders');

    await click('[data-test-done]');
    const didTransition = transitionStub.calledWith(
      'vault.cluster.secrets.backend.ldap.libraries.library.details.accounts'
    );
    assert.true(didTransition, 'Transitions to accounts route on done');
  });
});
