import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../../templates/components/confirm/trigger';

/**
 * @module Trigger
 * `Trigger` components trigger a confirmation message. They should only be used within a `Confirm` component.
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
 *     />
 * </Confirm>
 * </div>
 * ```
 *
 * @param id=null {ID} - A unique identifier used to bind a trigger to a confirmation message.
 * @param onTrigger {Func} - A function that displays the confirmation message. This must receive the `id` listed above.
 * @param onConfirm=null {Func} - The action to take when the user clicks the confirm button.
 * @param [title='Delete this?'] {String} - The header text to display in the confirmation message.
 * @param [triggerText='Delete'] {String} - The text on the trigger button.
 * @param [message='You will not be able to recover it later.'] {String} -
 * @param [confirmButtonText='Delete'] {String} - The text to display on the confirm button.
 * @param [cancelButtonText='Cancel'] {String} - The text to display on the cancel button.
 */

export default Component.extend({
  layout,
  tagName: '',
  renderedTrigger: null,
  id: null,
  onCancel() {},
  onConfirm() {},
  title: 'Delete this?',
  message: 'You will not be able to recover it later.',
  triggerText: 'Delete',
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
  showConfirm: computed('renderedTrigger', function() {
    return !!this.renderedTrigger;
  }),
});
