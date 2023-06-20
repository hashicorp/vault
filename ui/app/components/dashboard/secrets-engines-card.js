import Component from '@glimmer/component';

export default class DashboardSecretsEnginesCardComponent extends Component {
  constructor() {
    super(...arguments);

    if (this.args.secretsEngines.length) {
      this.secretsEngines = this.args.secretsEngines.slice(0, 5);
    }
  }
}
