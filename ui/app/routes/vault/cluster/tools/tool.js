import Route from '@ember/routing/route';
import { toolsActions } from 'vault/helpers/tools-actions';

export default Route.extend({
  model(params) {
    const supportedActions = toolsActions();
    if (supportedActions.includes(params.selected_action)) {
      return params.selected_action;
    }
    throw new Error('Given param is not a supported tool action');
  },

  setupController(controller, model) {
    this._super(...arguments);
    controller.set('selectedAction', model);
  },

  actions: {
    didTransition() {
      const params = this.paramsFor(this.routeName);
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      this.controller.setProperties(params);
      return true;
    },
  },
});
