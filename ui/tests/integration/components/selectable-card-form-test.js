import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import Service from '@ember/service';
import { click, render } from '@ember/test-helpers';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';
import { SELECTORS, TEST_LABEL } from 'vault/tests/helpers/components/selectable-card-form';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const TITLE = 'Select Card Form';
const PAGE_PATH = 'vault.cluster.secrets.backend.credentials';

const storeService = Service.extend({
  query(modelType) {
    return new Promise((resolve, reject) => {
      switch (modelType) {
        case 'database/role':
          resolve([{ id: 'my-role', backend: 'database' }]);
          break;
        default:
          reject({ httpStatus: 404, message: 'not found' });
          break;
      }
      reject({ httpStatus: 404, message: 'not found' });
    });
  },
});

module('Integration | Component | selectable-card-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.transitionTo = sinon.stub();

    this.owner.unregister('service:store');
    this.owner.register('service:store', storeService);
    this.set('title', TITLE);
    this.set('testLabel', TEST_LABEL);
    this.set('pagePath', PAGE_PATH);
  });

  hooks.afterEach(function () {
    this.router.transitionTo.reset();
  });

  test('it shows a disabled button when no item is selected', async function (assert) {
    await render(hbs`<SelectableCardForm @title={{this.title}} @testLabel={{this.testLabel}}/>`);
    assert.dom(SELECTORS.selectableCardFormButton).isDisabled();
  });

  test('it shows button that can be clicked to credentials route when an item is selected', async function (assert) {
    const models = ['database/role'];
    this.set('models', models);
    await render(
      hbs`<SelectableCardForm @title={{this.title}} @testLabel={{this.testLabel}} @placeholder="Search for a role..." @models={{this.models}} @pagePath={{this.pagePath}}/>`
    );
    assert.dom(SELECTORS.selectableCardFormInput).exists('renders search select component by default');
    assert
      .dom(SELECTORS.selectableCardFormInput)
      .hasText('Search for a role...', 'renders placeholder text passed to search select');
    await clickTrigger();
    await selectChoose('', 'my-role');
    assert.dom(SELECTORS.selectableCardFormButton).isEnabled();
    await click(SELECTORS.selectableCardFormButton);

    assert.propEqual(
      this.router.transitionTo.lastCall.args,
      ['vault.cluster.secrets.backend.credentials', 'my-role'],
      'transitionTo is called with correct route and role name'
    );
  });
});
