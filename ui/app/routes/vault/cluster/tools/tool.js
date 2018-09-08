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
  model() {},
  setupController(controller) {
    let { selected_action: selectedAction } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('selectedAction', selectedAction);
  },
  actions: {
    didTransition() {
      const params = this.paramsFor(this.routeName);
      if (this.get('wizard.currentMachine') === 'tools') {
        this.get('wizard').transitionFeatureMachine(
          this.get('wizard.featureState'),
          params.selectedAction.toUpperCase()
        );
      }
      this.controller.setProperties(params);
      return true;
    },
  },
});
