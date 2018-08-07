import Ember from 'ember';
import { Machine } from 'xstate';

const { Service, getOwner } = Ember;

import getStorage from 'vault/lib/token-storage';
const machinesDir = 'vault/machines/';

import TutorialMachineConfig from 'vault/machines/tutorial-machine';

const TutorialMachine = Machine(TutorialMachineConfig);
let FeatureMachine = null;

export default Service.extend({
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,

  init() {
    this._super(...arguments);
    if (!this.storageHasKey('vault-tutorial-state')) {
      this.saveExtState('vault-tutorial-state', 'idle');
      this.set('currentState', TutorialMachine.initialState);
    } else {
      this.set('currentState', this.getExtState('vault-tutorial-state'));
      if (this.storageHasKey('vault-feature-list')) {
        this.set('featureList', this.getExtState('vault-feature-list'));
        if (this.storageHasKey('vault-feature-state')) {
          this.set('featureState', this.getExtState('vault-feature-state'));
        } else {
          if (FeatureMachine !== null) {
            this.set('featureState', FeatureMachine.initialState);
          } else {
            this.buildFeatureMachine();
          }
        }
      }
    }
  },

  transitionMachine(currentState, event) {
    let { actions, value } = TutorialMachine.transition(currentState, event);
    this.set('currentState', value);
    for (let action in actions) {
      this.executeAction(action, event);
    }
  },

  saveExtState(key, value) {
    this.storage().setItem(key, value);
  },

  getExtState(key) {
    return this.storage().getItem(key);
  },

  storageHasKey(key) {
    return this.storage().keys().includes(key);
  },

  executeAction(action, event) {
    switch (action) {
      case 'saveFeatures':
        this.saveFeatures(event.features);
        break;
      case 'completeFeature':
        this.completeFeature();
        break;
      default:
        break;
    }
  },

  saveFeatures(features) {
    this.set('featuresList', features);
    this.saveExtState('vault-feature-list', this.get('featuresList'));
    this.buildFeatureMachine();
  },

  buildFeatureMachine() {
    const FeatureMachineConfig = getOwner(this).lookup(
      `machine:${this.get('featuresList').objectAt(0)}-machine`
    );
    FeatureMachine = Machine(FeatureMachineConfig);
    this.set('currentMachine', this.get('featuresList').objectAt(0));
    this.set('featureState', FeatureMachine.initialState);
  },

  completeFeature() {
    let features = this.get('featuresList');
    features.pop();
    this.saveExtState('vault-feature-list', this.get('featuresList'));
    if (features.length > 0) {
      FeatureMachine = Machine(JSON.loads(`${machinesDir}/${features.objectAt(0)}-machine`));
      this.set('currentMachine', features.objectAt(0));
      this.set('featureState', FeatureMachine.initialState);
    } else {
      this.completeTutorial();
      FeatureMachine = null;
      TutorialMachine.transition(this.get('currentState'), 'DONE');
    }
  },

  storage() {
    return getStorage();
  },
});
