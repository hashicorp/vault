import Component from '@ember/component';
import layout from '../templates/components/popup-menu';

/**
 * @module PopupMenu
 * `PopupMenu` displays a button that when pressed will toggle open a menu that is yielded in the content
 * block.
 *
 * @example
 * ```js
 * <PopupMenu><nav class="menu"> <ul class="menu-list"> <li class="action"> <button type="button">A menu!</button> </li> </ul> </nav></PopupMenu>
 * ```
 *
 * @param contentClass=''{String} A class that will be applied to the yielded content of the popup.
 */

export default Component.extend({
  layout,
  contentClass: '',
  tagName: 'span',
});
