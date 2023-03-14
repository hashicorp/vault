import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

export default class NamespaceReminder extends Component {
  @service namespace;

  get showMessage() {
    return !this.namespace.inRootNamespace;
  }

  get mode() {
    return this.args.mode || 'edit';
  }

  get noun() {
    return this.args.noun || null;
  }

  get modeVerb() {
    if (!this.mode) {
      return '';
    }
    return this.mode.endsWith('e') ? `${this.mode}d` : `${this.mode}ed`;
  }
}
