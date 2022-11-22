import { action } from '@ember/object';
import Component from '@glimmer/component';

interface Args {
  key: {
    backend: string;
    keyName: string;
    keyId: string;
  };
}

export default class PkiKeyDetails extends Component<Args> {
  get breadcrumbs() {
    return [
      { label: this.args.key.backend || 'pki', path: 'overview' },
      { label: 'keys', path: 'keys.index' },
      { label: this.args.key.keyId },
    ];
  }

  @action deleteKey() {
    // TODO handle delete
  }
}
