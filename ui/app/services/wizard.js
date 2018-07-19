import Ember from 'ember';
import { Machine } from 'xstate';

const { Service } = Ember;

const CubbyholeMachine = Machine({
  key: 'cubbyhole',
  initial: 'idle',
  states: {
    idle: {
      on: {
        DISMISS: 'dismissed',
        INTERACTION: 'active',
      },
    },
    active: {
      on: { DISMISS: 'dismissed' },
      key: 'feature',
      initial: 'create',
      states: {
        create: {
          on: { EDIT: 'edit' },
          onEntry: [{ type: 'render', component: 'cubbyHoleHelp' }],
        },
        edit: {
          on: { SAVE: 'details' },
          onEntry: [{ type: 'render', component: 'cubbyHoleEditHelp' }],
        },
        details: {
          onEntry: [{ type: 'render', component: 'cubbyHoleSuccess' }],
          on: { RESET: 'create' },
        },
      },
    },
    dismissed: {
      on: { RESET: 'idle' },
      onEntry: ['saveState'],
    },
  },
});

export default Service.extend({
  currentState: null,

  init() {
    this._super(...arguments);
    this.set('currentState', CubbyholeMachine.initialState);
  },

  transitionMachine(intendedState, extendedState) {
    let { actions, value } = CubbyholeMachine.transition(state);
    this.set('currentState', value);
  },
});
