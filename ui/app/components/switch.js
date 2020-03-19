/**
 * @module Switch
 * Switch components are used to indicate boolean values which can be toggled on or off.
 * They are a stylistic alternative to checkboxes, but use the input[type="checkbox"] under the hood.
 *
 * @example
 * ```js
 * <Switch @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} onChange - onChange is triggered on checkbox change (select, deselect)
 * @param {string} inputId - Input ID is needed to match the label to the input
 * @param {boolean} [disabled=false] - disabled makes the switch unclickable
 * @param {boolean} [isChecked=true] - isChecked is the checked status of the input, and must be passed and mutated from the parent
 * @param {boolean} [round=false] - default switch is squared off, this param makes it rounded
 * @param {string} [size='small'] - Sizing can be small, medium, or large
 * @param {string} [status='normal'] - Status can be normal or success, which makes the switch have a blue background when on
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  name: 'json',
  isChecked: false,
  disabled: false,
  size: 'small',
  round: false,
  status: 'normal',
  inputClasses: computed('size', 'round', 'status', function() {
    const roundClass = this.round ? 'is-rounded' : '';
    const sizeClass = `is-${this.size}`;
    const statusClass = `is-${this.status}`;
    return `switch ${statusClass} ${sizeClass} ${roundClass}`;
  }),
  actions: {
    handleChange(value) {
      this.onChange(value);
    },
  },
});
