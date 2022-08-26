import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../../templates/components/confirm/message';

/**
 * @module Message
 * `Message` components trigger and display a confirmation message. They should only be used within a `Confirm` component.
 *
 * @example
 * ```js
 * <div class="box">
 * <Confirm as |c|>
 *   <c.Message
 *     @id={{item.id}}
 *     @triggerText="Delete"
 *     @message="This will permanently delete this secret and all its versions."
 *     @onConfirm={{action "delete" item "secret"}}
 *     />
 * </Confirm>
 * </div>
 * ```
 *
 * @property id=null {ID} - A unique identifier used to bind a trigger to a confirmation message.
 * @property onConfirm=null {Func} - The action to take when the user clicks the confirm button.
 * @property [triggerText='Delete'] {String} - The text on the trigger button.
 * @property [title='Delete this?'] {String} - The header text to display in the confirmation message.
 * @property [message='You will not be able to recover it later.'] {String} - The message to display above the confirm and cancel buttons.
 * @property [confirmButtonText='Delete'] {String} - The text to display on the confirm button.
 * @property [cancelButtonText='Cancel'] {String} - The text to display on the cancel button.
 */

export default Component.extend({
  layout,
  tagName: '',
  renderedTrigger: null,
  id: null,
  onCancel() {},
  onConfirm() {},
  resetTrigger() {},
  title: 'Delete this?',
  message: 'You will not be able to recover it later.',
  triggerText: 'Delete',
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
  showConfirm: computed('id', 'renderedTrigger', function () {
    return this.renderedTrigger === this.id;
  }),
  actions: {
    onConfirm() {
      this.onConfirm();
      this.resetTrigger();
    },
  },
});
