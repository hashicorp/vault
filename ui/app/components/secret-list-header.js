import Component from '@glimmer/component';

export default class SecretListHeader extends Component {
  // api
  isCertTab = false;
  isConfigure = false;
  baseKey = null;
  backendCrumb = null;
  model = null;
  options = null;

  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
}
