import Component from '@glimmer/component';

export default class SecretVersionMenu extends Component {
  onRefresh() {}
  get useDefaultTrigger() {
    return this.args.useDefaultTrigger || false;
  }
}
