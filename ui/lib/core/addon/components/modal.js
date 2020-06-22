/**
 * @module Modal
 * Modal components are used to overlay content on top of the page. Has a darkened background,
 * a title, and in order to close it you must pass an onClose function.
 *
 * Every instance of the Modal component needs to be wrapped in a conditional to show it only if isActive or isModalActive
 * is true.  This is because the focus on the didInsertElement of the modal is set when the modal is loaded
 * and if there is more than one modal on the page, the focus is incorrectly set on the first modal the dom renders.
 *
 * @example
 * ```js
 * <Modal @title={'myTitle'} @showCloseButton={true} @onClose={() => {}} @modalId="uniqueId"/>
 * ```
 * @param {function} onClose - onClose is the action taken when someone clicks the modal background or close button (if shown).
 * @param {string} [title] - This text shows up in the header section of the modal.
 * @param {boolean} [showCloseButton=false] - controls whether the close button in the top right corner shows.
 * @param {string} type=null - The header type. This comes from the message-types helper.
 * @param {string} modalId=null - unique ID passed to each modal instance.  Used to set focus on the modal for accessibility.
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import { messageTypes } from 'core/helpers/message-types';
import layout from '../templates/components/modal';

export default Component.extend({
  layout,
  title: null,
  showCloseButton: false,
  type: null,
  modalId: null,
  didInsertElement() {
    this._super(...arguments);
    // if the modal is active, set the focus to the modal card
    // allows user to use keyboard on the modal
    if (this.isActive) {
      let modalCard = this.element.querySelector(`#modal-card-id-${this.modalId}`);
      modalCard.focus();
    }
  },
  glyph: computed('type', function() {
    const modalType = this.get('type');
    if (!modalType) {
      return;
    }
    return messageTypes([this.get('type')]);
  }),
  modalClass: computed('type', function() {
    const modalType = this.get('type');
    if (!modalType) {
      return 'modal';
    }
    return 'modal ' + messageTypes([this.get('type')]).class;
  }),
  onClose: () => {},
});
