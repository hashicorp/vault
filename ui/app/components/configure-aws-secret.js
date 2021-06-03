import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class ConfigureAwsSecretComponent extends Component {
  @action
  saveRootCreds(data) {
    this.args.saveAWSRoot(data);
  }

  @action
  saveLease(data) {
    this.args.saveAWSLease(data);
  }
}
