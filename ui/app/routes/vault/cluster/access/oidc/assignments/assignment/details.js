import Route from '@ember/routing/route';

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
  // ARG TODO test this.
  // Reset query params to default since query param values in Ember are "sticky"
  // and the latest query param is preserved,
  resetController(controller, isExiting) {
    if (isExiting) {
      // isExiting is false if only the route's model was changing
      controller.set('listEntities', false);
      controller.set('listGroups', false);
    }
  }
}
