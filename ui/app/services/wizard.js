import Ember from 'ember';
import { Machine } from 'xstate';

const { Service } = Ember;

import getStorage from 'vault/lib/token-storage';

import TutorialMachineConfig from 'vault/machines/tutorial-machine';
import SecretsMachineConfig from 'vault/machines/secrets-machine';

const TutorialMachine = Machine(TutorialMachineConfig);
let FeatureMachine = null;
const TUTORIAL_STATE = 'vault-tutorial-state';
const FEATURE_LIST = 'vault-feature-list';
const FEATURE_STATE = 'vault-feature-state';
const MACHINES = { secrets: SecretsMachineConfig };

export default Service.extend({
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,

  init() {
    this._super(...arguments);
    if (!this.storageHasKey(TUTORIAL_STATE)) {
      this.saveState('currentState', TutorialMachine.initialState);
      this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    } else {
      this.saveState('currentState', this.getExtState(TUTORIAL_STATE));
      if (this.storageHasKey(FEATURE_LIST)) {
        this.set('featureList', this.getExtState(FEATURE_LIST));
        if (this.storageHasKey(FEATURE_STATE)) {
          this.saveState('featureState', this.getExtState(FEATURE_STATE));
        } else {
          if (FeatureMachine !== null) {
            this.saveState('featureState', FeatureMachine.initialState);
            this.saveExtState(FEATURE_STATE, this.get('featureState'));
          } else {
            this.buildFeatureMachine();
          }
        }
      }
    }
  },

  saveState(stateType, state) {
    if (state.value) {
      state = state.value;
    }

    let stateKey = '';
    while (Ember.typeOf(state) === 'object') {
      let newState = Object.keys(state);
      stateKey += newState + '.';
      state = state[newState];
    }
    stateKey += state;
    this.set(stateType, stateKey);
  },

  transitionTutorialMachine(currentState, event) {
    let { actions, value } = TutorialMachine.transition(currentState, event);
    this.saveState('currentState', value);
    this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    for (let action in actions) {
      this.executeAction(action, event);
    }
  },

  transitionFeatureMachine(currentState, event) {
    let { actions, value } = FeatureMachine.transition(currentState, event);
    this.saveState('featureState', value);
    this.saveExtState(FEATURE_STATE, value);
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
    this.set('featureList', features);
    this.saveExtState(FEATURE_LIST, this.get('featureList'));
    this.buildFeatureMachine();
  },

  buildFeatureMachine() {
    if (this.get('featureList') === null) {
      return;
    }
    const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
    FeatureMachine = Machine(FeatureMachineConfig);
    this.set('currentMachine', this.get('featureList').objectAt(0));
    this.saveState('featureState', FeatureMachine.initialState);
    this.saveExtState(FEATURE_STATE, this.get('featureState'));
  },

  completeFeature() {
    let features = this.get('featureList');
    features.pop();
    this.saveExtState(FEATURE_LIST, this.get('featureList'));
    if (features.length > 0) {
      const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
      FeatureMachine = Machine(FeatureMachineConfig);
      this.set('currentMachine', features.objectAt(0));
      this.saveState('featureState', FeatureMachine.initialState);
      this.saveExtState(FEATURE_STATE, this.get('featureState'));
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
