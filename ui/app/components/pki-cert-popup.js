/**
 * @module PkiCertPopup
 * PkiCertPopup components
 *
 * @example
 * ```js
 * <PkiCertPopup @item=/>
 * ```
 * @param {class} item - the PKI cert in question.
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
export default class PkiCertPopup extends Component {
  get item() {
    return this.args.item || null;
  }

  @action
  delete(item) {
    item.save({ adapterOptions: { method: 'revoke' } });
  }
}
