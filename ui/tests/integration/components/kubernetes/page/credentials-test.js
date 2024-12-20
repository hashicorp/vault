/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import { Response } from 'miragejs';
import hbs from 'htmlbars-inline-precompile';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | kubernetes | Page::Credentials', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    this.backend = 'kubernetes-test';
    this.roleName = 'role-0';

    this.getCreateCredentialsError = (roleName, errorType = null) => {
      let errors;

      if (errorType === 'noNamespace') {
        errors = ["'kubernetes_namespace' is required"];
      } else {
        errors = [`role '${roleName}' does not exist`];
      }

      this.server.post(`/kubernetes-test/creds/${roleName}`, () => {
        return new Response(400, {}, { errors });
      });
    };
    this.breadcrumbs = [
      { label: this.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: this.roleName, route: 'roles.role.details' },
      { label: 'Credentials' },
    ];
    this.renderComponent = () => {
      return render(
        hbs`<Page::Credentials @backend={{this.backend}} @roleName={{this.roleName}} @breadcrumbs={{this.breadcrumbs}}/>`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should display generate credentials form', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-credentials-header]').hasText('Generate credentials');
    assert
      .dom('[data-test-generate-credentials] p')
      .hasText(`This will generate credentials using the role ${this.roleName}.`);
    assert.dom('[data-test-generate-credentials] label').hasText('Kubernetes namespace');
    assert
      .dom('[data-test-generate-credentials] .is-size-8')
      .hasText('The namespace in which to generate the credentials.');
    assert.dom('[data-test-toggle-label] .title').hasText('ClusterRoleBinding');
    assert
      .dom('[data-test-toggle-label] .description')
      .hasText(
        'Generate a ClusterRoleBinding to grant permissions across the whole cluster instead of within a namespace. This requires the Vault role to have kubernetes_role_type set to ClusterRole.'
      );
  });

  test('it should show errors states when generating credentials', async function (assert) {
    assert.expect(2);

    this.getCreateCredentialsError(this.roleName, 'noNamespace');
    await this.renderComponent();
    await click('[data-test-generate-credentials-button]');

    assert.dom('[data-test-message-error-description]').hasText("'kubernetes_namespace' is required");

    this.roleName = 'role-2';
    this.getCreateCredentialsError(this.roleName);

    await this.renderComponent();
    await click('[data-test-generate-credentials-button]');
    assert.dom('[data-test-message-error-description]').hasText(`role '${this.roleName}' does not exist`);
  });

  test('it should show correct credential information after generate credentials is clicked', async function (assert) {
    assert.expect(15);

    this.server.post('/kubernetes-test/creds/role-0', () => {
      assert.ok('POST request made to generate credentials');
      return {
        request_id: '58fefc6c-5195-c17a-94f2-8f889f3df57c',
        lease_id: 'kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38',
        renewable: false,
        lease_duration: 3600,
        data: {
          service_account_name: 'default',
          service_account_namespace: 'default',
          service_account_token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr',
        },
      };
    });

    await this.renderComponent();
    await fillIn('[data-test-kubernetes-namespace]', 'kubernetes-test');
    assert.dom('[data-test-kubernetes-namespace]').hasValue('kubernetes-test', 'kubernetes-test');

    await click('[data-test-toggle-input]');
    await click('[data-test-toggle-input="Time-to-Live (TTL)"]');
    await fillIn('[data-test-ttl-value="Time-to-Live (TTL)"]', 2);
    await click('[data-test-generate-credentials-button]');

    assert.dom('[data-test-credentials-header]').hasText('Credentials');
    assert.dom('[data-test-k8-alert-title]').hasText('Warning');
    assert
      .dom('[data-test-k8-alert-message]')
      .hasText("You won't be able to access these credentials later, so please copy them now.");
    assert.dom('[data-test-row-label="Service account token"]').hasText('Service account token');
    await click('[data-test-value-div="Service account token"] [data-test-button]');
    assert
      .dom('[data-test-value-div="Service account token"] .display-only')
      .hasText('eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr');
    assert.dom('[data-test-row-label="Namespace"]').hasText('Namespace');
    assert.dom('[data-test-value-div="Namespace"]').exists();
    assert.dom('[data-test-row-label="Service account name"]').hasText('Service account name');
    assert.dom('[data-test-value-div="Service account name"]').exists();

    assert.dom('[data-test-row-label="Lease expiry"]').hasText('Lease expiry');
    assert.dom('[data-test-value-div="Lease expiry"]').hasText('April 3rd 2018, 3:15:30 PM');
    assert.dom('[data-test-row-label="lease_id"]').hasText('lease_id');
    assert
      .dom('[data-test-value-div="lease_id"] [data-test-row-value="lease_id"]')
      .hasText('kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38');
  });
});
