import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import Service from '@ember/service';
import { render, settled } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import waitForError from 'vault/tests/helpers/wait-for-error';

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
    await settled();
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
    await settled();
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(
      component.options.objectAt(0).text,
      component.selectedOptionText,
      'first object in list is focused'
    );
  });

  test('it filters options and adds option to create new item when text is entered', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await settled();
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'shows all options');
    await typeInSearch('n');
    assert.equal(component.options.length, 3, 'list still shows three options, including the add option');
    await typeInSearch('ni');
    assert.equal(component.options.length, 2, 'list shows two options, including the add option');
    await typeInSearch('nine');
    assert.equal(component.options.length, 1, 'list shows one option');
  });

  test('it counts options when wildcard is used and displays the count', async function(assert) {
    const models = ['transform/transformation'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange wildcardLabel="role" }}`);
    await settled();
    await clickTrigger();
    await settled();
    await typeInSearch('*bar*');
    await settled();
    await component.selectOption();
    await settled();
    assert.dom('[data-test-count="2"]').exists('correctly counts with wildcard filter and shows the count');
  });

  test('it behaves correctly if new items not allowed', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange disallowNewItems=true}}`);
    await settled();
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    await typeInSearch('p');
    assert.equal(component.options.length, 1, 'list shows one option');
    assert.equal(component.options[0].text, 'No results found');
    await clickTrigger();
    assert.ok(this.onChange.notCalled, 'on change not called when empty state clicked');
  });

  test('it moves option from drop down to list when clicked', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await settled();
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'shows all options');
    await component.selectOption();
    await settled();
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith(['7']));
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 2, 'shows two options');
  });

  test('it pre-populates list with passed in selectedOptions', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await settled();
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 2, 'shows two options');
  });

  test('it adds discarded list items back into select', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await settled();
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    await component.deleteButtons.objectAt(0).click();
    await settled();
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith([]));
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'shows all options');
  });

  test('it adds created item to list items on create and removes without adding back to options on delete', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" models=models onChange=onChange}}`);
    await settled();
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'shows all options');
    await typeInSearch('n');
    assert.equal(component.options.length, 3, 'list still shows three options, including the add option');
    await typeInSearch('ni');
    await component.selectOption();
    await settled();
    assert.equal(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith(['ni']));
    await component.deleteButtons.objectAt(0).click();
    await settled();
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');
    assert.ok(this.onChange.calledWith([]));
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 3, 'does not add deleted option back to list');
  });

  test('it uses fallback component if endpoint 403s', async function(assert) {
    const models = ['policy/rgp'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    await settled();
    assert.ok(component.hasStringList);
  });

  test('it shows no results if endpoint 404s', async function(assert) {
    const models = ['test'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    await settled();
    await clickTrigger();
    await settled();
    assert.equal(component.options.length, 1, 'prompts for search to add new options');
    assert.equal(component.options.objectAt(0).text, 'Type to search', 'text of option shows Type to search');
  });

  test('it shows add suggestion if there are no options', async function(assert) {
    const models = [];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    await settled();
    await clickTrigger();
    await settled();

    await typeInSearch('new item');
    assert.equal(component.options.objectAt(0).text, 'Add new foo: new item', 'shows the create suggestion');
  });
  test('it shows items not in the returned response', async function(assert) {
    const models = ['test'];
    this.set('models', models);
    this.set('inputValue', ['test', 'two']);
    await render(
      hbs`{{search-select label="foo" inputValue=inputValue models=models fallbackComponent="string-list" onChange=onChange}}`
    );
    await settled();
    assert.equal(component.selectedOptions.length, 2, 'renders inputOptions as selectedOptions');
  });

  test('it shows both name and smaller id for identity endpoints', async function(assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await settled();
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(component.smallOptionIds.length, 3, 'shows the smaller id text and the name');
  });

  test('it does not show name and smaller id for non-identity endpoints', async function(assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    await settled();
    await clickTrigger();
    assert.equal(component.options.length, 3, 'shows all options');
    assert.equal(component.smallOptionIds.length, 0, 'only shows the regular sized id');
  });

  test('it throws an error if endpoint 500s', async function(assert) {
    const models = ['server/error'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    let promise = waitForError();
    render(hbs`{{search-select label="foo" inputValue=inputValue models=models onChange=onChange}}`);
    let err = await promise;
    assert.ok(err.message.includes('internal server error'), 'it throws an internal server error');
  });
});
