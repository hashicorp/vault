import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

export default class MfaLoginEnforcementHeaderComponent extends Component {
  @service store;

  constructor() {
    super(...arguments);
    if (!this.args.heading) {
      this.fetchEnforcements();
    }
  }

  @tracked enforcements = [];

  async fetchEnforcements() {
    try {
      this.enforcements = (await this.store.query('mfa-login-enforcement', {})).toArray();
    } catch (error) {
      this.enforcements = [];
    }
  }
}
