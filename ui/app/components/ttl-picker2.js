/**
 * @module TtlPicker2
 * TtlPicker2 components are used to...
 *
 * @example
 * ```js
 * <TtlPicker2 @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} [onChange] - requiredParam is...
 * @param {boolean} [defaultEnabled] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  enableTTL: false,
  unitOptions: computed(function() {
    return [
      { label: 'foo', value: 'bar' },
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),

  seconds: computed('time', 'unit', function() {
    return this.convertToSeconds(this.time, this.unit);
  }),

  helperTextUnset: 'Allow tokens to be used indefinitely',
  helperTextSet: 'Disable the use of the token after',
  helperText: computed('enableTTL', 'helperTextUnset', 'helperTextSet', function() {
    return this.enableTTL ? this.helperTextSet : this.helperTextUnset;
  }),
  errorMessage: null,
  actions: {
    toggleTTL(value) {
      this.set('enableTTL', value);
    },
    changedValue(data) {
      console.log({ data });
    },
    handleChange(unit) {
      console.log({ unit });
    },
  },
});
