/**
 * @module Modal
 * Modal components are used to overlay content on top of the page. Has a darkened background,
 * a title, and in order to close it you must pass an onClose function.
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
import layout from '../templates/components/modal';

export default Component.extend({
  layout,
  title: null,
  showCloseButton: false,
  onClose: () => {},
});
