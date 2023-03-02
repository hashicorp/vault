import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import config from '../config/environment';

export default Controller.extend({
  env: config.environment,
  auth: service(),
  store: service(),
});
