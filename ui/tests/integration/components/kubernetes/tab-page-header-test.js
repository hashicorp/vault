import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | kubernetes | TabPageHeader', function (hooks) {
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
    this.model = this.store.peekRecord('secret-engine', 'kubernetes-test');
    this.mount = this.model.path.slice(0, -1);
    this.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: this.mount }];
  });

  test('it should render breadcrumbs', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-breadcrumbs] li:nth-child(1) a').hasText('secrets', 'Secrets breadcrumb renders');

    assert
      .dom('[data-test-breadcrumbs] li:nth-child(2)')
      .containsText(this.mount, 'Mount path breadcrumb renders');
  });

  test('it should render title', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-header-title] svg')
      .hasClass('flight-icon-kubernetes', 'Correct icon renders in title');
    assert.dom('[data-test-header-title]').hasText(this.mount, 'Mount path renders in title');
  });

  test('it should render tabs', async function (assert) {
    await render(hbs`<TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-tab="overview"]').hasText('Overview', 'Overview tab renders');
    assert.dom('[data-test-tab="roles"]').hasText('Roles', 'Roles tab renders');
    assert.dom('[data-test-tab="config"]').hasText('Configuration', 'Configuration tab renders');
  });

  test('it should render filter for roles', async function (assert) {
    await render(
      hbs`<TabPageHeader @model={{this.model}} @filterRoles={{true}} @rolesFilterValue="test" @breadcrumbs={{this.breadcrumbs}} />`,
      { owner: this.engine }
    );
    assert.dom('[data-test-nav-input] input').hasValue('test', 'Filter renders with provided value');
  });

  test('it should yield block for toolbar actions', async function (assert) {
    await render(
      hbs`
      <TabPageHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}}>
        <span data-test-yield>It yields!</span>
      </TabPageHeader>
    `,
      { owner: this.engine }
    );

    assert
      .dom('.toolbar-actions [data-test-yield]')
      .hasText('It yields!', 'Block is yielded for toolbar actions');
  });
});
