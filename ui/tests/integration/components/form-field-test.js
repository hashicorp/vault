/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import sinon from 'sinon';
import formFields from '../../pages/components/form-field';
import { format, startOfDay } from 'date-fns';

import { GENERAL } from 'vault/tests/helpers/general-selectors';

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

  // ------------------
  // LEGACY FORM FIELDS
  // ------------------

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

  test('it renders: toggleButton', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('foobar', 'toggleButton', {
        defaultValue: false,
        editType: 'toggleButton',
        helperTextEnabled: 'Toggled on',
        helperTextDisabled: 'Toggled off',
      })
    );
    assert.ok(component.hasToggleButton, 'renders a toggle button');
    assert.dom('[data-test-toggle-input]').isNotChecked();
    assert.dom('[data-test-toggle-subtext]').hasText('Toggled off');

    await component.fields.objectAt(0).toggleButton();

    assert.true(model.get('foobar'));
    assert.ok(spy.calledWith('foobar', true), 'onChange called with correct args');
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

  // --- common elements (legacy) ---

  test('it uses a passed label', async function (assert) {
    await setup.call(this, createAttr('foo', 'string', { label: 'Not Foo' }));
    assert.strictEqual(component.fields.objectAt(0).labelValue, 'Not Foo', 'renders the label from options');
  });

  test('it renders a help tooltip', async function (assert) {
    await setup.call(
      this,
      createAttr('foo', 'string', { editType: 'stringArray', helpText: 'Here is some help text' })
    );
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
    assert.dom(GENERAL.validationWarningByAttr('path')).exists('Validation warning renders');
  });

  // ---------------
  // HDS FORM FIELDS
  // ---------------
  // Note: some tests may be duplicative of the generic tests above
  //

  // ––––– editType === 'radio' / possibleValues –––––

  test('it renders: editType=radio / possibleValues - as Hds::Form::Radio::Group', async function (assert) {
    const possibleValues = ['foo', 'bar', 'baz'];
    await setup.call(this, createAttr('myfield', '-', { editType: 'radio', possibleValues }));
    const labels = findAll(`${GENERAL.inputGroupByAttr('myfield')} label`);
    const inputs = findAll(`${GENERAL.inputGroupByAttr('myfield')} input[type="radio"]`);
    assert
      .dom('.field fieldset[class^="hds-form-group"] input[type="radio"].hds-form-radio')
      .exists('renders as Hds::Form::Radio::Group');
    assert.strictEqual(inputs.length, 3, 'renders a fieldset element with 3 radio elements');
    possibleValues.forEach((possibleValue, index) => {
      assert
        .dom(labels[index])
        .hasAttribute('id', `label-${possibleValue}`, 'label has correct id')
        .hasText(possibleValue, 'label has correct text');
      assert
        .dom(inputs[index])
        .hasAttribute('id', possibleValue, 'input[type="radio"] has correct id')
        .hasAttribute(
          'data-test-radio',
          possibleValue,
          'input[type="radio"] has correct `data-test-radio` attribute'
        );
    });
  });

  test('it renders: editType=radio / possibleValues - with no selected radio', async function (assert) {
    const possibleValues = ['foo', 'bar', 'baz'];
    await setup.call(this, createAttr('myfield', '-', { editType: 'radio', possibleValues }));
    possibleValues.forEach((possibleValue) => {
      assert
        .dom(GENERAL.radioByAttr(possibleValue))
        .isNotChecked(`input[type="radio"] "${possibleValue}" is not checked`);
    });
  });

  test('it renders: editType=radio / possibleValues - with selected value and changes it', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'radio',
        possibleValues: ['foo', 'bar', 'baz'],
        defaultValue: 'baz',
      })
    );
    assert.dom(GENERAL.radioByAttr('baz')).isChecked(`input[type="radio"] "baz" is checked`);
    await click(GENERAL.radioByAttr('foo'));
    assert.strictEqual(model.get('myfield'), 'foo');
    assert.ok(spy.calledWith('myfield', 'foo'), 'onChange called with correct args');
  });

  test('it renders: editType=radio / possibleValues - with `true/false` boolean values', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'radio',
        // we need to pass custom ID or the `true` value will not be assigned as `for` argument to the label
        // see bug in HDS: https://github.com/hashicorp/design-system/pull/2863
        // once the bug is fixed, we can change this to `possibleValues: [true, false],`
        possibleValues: [
          { value: true, id: 'true-option' },
          { value: false, id: 'false-option' },
        ],
        defaultValue: true,
      })
    );
    assert.dom(GENERAL.radioByAttr('true-option')).isChecked(`input[type="radio"] "true" is checked`);
    await click(GENERAL.radioByAttr('false-option'));
    // eslint-disable-next-line qunit/no-assert-equal-boolean
    assert.strictEqual(model.get('myfield'), false);
    assert.ok(spy.calledWith('myfield', false), 'onChange called with correct args');
  });

  test('it renders: editType=radio / possibleValues - with passed custom id, label, subtext, helptext', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'radio',
        possibleValues: [
          { value: 'foo', id: 'custom-id-1' },
          { value: 'bar', label: 'Custom label 2', subText: 'Some subtext 2' },
          { value: 'baz', label: 'Custom label 3', helpText: 'Some helptext 3' },
          { value: 'qux', label: 'Custom label 4', subText: 'Some subtext 4', helpText: 'Some helptext 2' },
        ],
      })
    );
    // first item should have custom ID, label `foo`, and no subText/helpText
    assert
      .dom(GENERAL.radioByAttr('custom-id-1'))
      .hasAttribute('id', 'custom-id-1', 'renders the radio input with a custom `id` attribute');
    assert.dom(GENERAL.labelByGroupControlIndex(1)).hasText('foo', 'renders default label from `foo` value');
    assert
      .dom(GENERAL.helpTextByGroupControlIndex(1))
      .doesNotExist('does not render subtext/helptext for `foo`');
    // second item should have custom label and subText but no helpText
    assert
      .dom(GENERAL.labelByGroupControlIndex(2))
      .hasText('Custom label 2', 'renders the custom label for `bar` from options');
    assert
      .dom(GENERAL.helpTextByGroupControlIndex(2))
      .hasText('Some subtext 2', 'renders the right subtext string for `bar` from options');
    // third item should have custom label and no subText/helpText (helpText is visible only if no subText is defined for any of the items)
    assert
      .dom(GENERAL.labelByGroupControlIndex(3))
      .hasText('Custom label 3', 'renders the custom label for `baz` from options');
    assert.dom(GENERAL.helpTextByGroupControlIndex(3)).doesNotExist('does not render the helptext for `baz`');
    // fourth item should have custom label and subText but no helpText (helpText is visible only if no subText is defined for any of the items)
    assert
      .dom(GENERAL.labelByGroupControlIndex(4))
      .hasText('Custom label 4', 'renders the custom label for `qux` from options');
    assert
      .dom(GENERAL.helpTextByGroupControlIndex(4))
      .exists({ count: 1 }, 'renders only the subtext for `qux` and not the helptext')
      .hasText('Some subtext 4', 'renders the right subtext string for `qux` from options');
  });

  test('it renders: editType=radio / possibleValues - with passed helptext', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'radio',
        possibleValues: [{ value: 'foo' }, { value: 'bar', helpText: 'Some helptext 2' }],
      })
    );
    // first item should not have helpText
    assert.dom(GENERAL.helpTextByGroupControlIndex(1)).doesNotExist('does not render helptext for `foo`');
    // second item should have helpText
    assert
      .dom(GENERAL.helpTextByGroupControlIndex(2))
      .hasText('Some helptext 2', 'renders the right helptext string for `bar` from options');
  });

  test('it renders: editType=radio / possibleValues - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', '-', { editType: 'radio', possibleValues: ['foo', 'bar', 'baz'] }),
      model: { myfield: 'bar' },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– editType === 'checkboxList' / possibleValues –––––

  test('it renders: editType=checkboxList / possibleValues - as Hds::Form::Checkbox::Group', async function (assert) {
    const possibleValues = ['foo', 'bar', 'baz'];
    await setup.call(this, createAttr('myfield', '-', { editType: 'checkboxList', possibleValues }));
    const labels = findAll(`${GENERAL.inputGroupByAttr('myfield')} label`);
    const inputs = findAll(`${GENERAL.inputGroupByAttr('myfield')} input[type="checkbox"]`);
    assert
      .dom('.field [class^="hds-form-group"] input[type="checkbox"].hds-form-checkbox')
      .exists('renders as Hds::Form::Checkbox::Group');
    assert.strictEqual(inputs.length, 3, 'renders a fieldset element with 3 checkbox elements');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input group label');
    possibleValues.forEach((possibleValue, index) => {
      assert
        .dom(labels[index])
        .hasAttribute('id', `label-${possibleValue}`, 'label has correct id')
        .hasText(possibleValue, 'label has correct text');
      assert
        .dom(inputs[index])
        .hasAttribute('id', possibleValue, 'input[type="checkbox"] has correct id')
        .hasAttribute(
          'data-test-checkbox',
          possibleValue,
          'input[type="checkbox"] has correct `data-test-checkbox` attribute'
        );
    });
  });

  test('it renders: editType=checkboxList / possibleValues - with no selected checkbox', async function (assert) {
    const possibleValues = ['foo', 'bar', 'baz'];
    await setup.call(this, createAttr('myfield', '-', { editType: 'checkboxList', possibleValues }));
    possibleValues.forEach((possibleValue) => {
      assert
        .dom(GENERAL.checkboxByAttr(possibleValue))
        .isNotChecked(`input[type="checkbox"] "${possibleValue}" is not checked`);
    });
  });

  test('it renders: editType=checkboxList / possibleValues - with selected value and changes it', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'checkboxList',
        possibleValues: ['foo', 'bar', 'baz'],
        defaultValue: ['baz'],
      })
    );
    assert.dom(GENERAL.checkboxByAttr('baz')).isChecked('input[type="checkbox"] "baz" is checked');
    // select the remaining items (they're appended to the model)
    await click(GENERAL.checkboxByAttr('foo'));
    await click(GENERAL.checkboxByAttr('bar'));
    // notice: we can't use `strictEqual` here because they're different objects
    assert.deepEqual(model.get('myfield'), ['baz', 'foo', 'bar']);
    assert.ok(spy.calledWith('myfield', ['baz', 'foo', 'bar']), 'onChange called with correct args');
  });

  test('it renders: editType=checkboxList / possibleValues - with passed label, subtext, helptext, doclink', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'checkboxList',
        possibleValues: ['foo', 'bar', 'baz'],
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText(
        'Some subtext See our documentation for help.',
        'renders the right subtext string from options'
      );
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right help text string from options');
  });

  test('it renders: editType=checkboxList / possibleValues - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', '-', { editType: 'checkboxList', possibleValues: ['foo', 'bar', 'baz'] }),
      model: { myfield: 'bar' },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– editType === 'select' / possibleValues –––––

  test('it renders: editType=select / possibleValues - as Hds::Form::Select', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'string', { editType: 'select', possibleValues: ['foo', 'bar', 'baz'] })
    );
    assert
      .dom('.field [class^="hds-form-field"] select.hds-form-select')
      .exists('renders as Hds::Form::Select');
    assert
      .dom('select')
      .hasAttribute('data-test-input', 'myfield', 'select has correct `data-test-input` attribute');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the select label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('foo', 'has first option value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'bar');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('bar', 'has selected option value');
    assert.strictEqual(model.get('myfield'), 'bar');
    assert.ok(spy.calledWith('myfield', 'bar'), 'onChange called with correct args');
  });

  test('it renders: editType=select / possibleValues - with no default', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'string', {
        editType: 'select',
        possibleValues: ['foo', 'bar', 'baz'],
        noDefault: true,
      })
    );

    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('', 'has no initial value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'foo');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('foo', 'has selected option value');
    assert.strictEqual(model.get('myfield'), 'foo');
    assert.ok(spy.calledWith('myfield', 'foo'), 'onChange called with correct args');
  });

  test('it renders: editType=select / possibleValues - with selected value', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'string', {
        editType: 'select',
        possibleValues: ['foo', 'bar', 'baz'],
        defaultValue: 'baz',
      })
    );

    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('baz', 'has initial value selected');
    await fillIn(GENERAL.inputByAttr('myfield'), 'foo');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('foo', 'has selected option value');
    assert.strictEqual(model.get('myfield'), 'foo');
    assert.ok(spy.calledWith('myfield', 'foo'), 'onChange called with correct args');
  });

  test('it renders: editType=select / possibleValues - with passed label, subtext, helptext, doclink', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', 'string', {
        editType: 'select',
        possibleValues: ['foo', 'bar', 'baz'],
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText(
        'Some subtext See our documentation for help.',
        'renders the right subtext string from options'
      );
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right help text string from options');
  });

  test('it renders: editType=select / possibleValues - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', 'string', { editType: 'select', possibleValues: ['foo', 'bar', 'baz'] }),
      model: { myfield: 'bar' },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– editType === 'password' –––––

  test('it renders: editType=password / type=string - as Hds::Form::TextInput [@type=password]', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'string', { editType: 'password', defaultValue: 'default' })
    );
    assert
      .dom('.field [class^="hds-form-field"] input.hds-form-text-input')
      .exists('renders as Hds::Form::TextInput');
    assert
      .dom(`input[type="password"]`)
      .exists('renders input with type=password')
      .hasAttribute(
        'data-test-input',
        'myfield',
        'input[type="password"] has correct `data-test-input` attribute'
      );
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('default', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'bar');
    assert.strictEqual(model.get('myfield'), 'bar');
    assert.ok(spy.calledWith('myfield', 'bar'), 'onChange called with correct args');
  });

  test('it renders: editType=password / type=number - as Hds::Form::TextInput [@type=password]', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'number', { editType: 'password', defaultValue: 123 })
    );
    assert
      .dom('.field [class^="hds-form-field"] input.hds-form-text-input')
      .exists('renders as Hds::Form::TextInput');
    assert
      .dom(`input${GENERAL.inputByAttr('myfield')}[type="password"]`)
      .exists('renders input with type=password');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('123', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), 987);
    assert.strictEqual(model.get('myfield'), '987');
    assert.ok(spy.calledWith('myfield', '987'), 'onChange called with correct args');
  });

  test('it renders: editType=password / type=string - with passed label, placeholder, subtext, helptext, doclink', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', 'string', {
        editType: 'password',
        placeholder: 'Custom placeholder',
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.inputByAttr('myfield'))
      .hasAttribute('placeholder', 'Custom placeholder', 'renders the placeholder from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText(
        'Some subtext See our documentation for help.',
        'renders the right subtext string from options'
      );
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right help text string from options');
  });

  test('it renders: editType=password / type=string - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', 'string', { editType: 'password' }),
      model: { myfield: 'bar' },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– editType === 'textarea' –––––

  test('it renders: editType=textarea / type=string - as Hds::Form::Textarea', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'string', { editType: 'textarea', defaultValue: 'default' })
    );
    assert
      .dom('.field [class^="hds-form-field"] textarea.hds-form-textarea')
      .exists('renders as Hds::Form::Textarea');
    assert
      .dom(`textarea`)
      .exists('renders textarea')
      .hasAttribute('data-test-input', 'myfield', 'textarea has correct `data-test-input` attribute');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('default', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'bar');
    assert.strictEqual(model.get('myfield'), 'bar');
    assert.ok(spy.calledWith('myfield', 'bar'), 'onChange called with correct args');
  });

  test('it renders: editType=textarea / type=number - as Hds::Form::Textarea', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', 'number', { editType: 'textarea', defaultValue: 123 })
    );
    assert
      .dom('.field [class^="hds-form-field"] textarea.hds-form-textarea')
      .exists('renders as Hds::Form::Textarea');
    assert
      .dom(`textarea`)
      .exists('renders textarea')
      .hasAttribute('data-test-input', 'myfield', 'textarea has correct `data-test-input` attribute');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('123', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'bar');
    assert.strictEqual(model.get('myfield'), 'bar');
    assert.ok(spy.calledWith('myfield', 'bar'), 'onChange called with correct args');
  });

  test('it renders: editType=textarea / type=string - with passed label, placeholder, subtext, helptext, doclink', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', 'string', {
        editType: 'textarea',
        placeholder: 'Custom placeholder',
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.inputByAttr('myfield'))
      .hasAttribute('placeholder', 'Custom placeholder', 'renders the placeholder from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText(
        'Some subtext See our documentation for help.',
        'renders the right subtext string from options'
      );
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right help text string from options');
  });

  test('it renders: editType=textarea / type=string - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', 'string', { editType: 'textarea' }),
      model: { myfield: 'bar' },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– type/editType === 'boolean' –––––

  test('it renders: type=boolean - as Hds::Form::Checkbox', async function (assert) {
    await setup.call(this, createAttr('myfield', 'boolean', { defaultValue: 'false' }));
    assert
      .dom('.field [class^="hds-form-field"] input[type="checkbox"].hds-form-checkbox')
      .exists('renders as Hds::Form::Checkbox::Field');
    assert
      .dom(`input[type=checkbox]`)
      .exists('renders input[type="checkbox"]')
      .hasAttribute(
        'data-test-input',
        'myfield',
        'input[type="checkbox"] has correct `data-test-input` attribute'
      );
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input[type="checkbox"] label');
  });

  test('it renders: editType=boolean - as Hds::Form::Checkbox', async function (assert) {
    await setup.call(this, createAttr('myfield', '-', { editType: 'boolean', defaultValue: 'false' }));
    assert
      .dom('.field [class^="hds-form-field"] input[type="checkbox"].hds-form-checkbox')
      .exists('renders as Hds::Form::Checkbox::Field');
    assert
      .dom(`input[type=checkbox]`)
      .exists('renders input[type="checkbox"]')
      .hasAttribute(
        'data-test-input',
        'myfield',
        'input[type="checkbox"] has correct `data-test-input` attribute'
      );
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input[type="checkbox"] label');
  });

  test('it renders: editType=boolean - unselected by default', async function (assert) {
    await setup.call(this, createAttr('myfield', '-', { editType: 'boolean' }));
    assert.dom(GENERAL.inputByAttr('myfield')).isNotChecked('input[type="checkbox"] is not checked');
  });

  test('it renders: editType=boolean - selected and changes it', async function (assert) {
    const [model, spy] = await setup.call(
      this,
      createAttr('myfield', '-', { editType: 'boolean', defaultValue: 'true' })
    );
    assert.dom(GENERAL.inputByAttr('myfield')).isChecked('input[type="checkbox"] is checked');
    await click(GENERAL.inputByAttr('myfield'));
    assert.false(model.get('myfield'));
    assert.ok(spy.calledWith('myfield', false), 'onChange called with correct args');
  });

  test('it renders: editType=boolean - with passed label, subtext, helptext, doclink', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', '-', {
        editType: 'boolean',
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText('Some subtext Learn more here.', 'renders the right subtext string from options');
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right help text string from options');
  });

  test('it renders: editType=boolean - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', '-', { editType: 'boolean' }),
      model: { myfield: false },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });

  // ––––– editType === undefined && (type === 'string' || type === 'number') –––––

  test('it renders: editType=undefined type=string - as Hds::Form::TextInput', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('myfield', 'string', { defaultValue: 'default' }));
    assert
      .dom('.field [class^="hds-form-field"] input[type="text"].hds-form-text-input')
      .exists('renders as Hds::Form::TextInput::Field');
    assert
      .dom(`input[type=text]`)
      .exists('renders input[type="text"]')
      .hasAttribute(
        'data-test-input',
        'myfield',
        'input[type="text"] has correct `data-test-input` attribute'
      )
      .hasAttribute('name', 'myfield', 'input[type="text"] has correct `name` attribute')
      .hasAttribute('id', 'myfield', 'input[type="text"] has correct `id` attribute');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input[type="text"] label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('default', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), 'bar');
    assert.strictEqual(model.get('myfield'), 'bar');
    assert.ok(spy.calledWith('myfield', 'bar'), 'onChange called with correct args');
  });

  test('it renders: editType=undefined type=number - as Hds::Form::TextInput', async function (assert) {
    const [model, spy] = await setup.call(this, createAttr('myfield', 'number', { defaultValue: 123 }));
    assert
      .dom('.field [class^="hds-form-field"] input[type="text"].hds-form-text-input')
      .exists('renders as Hds::Form::TextInput::Field');
    assert
      .dom(`input[type=text]`)
      .exists('renders input[type="text"]')
      .hasAttribute(
        'data-test-input',
        'myfield',
        'input[type="text"] has correct `data-test-input` attribute'
      )
      .hasAttribute('name', 'myfield', 'input[type="text"] has correct `name` attribute')
      .hasAttribute('id', 'myfield', 'input[type="text"] has correct `id` attribute');
    assert.dom(GENERAL.fieldLabel()).hasText('Myfield', 'renders the input[type="text"] label');
    assert.dom(GENERAL.inputByAttr('myfield')).hasValue('123', 'renders default value');
    await fillIn(GENERAL.inputByAttr('myfield'), '1234');
    assert.strictEqual(model.get('myfield'), '1234');
    assert.ok(spy.calledWith('myfield', '1234'), 'onChange called with correct args');
  });

  test('it renders: editType=undefined - with passed label, placeholder, subText, helpText, doclink, disabled, characterLimit', async function (assert) {
    await setup.call(
      this,
      createAttr('myfield', 'string', {
        placeholder: 'Custom placeholder',
        label: 'Custom label',
        subText: 'Some subtext',
        helpText: 'Some helptext',
        docLink: '/docs',
        editDisabled: true,
        characterLimit: 10,
      })
    );
    assert.dom(GENERAL.fieldLabel()).hasText('Custom label', 'renders the custom label from options');
    assert
      .dom(GENERAL.inputByAttr('myfield'))
      .hasAttribute('placeholder', 'Custom placeholder', 'renders the placeholder from options')
      .hasAttribute('disabled', '', 'renders the disabled attribute from options')
      .hasAttribute('maxlength', '10', 'renders the characterLimit from options');
    assert
      .dom(GENERAL.helpTextByAttr('Some subtext'))
      .exists('renders `subText` option as HelperText')
      .hasText(
        'Some subtext See our documentation for help.',
        'renders the right subText string from options'
      );
    assert
      .dom(`${GENERAL.helpTextByAttr('Some subtext')} ${GENERAL.docLinkByAttr('/docs')}`)
      .exists('renders `docLink` option as as link inside the subtext');
    assert
      .dom(GENERAL.helpTextByAttr('Some helptext'))
      .exists('renders `helptext` option as HelperText')
      .hasText('Some helptext', 'renders the right helpText string from options');
  });

  test('it renders: editType=undefined - with readOnly when mode=edit', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', 'string', { readOnly: true }),
      mode: 'edit',
      model: { myfield: false },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @mode={{this.mode}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.inputByAttr('myfield'))
      .hasAttribute('readonly', '', 'renders the readOnly attribute from options');
  });

  test('it renders: editType=undefined - with validation errors and warnings', async function (assert) {
    this.setProperties({
      attr: createAttr('myfield', 'string'),
      model: { myfield: false },
      modelValidations: {
        myfield: {
          isValid: false,
          errors: ['Error message #1', 'Error message #2'],
          warnings: ['Warning message #1', 'Warning message #2'],
        },
      },
      onChange: () => {},
    });

    await render(
      hbs`<FormField @attr={{this.attr}} @model={{this.model}} @modelValidations={{this.modelValidations}} @onChange={{this.onChange}} />`
    );
    assert
      .dom(GENERAL.validationErrorByAttr('myfield'))
      .exists('Validation error renders')
      .hasText('Error message #1 Error message #2', 'Validation errors are combined');
    assert
      .dom(GENERAL.validationWarningByAttr('myfield'))
      .exists('Validation warning renders')
      .hasText('Warning message #1 Warning message #2', 'Validation warnings are combined');
  });
});
