/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import hbs from 'htmlbars-inline-precompile';
import timestamp from 'core/utils/timestamp';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kubernetes | Page::Credentials', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    this.backend = 'kubernetes-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.roleName = 'role-0';

    this.generateStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'kubernetesGenerateCredentials')
      .resolves({
        lease_duration: 3600,
        lease_id: 'kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38',
        data: {
          service_account_name: 'default',
          service_account_namespace: 'default',
          service_account_token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr',
        },
      });

    this.getCreateCredentialsError = (roleName, errorType = null) => {
      let errors;

      if (errorType === 'noNamespace') {
        errors = ["'kubernetes_namespace' is required"];
      } else {
        errors = [`role '${roleName}' does not exist`];
      }

      this.generateStub.rejects(getErrorResponse({ errors }, 400));
    };

    this.breadcrumbs = [
      { label: this.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: this.roleName, route: 'roles.role.details' },
      { label: 'Credentials' },
    ];

    this.renderComponent = () => {
      return render(hbs`<Page::Credentials @roleName={{this.roleName}} @breadcrumbs={{this.breadcrumbs}}/>`, {
        owner: this.engine,
      });
    };
  });

  test('it should display generate credentials form', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Generate credentials');
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

    await this.renderComponent();
    await fillIn('[data-test-kubernetes-namespace]', 'kubernetes-test');
    assert.dom('[data-test-kubernetes-namespace]').hasValue('kubernetes-test', 'kubernetes-test');
    await click(GENERAL.toggleInput('kubernetes-clusterRoleBinding'));
    await click(GENERAL.toggleInput('Time-to-Live (TTL)'));
    await fillIn('[data-test-ttl-value="Time-to-Live (TTL)"]', 2);
    await click('[data-test-generate-credentials-button]');

    const payload = {
      kubernetes_namespace: 'kubernetes-test',
      cluster_role_binding: true,
      ttl: '2s',
    };
    assert.true(
      this.generateStub.calledWith(this.roleName, this.backend, payload),
      'Generate credentials request made'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Credentials');
    assert.dom('[data-test-k8-alert-title]').hasText('Warning');
    assert
      .dom('[data-test-k8-alert-message]')
      .hasText("You won't be able to access these credentials later, so please copy them now.");
    assert.dom('[data-test-row-label="Service account token"]').hasText('Service account token');
    await click(`${GENERAL.infoRowValue('Service account token')} [data-test-button]`);
    assert
      .dom(`${GENERAL.infoRowValue('Service account token')} .display-only`)
      .hasText('eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr');
    assert.dom('[data-test-row-label="Namespace"]').hasText('Namespace');
    assert.dom(GENERAL.infoRowValue('Namespace')).exists();
    assert.dom('[data-test-row-label="Service account name"]').hasText('Service account name');
    assert.dom(GENERAL.infoRowValue('Service account name')).exists();

    assert.dom('[data-test-row-label="Lease expiry"]').hasText('Lease expiry');
    assert.dom(GENERAL.infoRowValue('Lease expiry')).hasText('April 3rd 2018, 3:15:30 PM');
    assert.dom('[data-test-row-label="Lease ID"]').hasText('Lease ID');
    assert
      .dom(GENERAL.infoRowValue('Lease ID'))
      .hasText('kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38');
  });
});
