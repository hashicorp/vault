import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

const pgpKeyFileDefault = () => ({ value: '' });

/**
 * @module ChoosePgpKeyForm
 * ChoosePgpKeyForm component is used for DR Operation Token Generation workflow. It provides
 * an interface for the user to upload or paste a PGP key for use
 *
 * @example
 * ```js
 * <ChoosePgpKeyForm @onCancel={{this.reset}} @onSubmit={{handleGenerateWithPgpKey}}>
 * ```
 * @param {function} onCancel - This function will be triggered when the modal intends to be closed
 * @param {function} onSubmit - When the PGP key is confirmed, it will call this method with the pgpKey value as the only param
 */
export default class ChoosePgpKeyForm extends Component {
  @tracked pgpKeyFile = pgpKeyFileDefault();
  @tracked selectedPgp = '';

  get pgpKey() {
    return this.pgpKeyFile.value;
  }

  @action setKey(_, keyFile) {
    this.pgpKeyFile = keyFile;
  }

  // Form submit actions:
  @action usePgpKey(evt) {
    evt.preventDefault();
    this.selectedPgp = this.pgpKey;
  }
  @action handleSubmit(evt) {
    evt.preventDefault();
    this.args.onSubmit(this.pgpKey);
  }
}
