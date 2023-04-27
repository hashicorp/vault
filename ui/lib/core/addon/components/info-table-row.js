/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { typeOf } from '@ember/utils';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { convertFromSeconds, largestUnitFromSeconds } from 'core/utils/duration-utils';

/**
 * @module InfoTableRow
 * `InfoTableRow` displays a label and a value in a table-row style manner. The component is responsive so
 * that the value breaks under the label on smaller viewports.
 *
 * @example
 * ```js
 * <InfoTableRow @value={{5}} @label="TTL" @helperText="Some description"/>
 * ```
 *
 * @param label=null {string} - The display name for the value.
 * @param helperText=null {string} - Text to describe the value displayed beneath the label.
 * @param value=null {any} - The the data to be displayed - by default the content of the component will only show if there is a value. Also note that special handling is given to boolean values - they will render `Yes` for true and `No` for false. Overridden by block if exists
 * @param [alwaysRender=false] {Boolean} - Indicates if the component content should be always be rendered.  When false, the value of `value` will be used to determine if the component should render.
 * @param [defaultShown] {String} - Text that renders as value if alwaysRender=true. Eg. "Vault default"
 * @param [tooltipText] {String} - Text if a tooltip should display over the value.
 * @param [isTooltipCopyable] {Boolean} - Allows tooltip click to copy
 * @param [type=array] {string} - The type of value being passed in.  This is used for when you want to trim an array.  For example, if you have an array value that can equal length 15+ this will trim to show 5 and count how many more are there
 * @param [isLink=true] {Boolean} - Passed through to InfoTableItemArray. Indicates if the item should contain a link-to component.  Only setup for arrays, but this could be changed if needed.
 * @param [modelType=null] {string} - Passed through to InfoTableItemArray. Tells what model you want data for the allOptions to be returned from.  Used in conjunction with the the isLink.
 * @param [queryParam] {String} - Passed through to InfoTableItemArray. If you want to specific a tab for the View All XX to display to.  Ex= role
 * @param [backend] {String} - Passed through to InfoTableItemArray. To specify secrets backend to point link to  Ex= transformation
 * @param [viewAll] {String} - Passed through to InfoTableItemArray. Specify the word at the end of the link View all.
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
  get formattedTtl() {
    const { value } = this.args;
    if (Number.isInteger(value)) {
      const unit = largestUnitFromSeconds(value);
      return `${convertFromSeconds(value, unit)}${unit}`;
    }
    return value;
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
