/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import Service from '@ember/service';
import { click, render, settled } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import waitForError from 'vault/tests/helpers/wait-for-error';
import searchSelect from '../../pages/components/search-select';
import { isWildcardString } from 'vault/helpers/is-wildcard-string';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
        case 'some/model':
          resolve([
            { id: 'model-a-id', name: 'model-a', uuid: 'a123', type: 'a' },
            { id: 'model-b-id', name: 'model-b', uuid: 'b456', type: 'b' },
            { id: 'model-c-id', name: 'model-c', uuid: 'c789', type: 'c' },
          ]);
          break;
        case 'pki/issuer':
          resolve([
            { id: 'issuer-a-id', issuerName: 'my-first-issuer' },
            { id: 'issuer-b-id' },
            { id: 'issuer-c-id', issuerName: 'my-issuer-again' },
          ]);
          break;
        default:
          reject({ httpStatus: 404 });
          break;
      }
      reject({ httpStatus: 404 });
    });
  },
});

module('Integration | Component | search select', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const mockFunctionFromParent = (selection, dropdownOptions) => {
      const modelExists =
        !!dropdownOptions.find((opt) => opt.id === selection) ||
        !!dropdownOptions.find((opt) => opt.uuid === selection) ||
        isWildcardString([selection]);
      return !modelExists ? 'The model associated with this id no longer exists' : false;
    };
    this.set('renderTooltip', mockFunctionFromParent);
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
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    `);

    assert.ok(component.hasLabel, 'it renders the label');
    assert.strictEqual(component.labelText, 'foo', 'the label text is correct');
    assert.ok(component.hasTrigger, 'it renders the power select trigger');
    assert.strictEqual(component.selectedOptions.length, 0, 'there are no selected options');
  });

  test('it shows options when trigger is clicked', async function (assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    `);

    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    assert.strictEqual(
      component.options.objectAt(0).text,
      component.selectedOptionText,
      'first object in list is focused'
    );
  });

  test('it shows passed in options when trigger is clicked', async function (assert) {
    const options = [
      { name: 'namespace45', id: 'displayedName' },
      { name: 'name24', id: '1241' },
    ];
    this.set('options', options);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['third-option']);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @options={{this.options}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
      />
    `);
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 2, 'shows all options');
    assert.strictEqual(
      component.options.objectAt(0).text,
      component.selectedOptionText,
      'first object in list is focused'
    );

    // verify overflow styling on input field exists
    assert.dom('.list-item-text').exists('selected option field has width set');
    assert.dom('.text-overflow-ellipsis').exists('selected option text has overflow class');
  });

  test('it filters options and adds option to create new item when text is entered', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    `);

    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    assert.strictEqual(component.options.objectAt(0).text, 'seven 7', 'first option renders');
    await typeInSearch('n');
    assert.strictEqual(
      component.options.length,
      3,
      'list still shows three options, including the add option'
    );
    await typeInSearch('ni');
    assert.strictEqual(component.options.length, 2, 'list shows two options, including the add option');
    await typeInSearch('nine');
    assert.strictEqual(component.options.length, 1, 'list shows one option');
    assert.strictEqual(component.options.objectAt(0).text, 'nine 9', 'renders only matching option');
  });

  test('it counts options when wildcard is used and displays the count', async function (assert) {
    const models = ['transform/transformation'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @wildcardLabel="role"
      />
    `);

    await clickTrigger();
    await settled();
    await typeInSearch('*bar*');
    await settled();
    await component.selectOption();
    await settled();
    assert.dom('[data-test-count="2"]').exists('correctly counts with wildcard filter and shows the count');
    assert.strictEqual(
      component.selectedOptions.objectAt(0).text,
      '*bar* includes 2 roles',
      'renders correct selected text'
    );
  });

  test('it behaves correctly if new items not allowed', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @disallowNewItems={{true}}
      />
    `);

    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    await typeInSearch('p');
    assert.strictEqual(component.options.length, 1, 'list shows one option');
    assert.strictEqual(component.options[0].text, 'No results found');
    await clickTrigger();
    assert.ok(this.onChange.notCalled, 'on change not called when empty state clicked');
  });

  test('it moves option from drop down to list when clicked', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
    <div class='box'>
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    </div>
    `);
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    await component.selectOption();
    await settled();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith(['7']));
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 2, 'shows two options');
  });

  test('it pre-populates list with passed in selectedOptions', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
      />
    `);
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.strictEqual(component.selectedOptions.objectAt(0).text, 'eight 8', 'selected option renders');
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 2, 'shows two options');
  });

  test('it adds discarded list items back into select', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('inputValue', ['8']);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
      />
    `);

    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    await component.deleteButtons.objectAt(0).click();
    await settled();
    assert.strictEqual(component.selectedOptions.length, 0, 'there are no selected options');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith([]));
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    assert.strictEqual(
      component.options.objectAt(2).text,
      'eight 8',
      'previously selected option returns to dropdown and renders properly'
    );
  });

  test('it adds created item to list items on create and removes without adding back to options on delete', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    `);
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    await typeInSearch('n');
    assert.strictEqual(
      component.options.length,
      3,
      'list still shows three options, including the add option'
    );
    await typeInSearch('ni');
    await component.selectOption();
    await settled();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(this.onChange.calledWith(['ni']));
    await component.deleteButtons.objectAt(0).click();
    await settled();
    assert.strictEqual(component.selectedOptions.length, 0, 'there are no selected options');
    assert.ok(this.onChange.calledWith([]));
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 3, 'does not add deleted option back to list');
  });

  test('it uses fallback component if endpoint 403s', async function (assert) {
    const models = ['policy/rgp'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="Policy"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @fallbackComponent="string-list"
      />
    `);
    assert.ok(component.hasStringList);
  });

  test('it uses pre-populates fallback component with inputValue if endpoint 403s', async function (assert) {
    const models = ['policy/rgp'];
    this.set('models', models);
    this.set('inputValue', ['1']);
    await render(hbs`
      <SearchSelect
        @label="Policy"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @fallbackComponent="string-list"
      />
    `);
    assert.dom('[data-test-string-list-input="0"]').hasValue('1');
  });

  test('it shows no results if endpoint 404s', async function (assert) {
    const models = ['test'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @fallbackComponent="string-list"
      />
    `);

    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 1, 'prompts for search to add new options');
    assert.strictEqual(
      component.options.objectAt(0).text,
      'Type to search',
      'text of option shows Type to search'
    );
  });

  test('it shows add suggestion if there are no models', async function (assert) {
    const models = [];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @fallbackComponent="string-list"
      />
    `);
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.length, 1);
    assert.strictEqual(
      component.options.objectAt(0).text,
      'Type to search',
      'no options in dropdown, just Type to search prompt'
    );
    await typeInSearch('new-model');
    assert.strictEqual(
      component.options.objectAt(0).text,
      'Click to add new item: new-model',
      'shows the create suggestion'
    );
  });

  test('it shows selected items not in the returned response and if one model 404s', async function (assert) {
    const models = ['test', 'policy/acl'];
    this.set('models', models);
    this.set('inputValue', ['test-1', 'test-2']);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @fallbackComponent="string-list"
      />
    `);
    assert.strictEqual(component.selectedOptions.objectAt(0).text, 'test-1', 'renders first selected option');
    assert.strictEqual(
      component.selectedOptions.objectAt(1).text,
      'test-2',
      'renders second selected option'
    );
    await clickTrigger();
    assert.strictEqual(component.options.objectAt(0).text, '1', 'renders options from successful query');
    await typeInSearch('new-item');
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.objectAt(2).text, 'new-item', 'renders newly added item');
    assert.ok(
      this.onChange.calledWith(['test-1', 'test-2', 'new-item']),
      'onChange called with all three items'
    );
  });

  test('it shows both name and smaller id for identity endpoints', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
      />
    `);
    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'shows three options');
    assert.strictEqual(component.options.objectAt(0).text, 'seven 7', 'renders correct dropdown text');
    assert.strictEqual(component.smallOptionIds.length, 3, 'shows 3 smaller id text and the name');
  });

  test('it renders correctly when model keys are not standardized', async function (assert) {
    const models = ['pki/issuer'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('disallowNewItems', true);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @shouldRenderName={{true}}
        @nameKey="issuerName"
        @disallowNewItems={{this.disallowNewItems}}
      />
    `);
    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'shows three options');
    assert.strictEqual(
      component.options.objectAt(0).text,
      'my-first-issuer issuer-a-id',
      'first option renders custom ID and name'
    );
    assert.strictEqual(
      component.options.objectAt(1).text,
      'issuer-b-id',
      `second option renders only id at custom key`
    );
    await typeInSearch('issuer-a');
    await settled();
    assert.strictEqual(
      component.options.length,
      2,
      'shows two options after filter, filtering on both name and id keys'
    );
    this.set('disallowNewItems', false);
    await typeInSearch('new-issuer');
    await settled();
    assert.strictEqual(component.options.length, 1, 'shows suggestion');
    assert.strictEqual(
      component.options.objectAt(0).text,
      'Click to add new item: new-issuer',
      'Prompts to add new item'
    );
  });

  test('it does not show name and smaller id for non-identity endpoints', async function (assert) {
    const models = ['policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    await render(hbs`
    <SearchSelect
    @label="foo"
    @models={{this.models}}
    @onChange={{this.onChange}}
    @inputValue={{this.inputValue}}
    @fallbackComponent="string-list"
    />
    `);
    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'shows all options');
    assert.strictEqual(component.options.objectAt(0).text, '1', 'renders just id');
    assert.strictEqual(component.smallOptionIds.length, 0, 'only shows the regular sized id');
  });

  test('it throws an error if endpoint 500s', async function (assert) {
    const models = ['server/error'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    const promise = waitForError();
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
      />
    `);
    const err = await promise;
    assert.ok(err.message.includes('internal server error'), 'it throws an internal server error');
  });

  test('it queries multiple models', async function (assert) {
    const models = ['identity/entity', 'policy/acl'];
    this.set('models', models);
    this.set('onChange', sinon.spy());

    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
      />
    `);
    await clickTrigger();
    assert.strictEqual(component.options.length, 6, 'shows options from both models');
    assert.strictEqual(component.options.objectAt(0).text, 'seven 7', 'first dropdown item renders');
    assert.strictEqual(component.options.objectAt(5).text, '3 3', 'last dropdown item renders');
  });

  test('it returns array with objects instead of strings if passObject=true', async function (assert) {
    const models = ['identity/entity'];
    this.set('models', models);
    this.set('onChange', sinon.spy());
    this.set('passObject', true);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @passObject={{this.passObject}}
      />
    `);
    await clickTrigger();
    await settled();
    // First select existing option
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert.ok(this.onChange.calledOnce);
    assert.ok(
      this.onChange.calledWith([{ id: '7', isNew: false }]),
      'onClick is called with array of single object with isNew false'
    );
    // Then create a new item and select it
    await clickTrigger();
    await settled();
    await typeInSearch('newItem');
    await component.selectOption();
    await settled();
    assert.ok(
      this.onChange.calledWith([
        { id: '7', isNew: false },
        { id: 'newItem', isNew: true },
      ]),
      'onClick is called with array of objects with isNew true on new item'
    );
  });

  test(`it returns custom object if passObject=true and multiple objectKeys with objectKeys[0]='id'`, async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', ['id', 'uuid']);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @passObject={{true}}
        @objectKeys={{this.objectKeys}}
      />
    `);

    await clickTrigger();
    await settled();

    // First select existing option
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .hasText('model-a-id', 'does not render name if first objectKey is id');
    assert.ok(this.onChange.calledOnce);
    assert.ok(
      this.onChange.calledWith([{ id: 'model-a-id', isNew: false, uuid: 'a123' }]),
      'onClick is called with array of single object with keys: id, uuid'
    );
    // Then create a new item and select it
    await clickTrigger();
    await settled();
    await typeInSearch('newItem');
    await component.selectOption();
    await settled();
    assert.propEqual(
      spy.args[1][0],
      [
        {
          id: 'model-a-id',
          isNew: false,
          uuid: 'a123',
        },
        {
          id: 'newItem',
          isNew: true,
        },
      ],
      'onClick is called with array of objects with isNew=true (and no additional keys) on new item'
    );
  });

  test('it returns custom object and renders name if passObject=true and multiple objectKeys', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const objectKeys = ['uuid', 'name'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @passObject={{true}}
        @objectKeys={{this.objectKeys}}
      />
    `);

    await clickTrigger();
    await settled();

    // First select existing option
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .hasText('model-a a123', `renders name and ${objectKeys[0]} if first objectKey is not id`);
    assert.dom('[data-test-smaller-id]').exists();
    assert.propEqual(
      spy.args[0][0],
      [
        {
          id: 'model-a-id',
          isNew: false,
          name: 'model-a',
          uuid: 'a123',
        },
      ],
      `onClick is called with array of single object: isNew=false, and has keys: ${objectKeys.join(', ')}`
    );
  });

  test('it renders ids if model does not have the passed objectKeys as an attribute', async function (assert) {
    const models = ['policy/acl'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @objectKeys={{this.objectKeys}}
      />
    `);

    await clickTrigger();
    await settled();

    // First select existing option
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .hasText('1', 'renders model id if does not have objectKey as an attribute');
    assert.propEqual(spy.args[0][0], ['1'], 'onClick is called with array of single id string');
  });

  test('it renders when passObject=true and model does not have the passed objectKeys as an attr', async function (assert) {
    const models = ['policy/acl'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @passObject={{true}}
        @objectKeys={{this.objectKeys}}
      />
    `);

    await clickTrigger();
    await settled();
    // First select existing option
    await component.selectOption();
    assert.strictEqual(component.selectedOptions.length, 1, 'there is 1 selected option');
    assert
      .dom(GENERAL.searchSelect.selectedOption(0))
      .hasText('1', 'renders model id if does not have objectKey');
    assert.propEqual(
      spy.args[0][0],
      [
        {
          id: '1',
          isNew: false,
        },
      ],
      'onClick is called with array of single object with correct keys'
    );
  });

  test('it renders when passed multiple models, passObject=true and one model does not have the attr in objectKeys', async function (assert) {
    const models = ['policy/acl', 'some/model'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @passObject={{true}}
        @objectKeys={{this.objectKeys}}
      />
    `);
    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.objectAt(0).text, '1', 'first option renders just id as name');
    assert.strictEqual(
      component.options.objectAt(3).text,
      'model-a a123',
      `4 option renders both name and ${objectKeys[0]}`
    );

    // First select options with and without id
    await component.selectOption();
    await clickTrigger();
    await settled();
    await click('[data-option-index="2"]');
    const expectedArray = [
      {
        id: '1',
        isNew: false,
      },
      {
        id: 'model-a-id',
        isNew: false,
        uuid: 'a123',
      },
    ];
    assert.propEqual(
      spy.args[1][0],
      expectedArray,
      `onClick is called with array of objects and correct keys.
      first object: ${Object.keys(expectedArray[0]).join(', ')},
      second object: ${Object.keys(expectedArray[1]).join(', ')}`
    );
  });

  test('it renders when passed multiple models, passedObject=false and one model does not have the attr in objectKeys', async function (assert) {
    const models = ['policy/acl', 'some/model'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @objectKeys={{this.objectKeys}}
      />
    `);

    await clickTrigger();
    await settled();
    assert.strictEqual(component.options.objectAt(0).text, '1', 'first option is just id as name');
    assert.strictEqual(
      component.options.objectAt(3).text,
      'model-a a123',
      `4th option has both name and ${objectKeys[0]}`
    );

    // First select options with and without id
    await component.selectOption();
    await clickTrigger();
    await settled();
    await click('[data-option-index="2"]');
    assert.propEqual(spy.args[1][0], ['1', 'model-a-id'], 'onClick is called with array of id strings');
  });

  test('it renders a tooltip beside selection if does not match a record returned from query when passObject=false, passed objectKeys', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    const inputValue = ['a123', 'non-existent-model'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    this.set('inputValue', inputValue);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @objectKeys={{this.objectKeys}}
        @inputValue={{this.inputValue}}
        @renderTooltip={{this.renderTooltip}}
      />
      `);
    assert.strictEqual(component.selectedOptions.length, 2, 'there are two selected options');
    assert.dom(GENERAL.searchSelect.selectedOption(0)).hasText('model-a a123');
    assert.dom(GENERAL.searchSelect.selectedOption(1)).hasText('non-existent-model');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(0)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model that exists');

    assert
      .dom(`${GENERAL.searchSelect.selectedOption(1)} ${GENERAL.tooltip('search-select')}`)
      .exists('renders info tooltip for model not returned from query');
  });

  test('it renders a tooltip beside selection if does not match a record returned from query when passObject=true, passed objectKeys', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const objectKeys = ['uuid'];
    const inputValue = ['a123', 'non-existent-model'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('objectKeys', objectKeys);
    this.set('inputValue', inputValue);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @objectKeys={{this.objectKeys}}
        @inputValue={{this.inputValue}}
        @passObject={{true}}
        @renderTooltip={{this.renderTooltip}}
      />
    `);

    assert.strictEqual(component.selectedOptions.length, 2, 'there are two selected options');
    assert.dom(GENERAL.searchSelect.selectedOption(0)).hasText('model-a a123');
    assert.dom(GENERAL.searchSelect.selectedOption(1)).hasText('non-existent-model');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(0)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model that exists');

    assert
      .dom(`${GENERAL.searchSelect.selectedOption(1)} ${GENERAL.tooltip('search-select')}`)
      .exists('renders info tooltip for model not returned from query');
  });

  test('it renders a tooltip beside selection if does not match a record returned from query when passObject=true and idKey=id', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const inputValue = ['model-a-id', 'non-existent-model'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('inputValue', inputValue);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @passObject={{true}}
        @renderTooltip={{this.renderTooltip}}
      />
    `);

    assert.strictEqual(component.selectedOptions.length, 2, 'there are two selected options');
    assert.dom(GENERAL.searchSelect.selectedOption(0)).hasText('model-a-id');
    assert.dom(GENERAL.searchSelect.selectedOption(1)).hasText('non-existent-model');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(0)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model that exists');

    assert
      .dom(`${GENERAL.searchSelect.selectedOption(1)} ${GENERAL.tooltip('search-select')}`)
      .exists('renders info tooltip for model not returned from query');
  });

  test('it renders a tooltip beside selection if does not match a record returned from query when passObject=false and idKey=id', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const inputValue = ['model-a-id', 'non-existent-model', 'wildcard*'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('inputValue', inputValue);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @passObject={{false}}
        @renderTooltip={{this.renderTooltip}}
      />
    `);
    assert.strictEqual(component.selectedOptions.length, 3, 'there are three selected options');
    assert.dom(GENERAL.searchSelect.selectedOption(0)).hasText('model-a-id');
    assert.dom(GENERAL.searchSelect.selectedOption(1)).hasText('non-existent-model');
    assert.dom(GENERAL.searchSelect.selectedOption(2)).hasText('wildcard*');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(0)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model that exists');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(1)} ${GENERAL.tooltip('search-select')}`)
      .exists('renders a tooltip for model not returned from query');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(2)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for wildcard option');
  });

  test('it does not render a tooltip beside selection if not passed @renderTooltip', async function (assert) {
    const models = ['some/model'];
    const spy = sinon.spy();
    const inputValue = ['model-a-id', 'non-existent-model', 'wildcard*'];
    this.set('models', models);
    this.set('onChange', spy);
    this.set('inputValue', inputValue);
    await render(hbs`
      <SearchSelect
        @label="foo"
        @models={{this.models}}
        @onChange={{this.onChange}}
        @inputValue={{this.inputValue}}
        @passObject={{false}}
      />
    `);
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(0)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model that exists');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(1)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for model not returned from query');
    assert
      .dom(`${GENERAL.searchSelect.selectedOption(2)} ${GENERAL.tooltip('search-select')}`)
      .doesNotExist('does not render info tooltip for wildcard option');
  });
});
