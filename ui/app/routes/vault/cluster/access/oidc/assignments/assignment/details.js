import Route from '@ember/routing/route';
import RSVP from 'rsvp';

export default class OidcAssignmentDetailsRoute extends Route {
  // we want to trigger a full transition so that the model hood is called.
  // @service store;
  queryParams = {
    listEntities: {
      refreshModel: true,
    },
    listGroups: {
      refreshModel: true,
    },
  };

  model() {
    // because the Read Assignment by Name API only returns an array of IDs and we need to display the entity &|| group names
    // we make an API request for each entity and each group to the `identity/entity/id/the-id` endpoint and replace the entityIds array of strings with an array of objects
    // that contains both the name and id.
    // ex: {name: entity-1, id: 54ab8e70-1b05-1f83-5577-d18575ddf759}
    const model = this.modelFor('vault.cluster.access.oidc.assignments.assignment');
    //  ARG TODO stopped here. Need to map through and make a request on each one.
    // if (model.entityIds.length > 0) {
    //   model.entityIds.map((entity) => {
    //     return this.store.findRecord('');
    //   });
    // }
    return RSVP.hash({
      name: model.name,
      // enitityIds:
    });
  }

  // Reset query params to default since query param values in Ember are "sticky"
  // and the latest query param is preserved,
  resetController(controller, isExiting) {
    if (isExiting) {
      // isExiting is false if only the route's model was changing
      controller.set('listEntities', 'false');
      controller.set('listGroups', 'false');
    }
  }
}
