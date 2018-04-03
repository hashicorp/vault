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
import { getter } from 'ember-cli-page-object/macros';

export default {
  hasStringList: isPresent('[data-test-component=string-list]'),
  hasTextFile: isPresent('[data-test-component=text-file]'),
  hasTTLPicker: isPresent('[data-test-component=ttl-picker]'),
  hasJSONEditor: isPresent('[data-test-component=json-editor]'),
  hasSelect: isPresent('select'),
  hasInput: isPresent('input'),
  hasCheckbox: isPresent('input[type=checkbox]'),
  hasTextarea: isPresent('textarea'),
  hasTooltip: isPresent('[data-test-component=info-tooltip]'),
  tooltipTrigger: focusable('[data-test-tool-tip-trigger]'),
  tooltipContent: text('[data-test-help-text]'),

  fields: collection({
    itemScope: '[data-test-field]',
    item: {
      clickLabel: clickable('label'),
      for: attribute('for', 'label'),
      labelText: text('label', { multiple: true }),
      input: fillable('input'),
      select: fillable('select'),
      textarea: fillable('textarea'),
      change: triggerable('keyup', 'input'),
      inputValue: value('input'),
      textareaValue: value('textarea'),
      inputChecked: attribute('checked', 'input[type=checkbox]'),
      selectValue: value('select'),
    },
    findByName(name) {
      // we use name in the label `for` attribute
      // this is consistent across all types of fields
      //(otherwise we'd have to use name on select or input or textarea)
      return this.toArray().findBy('for', name);
    },
    fillIn(name, value) {
      return this.findByName(name).input(value);
    },
  }),
  field: getter(function() {
    return this.fields(0);
  }),
};
