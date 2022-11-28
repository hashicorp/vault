import { action } from '@ember/object';
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
interface Args {
  key: {
    backend: string;
    keyName: string;
    keyId: string;
  };
}

export default class PkiKeyDetails extends Component<Args> {
  @service declare secretMountPath: { currentPath: string };

  get breadcrumbs() {
    return [
      { label: 'secrets', path: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath || 'pki', path: 'overview' },
      { label: 'keys', path: 'keys.index' },
      { label: this.args.key.keyId },
    ];
  }

  @action deleteKey() {
    // TODO handle delete
  }
}
