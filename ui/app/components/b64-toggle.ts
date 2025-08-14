/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { isBlank } from '@ember/utils';
import { encodeString, decodeString } from 'core/utils/b64';

const B64 = 'base64';
const UTF8 = 'utf-8';

type EncodingType = 'base64' | 'utf-8';

interface B64ToggleSignature {
  Args: {
    /**
     * The value that will be mutated when the encoding is toggled
     */
    value: string;
    /**
     * The encoding of `value` when the component is initialized.
     * Defaults to 'utf-8'.
     * Possible values: 'utf-8' and 'base64'
     */
    initialEncoding?: EncodingType;
    /**
     * Whether or not the toggle is associated with an input.
     * Also bound to `is-input` and `is-textarea` classes
     * Defaults to true
     */
    isInput?: boolean;
    /**
     * Callback to update the value when encoding is toggled
     */
    onUpdate?: (newValue: string) => void;
  };
  Element: HTMLButtonElement;
}

/**
 * @module B64Toggle
 * B64Toggle component provides base64 encoding/decoding functionality.
 * It toggles between UTF-8 and base64 encoding of a value.
 *
 * @example
 * <B64Toggle @value={{this.data}} @initialEncoding="utf-8" @isInput={{true}} />
 */
export default class B64Toggle extends Component<B64ToggleSignature> {
  @tracked _value = '';
  @tracked _b64Value = '';

  get valuesMatch(): boolean {
    if (isBlank(this.args.value) || isBlank(this._b64Value)) {
      return false;
    }
    return this.args.value === this._b64Value;
  }

  get currentEncoding(): EncodingType {
    return isBlank(this.args.value) ? UTF8 : this.valuesMatch ? B64 : UTF8;
  }

  constructor(owner: unknown, args: B64ToggleSignature['Args']) {
    super(owner, args);

    if (this.initialEncoding) {
      if (this.initialEncoding === B64) {
        this._b64Value = this.args.value;
      }
    }
  }

  get initialEncoding(): EncodingType {
    return this.args.initialEncoding || UTF8;
  }

  /**
   * Is the value known to be base64-encoded.
   */
  get isBase64(): boolean {
    return this.currentEncoding === B64;
  }

  @action
  handleClick(): void {
    const val = this.args.value;
    const isUTF8 = this.currentEncoding === UTF8;
    if (!val) {
      return;
    }
    const newVal = isUTF8 ? encodeString(val) : decodeString(val);

    // if the current value is UTF-8, store the base64 for the newVal which is base64
    if (isUTF8) {
      this._b64Value = newVal;
    }
    // Update internal state
    this._value = newVal;

    // Call the update callback if provided
    if (this.args.onUpdate) {
      this.args.onUpdate(newVal);
    }
  }
}
