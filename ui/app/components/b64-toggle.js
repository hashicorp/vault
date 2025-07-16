/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { equal } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import Component from '@ember/component';
import { set, computed } from '@ember/object';
import { encodeString, decodeString } from 'vault/utils/b64';

const B64 = 'base64';
const UTF8 = 'utf-8';
export default Component.extend({
  tagName: 'button',
  attributeBindings: ['type', 'data-test-transit-b64-toggle'],
  type: 'button',
  classNames: ['button', 'b64-toggle'],
  classNameBindings: ['isInput:is-input:is-textarea'],

  /*
   * Whether or not the toggle is associated with an input.
   * Also bound to `is-input` and `is-textarea` classes
   * Defaults to true
   *
   * @public
   * @type boolean
   */

  isInput: true,

  /*
   * The value that will be mutated when the encoding is toggled
   *
   * @public
   * @type string
   */
  value: null,

  /*
   * The encoding of `value` when the component is initialized.
   * Defaults to 'utf-8'.
   * Possible values: 'utf-8' and 'base64'
   *
   * @public
   * @type string
   */
  initialEncoding: UTF8,

  /*
   * A cached version of value - used to determine if the input has changed since encoding.
   *
   * @private
   * @type string
   */
  _value: '',

  /*
   * The current encoding of `value`.
   * Possible values: 'utf-8' and 'base64'
   *
   * @private
   * @type string
   */
  currentEncoding: '',

  /*
   * The encoding when we last mutated `value`.
   * Possible values: 'utf-8' and 'base64'
   *
   * @private
   * @type string
   */
  lastEncoding: '',

  /*
   * Is the value known to be base64-encoded.
   *
   * @private
   * @type boolean
   */
  isBase64: equal('currentEncoding', B64),

  /*
   * Does the current value match the cached _value, i.e. has the input been mutated since we encoded.
   *
   * @private
   * @type boolean
   */
  valuesMatch: computed('value', '_value', function () {
    const { value, _value } = this;
    const anyBlank = isBlank(value) || isBlank(_value);
    return !anyBlank && value === _value;
  }),

  init() {
    this._super(...arguments);
    const initial = this.initialEncoding;
    set(this, 'currentEncoding', initial);
    if (initial === B64) {
      set(this, '_value', this.value);
      set(this, 'lastEncoding', B64);
    }
  },

  didReceiveAttrs() {
    this._super();
    // if there's no value, reset encoding
    if (this.value === '') {
      set(this, 'currentEncoding', UTF8);
      return;
    }
    // the value has changed after we transformed it so we reset currentEncoding
    if (this.isBase64 && !this.valuesMatch) {
      set(this, 'currentEncoding', UTF8);
    }
    // the value changed back to one we previously had transformed
    if (this.lastEncoding === B64 && this.valuesMatch) {
      set(this, 'currentEncoding', B64);
    }
  },

  click() {
    const val = this.value;
    const isUTF8 = this.currentEncoding === UTF8;
    if (!val) {
      return;
    }
    const newVal = isUTF8 ? encodeString(val) : decodeString(val);
    const encoding = isUTF8 ? B64 : UTF8;
    set(this, 'value', newVal);
    set(this, '_value', newVal);
    set(this, 'lastEncoding', encoding);
    set(this, 'currentEncoding', encoding);
  },
});
