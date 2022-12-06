import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { typeInSearch, clickTrigger, selectChoose } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | kubernetes | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      },
    });
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-config'),
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-role'),
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      ...this.server.create('kubernetes-role'),
    });
    this.model = {
      backend: this.store.peekRecord('secret-engine', 'kubernetes-test'),
      config: this.store.peekRecord('kubernetes/config', 'kubernetes-test'),
      roles: this.store.peekAll('kubernetes/role'),
    };
  });

  test('it should display role card', async function (assert) {
    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    assert.dom('[data-test-roles-card] .title').hasText('Roles');
    assert
      .dom('[data-test-roles-card] p')
      .hasText('The number of Vault roles being used to generate Kubernetes credentials.');
    assert.dom('[data-test-roles-card] a').hasText('View Roles');

    this.model.roles = [];

    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });

    assert.dom('[data-test-roles-card] a').hasText('Create Role');
  });

  test('it should display correct number of roles in role card', async function (assert) {
    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    assert.dom('[data-test-roles-card] .has-font-weight-normal').hasText('2');

    this.model.roles = [];

    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    assert.dom('[data-test-roles-card] .has-font-weight-normal').hasText('None');
  });

  test('it should display generate credentials card', async function (assert) {
    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    assert.dom('[data-test-generate-credential-card] .title').hasText('Generate credentials');
    assert
      .dom('[data-test-generate-credential-card] p')
      .hasText('Quickly generate credentials by typing the role name.');
  });

  test('it should show options for SearchSelect', async function (assert) {
    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    await clickTrigger();
    assert.strictEqual(this.element.querySelectorAll('.ember-power-select-option').length, 2);
    await typeInSearch('role-0');
    assert.strictEqual(this.element.querySelectorAll('.ember-power-select-option').length, 1);
    assert.dom('[data-test-generate-credential-card] button').isDisabled();
    await selectChoose('', '.ember-power-select-option', 2);
    assert.dom('[data-test-generate-credential-card] button').isNotDisabled();
  });

  test('it should show ConfigCta when no config is set up', async function (assert) {
    this.model.config = null;

    await render(hbs`<Page::Overview @model={{this.model}} />`, { owner: this.engine });
    assert.dom('.empty-state .empty-state-title').hasText('Kubernetes not configured');
    assert
      .dom('.empty-state .empty-state-message')
      .hasText(
        'Get started by establishing the URL of the Kubernetes API to connect to, along with some additional options.'
      );
    assert.dom('.empty-state .empty-state-actions').hasText('Configure Kubernetes');
  });
});
