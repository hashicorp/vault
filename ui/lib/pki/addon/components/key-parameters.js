import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module KeyParameters
 * KeyParameters components are used to set the default and update the key_bits pki role api param whenever the key_type changes.
 * key_bits is conditional on key_type and should be set as a default value whenever key_type changes.
 * @example
 * ```js
 * <KeyParameters @model={@model} @group={group}/>
 * ```
 * @param {class} model - The pki/role model.
 * @param {string} group - The name of the group created in the model. In this case, it's the "Key parameters" group.
 */

const KEY_BITS_OPTIONS = {
  rsa: [2048, 3072, 4096],
  ec: [256, 224, 384, 521],
  ed25519: [0],
  any: [0],
};

export default class KeyParameters extends Component {
  get keyBitOptions() {
    return KEY_BITS_OPTIONS[this.args.model.keyType];
  }

  get keyBitsDefault() {
    return Number(KEY_BITS_OPTIONS[this.args.model.keyType][0]);
  }

  @action onKeyBitsChange(selection) {
    this.args.model.set('keyBits', Number(selection.target.value));
  }

  @action onSignatureBitsOrKeyTypeChange(name, selection) {
    if (name === 'signatureBits') {
      this.args.model.set(name, Number(selection.target.value));
    }
    if (name === 'keyType') {
      this.args.model.set(name, selection.target.value);
      this.args.model.set('keyBits', this.keyBitsDefault);
    }
  }
}
