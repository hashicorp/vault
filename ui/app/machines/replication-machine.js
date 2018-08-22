export default {
  key: 'replication',
  initial: 'setup',
  states: {
    setup: {
      on: {
        CONTINUE: 'details',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/replication-setup' }],
    },
    details: {
      on: {
        CONTINUE: 'complete',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/replication-details' }],
    },
    complete: {
      onEntry: ['completeFeature'],
      on: { RESET: 'idle' },
    },
  },
};
