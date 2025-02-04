/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | kubernetes | Page::Role::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    const routerStub = sinon.stub(router, 'transitionTo');
    this.transitionCalledWith = (routeName, name) => {
      const route = `vault.cluster.secrets.backend.kubernetes.${routeName}`;
      const args = name ? [route, name] : [route];
      return routerStub.calledWith(...args);
    };

    const store = this.owner.lookup('service:store');
    this.getRole = (trait) => {
      const role = this.server.create('kubernetes-role', trait);
      store.pushPayload('kubernetes/role', {
        modelName: 'kubernetes/role',
        backend: 'kubernetes-test',
        ...role,
      });
      return store.peekRecord('kubernetes/role', role.name);
    };

    this.newModel = store.createRecord('kubernetes/role', { backend: 'kubernetes-test' });
    this.breadcrumbs = [
      { label: this.newModel.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: 'Create' },
    ];
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should display placeholder when generation preference is not selected', async function (assert) {
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('Choose an option above', 'Empty state title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'To configure a Vault role, choose what should be generated in Kubernetes by Vault.',
        'Empty state message renders'
      );
    assert.dom('[data-test-save]').isDisabled('Save button is disabled');
  });

  test('it should display different form fields based on generation preference selection', async function (assert) {
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );
    const commonFields = [
      'name',
      'allowedKubernetesNamespaces',
      'tokenMaxTtl',
      'tokenDefaultTtl',
      'annotations',
    ];

    await click('[data-test-radio-card="basic"]');
    ['serviceAccountName', ...commonFields].forEach((field) => {
      assert.dom(`[data-test-field="${field}"]`).exists(`${field} field renders`);
    });

    await click('[data-test-radio-card="expanded"]');
    ['kubernetesRoleType', 'kubernetesRoleName', 'nameTemplate', ...commonFields].forEach((field) => {
      assert.dom(`[data-test-field="${field}"]`).exists(`${field} field renders`);
    });

    await click('[data-test-radio-card="full"]');
    ['kubernetesRoleType', 'nameTemplate', ...commonFields].forEach((field) => {
      assert.dom(`[data-test-field="${field}"]`).exists(`${field} field renders`);
    });
    assert.dom('[data-test-generated-role-rules]').exists('Generated role rules section renders');
  });

  test('it should clear specific form fields when switching generation preference', async function (assert) {
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );

    await click('[data-test-radio-card="basic"]');
    await fillIn('[data-test-input="serviceAccountName"]', 'test');
    await click('[data-test-radio-card="expanded"]');
    assert.strictEqual(
      this.newModel.serviceAccountName,
      null,
      'Service account name cleared when switching from basic to expanded'
    );

    await fillIn('[data-test-input="kubernetesRoleName"]', 'test');
    await click('[data-test-radio-card="full"]');
    assert.strictEqual(
      this.newModel.kubernetesRoleName,
      null,
      'Kubernetes role name cleared when switching from expanded to full'
    );

    await click('[data-test-input="kubernetesRoleType"] input');
    await click('[data-test-toggle-input="show-nameTemplate"]');
    await fillIn('[data-test-input="nameTemplate"]', 'bar');
    await fillIn('[data-test-select-template]', '6');
    await click('[data-test-radio-card="expanded"]');
    assert.strictEqual(
      this.newModel.generatedRoleRules,
      null,
      'Role rules cleared when switching from full to expanded'
    );

    await click('[data-test-radio-card="basic"]');
    assert.strictEqual(
      this.newModel.kubernetesRoleType,
      null,
      'Kubernetes role type cleared when switching from expanded to basic'
    );
    assert.strictEqual(
      this.newModel.kubernetesRoleName,
      null,
      'Kubernetes role name cleared when switching from expanded to basic'
    );
    assert.strictEqual(
      this.newModel.nameTemplate,
      null,
      'Name template cleared when switching from expanded to basic'
    );
  });

  test('it should create new role', async function (assert) {
    assert.expect(3);

    this.server.post('/kubernetes-test/roles/role-1', () => assert.ok('POST request made to save role'));

    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );
    await click('[data-test-radio-card="basic"]');
    await click('[data-test-save]');
    assert.dom('[data-test-inline-error-message]').hasText('Name is required', 'Validation error renders');
    await fillIn('[data-test-input="name"]', 'role-1');
    await fillIn('[data-test-input="serviceAccountName"]', 'default');
    await click('[data-test-save]');
    assert.ok(
      this.transitionCalledWith('roles.role.details', this.newModel.name),
      'Transitions to details route on save'
    );
  });

  test('it should populate fields when editing role', async function (assert) {
    assert.expect(15);

    this.server.post('/kubernetes-test/roles/:name', () => assert.ok('POST request made to save role'));

    for (const pref of ['basic', 'expanded', 'full']) {
      const trait = { expanded: 'withRoleName', full: 'withRoleRules' }[pref];
      this.role = this.getRole(trait);
      await render(
        hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}} />`,
        { owner: this.engine }
      );
      assert.dom(`[data-test-radio-card="${pref}"] input`).isChecked('Correct radio card is checked');
      assert.dom('[data-test-input="name"]').hasValue(this.role.name, 'Role name is populated');
      const selector = {
        basic: { name: '[data-test-input="serviceAccountName"]', method: 'hasValue', value: 'default' },
        expanded: {
          name: '[data-test-input="kubernetesRoleName"]',
          method: 'hasValue',
          value: 'vault-k8s-secrets-role',
        },
        full: {
          name: '[data-test-select-template]',
          method: 'hasValue',
          value: '5',
        },
      }[pref];
      assert.dom(selector.name)[selector.method](selector.value);
      await click('[data-test-save]');
      assert.ok(
        this.transitionCalledWith('roles.role.details', this.role.name),
        'Transitions to details route on save'
      );
    }
  });

  test('it should show and hide annotations and labels', async function (assert) {
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );
    await click('[data-test-radio-card="basic"]');
    assert.dom('[data-test-annotations]').doesNotExist('Annotations and labels are hidden');

    await click('[data-test-field="annotations"]');
    await fillIn('[data-test-kv="annotations"] [data-test-kv-key]', 'foo');
    await fillIn('[data-test-kv="annotations"] [data-test-kv-value]', 'bar');
    await click('[data-test-kv="annotations"] [data-test-kv-add-row]');
    assert.deepEqual(this.newModel.extraAnnotations, { foo: 'bar' }, 'Annotations set');

    await fillIn('[data-test-kv="labels"] [data-test-kv-key]', 'bar');
    await fillIn('[data-test-kv="labels"] [data-test-kv-value]', 'baz');
    await click('[data-test-kv="labels"] [data-test-kv-add-row]');
    assert.deepEqual(this.newModel.extraLabels, { bar: 'baz' }, 'Labels set');
  });

  test('it should expand annotations and labels when editing if they were populated', async function (assert) {
    this.role = this.getRole();
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}}/>`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-annotations]')
      .doesNotExist('Annotations and labels are collapsed initially when not defined');
    this.role = this.getRole('withRoleRules');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}}/>`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-annotations]')
      .exists('Annotations and labels are expanded initially when defined');
  });

  test('it should restore role rule example', async function (assert) {
    this.role = this.getRole('withRoleRules');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}}/>`, {
      owner: this.engine,
    });
    const addedText = 'this will be add to the start of the first line in the JsonEditor';
    await fillIn('[data-test-component="code-mirror-modifier"] textarea', addedText);
    await click('[data-test-restore-example]');
    assert.dom('.CodeMirror-code').doesNotContainText(addedText, 'Role rules example restored');
  });

  test('it should set generatedRoleRoles model prop on save', async function (assert) {
    assert.expect(1);

    this.server.post('/kubernetes-test/roles/role-1', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      const role = this.server.create('kubernetes-role', 'withRoleRules');
      assert.strictEqual(
        payload.generated_role_rules,
        role.generated_role_rules,
        'Generated roles rules are passed in save request'
      );
    });

    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}}/>`,
      { owner: this.engine }
    );
    await click('[data-test-radio-card="full"]');
    await fillIn('[data-test-input="name"]', 'role-1');
    await fillIn('[data-test-select-template]', '5');
    await click('[data-test-save]');
  });

  test('it should unset selectedTemplateId when switching from full generation preference', async function (assert) {
    assert.expect(1);

    this.server.post('/kubernetes-test/roles/role-1', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.generated_role_rules, null, 'Generated roles rules are not set');
    });

    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}}/>`,
      { owner: this.engine }
    );
    await click('[data-test-radio-card="full"]');
    await fillIn('[data-test-input="name"]', 'role-1');
    await fillIn('[data-test-select-template]', '5');
    await click('[data-test-radio-card="basic"]');
    await fillIn('[data-test-input="serviceAccountName"]', 'default');
    await click('[data-test-save]');
  });

  test('it should go back to list route and clean up model', async function (assert) {
    const unloadSpy = sinon.spy(this.newModel, 'unloadRecord');
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}}/>`,
      { owner: this.engine }
    );
    await click('[data-test-cancel]');
    assert.ok(unloadSpy.calledOnce, 'New model is unloaded on cancel');
    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list on cancel');

    this.role = this.getRole();
    const rollbackSpy = sinon.spy(this.role, 'rollbackAttributes');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    await click('[data-test-cancel]');
    assert.ok(rollbackSpy.calledOnce, 'Attributes are rolled back for existing model on cancel');
    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list on cancel');
  });

  test('it should check for form errors', async function (assert) {
    await render(
      hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} @breadcrumbs={{this.breadcrumbs}}/>`,
      { owner: this.engine }
    );
    await click('[data-test-radio-card="basic"]');
    await click('[data-test-save]');
    assert
      .dom('[data-test-input="name"]')
      .hasClass('has-error-border', 'shows border error on input with error');
    assert.dom('[data-test-inline-error-message]').hasText('Name is required');
    assert
      .dom('[data-test-invalid-form-alert] [data-test-inline-error-message]')
      .hasText('There is an error with this form.');
  });

  test('it should save edited role with correct properties', async function (assert) {
    assert.expect(1);

    this.role = this.getRole();

    this.server.post('/kubernetes-test/roles/:name', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = {
        name: 'role-0',
        service_account_name: 'demo',
        kubernetes_role_type: 'Role',
        allowed_kubernetes_namespaces: '*',
        token_max_ttl: 86400,
        token_default_ttl: 600,
      };
      assert.deepEqual(expected, data, 'POST request made to save role with correct properties');
    });

    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    await fillIn('[data-test-input="serviceAccountName"]', 'demo');
    await click('[data-test-save]');
  });
});
