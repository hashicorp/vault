import Ember from 'ember';
import { Machine } from 'xstate';

const { Service, inject } = Ember;

import getStorage from 'vault/lib/token-storage';

import TutorialMachineConfig from 'vault/machines/tutorial-machine';
import SecretsMachineConfig from 'vault/machines/secrets-machine';
import PoliciesMachineConfig from 'vault/machines/policies-machine';
import ReplicationMachineConfig from 'vault/machines/replication-machine';
import ToolsMachineConfig from 'vault/machines/tools-machine';
import AuthMachineConfig from 'vault/machines/auth-machine';

const TutorialMachine = Machine(TutorialMachineConfig);
let FeatureMachine = null;
const TUTORIAL_STATE = 'vault-tutorial-state';
const FEATURE_LIST = 'vault-feature-list';
const FEATURE_STATE = 'vault-feature-state';
const COMPLETED_FEATURES = 'vault-completed-list';
const MACHINES = {
  secrets: SecretsMachineConfig,
  policies: PoliciesMachineConfig,
  replication: ReplicationMachineConfig,
  tools: ToolsMachineConfig,
  authentication: AuthMachineConfig,
};

const DEFAULTS = {
  currentState: null,
  featureList: null,
  featureState: null,
  currentMachine: null,
  tutorialComponent: null,
  featureComponent: null,
  stepComponent: null,
  detailsComponent: null,
  componentState: null,
  nextFeature: null,
  nextStep: null,
};

export default Service.extend(DEFAULTS, {
  router: inject.service(),
  showWhenUnauthenticated: false,

  init() {
    this._super(...arguments);
    this.initializeMachines();
  },

  initializeMachines() {
    if (!this.storageHasKey(TUTORIAL_STATE)) {
      let state = TutorialMachine.initialState;
      this.saveState('currentState', state.value);
      this.saveExtState(TUTORIAL_STATE, state.value);
    }
    this.saveState('currentState', this.getExtState(TUTORIAL_STATE));
    let stateNodes = TutorialMachine.getStateNodes(this.get('currentState'));
    this.executeActions(stateNodes.reduce((acc, node) => acc.concat(node.onEntry), []));
    if (this.storageHasKey(FEATURE_LIST)) {
      this.set('featureList', this.getExtState(FEATURE_LIST));
      if (this.storageHasKey(FEATURE_STATE)) {
        this.saveState('featureState', this.getExtState(FEATURE_STATE));
      } else {
        if (FeatureMachine != null) {
          this.saveState('featureState', FeatureMachine.initialState);
          this.saveExtState(FEATURE_STATE, this.get('featureState'));
        }
      }
      this.buildFeatureMachine();
    }
  },

  restartGuide() {
    let storage = this.storage();
    // empty storage
    [TUTORIAL_STATE, FEATURE_LIST, FEATURE_STATE, COMPLETED_FEATURES].forEach(key => storage.removeItem(key));
    // reset wizard state
    this.setProperties(DEFAULTS);
    // restart machines from blank state
    this.initializeMachines();
    // progress machine to 'active.select'
    this.transitionTutorialMachine('idle', 'AUTH');
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

  transitionTutorialMachine(currentState, event, extendedState) {
    if (extendedState) {
      this.set('componentState', extendedState);
    }
    let { actions, value } = TutorialMachine.transition(currentState, event);
    this.saveState('currentState', value);
    this.saveExtState(TUTORIAL_STATE, this.get('currentState'));
    this.executeActions(actions, event);
  },

  transitionFeatureMachine(currentState, event, extendedState) {
    if (!this.get('currentState').includes('active')) {
      return;
    }
    if (extendedState) {
      this.set('componentState', extendedState);
    }

    let { actions, value } = FeatureMachine.transition(currentState, event, extendedState);
    this.saveState('featureState', value);
    this.saveExtState(FEATURE_STATE, value);
    this.executeActions(actions, event);
    let next = FeatureMachine.transition(value, 'CONTINUE', extendedState);
    this.saveState('nextStep', next.value);
  },

  saveExtState(key, value) {
    this.storage().setItem(key, value);
  },

  getExtState(key) {
    return this.storage().getItem(key);
  },

  storageHasKey(key) {
    return Boolean(this.getExtState(key));
  },

  executeActions(actions, event) {
    for (let action of actions) {
      let type = action;
      if (action.type) {
        type = action.type;
      }
      switch (type) {
        case 'render':
          this.set(`${action.level}Component`, action.component);
          break;
        case 'routeTransition':
          Ember.run.next(() => {
            this.get('router').transitionTo(...action.params);
          });
          break;
        case 'saveFeatures':
          this.saveFeatures(event.features);
          break;
        case 'completeFeature':
          this.completeFeature();
          break;
        case 'handleDismissed':
          this.handleDismissed();
          break;
        case 'showTutorialWhenAuthenticated':
          this.set('showWhenUnauthenticated', false);
          break;
        case 'showTutorialAlways':
          this.set('showWhenUnauthenticated', true);
          break;
        case 'continueFeature':
          this.transitionFeatureMachine(this.get('featureState'), 'CONTINUE', this.get('componentState'));
          break;
        default:
          break;
      }
    }
  },

  handleDismissed() {
    this.storage().removeItem(FEATURE_STATE);
    this.storage().removeItem(FEATURE_LIST);
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
    this.startFeature();
    if (this.storageHasKey(FEATURE_STATE)) {
      this.saveState('featureState', this.getExtState(FEATURE_STATE));
    }
    this.saveExtState(FEATURE_STATE, this.get('featureState'));
    let nextFeature =
      this.get('featureList').length > 1 ? this.get('featureList').objectAt(1).capitalize() : 'Finish';
    this.set('nextFeature', nextFeature);
    let next = FeatureMachine.transition(this.get('featureState'), 'CONTINUE', this.get('componentState'));
    this.saveState('nextStep', next.value);
    let stateNodes = FeatureMachine.getStateNodes(this.get('featureState'));
    this.executeActions(stateNodes.reduce((acc, node) => acc.concat(node.onEntry), []));
  },

  startFeature() {
    const FeatureMachineConfig = MACHINES[this.get('featureList').objectAt(0)];
    FeatureMachine = Machine(FeatureMachineConfig);
    this.set('currentMachine', this.get('featureList').objectAt(0));
    this.saveState('featureState', FeatureMachine.initialState);
  },

  completeFeature() {
    let features = this.get('featureList');
    let done = features.shift();
    if (!this.getExtState(COMPLETED_FEATURES)) {
      let completed = [];
      completed.push(done);
      this.saveExtState(COMPLETED_FEATURES, completed);
    } else {
      this.saveExtState(COMPLETED_FEATURES, this.getExtState(COMPLETED_FEATURES).toArray().addObject(done));
    }

    this.saveExtState(FEATURE_LIST, this.get('featureList'));
    this.storage().removeItem(FEATURE_STATE);
    if (features.length > 0) {
      this.buildFeatureMachine();
    } else {
      FeatureMachine = null;
      TutorialMachine.transition(this.get('currentState'), 'DONE');
    }
  },

  storage() {
    return getStorage();
  },
});
