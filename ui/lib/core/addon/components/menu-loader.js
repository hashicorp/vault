/**
 * @module MenuLoader
 * MenuLoader components are used to...
 *
 * @example
 * ```js
 * <MenuLoader @loadingParam={model.updatePath.isPending} />
 * ```
 *
 * @param loadingParam {Boolean} - If the value of this param is true, the loading state will be rendered,
 * else the component will yield.
 */
import Component from '@ember/component';
import layout from '../templates/components/menu-loader';

export default Component.extend({
  tagName: 'li',
  classNames: 'action',
  layout,
  loadingParam: null,
});
