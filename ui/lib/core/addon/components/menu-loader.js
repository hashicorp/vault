/**
 * @module MenuLoader
 * MenuLoader components are used to show a loading state when fetching data is triggered by opening a
 * popup menu.
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
