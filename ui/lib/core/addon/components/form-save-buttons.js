/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module FormSaveButtons
 * `FormSaveButtons` displays a button save and a cancel button at the bottom of a form.
 * To show an overall inline error message, use the :error yielded block like shown below.
 *
 * @example
 * ```js
 * <FormSaveButtons @saveButtonText="Save" @isSaving={{isSaving}} @cancelLinkParams={{array
 * "foo.route"}}>
 *   <:error>This is an error</:error>
 * </FormSaveButtons>
 * ```
 *
 * @param [saveButtonText="Save" {String}] - The text that will be rendered on the Save button.
 * @param [cancelButtonText="Cancel" {String}] - The text that will be rendered on the Cancel button.
 * @param [isSaving=false {Boolean}] - If the form is saving, this should be true. This will disable the save button and render a spinner on it;
 * @param [cancelLinkParams=[] {Array}] - An array of arguments used to construct a link to navigate back to when the Cancel button is clicked.
 * @param [onCancel=null {Function}] - If the form should call an action on cancel instead of route somewhere, the function can be passed using onCancel instead of passing an array to cancelLinkParams.
 * @param [includeBox=true {Boolean}] - By default we include padding around the form with underlines. Passing this value as false will remove that padding.
 *
 */

export default class FormSaveButtons extends Component {
  get cancelLink() {
    const { cancelLinkParams } = this.args;
    if (!Array.isArray(cancelLinkParams) || !cancelLinkParams.length) return null;
    const [route, ...models] = cancelLinkParams;
    return { route, models };
  }
}
