import Component from '@glimmer/component';
import { setComponentTemplate } from '@ember/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

// import Component from '@ember/component';
import layout from '../templates/components/masked-input';

/**
 * @module MaskedInput
 * `MaskedInput` components are textarea inputs where the input is hidden. They are used to enter sensitive information like passwords.
 * If the field needs to be something other than displayOnly or a input field, you should use the component TextFile nested in the FormField component.
 *
 * @example
 * <MaskedInput
 *  @value={{attr.options.defaultValue}}
 *  @placeholder="secret"
 *  @allowCopy={{true}}
 *  @onChange={{action "someAction"}}
 * />
 *
 * @param [value] {String} - The value to display in the input.
 * @param [placeholder=value] {String} - The placeholder to display before the user has entered any input.
 * @param [allowCopy=null] {bool} - Whether or not the input should render with a copy button.
 * @param [displayOnly=false] {bool} - Whether or not to display the value as a display only `pre` element or as an input.
 * @param [onChange=Function.prototype] {Function|action} - A function to call when the value of the input changes.
 * @param [isCertificate=false] {bool} - If certificate display the label and icons differently.
 */
class MaskedInput extends Component {
  // export default Component.extend({
  layout;

  placeholder = 'value';
  displayOnly = false;
  onKeyDown() {}
  onChange() {}

  @tracked
  showValue = false;
  @tracked
  value = null;

  @action
  toggleMask() {
    this.showValue = !this.showValue;
  }
  @action
  updateValue(e) {
    e.preventDefault();
    let value = e.target.value;
    console.log(value, 'value');
    this.value = value;
    this.onChange(value);
  }
}

export default setComponentTemplate(layout, MaskedInput);
