/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import sinon from 'sinon';
import formFields from '../../pages/components/form-field';
import { format, startOfDay } from 'date-fns';

const component = create(formFields);

module('Integration | Component | form field', function (hooks) {
  setupRenderingTest(hooks);

  const createAttr = (name, type, options) => {
    return {
      name,
      type,
      options,
    };
  };

  const setup = async function (attr) {
    // ember sets model attrs from the defaultValue key, mimicking that behavior here
    const model = EmberObject.create({ [attr.name]: attr.options?.defaultValue });
    const spy = sinon.spy();
    this.set('onChange', spy);
    this.set('model', model);
    this.set('attr', attr);
    await render(hbs`<FormField @attr={{this.attr}} @model={{this.model}} @onChange={{this.onChange}} />`);
    return [model, spy];
  };

  test('it renders', async function (assert) {
    const model = EmberObject.create({});
    this.attr = { name: 'foo' };
    this.model = model;
    await render(hbs`<FormField @attr={{this.attr}} @model={{this.model}} />`);
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Foo', 'renders a label');
    assert.notOk(component.hasInput, 'renders only the label');
  });

  test('it renders: string', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('foo', 'string', { defaultValue: 'default' }));
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Foo', 'renders a label');
    assert.strictEqual(component.fields.objectAt(0).inputValue, 'default', 'renders default value');
    assert.ok(component.hasInput, 'renders input for string');
    await component.fields.objectAt(0).input('bar').change();

    assert.strictEqual(model.get('foo'), 'bar');
    assert.ok(spy.calledWith('foo', 'bar'), 'onChange called with correct args');
  });

  test('it renders: boolean', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('foo', 'boolean', { defaultValue: false }));
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Foo', 'renders a label');
    assert.notOk(component.fields.objectAt(0).inputChecked, 'renders default value');
    assert.ok(component.hasCheckbox, 'renders a checkbox for boolean');
    await component.fields.objectAt(0).clickLabel();

    assert.true(model.get('foo'));
    assert.ok(spy.calledWith('foo', true), 'onChange called with correct args');
  });

  test('it renders: number', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('foo', 'number', { defaultValue: 5 }));
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Foo', 'renders a label');
    assert.strictEqual(component.fields.objectAt(0).inputValue, '5', 'renders default value');
    assert.ok(component.hasInput, 'renders input for number');
    await component.fields.objectAt(0).input(8).change();

    assert.strictEqual(model.get('foo'), '8');
    assert.ok(spy.calledWith('foo', '8'), 'onChange called with correct args');
  });

  test('it renders: object', async function (assert) {
    await setup.call(this, createAttr('foo', 'object'));
    assert.dom('[data-test-component="json-editor-title"]').hasText('Foo', 'renders a label');
    assert.ok(component.hasJSONEditor, 'renders the json editor');
  });

  test('it renders: string as json with clear button', async function (assert) {
    await setup.call(this, createAttr('foo', 'string', { editType: 'json', allowReset: true }));
    assert.dom('[data-test-component="json-editor-title"]').hasText('Foo', 'renders a label');
    assert.ok(component.hasJSONEditor, 'renders the json editor');
    assert.ok(component.hasJSONClearButton, 'renders button that will clear the JSON value');
  });

  test('it renders: editType textarea', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', 'string', { defaultValue: 'goodbye', editType: 'textarea' })
    );
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Foo', 'renders a label');
    assert.ok(component.hasTextarea, 'renders a textarea');
    assert.strictEqual(component.fields.objectAt(0).textareaValue, 'goodbye', 'renders default value');
    await component.fields.objectAt(0).textarea('hello');

    assert.strictEqual(model.get('foo'), 'hello');
    assert.ok(spy.calledWith('foo', 'hello'), 'onChange called with correct args');
  });

  test('it renders: editType file', async function (assert) {
    const subText = 'My subtext.';
    await setup.call(this, createAttr('foo', 'string', { editType: 'file', subText, docLink: '/docs' }));
    assert.ok(component.hasTextFile, 'renders the text-file component');
    assert
      .dom('.hds-form-helper-text')
      .hasText(
        `Select a file from your computer. ${subText} See our documentation for help.`,
        'renders subtext'
      );
    assert.dom('.hds-form-helper-text a').exists('renders doc link');
    await click('[data-test-text-toggle]');
    // assert again after toggling because subtext is rendered differently for each input
    assert
      .dom('.hds-form-helper-text')
      .hasText(`Enter the value as text. ${subText} See our documentation for help.`, 'renders subtext');
    assert.dom('.hds-form-helper-text a').exists('renders doc link');
    await fillIn('[data-test-text-file-textarea]', 'hello world');
  });

  test('it renders: editType ttl', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', null, {
        editType: 'ttl',
        helperTextDisabled: 'TTL is disabled',
        helperTextEnabled: 'TTL is enabled',
      })
    );
    assert.ok(component.hasTTLPicker, 'renders the ttl-picker component');
    assert.dom('[data-test-ttl-form-subtext]').hasText('TTL is disabled');
    assert.dom('[data-test-ttl-toggle]').isNotChecked();
    await component.fields.objectAt(0).toggleTtl();
    await component.fields.objectAt(0).select('h').change();
    await component.fields.objectAt(0).ttlTime('3');
    const expectedSeconds = `${3 * 3600}s`;
    assert.strictEqual(model.get('foo'), expectedSeconds);
    assert.ok(spy.calledWith('foo', expectedSeconds), 'onChange called with correct args');
    await component.fields.objectAt(0).toggleTtl();
    assert.ok(spy.calledWith('foo', '0'), 'onChange called with 0 when toggle off');
  });

  test('it renders: editType ttl with special settings', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', null, {
        editType: 'ttl',
        setDefault: '3600s',
        ttlOffValue: '',
      })
    );
    assert.ok(component.hasTTLPicker, 'renders the ttl-picker component');
    assert.dom('[data-test-ttl-toggle]').isChecked();
    await component.fields.objectAt(0).toggleTtl();
    assert.strictEqual(model.get('foo'), '');
    assert.ok(spy.calledWith('foo', ''), 'onChange called with correct args');
  });

  test('it renders: editType ttl without toggle', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', null, { editType: 'ttl', hideToggle: true })
    );
    await component.fields.objectAt(0).select('h').change();
    await component.fields.objectAt(0).ttlTime('3');
    const expectedSeconds = `${3 * 3600}s`;
    assert.strictEqual(model.get('foo'), expectedSeconds);
    assert.ok(spy.calledWith('foo', expectedSeconds), 'onChange called with correct args');
  });

  test('it renders: radio buttons for possible values', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', null, { editType: 'radio', possibleValues: ['SHA1', 'SHA256'] })
    );
    assert.ok(component.hasRadio, 'renders radio buttons');
    const selectedValue = 'SHA256';
    await component.selectRadioInput(selectedValue);
    assert.strictEqual(model.get('foo'), selectedValue);
    assert.ok(spy.calledWith('foo', selectedValue), 'onChange called with correct args');
  });
  test('it renders: radio buttons for possible values, labels, and subtext', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foo', null, {
        editType: 'radio',
        possibleValues: [
          { label: 'Label 1', subText: 'Some subtext 1', value: 'SHA1' },
          { label: 'Label 2', subText: 'Some subtext 2', value: 'SHA256' },
          { subText: 'Some subtext 3', value: 'SHA256' },
        ],
      })
    );
    assert.ok(component.hasRadio, 'renders radio buttons');
    const selectedValue = 'SHA256';
    await component.selectRadioInput(selectedValue);
    assert.dom('[data-test-radio-label="Label 1"]').hasTextContaining('Label 1');
    assert.dom('[data-test-radio-label="Label 2"]').hasTextContaining('Label 2');
    assert.dom('[data-test-radio-label="SHA256"]').hasTextContaining('SHA256');
    assert.dom('[data-test-radio-subText="Some subtext 1"]').hasText('Some subtext 1');
    assert.dom('[data-test-radio-subText="Some subtext 2"]').hasText('Some subtext 2');
    assert.dom('[data-test-radio-subText="Some subtext 3"]').hasText('Some subtext 3');
    assert.strictEqual(model.get('foo'), selectedValue);
    assert.ok(spy.calledWith('foo', selectedValue), 'onChange called with correct args');
  });
  test('it renders: datetimelocal', async function (assert) {
    const [model] = await setup.call(
      this,
      createAttr('bar', null, {
        editType: 'dateTimeLocal',
      })
    );
    assert.dom("[data-test-input='bar']").exists();
    await fillIn(
      "[data-test-input='bar']",
      format(startOfDay(new Date('2023-12-17T03:24:00')), "yyyy-MM-dd'T'HH:mm")
    );
    // add a click label to focus out the date we filled in above
    await click('.is-label');
    assert.deepEqual(model.get('bar'), '2023-12-17T00:00', 'sets the value on the model');
  });

  test('it renders: editType stringArray', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('foo', 'string', { editType: 'stringArray' }));
    assert.ok(component.hasStringList, 'renders the string-list component');

    await component.fields.objectAt(0).textarea('array').change();
    assert.deepEqual(model.get('foo'), ['array'], 'sets the value on the model');
    assert.deepEqual(spy.args[0], ['foo', ['array']], 'onChange called with correct args');
  });

  test('it renders: sensitive', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('password', 'string', { sensitive: true }));
    assert.ok(component.hasMaskedInput, 'renders the masked-input component');
    await component.fields.objectAt(0).textarea('secret');
    assert.strictEqual(model.get('password'), 'secret');
    assert.ok(spy.calledWith('password', 'secret'), 'onChange called with correct args');
  });

  test('it uses a passed label', async function (assert) {
    await setup.call(this, createAttr('foo', 'string', { label: 'Not Foo' }));
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Not Foo', 'renders the label from options');
  });

  test('it renders a help tooltip', async function (assert) {
    await setup.call(this, createAttr('foo', 'string', { helpText: 'Here is some help text' }));
    await component.tooltipTrigger();
    assert.ok(component.hasTooltip, 'renders the tooltip component');
  });

  test('it should not expand and toggle ttl when default 0s value is present', async function (assert) {
    assert.expect(2);

    this.setProperties({
      model: EmberObject.create({ foo: '0s' }),
      attr: createAttr('foo', null, { editType: 'ttl' }),
      onChange: () => {},
    });

    await render(hbs`<FormField @attr={{this.attr}} @model={{this.model}} @onChange={{this.onChange}} />`);
    assert
      .dom('[data-test-toggle-input="Foo"]')
      .isNotChecked('Toggle is initially unchecked when given default value');
    assert.dom('[data-test-ttl-picker-group="Foo"]').doesNotExist('Ttl input is hidden');
  });

  test('it should toggle and expand ttl when initial non default value is provided', async function (assert) {
    assert.expect(2);

    this.setProperties({
      model: EmberObject.create({ foo: '1s' }),
      attr: createAttr('foo', null, { editType: 'ttl' }),
      onChange: () => {},
    });

    await render(hbs`<FormField @attr={{this.attr}} @model={{this.model}} @onChange={{this.onChange}} />`);
    assert.dom('[data-test-toggle-input="Foo"]').isChecked('Toggle is initially checked when given value');
    assert.dom('[data-test-ttl-value="Foo"]').hasValue('1', 'Ttl input displays with correct value');
  });

  test('it should show validation warning', async function (assert) {
    const model = this.owner.lookup('service:store').createRecord('auth-method');
    model.path = 'foo bar';
    this.validations = model.validate().state;
    this.setProperties({
      model,
      attr: createAttr('path', 'string'),
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.validations}} @onChange={{this.onChange}} />`
    );
    assert.dom('[data-test-validation-warning]').exists('Validation warning renders');
  });
});
