import { typeOf } from '@ember/utils';
import { computed } from '@ember/object';
import { or } from '@ember/object/computed';
import Component from '@ember/component';

export default Component.extend({
  'data-test-component': 'info-table-row',
  classNames: ['info-table-row'],
  isVisible: or('alwaysRender', 'value'),

  /*
   * @param boolean
   * indicates if the component content should be always be rendered.
   * when false, the value of `value` will be used to determine if the component should render
   */
  alwaysRender: false,

  /*
   * @param string
   * the display name for the value
   *
   */
  label: null,

  /*
   *
   * the value of the data passed in - by default the content of the component will only show if there is a value
   */
  value: null,

  valueIsBoolean: computed('value', function() {
    return typeOf(this.get('value')) === 'boolean';
  }),
});
