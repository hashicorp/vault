/**
 * @module Icon
 * Icon components are used to...
 *
 * @example
 * ```js
 * <IcOn @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */
import Component from '@ember/component';
import layout from '../templates/components/icon';

export default Component.extend({
  tagName: '',
  layout,
});
