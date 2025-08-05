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

const allFields = [
  { label: 'Role name', key: 'name' },
  { label: 'Kubernetes role type', key: 'kubernetesRoleType' },
  { label: 'Kubernetes role name', key: 'kubernetesRoleName' },
  { label: 'Service account name', key: 'serviceAccountName' },
  { label: 'Allowed Kubernetes namespaces', key: 'allowedKubernetesNamespaces' },
  { label: 'Max Lease TTL', key: 'tokenMaxTtl' },
  { label: 'Default Lease TTL', key: 'tokenDefaultTtl' },
  { label: 'Name template', key: 'nameTemplate' },
];

module('Integration | Component | kubernetes | Page::Role::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));
    this.renderComponent = (trait) => {
      const data = this.server.create('kubernetes-role', trait);
      store.pushPayload('kubernetes/role', {
        modelName: 'kubernetes/role',
        backend: 'kubernetes-test',
        ...data,
      });
      this.model = store.peekRecord('kubernetes/role', data.name);
      this.breadcrumbs = [
        { label: this.model.backend, route: 'overview' },
        { label: 'Roles', route: 'roles' },
        { label: this.model.name },
      ];
      return render(hbs`<Page::Role::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
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
        const modelValue = this.model[field.key];
        const value = field.key.includes('Ttl') ? duration([modelValue]) : modelValue;
        assert.dom(`[data-test-row-value="${field.label}"]`).hasText(value, `${field.label} value renders`);
      });
    };

    this.assertExtraFields = (modelKeys, assert) => {
      modelKeys.forEach((modelKey) => {
        for (const key in this.model[modelKey]) {
          assert.dom(`[data-test-row-label="${key}"]`).hasText(key, `${modelKey} key renders`);
          assert
            .dom(`[data-test-row-value="${key}"]`)
            .hasText(this.model[modelKey][key], `${modelKey} value renders`);
        }
      });
    };
  });

  test('it should render header with role name and breadcrumbs', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-header-title]').hasText(this.model.name, 'Role name renders in header');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(1)')
      .containsText(this.model.backend, 'Overview breadcrumb renders');
    assert.dom('[data-test-breadcrumbs] li:nth-child(2) a').containsText('Roles', 'Roles breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText(this.model.name, 'Role breadcrumb renders');
  });

  test('it should render toolbar actions', async function (assert) {
    assert.expect(5);

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();

    this.server.delete(`/${this.model.backend}/roles/${this.model.name}`, () => {
      assert.ok(true, 'Request made to delete role');
      return;
    });

    assert.dom('[data-test-delete]').hasText('Delete role', 'Delete action renders');
    assert
      .dom('[data-test-generate-credentials]')
      .hasText('Generate credentials', 'Generate credentials action renders');
    assert.dom('[data-test-edit]').hasText('Edit role', 'Edit action renders');

    await click('[data-test-delete]');
    await click('[data-test-confirm-button]');
    assert.ok(
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
    this.assertExtraFields(['extraAnnotations'], assert);
    assert.dom('[data-test-extra-fields="Labels"]').doesNotExist('Labels do not render');
  });

  test('it should render fields that correspond to full creation', async function (assert) {
    assert.expect(22);
    await this.renderComponent('withRoleRules');
    this.assertFilteredFields([2, 3], assert);
    assert.dom('[data-test-generated-role-rules]').exists('Generated role rules render');
    this.assertExtraFields(['extraAnnotations', 'extraLabels'], assert);
  });
});
