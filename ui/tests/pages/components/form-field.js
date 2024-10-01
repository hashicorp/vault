/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  attribute,
  focusable,
  value,
  clickable,
  isPresent,
  collection,
  fillable,
  triggerable,
  text,
} from 'ember-cli-page-object';

export default {
  hasStringList: isPresent('[data-test-component=string-list]'),
  hasTextFile: isPresent('[data-test-component=text-file]'),
  hasTTLPicker: isPresent('[data-test-toggle-input="Foo"]'),
  hasJSONEditor: isPresent('[data-test-component="code-mirror-modifier"]'),
  hasJSONClearButton: isPresent('[data-test-json-clear-button]'),
  hasInput: isPresent('input'),
  hasCheckbox: isPresent('input[type=checkbox]'),
  hasTextarea: isPresent('textarea'),
  hasMaskedInput: isPresent('[data-test-masked-input]'),
  hasTooltip: isPresent('[data-test-component=info-tooltip]'),
  tooltipTrigger: focusable('[data-test-tool-tip-trigger]'),
  hasRadio: isPresent('[data-test-radio-input]'),
  radioButtons: collection('input[type=radio]', {
    select: clickable(),
    id: attribute('id'),
  }),

  fields: collection('[data-test-field]', {
    clickLabel: clickable('label'),
    toggleTtl: clickable('[data-test-toggle-input="Foo"]'),
    labelValue: text('[data-test-form-field-label]'),
    input: fillable('input'),
    ttlTime: fillable('[data-test-ttl-value]'),
    select: fillable('select'),
    textarea: fillable('textarea'),
    change: triggerable('keyup', '.input'),
    inputValue: value('input'),
    textareaValue: value('textarea'),
    inputChecked: attribute('checked', 'input[type=checkbox]'),
  }),
  selectRadioInput: async function (value) {
    return this.radioButtons.filterBy('id', value)[0].select();
  },
};
