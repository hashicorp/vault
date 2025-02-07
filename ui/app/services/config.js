import Service from '@ember/service';
import { getOwner } from '@ember/owner';

export default class ConfigService extends Service {
  configFor(key) {
    return getOwner(this).resolveRegistration('config:environment')[key];
  }

  constructor() {
    super();

    this.app = Object.freeze(this.configFor('APP'));
    this.environment = this.configFor('environment');
    this.host = this.configFor('host');
  }
}
