import { typeOf } from '@ember/utils';
import { computed } from '@ember/object';
import { or } from '@ember/object/computed';
import Component from '@ember/component';
import layout from '../templates/components/info-table-row';

/**
 * @module InfoTableRow
 * `InfoTableRow` displays a label and a value in a table-row style manner. The component is responsive so
 * that the value breaks under the label on smaller viewports.
 *
 * @example
 * ```js
 * <InfoTableRow @value={{5}} @label="TTL" />
 * ```
 *
 * @param value=null {any} - The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false.
 * @param label=null {string} - The display name for the value.
 * @param alwaysRender=false {Boolean} - Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render.
 *
 */
export default Component.extend({
  layout,
  'data-test-component': 'info-table-row',
  classNames: ['info-table-row'],
  isVisible: or('alwaysRender', 'value'),

  alwaysRender: false,
  label: null,
  value: null,

  valueIsBoolean: computed('value', function() {
    return typeOf(this.get('value')) === 'boolean';
  }),
});
