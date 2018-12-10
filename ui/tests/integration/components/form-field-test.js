import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import sinon from 'sinon';
import formFields from '../../pages/components/form-field';

const component = create(formFields);

module('Integration | Component | form field', function(hooks) {
  setupRenderingTest(hooks);

  const createAttr = (name, type, options) => {
    return {
      name,
      type,
      options,
    };
  };

  const setup = async function(attr) {
    let model = EmberObject.create({});
    let spy = sinon.spy();
    this.set('onChange', spy);
    this.set('model', model);
    this.set('attr', attr);
    await render(hbs`{{form-field attr=attr model=model onChange=onChange}}`);
    return [model, spy];
  };

  test('it renders', async function(assert) {
    let model = EmberObject.create({});
    this.set('attr', { name: 'foo' });
    this.set('model', model);
    await render(hbs`{{form-field attr=attr model=model}}`);

    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.notOk(component.hasInput, 'renders only the label');
  });

  test('it renders: string', async function(assert) {
    let [model, spy] = await setup.call(this, createAttr('foo', 'string', { defaultValue: 'default' }));
    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.equal(component.fields[0].inputValue, 'default', 'renders default value');
    assert.ok(component.hasInput, 'renders input for string');
    await component.fields[0].input('bar').change();

    assert.equal(model.get('foo'), 'bar');
    assert.ok(spy.calledWith('foo', 'bar'), 'onChange called with correct args');
  });

  test('it renders: boolean', async function(assert) {
    let [model, spy] = await setup.call(this, createAttr('foo', 'boolean', { defaultValue: false }));
    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.notOk(component.fields[0].inputChecked, 'renders default value');
    assert.ok(component.hasCheckbox, 'renders a checkbox for boolean');
    await component.fields[0].clickLabel();

    assert.equal(model.get('foo'), true);
    assert.ok(spy.calledWith('foo', true), 'onChange called with correct args');
  });

  test('it renders: number', async function(assert) {
    let [model, spy] = await setup.call(this, createAttr('foo', 'number', { defaultValue: 5 }));
    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.equal(component.fields[0].inputValue, 5, 'renders default value');
    assert.ok(component.hasInput, 'renders input for number');
    await component.fields[0].input(8).change();

    assert.equal(model.get('foo'), 8);
    assert.ok(spy.calledWith('foo', '8'), 'onChange called with correct args');
  });

  test('it renders: object', async function(assert) {
    await setup.call(this, createAttr('foo', 'object'));
    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.ok(component.hasJSONEditor, 'renders the json editor');
  });

  test('it renders: editType textarea', async function(assert) {
    let [model, spy] = await setup.call(
      this,
      createAttr('foo', 'string', { defaultValue: 'goodbye', editType: 'textarea' })
    );
    assert.equal(component.fields[0].labelText, 'Foo', 'renders a label');
    assert.ok(component.hasTextarea, 'renders a textarea');
    assert.equal(component.fields[0].textareaValue, 'goodbye', 'renders default value');
    await component.fields[0].textarea('hello');

    assert.equal(model.get('foo'), 'hello');
    assert.ok(spy.calledWith('foo', 'hello'), 'onChange called with correct args');
  });

  test('it renders: editType file', async function(assert) {
    await setup.call(this, createAttr('foo', 'string', { editType: 'file' }));
    assert.ok(component.hasTextFile, 'renders the text-file component');
  });

  test('it renders: editType ttl', async function(assert) {
    let [model, spy] = await setup.call(this, createAttr('foo', null, { editType: 'ttl' }));
    assert.ok(component.hasTTLPicker, 'renders the ttl-picker component');
    await component.fields[0].input('3');
    await component.fields[0].select('h').change();

    assert.equal(model.get('foo'), '3h');
    assert.ok(spy.calledWith('foo', '3h'), 'onChange called with correct args');
  });

  test('it renders: editType stringArray', async function(assert) {
    let [model, spy] = await setup.call(this, createAttr('foo', 'string', { editType: 'stringArray' }));
    assert.ok(component.hasStringList, 'renders the string-list component');

    await component.fields[0].input('array').change();
    assert.deepEqual(model.get('foo'), ['array'], 'sets the value on the model');
    assert.deepEqual(spy.args[0], ['foo', ['array']], 'onChange called with correct args');
  });

  test('it renders: sensitive', async function(assert) {
    await setup.call(this, createAttr('password', 'string', { sensitive: true }));
    assert.ok(component.hasMaskedInput, 'renders the masked-input component');
  });

  test('it uses a passed label', async function(assert) {
    await setup.call(this, createAttr('foo', 'string', { label: 'Not Foo' }));
    assert.equal(component.fields[0].labelText, 'Not Foo', 'renders the label from options');
  });

  test('it renders a help tooltip', async function(assert) {
    await setup.call(this, createAttr('foo', 'string', { helpText: 'Here is some help text' }));
    await component.tooltipTrigger();
    assert.ok(component.hasTooltip, 'renders the tooltip component');
  });
});
