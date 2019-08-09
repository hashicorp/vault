import Component from '@ember/component';
import { computed } from '@ember/object';
import { htmlSafe } from '@ember/template';
import layout from '../templates/components/confirm';

/**
 * @module Confirm
 * `Confirm` components prevent users from performing actions they do not intend to by showing a confirmation message as an overlay. This is a contextual component that should always be rendered with a `Trigger` which triggers the message.
 *
 * @example
 * ```js
 * <div class="box">
 * <Confirm as |c|>
 *   <c.Trigger
 *     @id={{item.id}}
 *     @onTrigger={{action c.onTrigger item.id}}
 *     @triggerText="Delete"
 *     @message="This will permanently delete this secret and all its vesions."
 *     @onConfirm={{action "delete" item "secret"}}
 *     @onCancel={{action c.onCancel}}
 *     />
 * </Confirm>
 * </div>
 * ```
 */

export default Component.extend({
  layout,
  itemId: null,
  height: 0,
  style: computed('height', function() {
    return htmlSafe(`height: ${this.height}px`);
  }),
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
