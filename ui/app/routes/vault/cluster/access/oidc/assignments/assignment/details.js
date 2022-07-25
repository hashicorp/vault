import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { inject as service } from '@ember/service';

export default class OidcAssignmentDetailsRoute extends Route {
  @service store;
  model() {
    // get parent dynamic route name. This is the child and cannot access :name from params.
    let { name } = this.paramsFor('vault.cluster.access.oidc.assignments.assignment');
    return hash({
      details: this.store.findRecord('oidc/assignment', name).then((data) => data),
      entities: this.store
        .query('oidc/assignment', {})
        .then((data) => {
          let filteredEntities;
          data.filter((item) => {
            let assignments = item.id;
            if (assignments === name) {
              filteredEntities = item.hasMany('entity_ids').ids();
            }
          });
          return filteredEntities;
        })
        .catch(() => {
          // Do nothing
        }),
      groups: this.store
        .query('oidc/assignment', {})
        .then((data) => {
          let filteredGroups;
          data.filter((item) => {
            let assignments = item.id;
            if (assignments === name) {
              filteredGroups = item.hasMany('group_ids').ids();
            }
          });
          return filteredGroups;
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
