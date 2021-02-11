import { module, test } from 'qunit';
import { run } from '@ember/runloop';
import { setupRenderingTest } from 'ember-qunit';
import Service from '@ember/service';
import { render } from '@ember/test-helpers';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';

const TITLE = 'Get Credentials';
const SEARCH_LABEL = 'Role to use';

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

module('Integration | Component | get-credentials-card', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
      this.set('title', TITLE);
      this.set('searchLabel', SEARCH_LABEL);
    });
  });

  test('it shows a disabled button when no item is selected', async function(assert) {
    await render(hbs`<GetCredentialsCard @title={{title}} @searchLabel={{searchLabel}}/>`);
    assert.dom('[data-test-get-credentials]').isDisabled();
  });

  test('it shows button that can be clicked to credentials route when an item is selected', async function(assert) {
    const models = ['database/role'];
    this.set('models', models);
    await render(hbs`<GetCredentialsCard @title={{title}} @searchLabel={{searchLabel}} @models={{models}}/>`);
    await clickTrigger();
    await selectChoose('', 'my-role');
    assert.dom('[data-test-get-credentials]').isEnabled();
  });
});
