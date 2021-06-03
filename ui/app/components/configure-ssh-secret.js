import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class ConfigureAwsSecretComponent extends Component {
  @action
  saveConfig(data) {
    this.args.saveConfig(data);
  }
}
