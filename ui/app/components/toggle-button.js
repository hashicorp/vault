import Component from '@ember/component';
import { set, get, defineProperty, computed } from '@ember/object';

/**
 * @module ToggleButton
 * `ToggleButton` components are used to expand and collapse content with a toggle.
 *
 * @example
 * ```js
 *   <ToggleButton @openLabel="Encrypt Output with PGP" @closedLabel="Encrypt Output with PGP" @toggleTarget={{this}} @toggleAttr="showOptions"/>
 *  {{#if showOptions}}
 *     <div>
 *       <p>
 *         I will be toggled!
 *       </p>
 *     </div>
 *   {{/if}}
 * ```
 *
 * @param toggleAttr=null {String} - The attribute upon which to toggle.
 * @param openLabel=Hide options {String} - The message to display when the toggle is open.
 * @param closedLabel=More options {String} - The message to display when the toggle is closed.
 */
export default Component.extend({
  tagName: 'button',
  type: 'button',
  toggleTarget: null,
  toggleAttr: null,
  classNameBindings: ['buttonClass'],
  attributeBindings: ['type'],
  buttonClass: 'has-text-info',
  classNames: ['button', 'is-transparent'],
  openLabel: 'Hide options',
  closedLabel: 'More options',
  init() {
    this._super(...arguments);
    const toggleAttr = this.get('toggleAttr');
    defineProperty(
      this,
      'isOpen',
      computed(`toggleTarget.${toggleAttr}`, () => {
        const props = this.getProperties('toggleTarget', 'toggleAttr');
        return get(props.toggleTarget, props.toggleAttr);
      })
    );
  },
  click() {
    const target = this.get('toggleTarget');
    const attr = this.get('toggleAttr');
    const current = get(target, attr);
    set(target, attr, !current);
  },
});
