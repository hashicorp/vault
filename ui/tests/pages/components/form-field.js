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
  hasTTLPicker: isPresent('[data-test-component=ttl-picker]'),
  hasJSONEditor: isPresent('[data-test-component=json-editor]'),
  hasSelect: isPresent('select'),
  hasInput: isPresent('input'),
  hasCheckbox: isPresent('input[type=checkbox]'),
  hasTextarea: isPresent('textarea'),
  hasMaskedInput: isPresent('[data-test-masked-input]'),
  hasTooltip: isPresent('[data-test-component=info-tooltip]'),
  tooltipTrigger: focusable('[data-test-tool-tip-trigger]'),
  tooltipContent: text('[data-test-help-text]'),

  fields: collection('[data-test-field]', {
    clickLabel: clickable('label'),
    for: attribute('for', 'label', { multiple: true }),
    labelText: text('label', { multiple: true }),
    input: fillable('input'),
    select: fillable('select'),
    textarea: fillable('textarea'),
    change: triggerable('keyup', 'input'),
    inputValue: value('input'),
    textareaValue: value('textarea'),
    inputChecked: attribute('checked', 'input[type=checkbox]'),
    selectValue: value('select'),
  }),
  fillInTextarea: async function(name, value) {
    return this.fields
      .filter(field => {
        return field.for.includes(name);
      })[0]
      .textarea(value);
  },
  fillIn: async function(name, value) {
    return this.fields
      .filter(field => {
        return field.for.includes(name);
      })[0]
      .input(value);
  },
};
