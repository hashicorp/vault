import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module ConfirmAction
 * `ConfirmAction` is a button followed by a pop up confirmation message and button used to prevent users from performing actions they do not intend to.
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
 * @property {String} [confirmMessage=You will not be able to recover it later.] - The message to display when confirming.
 * @property {String} [confirmButtonText=Delete] - The confirm button text.
 * @property {String} [cancelButtonText=Cancel] - The cancel button text.
 *
 */

export default class ConfirmActionComponent extends Component {
  @tracked supportsDataTestProperties = true;
  @tracked horizontalPosition = 'auto-right';
  @tracked verticalPosition = 'below';
  @tracked isRunning = false;
  @tracked disabled = false;
  @tracked showConfirm = false;

  get confirmTitle() {
    return this.args.confirmTitle || 'Delete this?';
  }

  get confirmMessage() {
    return this.args.confirmMessage || 'You will not be able to recover it later.';
  }

  get confirmButtonText() {
    return this.args.confirmButtonText || 'Delete';
  }

  get cancelButtonText() {
    return this.args.cancelButtonText || 'Cancel';
  }

  @action
  closeButton(d) {
    d.actions.close();
  }

  @action
  toggleConfirm() {
    // toggle
    this.showConfirm = !this.showConfirm;
  }

  @action
  onConfirm() {
    const confirmAction = this.args.onConfirmAction;

    if (typeof confirmAction !== 'function') {
      throw new Error('confirm-action components expects `onConfirmAction` attr to be a function');
    } else {
      confirmAction();
      // toggle
      this.showConfirm = !this.showConfirm;
    }
  }
}
