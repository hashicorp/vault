import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/confirm';

/**
 * @module Confirm
 * `Confirm` components prevent users from performing actions they do not intend to. This component should always be rendered with a Trigger (usually a link or button) and Message.
 *
 * @example
 * ```js
 * <Confirm as |c|>
 * <c.Trigger>
 *   <button
 *     type="button"
 *     class="link is-destroy"
 *     onclick={{action c.onTrigger item.id}}>
 *     Delete
 *   </button>
 * </c.Trigger>
 * <c.Message
 *   @id={{item.id}}
 *   @onCancel={{action c.onCancel}}
 *   @onConfirm={{action "delete" item "secret"}}
 *   @message="This will permanently delete this secret and all its vesions.">
 * </c.Message>
 * </Confirm>
 * ```
 */

export default Component.extend({
  layout,
  itemId: null,
  height: 0,
  wormholeReference: null,
  wormholeId: computed(function() {
    return `confirm-${this.elementId}`;
  }),
  didInsertElement() {
    this.set('wormholeReference', this.element.querySelector(`#${this.wormholeId}`));
  },
  didRender() {
    this.updateHeight();
  },
  updateHeight: function() {
    let height;
    height = this.openTrigger
      ? this.element.querySelector('.confirm-overlay').clientHeight
      : this.element.querySelector('.confirm').clientHeight;
    this.set('height', height);
  },
  actions: {
    onTrigger: function(itemId) {
      this.set('openTrigger', itemId);
      this.updateHeight();
    },
    onCancel: function() {
      this.set('openTrigger', '');
      this.updateHeight();
    },
  },
});
