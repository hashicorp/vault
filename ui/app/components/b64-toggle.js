import { isBlank } from '@ember/utils';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { encodeString, decodeString } from 'vault/utils/b64';

/**
 * @module B64Toggle2
 * B64Toggle2 components are used to...
 *
 * @example
 * ```js
 * <B64Toggle2 @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 * @param {function} onChange - Function to handle the changing of the value passed in.
 */
export const B64 = 'base64';
export const UTF8 = 'utf-8';

export default class B64Toggle2 extends Component {
  @tracked _value; // internal tracker of encoded value
  @tracked lastEncoding = ''; // only becomes value once the action has been hit

  constructor() {
    super(...arguments);
  }

  get currentEncoding() {
    // can only be two values: B64 or UTF8
    // only set to B64 if lastEncoding was B64 and the valuesMatch
    if (this.lastEncoding === B64 && this.valuesMatch) {
      return B64;
    }
    // otherwise in all cases it is UTF8 encoding
    return UTF8;
  }

  get isBase64() {
    return this.currentEncoding === B64 ? true : false;
  }

  get isInput() {
    // isInput is either false or undefined so need to use equals instead of simple or statement
    if (this.args.isInput === false) {
      return false;
    } else {
      return true;
    }
  }

  get valuesMatch() {
    const anyBlank = isBlank(this.value);
    return !anyBlank && this.args.value === this._value;
  }

  get value() {
    return this._value || this.args.value;
  }

  @action
  encodeDecodeValue() {
    const isUTF8 = this.currentEncoding === UTF8;
    if (!this.args.value) {
      return;
    }
    let newVal = isUTF8 ? encodeString(this.args.value) : decodeString(this.args.value);
    let toggleEncodingType = isUTF8 ? B64 : UTF8;

    this._value = newVal; // the encoded value
    this.lastEncoding = toggleEncodingType;

    if (this.args.onChange) {
      // function passed to component to change the changed value
      this.args.onChange(newVal);
    }
  }
}
