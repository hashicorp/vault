export default {
  key: 'tools',
  initial: 'wrap',
  states: {
    wrap: {
      onEntry: [
        { type: 'routeTransition', params: ['vault.cluster.tools'] },
        { type: 'render', level: 'feature', component: 'wizard/tools-wrap' },
      ],
      on: {
        CONTINUE: 'wrapped',
      },
    },
    wrapped: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-wrapped' }],
      on: {
        LOOKUP: 'lookup',
      },
    },
    lookup: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-lookup' }],
      on: {
        CONTINUE: 'info',
      },
    },
    info: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-info' }],
      on: {
        REWRAP: 'rewrap',
      },
    },
    rewrap: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-rewrap' }],
      on: {
        CONTINUE: 'rewrapped',
      },
    },
    rewrapped: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-rewrapped' }],
      on: {
        UNWRAP: 'unwrap',
      },
    },
    unwrap: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-unwrap' }],
      on: {
        CONTINUE: 'unwrapped',
      },
    },
    unwrapped: {
      onEntry: [{ type: 'render', level: 'feature', component: 'wizard/tools-unwrapped' }],
      on: {
        CONTINUE: 'complete',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};
