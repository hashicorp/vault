import { module, test, skip } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import Service from '@ember/service';
import { render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

import searchSelect from '../../pages/components/search-select';

const component = create(searchSelect);

const storeService = Service.extend({
  query(modelType) {
    return new Promise((resolve, reject) => {
      switch (modelType) {
        case 'policy/acl':
          resolve([{ id: '1', name: '1' }, { id: '2', name: '2' }, { id: '3', name: '3' }]);
          break;
        case 'policy/rgp':
          reject({ httpStatus: 403, message: 'permission denied' });
          break;
        case 'identity/entity':
          resolve([{ id: '7', name: 'seven' }, { id: '8', name: 'eight' }, { id: '9', name: 'nine' }]);
          break;
        case 'server/error':
          var error = new Error('internal server error');
          error.httpStatus = 500;
          reject(error);
          break;
        default:
          reject({ httpStatus: 404, message: 'not found' });
          break;
      }
      reject({ httpStatus: 404, message: 'not found' });
    });
  },
});

module('Integration | Component | search select', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
  });

  test('it renders', async function(assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    assert.ok(component.hasLabel, 'it renders the label');
    assert.equal(component.labelText, 'foo', 'the label text is correct');
    assert.ok(component.hasTrigger, 'it renders the power select trigger');
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');
  });

  test('it shows options when trigger is clicked', async function(assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(component.options[0].text, component.selectedOptionText, 'first object in list is focused');
  });

  test('it filters options when text is entered', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    await typeInSearch('n');
    assert.equal(component.options.length, 2, 'shows two options');
  });

  test('it moves option from drop down to list when clicked', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    await component.selectOption();
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith(['7']));
    await clickTrigger();
    assert.equal(component.options.length, 2, 'shows two options');
  });

  test('it pre-populates list with passed in selectedOptions', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    await clickTrigger();
    assert.equal(component.options.length, 2, 'shows two options');
  });

  test('it adds discarded list items back into select', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    await component.deleteButtons[0].click();
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith([]));
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
  });

  test('it uses fallback component if endpoint 403s', async function(assert) {
    const models = ['policy/rgp'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    assert.ok(component.hasStringList);
  });

  test('it shows no results if endpoint 404s', async function(assert) {
    const models = ['test'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    await clickTrigger();
    assert.equal(component.options.length, 1, 'has the disabled no results option');
    assert.equal(component.options[0].text, 'No results found', 'text of option shows No results found');
  });

  test('it shows both name and smaller id for identity endpoints', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(component.smallOptionIds.length, 3, 'shows the smaller id text and the name');
  });

  test('it does not show name and smaller id for non-identity endpoints', async function(assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(component.smallOptionIds.length, 0, 'only shows the regular sized id');
  });

  skip('it throws an error if endpoint 500s', async function(assert) {
    const models = ['server/error'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    assert.throws(
      await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`)
    );
  });
});
