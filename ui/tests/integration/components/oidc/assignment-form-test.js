import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | oidc/assignment-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should save new assignment meep', async function (assert) {
    assert.expect(5);

    this.server.post('/identity/oidc/assignment/test', (schema, req) => {
      assert.ok(true, 'Request made to save assignment');
      return JSON.parse(req.requestBody);
    });

    this.model = this.store.createRecord('oidc/assignment');
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
    await fillIn('[data-test-input="name"]', 'test');
    await click('[data-test-oidc-assignment-save]');
  });
});
