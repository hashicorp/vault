import { action } from '@ember/object';
import Component from '@glimmer/component';
import PkiConfigModel from 'vault/models/pki/config';

interface File {
  value: string;
  fileName?: string;
  enterAsText: boolean;
}

interface Args {
  config: PkiConfigModel;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

/**
 * Pki Config Import Component shows the relevant form fields for configuring a PKI mount via import
 */
export default class PkiConfigImportComponent extends Component<Args> {
  @action
  onFileUploaded(file: File) {
    if (!this.args.config) return;
    this.args.config.pemBundle = file.value;
  }
}
