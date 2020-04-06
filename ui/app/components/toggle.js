/**
 * @module Toggle
 * Toggle components are used to indicate boolean values which can be toggled on or off.
 * They are a stylistic alternative to checkboxes, but still use the input[type="checkbox"] under the hood.
 *
 * @example
 * ```js
 * <Toggle @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} onChange - onChange is triggered on checkbox change (select, deselect). Must manually mutate checked value
 * @param {string} name - name is passed along to the form field, as well as to generate the ID of the input & "for" value of the label
 * @param {boolean} [checked=false] - checked status of the input, and must be passed in and mutated from the parent
 * @param {boolean} [disabled=false] - disabled makes the switch unclickable
 * @param {string} [size='medium'] - Sizing can be small or medium
 * @param {string} [status='normal'] - Status can be normal or success, which makes the switch have a blue background when checked=true
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',
  checked: false,
  disabled: false,
  size: 'normal',
  status: 'normal',
  safeId: computed('name', function() {
    return `toggle-${this.name.replace(/\W/g, '')}`;
  }),
  inputClasses: computed('size', 'status', function() {
    const sizeClass = `is-${this.size}`;
    const statusClass = `is-${this.status}`;
    return `toggle ${statusClass} ${sizeClass}`;
  }),
  actions: {
    handleChange(value) {
      this.onChange(value);
    },
  },
});
