import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const generated_role_rules = `rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - <list of policies to authorize>
`;

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
    this.getRole = (prefType) => {
      const data = {
        expanded: { kubernetes_role_name: 'test' },
        full: { generated_role_rules },
      }[prefType];
      const role = this.server.create('kubernetes-role', data);
      store.pushPayload('kubernetes/role', {
        modelName: 'kubernetes/role',
        backend: 'kubernetes-test',
        ...role,
      });
      return store.peekRecord('kubernetes/role', role.name);
    };

    this.newModel = store.createRecord('kubernetes/role', { backend: 'kubernetes-test' });
  });

  test('it should display placeholder when generation preference is not selected', async function (assert) {
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });
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
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });
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
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });

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

    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });
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
      this.role = this.getRole(pref);
      await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} />`, { owner: this.engine });
      assert.dom(`[data-test-radio-card="${pref}"] input`).isChecked('Correct radio card is checked');
      assert.dom('[data-test-input="name"]').hasValue(this.role.name, 'Role name is populated');
      const selector = {
        basic: { name: '[data-test-input="serviceAccountName"]', method: 'hasValue', value: 'default' },
        expanded: { name: '[data-test-input="kubernetesRoleName"]', method: 'hasValue', value: 'test' },
        full: {
          name: '[data-test-select-template]',
          method: 'hasValue',
          value: '6',
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
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });
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

  test('it should restore role rule example', async function (assert) {
    this.role = this.getRole('full');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} />`, { owner: this.engine });
    const addedText = 'this will be add to the start of the first line in the JsonEditor';
    await fillIn('[data-test-component="code-mirror-modifier"] textarea', addedText);
    await click('[data-test-restore-example]');
    assert.dom('.CodeMirror-code').doesNotContainText(addedText, 'Role rules example restored');
  });

  test('it should go back to list route and clean up model', async function (assert) {
    const unloadSpy = sinon.spy(this.newModel, 'unloadRecord');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.newModel}} />`, { owner: this.engine });
    await click('[data-test-cancel]');
    assert.ok(unloadSpy.calledOnce, 'New model is unloaded on cancel');
    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list on cancel');

    this.role = this.getRole('basic');
    const rollbackSpy = sinon.spy(this.role, 'rollbackAttributes');
    await render(hbs`<Page::Role::CreateAndEdit @model={{this.role}} />`, { owner: this.engine });
    await click('[data-test-cancel]');
    assert.ok(rollbackSpy.calledOnce, 'Attributes are rolled back for existing model on cancel');
    assert.ok(this.transitionCalledWith('roles'), 'Transitions to roles list on cancel');
  });
});
