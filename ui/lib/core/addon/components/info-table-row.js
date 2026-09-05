/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { typeOf } from '@ember/utils';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module InfoTableRow
 * `InfoTableRow` displays a label and a value in a table-row style manner. The component is responsive so
 * that the value breaks under the label on smaller viewports.
 *
 * @example
 * <InfoTableRow @value={{5}} @label="TTL" @helperText="Some description"/>
 *
 * @param {string} label=null - The display name for the value.
 * @param {string} helperText=null - Text to describe the value displayed beneath the label.
 * @param {any} value=null  - The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false. Overridden by block if exists
 * @param {boolean} [alwaysRender=false] - Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render.
 * @param {boolean} [truncateValue=false] - Indicates if the value should be truncated.
 * @param {string} [defaultShown] - Text that renders as value if alwaysRender=true. Eg. "Vault default"
 * @param {string} [tooltipText] - Text if a tooltip should display over the value.
 * @param {string} [formatDate] - A string of the desired date format that's passed to the date-format helper to render timestamps (ex. "MMM d yyyy, h:mm:ss aaa", see: https://date-fns.org/v2.30.0/docs/format)
 * @param {boolean} [formatTtl=false] - When true, value is passed to the format-duration helper, useful for TTL values
 * @param {string} [type=array] - The type of value being passed in.  This is used for when you want to trim an array.  For example, if you have an array value that can equal length 15+ this will trim to show 5 and count how many more are there
 *
 * @param {boolean} [isLink=true] - Passed through to InfoTableItemArray. Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param {string} [queryParam] - Passed through to InfoTableItemArray. If you want to specific a tab for the View All XX to display to.  Ex= role
 * @param {array} [arrayOptions] - Passed through to InfoTableItemArray for wildcard match counts.
 * @param {string} [wildcardLabel] - Passed through to InfoTableItemArray for wildcard badge text.
 * @param {string} [backend] - Passed through to InfoTableItemArray. To specify secrets backend to point link to  Ex= transformation
 * @param {string | array} [rootRoute] - Passed through to InfoTableItemArray for the View all link target.
 * @param {string | array} [itemRoute] - Passed through to InfoTableItemArray for individual array item links.
 * @param {boolean} [doNotTruncate=false] - Passed through to InfoTableItemArray to disable truncation and use ReadMore instead.
 */

export default class InfoTableRowComponent extends Component {
  @tracked
  hasLabelOverflow = false; // is calculated and set in didInsertElement

  get isVisible() {
    return this.args.alwaysRender || !this.valueIsEmpty;
  }

  get valueIsBoolean() {
    return typeOf(this.args.value) === 'boolean';
  }

  get valueIsEmpty() {
    const { value } = this.args;
    if (typeOf(value) === 'array' && value.length === 0) {
      return true;
    }
    switch (value) {
      case undefined:
        return true;
      case null:
        return true;
      case '':
        return true;
      default:
        return false;
    }
  }

  @action
  calculateLabelOverflow(el) {
    const labelDiv = el;
    const labelText = el.querySelector('.is-label');
    if (labelDiv && labelText) {
      if (labelText.offsetWidth > labelDiv.offsetWidth) {
        this.hasLabelOverflow = true;
      }
    }
  }
}
