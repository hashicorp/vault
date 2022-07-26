import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
// import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import Service from '@ember/service';
import { render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
// import waitForError from 'vault/tests/helpers/wait-for-error';

import searchSelect from '../../pages/components/search-select';

const component = create(searchSelect);

const storeService = Service.extend({
  query(modelType) {
    return new Promise((resolve, reject) => {
      switch (modelType) {
        case 'policy/acl':
          resolve([
            { id: '1', name: '1' },
            { id: '2', name: '2' },
            { id: '3', name: '3' },
          ]);
          break;
        case 'policy/rgp':
          reject({ httpStatus: 403, message: 'permission denied' });
          break;
        case 'identity/entity':
          resolve([
            { id: '7', name: 'seven' },
            { id: '8', name: 'eight' },
            { id: '9', name: 'nine' },
          ]);
          break;
        case 'server/error':
          var error = new Error('internal server error');
          error.httpStatus = 500;
          reject(error);
          break;
        case 'transform/transformation':
          resolve([
            { id: 'foo', name: 'bar' },
            { id: 'foobar', name: '' },
            { id: 'barfoo1', name: 'different' },
          ]);
          break;
        default:
          reject({ httpStatus: 404, message: 'not found' });
          break;
      }
      reject({ httpStatus: 404, message: 'not found' });
    });
  },
});

module('Integration | Component | search select with modal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
  });

  test('it renders', async function (assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
    <SearchSelectWithModal
      @id="assignments"
      @label="assignment"
      @subText="Search for an existing assignment, or type a new name to create it."
      @model="oidc/assignment"
      @inputValue={{map-by "id" @model.assignments}}
      @onChange={{this.handleAssignmentSelection}}
      @excludeOptions={{array "allow_all"}}
      @fallbackComponent="string-list"
      @modalFormComponent="oidc/assignment-form"
      @modalSubtext="Use assignment to specify which Vault entities and groups are allowed to authenticate."
      />
`);

    assert.ok(component.hasLabel, 'it renders the label');
    assert.equal(component.labelText, 'foo', 'the label text is correct');
    assert.ok(component.hasTrigger, 'it renders the power select trigger');
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');
  });
});
