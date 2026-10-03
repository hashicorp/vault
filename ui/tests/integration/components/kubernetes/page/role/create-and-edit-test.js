/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn, waitFor, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror, { getCodeEditorValue, setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import KubernetesRoleForm from 'vault/forms/secrets/kubernetes/role';
import { getRules } from 'kubernetes/utils/generated-role-rules';

module('Integration | Component | kubernetes | Page::Role::CreateAndEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kubernetes-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.transitionCalledWith = (routeName, name) => {
      const route = `vault.cluster.secrets.backend.kubernetes.${routeName}`;
      const args = name ? [route, name] : [route];
      return routerStub.calledWith(...args);
    };

    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'kubernetesWriteRole').resolves();

    this.setupEdit = (trait) => {
      const role = this.server.create('kubernetes-role', trait);
      this.form = new KubernetesRoleForm(role);
      return role;
    };

    this.form = new KubernetesRoleForm({}, { isNew: true });

    this.breadcrumbs = [
      { label: this.backend, route: 'overview' },
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

    this.renderComponent = () =>
      render(hbs`<Page::Role::CreateAndEdit @form={{this.form}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
  });

  test('it should display placeholder when generation preference is not selected', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('Choose an option above', 'Empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'To configure a Vault role, choose what should be generated in Kubernetes by Vault.',
        'Empty state message renders'
      );
    assert.dom(GENERAL.submitButton).isDisabled('Save button is disabled');
  });

  test('it should display different form fields based on generation preference selection', async function (assert) {
    const commonFields = [
      'name',
      'allowed_kubernetes_namespaces',
      'token_max_ttl',
      'token_default_ttl',
      'annotations',
    ];

    await this.renderComponent();

    await click('[data-test-radio-card="basic"]');
    ['service_account_name', ...commonFields].forEach((field) => {
      assert.dom(GENERAL.fieldByAttr(field)).exists(`${field} field renders`);
    });

    await click('[data-test-radio-card="expanded"]');
    ['kubernetes_role_type', 'kubernetes_role_name', 'name_template', ...commonFields].forEach((field) => {
      assert.dom(GENERAL.fieldByAttr(field)).exists(`${field} field renders`);
    });

    await click('[data-test-radio-card="full"]');
    ['kubernetes_role_type', 'name_template', ...commonFields].forEach((field) => {
      assert.dom(GENERAL.fieldByAttr(field)).exists(`${field} field renders`);
    });
    assert.dom('[data-test-generated-role-rules]').exists('Generated role rules section renders');
  });

  test('it should clear specific form fields when switching generation preference', async function (assert) {
    await this.renderComponent();

    await click('[data-test-radio-card="basic"]');
    await fillIn(GENERAL.inputByAttr('service_account_name'), 'test');
    await click('[data-test-radio-card="expanded"]');
    assert.strictEqual(
      this.form.data.service_account_name,
      undefined,
      'Service account name cleared when switching from basic to expanded'
    );

    await fillIn(GENERAL.inputByAttr('kubernetes_role_name'), 'test');
    await click('[data-test-radio-card="full"]');
    assert.strictEqual(
      this.form.data.kubernetes_role_name,
      undefined,
      'Kubernetes role name cleared when switching from expanded to full'
    );

    await click(`${GENERAL.inputGroupByAttr('kubernetes_role_type')} input`);
    await click(GENERAL.toggleInput('show-name_template'));
    await fillIn(GENERAL.inputByAttr('name_template'), 'bar');
    await fillIn('[data-test-select-template]', '6');
    await click('[data-test-radio-card="expanded"]');
    assert.strictEqual(
      this.form.data.generated_role_rules,
      undefined,
      'Role rules cleared when switching from full to expanded'
    );

    await click('[data-test-radio-card="basic"]');
    assert.strictEqual(
      this.form.data.kubernetes_role_type,
      undefined,
      'Kubernetes role type cleared when switching from expanded to basic'
    );
    assert.strictEqual(
      this.form.data.kubernetes_role_name,
      undefined,
      'Kubernetes role name cleared when switching from expanded to basic'
    );
    assert.strictEqual(
      this.form.data.name_template,
      undefined,
      'Name template cleared when switching from expanded to basic'
    );
  });

  test('it should update code editor when template selection changes', async function (assert) {
    await this.renderComponent();

    await click('[data-test-radio-card="full"]');
    await waitFor('.cm-editor');
    const editor = codemirror();
    const expectedInitialValue = `# The below is an example that you can use as a starting point.
#
# rules:
#   - apiGroups: [""]
#     resources: ["serviceaccounts", "serviceaccounts/token"]
#     verbs: ["create", "update", "delete"]
#   - apiGroups: ["rbac.authorization.k8s.io"]
#     resources: ["rolebindings", "clusterrolebindings"]
#     verbs: ["create", "update", "delete"]
#   - apiGroups: ["rbac.authorization.k8s.io"]
#     resources: ["roles", "clusterroles"]
#     verbs: ["bind", "escalate", "create", "update", "delete"]
`;
    assert.strictEqual(
      getCodeEditorValue(editor),
      expectedInitialValue,
      'editor initially renders rules from example template'
    );
    // Select a different template
    await fillIn('[data-test-select-template]', '6');
    const expectedRule = `rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - <list of policies to authorize>
`;
    assert.strictEqual(
      getCodeEditorValue(editor),
      expectedRule,
      'code editor updates and renders rules from selected template'
    );
  });

  test('it should create new role', async function (assert) {
    assert.expect(3);

    await this.renderComponent();
    await click('[data-test-radio-card="basic"]');

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.validationErrorByAttr('name')).hasText('Name is required', 'Validation error renders');
    const name = 'role-1';
    await fillIn(GENERAL.inputByAttr('name'), name);
    await fillIn(GENERAL.inputByAttr('service_account_name'), 'default');
    await click(GENERAL.submitButton);

    assert.true(
      this.writeStub.calledWith(name, this.backend, { ...this.payload, service_account_name: 'default' }),
      'Write role request made with correct params'
    );
    assert.true(
      this.transitionCalledWith('roles.role.details', name),
      'Transitions to details route on save'
    );
  });

  test('it should populate fields when editing role', async function (assert) {
    assert.expect(15);

    for (const pref of ['basic', 'expanded', 'full']) {
      const trait = { expanded: 'withRoleName', full: 'withRoleRules' }[pref];

      this.role = this.setupEdit(trait);
      await this.renderComponent();

      assert.dom(`[data-test-radio-card="${pref}"] input`).isChecked('Correct radio card is checked');
      assert.dom(GENERAL.inputByAttr('name')).hasValue(this.role.name, 'Role name is populated');
      const selector = {
        basic: { name: GENERAL.inputByAttr('service_account_name'), method: 'hasValue', value: 'default' },
        expanded: {
          name: GENERAL.inputByAttr('kubernetes_role_name'),
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
      await click(GENERAL.submitButton);

      assert.true(
        this.writeStub.calledWith(this.role.name, this.backend),
        'Write role request made with correct params'
      );
      assert.true(
        this.transitionCalledWith('roles.role.details', this.role.name),
        'Transitions to details route on save'
      );
    }
  });

  test('it should show and hide annotations and labels', async function (assert) {
    await this.renderComponent();

    await click('[data-test-radio-card="basic"]');
    assert.dom('[data-test-annotations]').doesNotExist('Annotations and labels are hidden');

    await click(GENERAL.fieldByAttr('annotations'));
    await fillIn('[data-test-kv="annotations"] [data-test-kv-key]', 'foo');
    await fillIn('[data-test-kv="annotations"] [data-test-kv-value]', 'bar');
    await click(`[data-test-kv="annotations"] ${GENERAL.kvObjectEditor.addRow}`);
    assert.deepEqual(this.form.data.extra_annotations, { foo: 'bar' }, 'Annotations set');

    await fillIn('[data-test-kv="labels"] [data-test-kv-key]', 'bar');
    await fillIn('[data-test-kv="labels"] [data-test-kv-value]', 'baz');
    await click(`[data-test-kv="labels"]  ${GENERAL.kvObjectEditor.addRow}`);
    assert.deepEqual(this.form.data.extra_labels, { bar: 'baz' }, 'Labels set');
  });

  test('it should expand annotations and labels when editing if they were populated', async function (assert) {
    this.setupEdit();
    await this.renderComponent();
    assert
      .dom('[data-test-annotations]')
      .doesNotExist('Annotations and labels are collapsed initially when not defined');

    this.setupEdit('withRoleRules');
    await this.renderComponent();
    assert
      .dom('[data-test-annotations]')
      .exists('Annotations and labels are expanded initially when defined');
  });

  test('it should restore role rule example', async function (assert) {
    this.role = this.setupEdit('withRoleRules');
    await this.renderComponent();

    const addedText = 'this will be add to the start of the first line in the JsonEditor';
    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, addedText);
    await settled();
    assert.strictEqual(getCodeEditorValue(editor), addedText, 'code editor contains addedText');
    await click('[data-test-restore-example]');
    const expectedValue = `rules:
- apiGroups: [""]
  resources: ["secrets", "services"]
  verbs: ["get", "watch", "list", "create", "delete", "deletecollection", "patch", "update"]
`;
    assert.strictEqual(getCodeEditorValue(editor), expectedValue, 'code editor is reset to initial value');
    assert.strictEqual(this.role.generated_role_rules, expectedValue, 'model value matches code editor');
    assert.dom('.cm-content').doesNotContainText(addedText, 'editor does not contain added text');
  });

  test('it should set generated_role_rules prop on save', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    await click('[data-test-radio-card="full"]');
    await fillIn(GENERAL.inputByAttr('name'), 'role-1');
    await fillIn('[data-test-select-template]', '5');
    await click(GENERAL.submitButton);

    const { rules } = getRules().find((r) => r.id === '5');
    assert.true(
      this.writeStub.calledWith('role-1', this.backend, { generated_role_rules: rules }),
      'Generated roles rules are passed in save request'
    );
  });

  test('it should unset selectedTemplateId when switching from full generation preference', async function (assert) {
    assert.expect(1);

    this.server.post('/kubernetes-test/roles/role-1', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.generated_role_rules, null, 'Generated roles rules are not set');
    });

    await this.renderComponent();
    await click('[data-test-radio-card="full"]');
    await fillIn(GENERAL.inputByAttr('name'), 'role-1');
    await fillIn('[data-test-select-template]', '5');
    await click('[data-test-radio-card="basic"]');
    await fillIn(GENERAL.inputByAttr('service_account_name'), 'default');
    await click(GENERAL.submitButton);
    assert.true(
      this.writeStub.calledWith('role-1', this.backend, { service_account_name: 'default' }),
      'Save request made successfully without generated role rules'
    );
  });

  test('it should go back to list route on cancel', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.cancelButton);
    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list on cancel');
  });

  test('it should check for form errors', async function (assert) {
    await this.renderComponent();
    await click('[data-test-radio-card="basic"]');

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.validationErrorByAttr('name')).hasText('Name is required');

    assert
      .dom(`[data-test-invalid-form-alert] ${GENERAL.inlineError}`)
      .hasText('There is an error with this form.');
  });

  test('it should save edited role with correct properties', async function (assert) {
    assert.expect(1);

    this.role = this.setupEdit();
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('service_account_name'), 'demo');
    await click(GENERAL.submitButton);

    const payload = {
      ...this.role,
      service_account_name: 'demo',
    };
    delete payload.name;
    assert.true(
      this.writeStub.calledWith(this.role.name, this.backend, payload),
      'Save request made with correct params'
    );
  });
});
