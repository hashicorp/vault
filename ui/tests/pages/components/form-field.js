import {
  attribute,
  focusable,
  value,
  clickable,
  isPresent,
  collection,
  fillable,
  text,
  triggerable,
} from 'ember-cli-page-object';

export default {
  hasStringList: isPresent('[data-test-component=string-list]'),
  hasSearchSelect: isPresent('[data-test-component=search-select]'),
  hasTextFile: isPresent('[data-test-component=text-file]'),
  hasTTLPicker: isPresent('[data-test-toggle-input="Foo"]'),
  hasJSONEditor: isPresent('[data-test-component="code-mirror-modifier"]'),
  hasJSONClearButton: isPresent('[data-test-json-clear-button]'),
  hasSelect: isPresent('select'),
  hasInput: isPresent('input'),
  hasCheckbox: isPresent('input[type=checkbox]'),
  hasTextarea: isPresent('textarea'),
  hasMaskedInput: isPresent('[data-test-masked-input]'),
  hasTooltip: isPresent('[data-test-component=info-tooltip]'),
  tooltipTrigger: focusable('[data-test-tool-tip-trigger]'),
  tooltipContent: text('[data-test-help-text]'),
  hasRadio: isPresent('[data-test-radio-input]'),
  radioButtons: collection('input[type=radio]', {
    select: clickable(),
    id: attribute('id'),
  }),

  fields: collection('[data-test-field]', {
    clickLabel: clickable('label'),
    toggleTtl: clickable('[data-test-toggle-input="Foo"]'),
    for: attribute('for', 'label', { multiple: true }),
    labelText: text('label', { multiple: true }),
    input: fillable('input'),
    ttlTime: fillable('[data-test-ttl-value]'),
    select: fillable('select'),
    textarea: fillable('textarea'),
    change: triggerable('keyup', '.input'),
    inputValue: value('input'),
    textareaValue: value('textarea'),
    inputChecked: attribute('checked', 'input[type=checkbox]'),
    selectValue: value('select'),
  }),
  selectRadioInput: async function (value) {
    return this.radioButtons.filterBy('id', value)[0].select();
  },
  fillInTextarea: async function (name, value) {
    return this.fields
      .filter((field) => {
        return field.for.includes(name);
      })[0]
      .textarea(value);
  },
  fillIn: async function (name, value) {
    return this.fields
      .filter((field) => {
        return field.for.includes(name);
      })[0]
      .input(value);
  },
};
