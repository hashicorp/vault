import Component from '@ember/component';
import layout from '../templates/components/form-save-buttons';

/**
 * @module FormSaveButtons
 * `FormSaveButtons` displays a button save and a cancel button at the bottom of a form.
 *
 * @example
 * ```js
 * <FormSaveButtons @saveButtonText="Save" @isSaving={{isSaving}} @cancelLinkParams={{array
 * "foo.route"}} />
 * ```
 *
 * @param [saveButtonText="Save" {String}] - The text that will be rendered on the Save button.
 * @param [isSaving=false {Boolean}] - If the form is saving, this should be true. This will disable the save button and render a spinner on it;
 * @param [cancelLinkParams=[] {Array}] - An array of arguments used to construct a link to navigate back to when the Cancel button is clicked.
 *
 */

export default Component.extend({
  layout,
  tagName: '',
});
