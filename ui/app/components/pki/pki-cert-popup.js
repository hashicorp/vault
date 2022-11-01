import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module PkiCertPopup
 * PkiCertPopup component is the hotdog menu button that allows you to see details or revoke a certificate.
 *
 * @example
 * ```js
 * <PkiCertPopup @item={{@item}}/>
 * ```
 * @param {class} item - the PKI cert in question.
 */

export default class PkiCertPopup extends Component {
  get item() {
    return this.args.item || null;
  }

  @action
  delete(item) {
    item.save({ adapterOptions: { method: 'revoke' } });
  }
}
