import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { toolsActions } from 'vault/helpers/tools-actions';

export default Route.extend({
  wizard: service(),

  beforeModel(transition) {
    const supportedActions = toolsActions();
    const { selected_action: selectedAction } = this.paramsFor(this.routeName);
    if (!selectedAction || !supportedActions.includes(selectedAction)) {
      transition.abort();
      return this.transitionTo(this.routeName, supportedActions[0]);
    }
  },

  model(params) {
    return params.selected_action;
  },

  setupController(controller, model) {
    this._super(...arguments);
    controller.set('selectedAction', model);
  },

  actions: {
    didTransition() {
      const params = this.paramsFor(this.routeName);
      if (this.wizard.currentMachine === 'tools') {
        this.wizard.transitionFeatureMachine(this.wizard.featureState, params.selected_action.toUpperCase());
      }
      this.controller.setProperties(params);
      return true;
    },
  },
});
