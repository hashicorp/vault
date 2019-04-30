import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';

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
 * @property {Func} onConfirmAction=null - The action to take upon confirming.
 * @property {String} [confirmMessage=Are you sure you want to do this?] - The message to display upon confirming.
 * @property {String} [confirmButtonText=Delete] - The confirm button text.
 * @property {String} [cancelButtonText=Cancel] - The cancel button text.
 * @property {String} [disabledMessage=Complete the form to complete this action] - The message to display when the button is disabled.
 *
 */

export default Component.extend({
  tagName: 'span',
  classNames: ['confirm-action'],
  layout: hbs`
    {{#if showConfirm ~}}
      <span class={{containerClasses}}>
        <span class={{concat 'confirm-action-text ' messageClasses}}>{{if disabled disabledMessage confirmMessage}}</span>
        <button {{action 'onConfirm'}} disabled={{disabled}} class={{confirmButtonClasses}} type="button" data-test-confirm-button=true>{{confirmButtonText}}</button>
        <button {{action 'toggleConfirm'}} type="button" class={{cancelButtonClasses}} data-test-confirm-cancel-button=true>{{cancelButtonText}}</button>
      </span>
    {{else}}
      <button
        class={{buttonClasses}}
        type="button"
        disabled={{disabled}}
        data-test-confirm-action-trigger=true
        {{action 'toggleConfirm'}}
      >
        {{yield}}
      </button>
    {{~/if}}
  `,

  disabled: false,
  disabledMessage: 'Complete the form to complete this action',
  showConfirm: false,
  messageClasses: 'is-size-8 has-text-grey',
  confirmButtonClasses: 'is-danger is-outlined button',
  containerClasses: '',
  buttonClasses: 'button',
  buttonText: 'Delete',
  confirmMessage: 'Are you sure you want to do this?',
  confirmButtonText: 'Delete',
  cancelButtonClasses: 'button',
  cancelButtonText: 'Cancel',
  // the action to take when we confirm
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
