/**
 * @module FieldGroupShow
 * FieldGroupShow components are used to...
 *
 * @example
 * ```js
 * <FieldGroupShow @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */
import Component from '@ember/component';
import layout from '../templates/components/field-group-show';

export default Component.extend({
  layout,
  model: null,
  showAllFields: false,
});
