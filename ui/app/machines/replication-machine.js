export default {
  key: 'replication',
  initial: 'setup',
  states: {
    setup: {
      on: {
        ENABLEREPLICATION: 'details',
      },
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.replication'] },
        { type: 'render', level: 'feature', component: 'wizard/replication-setup' },
      ],
    },
    details: {
      on: {
        CONTINUE: 'complete',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/replication-details' }],
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};
