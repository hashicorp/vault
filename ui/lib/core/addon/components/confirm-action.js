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
    {{#basic-dropdown class="popup-menu" horizontalPosition=horizontalPosition verticalPosition=verticalPosition onOpen=(action "toggleConfirm") onClose=(action "toggleConfirm") as |d|}}
      {{#d.trigger
        tagName="button"
        class=(concat buttonClasses " popup-menu-trigger" (if d.isOpen " is-active"))
        disabled=disabled
        data-test-confirm-action-trigger="true"
      }}
        {{yield}}
        {{#if (eq buttonClasses 'toolbar-link') ~}}
          <Icon @glyph="chevron-{{if showConfirm 'up' 'down'}}" />
        {{~/if}}
      {{/d.trigger}}
      {{#d.content class=(concat "popup-menu-content")}}
        <div class="box confirm-action-message">
          <div class="message is-highlight">
            <div class="message-title">
              <Icon @glyph="alert-triangle" />
              {{if disabled disabledTitle confirmTitle}}
            </div>
            <p>
              {{if disabled disabledMessage confirmMessage}}
            </p>
          </div>
          <div class="confirm-action-options">
            <button
              type="button"
              disabled={{disabled}}
              class="link is-destroy"
              data-test-confirm-button="true"
              {{action 'onConfirm'}}
            >
              {{confirmButtonText}}
            </button>
            <button
              type="button"
              class="link"
              data-test-confirm-cancel-button="true"
              {{action d.actions.close}}
            >
              {{cancelButtonText}}
            </button>
          </div>
        </div>
      {{/d.content}}
    {{/basic-dropdown}}
  `,
  buttonText: 'Delete',
  confirmTitle: 'Delete this?',
  confirmMessage: 'This data will be permenantly deleted.',
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
