import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | oidc/assignment-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.get('/identity/entity/id', () => ({
      data: {
        key_info: { '1234-12345': { name: 'test-entity' } },
        keys: ['1234-12345'],
      },
    }));
    this.server.get('/identity/group/id', () => ({
      data: {
        key_info: { 'abcdef-123': { name: 'test-group' } },
        keys: ['abcdef-123'],
      },
    }));
  });

  test('it should save new assignment', async function (assert) {
    assert.expect(6);
    this.model = this.store.createRecord('oidc/assignment');
    this.server.post('/identity/oidc/assignment/test', (schema, req) => {
      assert.ok(true, 'Request made to save assignment');
      return JSON.parse(req.requestBody);
    });
    // override capability getters
    Object.defineProperties(this.model, {
      canListGroups: { value: true },
      canListEntities: { value: true },
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-assignment-title]').hasText('Create assignment', 'Form title renders');
    assert.dom('[data-test-oidc-assignment-save]').hasText('Create', 'Save button has correct label');
    await click('[data-test-oidc-assignment-save]');
    assert
      .dom('[data-test-inline-alert]')
      .hasText('Name is required.', 'Validation message is shown for name');
    assert.equal(findAll('[data-test-inline-error-message]').length, 2, `there are two validations errors.`);
    await fillIn('[data-test-input="name"]', 'test');
    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await click('[data-test-oidc-assignment-save]');
  });

  test('it should populate fields with model data on edit view and update an assignment', async function (assert) {
    assert.expect(5);

    this.store.pushPayload('oidc/assignment', {
      modelName: 'oidc/assignment',
      name: 'test',
    });
    this.model = this.store.peekRecord('oidc/assignment', 'test');
    // override capability getters
    Object.defineProperties(this.model, {
      canListGroups: { value: true },
      canListEntities: { value: true },
    });
    const [group] = (await this.store.query('identity/group', {})).toArray();
    this.model.group_ids.addObject(group);
    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-assignment-title]').hasText('Edit assignment', 'Form title renders');
    assert.dom('[data-test-oidc-assignment-save]').hasText('Update', 'Save button has correct label');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test', 'Name input is populated with model value');
    assert.dom('[data-test-smaller-id="true"]').hasText('abcdef-123', 'group id renders in selected option');

    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
  });

  test('it should not show an error message if permissions shows at least entities or groups', async function (assert) {
    assert.expect(1);
    this.model = this.store.createRecord('oidc/assignment');

    Object.defineProperties(this.model, {
      canListGroups: { value: false },
      canListEntities: { value: true },
    });

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-empty-state-title]').doesNotExist();
  });

  test('it should show error message if permissions do not show group and entity', async function (assert) {
    assert.expect(1);
    this.model = this.store.createRecord('oidc/assignment');

    Object.defineProperties(this.model, {
      canListGroups: { value: false },
      canListEntities: { value: false },
    });

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-empty-state-title]').hasText('Permissions error');
  });
});
