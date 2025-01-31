import Service from '@ember/service';
import { getOwner } from '@ember/owner';

export default class ConfigService extends Service {
  configFor(key) {
    return getOwner(this).resolveRegistration('config:environment')[key];
  }
  app = Object.freeze(this.configFor('APP'));
  environment = this.configFor('environment');
  host = this.configFor('host');
}
