/**
 * @module Switch
 * Switch components are used to...
 *
 * @example
 * ```js
 * <Switch @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} onChange - requiredParam is...
 * @param {string} id - id is for the input
 * @param {string} [name='json'] - param1 is...
 * @param {boolean} [disabled=false] - param1 is...
 * @param {boolean} [isChecked=true] - param1 is...
 * @param {string} [classNames] - param1 is...
 * @param {boolean} [round=false] - default switch is squared off
 * @param {string} [size='small'] - Sizing which is
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  id: null,
  name: 'json',
  classNames: ['switch', 'is-small'],
  classNameBindings: ['round:is-round'],
  isChecked: true,
  onChange: () => {},
  disabled: false,
  size: 'small',
  isSize: computed('size', function() {
    return `is-${this.size}`;
  }),
  round: false,
  barfoo() {},
  actions: {
    handleChange(key, value) {
      console.log('switch clicked');
      console.log(key, value);
      this.onChange(key, value);
    },
    foobar(value) {
      console.log(`clicked! ${value}`);
      // this.onChange();
    },
  },
});
