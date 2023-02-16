export default {
  key: 'policies',
  initial: 'idle',
  states: {
    idle: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.policies.index', 'acl'] },
        { type: 'render', level: 'feature', component: 'wizard/policies-intro' },
      ],
      on: {
        CONTINUE: 'create',
      },
    },
    create: {
      on: {
        CONTINUE: 'details',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/policies-create' }],
    },
    details: {
      on: {
        CONTINUE: 'delete',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/policies-details' }],
    },
    delete: {
      on: {
        CONTINUE: 'others',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/policies-delete' }],
    },
    others: {
      on: {
        CONTINUE: 'complete',
      },
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/policies-others' }],
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};
