/**
 * @module Modal
 * Modal components are used to...
 *
 * @example
 * ```js
 * <Modal @title={'myTitle'} @showCloseButton={true} @onClose={() => {}}/>
 * ```
 * @param {function} onClose - onClose is the action taken when someone clicks the modal background or close button (if shown).
 * @param {string} [title] - This text shows up in the header section of the modal.
 * @param {boolean} [showCloseButton=false] - controls whether the close button in the top right corner shows.
 */

import Component from '@ember/component';

export default Component.extend({
  title: null,
  showCloseButton: false,
  onClose: () => {},
});
