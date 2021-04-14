import Component from '@ember/component';
import autosize from 'autosize';
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
export default Component.extend({
  layout,
  value: null,
  placeholder: 'value',
  maskWhileTyping: false,
  showValue: false,
  didInsertElement() {
    this._super(...arguments);
    autosize(this.element.querySelector('textarea'));
  },
  didUpdate() {
    this._super(...arguments);
    autosize.update(this.element.querySelector('textarea'));
  },
  willDestroyElement() {
    this._super(...arguments);
    autosize.destroy(this.element.querySelector('textarea'));
  },
  displayOnly: false,
  onKeyDown() {},
  onChange() {},
  actions: {
    toggleMask() {
      this.toggleProperty('showValue');
    },
    updateValue(e) {
      let value = e.target.value;
      this.set('value', value);
      this.onChange(value);
    },
  },
});
