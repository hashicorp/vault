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
 * Pki Config Import Component creates a PKI Config record on mount, and cleans it up if dirty on unmount
 */
export default class PkiConfigImportComponent extends Component<Args> {
  @action
  onFileUploaded(file: File) {
    if (!this.args.config) return;
    this.args.config.pemBundle = file.value;
  }
}
