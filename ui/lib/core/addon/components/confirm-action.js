import Component from '@ember/component';
import layout from '../templates/components/confirm-action';

/**
 * @module ConfirmAction
 * `ConfirmAction` is a button followed by a confirmation message and button used to prevent users from performing actions they do not intend to.
 *
 * @example
 * ```js
 *  <ConfirmAction
 *    @onConfirmAction={{ () => { console.log('Action!') } }}
 *    @confirmMessage="Are you sure you want to delete this config?">
 *    Delete
 *  </ConfirmAction>
 * ```
 *
 * @property {Func} [onConfirmAction=null] - The action to take upon confirming.
 * @property {String} [confirmTitle=Delete this?] - The title to display when confirming.
 * @property {String} [confirmMessage=Are you sure you want to do this?] - The message to display when confirming.
 * @property {String} [confirmButtonText=Delete] - The confirm button text.
 * @property {String} [cancelButtonText=Cancel] - The cancel button text.
 * @property {String} [disabledTitle=Can't delete this yet] - The title to display when the button is disabled.
 * @property {String} [disabledMessage=Complete the form to complete this action] - The message to display when the button is disabled.
 *
 */

export default Component.extend({
  layout,
  tagName: '',
  supportsDataTestProperties: true,
  buttonText: 'Delete',
  confirmTitle: 'Delete this?',
  confirmMessage: 'You will not be able to recover it later.',
  confirmButtonText: 'Delete',
  cancelButtonText: 'Cancel',
  disabledTitle: "Can't delete this yet",
  disabledMessage: 'Complete the form to complete this action',
  horizontalPosition: 'auto-right',
  verticalPosition: 'below',
  disabled: false,
  showConfirm: false,
  onConfirmAction: null,

  actions: {
    toggleConfirm() {
      this.toggleProperty('showConfirm');
    },

    onConfirm() {
      const confirmAction = this.get('onConfirmAction');

      if (typeof confirmAction !== 'function') {
        throw new Error('confirm-action components expects `onConfirmAction` attr to be a function');
      } else {
        confirmAction();
        this.toggleProperty('showConfirm');
      }
    },
  },
});
