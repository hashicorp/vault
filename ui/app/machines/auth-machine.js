export default {
  key: 'auth',
  initial: 'idle',
  on: {
    RESET: 'idle',
    DONE: 'complete',
  },
  states: {
    idle: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.settings.auth.enable'] },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/auth-idle' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    enable: {
      onEntry: [
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
        { type: 'render', level: 'step', component: 'wizard/auth-enable' },
      ],
      on: {
        CONTINUE: 'list',
      },
    },
    list: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/auth-list' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        DETAILS: 'details',
      },
    },
    details: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/auth-details' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'complete',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};
