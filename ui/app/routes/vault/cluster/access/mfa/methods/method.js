import Route from '@ember/routing/route';
import { hash } from 'rsvp';
export default class MfaMethodRoute extends Route {
  model({ id }) {
    return hash({
      method: this.store.findRecord('mfa-method', id).then((data) => data),
      enforcements: this.store
        .query('mfa-login-enforcement', {})
        .then((data) => {
          let filteredEnforcements = data.filter((item) => {
            let results = item.hasMany('mfa_methods').ids();
            return results.includes(id);
          });
          return filteredEnforcements;
        })
        .catch(() => {
          // Do nothing
        }),
    });
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
