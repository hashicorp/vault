import Ember from 'ember';
import { toolsActions } from 'vault/helpers/tools-actions';

export default Ember.Route.extend({
  beforeModel(transition) {
    const supportedActions = toolsActions();
    const { selectedAction } = this.paramsFor(this.routeName);
    if (!selectedAction || !supportedActions.includes(selectedAction)) {
      transition.abort();
      return this.transitionTo(this.routeName, supportedActions[0]);
    }
  },
  model() {},
  actions: {
    didTransition() {
      const params = this.paramsFor(this.routeName);
      this.controller.setProperties(params);
    },
  },
});
