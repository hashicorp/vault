import Component from '@ember/component';
import autosize from 'autosize';
import layout from '../templates/components/masked-input';

/**
 * @module MaskedInput
 * `MaskedInput` components are textarea inputs where the input is hidden. They are used to enter sensitive information like passwords.
 *
 * @example
 * <MaskedInput
 *  @value={{attr.options.defaultValue}}
 *  @allowCopy={{true}}
 *  @onChange={{action "someAction"}}
 *  @onKeyUp={{action "onKeyUp"}}
 * />
 *
 * @param [value] {String} - The value to display in the input.
 * @param [allowCopy=null] {bool} - Whether or not the input should render with a copy button.
 * @param [displayOnly=false] {bool} - Whether or not to display the value as a display only `pre` element or as an input.
 * @param [onChange=Function.prototype] {Function|action} - A function to call when the value of the input changes.
 * @param [onKeyUp=Function.prototype] {Function|action} - A function to call whenever on the dom event onkeyup. Generally passed down from higher level parent.
 * @param [isCertificate=false] {bool} - If certificate display the label and icons differently.
 *
 */
export default Component.extend({
  layout,
  value: null,
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
  onKeyUp() {},
  onChange() {},
  actions: {
    toggleMask() {
      this.toggleProperty('showValue');
    },
    updateValue(e) {
      const value = e.target.value;
      this.set('value', value);
      this.onChange(value);
    },
    handleKeyUp(name, value) {
      if (this.onKeyUp) {
        this.onKeyUp(name, value);
      }
    },
  },
});
