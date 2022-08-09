import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module FormFieldLabel
 * FormFieldLabel components add labels and descriptions to inputs
 *
 * @example
 * ```js
 * <FormFieldLabel for="input-name" @label={{this.label}} @helpText={{this.helpText}} @subText={{this.subText}} @dockLink={{this.docLink}} />
 * ```
 * @param {string} [label] - label text -- component attributes are spread on label element
 * @param {string} [helpText] - adds a tooltip
 * @param {string} [subText] - de-emphasized text rendered below the label
 * @param {string} [docLink] - url to documentation rendered after the subText
 */

export default class FormFieldLabel extends Component {
  @action
  preventSubmit(e) {
    e.preventDefault();
  }
}
