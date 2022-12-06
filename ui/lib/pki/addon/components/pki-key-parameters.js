import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module PkiKeyParameters
 * PkiKeyParameters components are used to display a list of key bit options depending on the selected key type.
 * If the component renders in a group, other attrs may be passed in and will be rendered using the <FormField> component
 * @example
 * ```js
 * <PkiKeyParameters @model={{@model}} @fields={{fields}}/>
 * ```
 * @param {class} model - The pki/role model.
 * @param {string} fields - The name of the fields created in the model. In this case, it's the "Key parameters" fields.
 */

const KEY_BITS_OPTIONS = {
  rsa: [2048, 3072, 4096],
  ec: [256, 224, 384, 521],
  ed25519: [0],
  any: [0],
};

export default class PkiKeyParameters extends Component {
  // TODO clarify types here
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
      this.args.model.set(name, Number(selection));
    }
    if (name === 'keyType') {
      this.args.model.set('keyBits', this.keyBitsDefault);
    }
  }
}
