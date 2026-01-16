/**
 * Copyright IBM Corp. 2016, 2025
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
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const allFields = [
  { label: 'Role name', key: 'name' },
  { label: 'Kubernetes role type', key: 'kubernetes_role_type' },
  { label: 'Kubernetes role name', key: 'kubernetes_role_name' },
  { label: 'Service account name', key: 'service_account_name' },
  { label: 'Allowed Kubernetes namespaces', key: 'allowed_kubernetes_namespaces' },
  { label: 'Max lease TTL', key: 'token_max_ttl' },
  { label: 'Default lease TTL', key: 'token_default_ttl' },
  { label: 'Name template', key: 'name_template' },
];

module('Integration | Component | kubernetes | Page::Role::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kubernetes-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.renderComponent = (trait) => {
      this.role = this.server.create('kubernetes-role', trait);
      this.capabilities = { canUpdate: true, canDelete: true, canGenerateCreds: true };
      this.breadcrumbs = [
        { label: this.backend, route: 'overview' },
        { label: 'Roles', route: 'roles' },
        { label: this.role.name },
      ];

      return render(
        hbs`<Page::Role::Details @role={{this.role}} @capabilities={{this.capabilities}} @breadcrumbs={{this.breadcrumbs}} />`,
        {
          owner: this.engine,
        }
      );
    };

    this.assertFilteredFields = (hiddenIndices, assert) => {
      const fields = allFields.filter((field, index) => !hiddenIndices.includes(index));
      assert
        .dom('[data-test-filtered-field]')
        .exists({ count: fields.length }, 'Correct number of filtered fields render');
      fields.forEach((field) => {
        assert
          .dom(`[data-test-row-label="${field.label}"]`)
          .hasText(field.label, `${field.label} label renders`);
        const fieldValue = this.role[field.key];
        const value = field.key.includes('ttl') ? duration([fieldValue]) : fieldValue;
        assert.dom(`[data-test-row-value="${field.label}"]`).hasText(value, `${field.label} value renders`);
      });
    };

    this.assertExtraFields = (roleKeys, assert) => {
      roleKeys.forEach((roleKey) => {
        for (const key in this.role[roleKey]) {
          assert.dom(`[data-test-row-label="${key}"]`).hasText(key, `${roleKey} key renders`);
          assert
            .dom(`[data-test-row-value="${key}"]`)
            .hasText(this.role[roleKey][key], `${roleKey} value renders`);
        }
      });
    };
  });

  test('it should render header with role name and breadcrumbs', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(this.role.name, 'Role name renders in header');
    assert.dom(GENERAL.breadcrumbAtIdx(0)).containsText(this.backend, 'Overview breadcrumb renders');
    assert.dom(GENERAL.breadcrumbAtIdx(1)).containsText('Roles', 'Roles breadcrumb renders');
    assert
      .dom(GENERAL.currentBreadcrumb(this.role.name))
      .containsText(this.role.name, 'Role name breadcrumb renders');
  });

  test('it should render role page header dropdown', async function (assert) {
    assert.expect(5);

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    const deleteStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'kubernetesDeleteRole')
      .resolves();

    await this.renderComponent();
    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Delete role')).hasText('Delete role', 'Delete action renders in dropdown');
    assert
      .dom(GENERAL.menuItem('Generate credentials'))
      .hasText('Generate credentials', 'Generate credentials action renders');
    assert.dom(GENERAL.menuItem('Edit role')).hasText('Edit role', 'Edit action renders');
    await click(GENERAL.menuItem('Delete role'));
    await click(GENERAL.confirmButton);

    assert.true(deleteStub.calledWith(this.role.name, this.backend), 'Request made to delete role');
    assert.true(
      transitionStub.calledWith('vault.cluster.secrets.backend.kubernetes.roles'),
      'Transitions to roles route on delete success'
    );
  });

  test('it should render fields that correspond to basic creation', async function (assert) {
    assert.expect(13);

    await this.renderComponent();
    this.assertFilteredFields([1, 2, 7], assert);

    assert.dom('[data-test-generated-role-rules]').doesNotExist('Generated role rules do not render');
    assert.dom('[data-test-extra-fields]').doesNotExist('Annotations and labels do not render');
  });

  test('it should render fields that correspond to expanded creation', async function (assert) {
    assert.expect(21);

    await this.renderComponent('withRoleName');
    this.assertFilteredFields([3], assert);

    assert.dom('[data-test-generated-role-rules]').doesNotExist('Generated role rules do not render');
    this.assertExtraFields(['extra_annotations'], assert);
    assert.dom('[data-test-extra-fields="Labels"]').doesNotExist('Labels do not render');
  });

  test('it should render fields that correspond to full creation', async function (assert) {
    assert.expect(22);

    await this.renderComponent('withRoleRules');
    this.assertFilteredFields([2, 3], assert);

    assert.dom('[data-test-generated-role-rules]').exists('Generated role rules render');
    this.assertExtraFields(['extra_annotations', 'extra_labels'], assert);
  });
});
