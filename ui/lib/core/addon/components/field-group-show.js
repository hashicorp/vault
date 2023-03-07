import Component from '@glimmer/component';

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
export default class FieldGroupShow extends Component {
  model = null;
  showAllFields = false;
}
