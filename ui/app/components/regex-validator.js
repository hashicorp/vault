/**
 * @module RegexValidator
 * RegexValidator components are used to provide input forms for regex values, along with a toggle-able validation input which does not get saved to the model.
 *
 * @example
 * ```js
 * const attrExample = {
 *    name: 'valName',
 *    options: {
 *      helpText: 'Shows in tooltip',
 *      subText: 'Shows underneath label',
 *      docLink: 'Adds docs link to subText if present',
 *      defaultValue: 'hello', // Shows if no value on model
 *    }
 * }
 * <RegexValidator @onChange={action 'myAction'} @attr={attrExample} @labelString="Label String" @value="initial value" />
 * ```
 * @param {func} onChange - the action that should trigger when the main input is changed.
 * @param {string} value - the value of the main input which will be updated in onChange
 * @param {string} labelString - Form label. Anticipated from form-field
 * @param {object} attr - attribute from model. Anticipated from form-field. Example of attribute shape above
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class RegexValidator extends Component {
  @tracked testValue = '';
  @tracked showTestValue = false;

  get regexError() {
    const testString = this.testValue;
    if (!testString || !this.args.value) return false;
    const regex = new RegExp(this.args.value, 'g');
    const matchArray = testString.toString().match(regex);
    return testString !== matchArray?.join('');
  }

  @action
  updateTestValue(evt) {
    this.testValue = evt.target.value;
  }

  @action
  toggleTestValue() {
    this.showTestValue = !this.showTestValue;
  }
}
