export default {
  key: 'secrets',
  initial: 'idle',
  states: {
    idle: {
      on: {
        CONTINUE: {
          aws: {
            cond: (extState, event) => event.selected === 'aws',
          },
          cubbyhole: {
            cond: (extState, event) => event.selected === 'ch',
          },
          kv: {
            cond: (extState, event) => event.selected === 'kv',
          },
        },
      },
    },
    aws: {
      on: {
        RESET: 'idle',
        DONE: 'complete',
        PAUSE: 'paused',
      },
      key: 'aws',
      initial: 'credentials',
      states: {
        credentials: {
          on: {
            CONTINUE: 'role',
          },
        },
        role: {
          on: {
            REPEAT: 'role',
          },
        },
      },
    },
    cubbyhole: {
      on: {
        RESET: 'idle',
        DONE: 'complete',
        PAUSE: 'paused',
      },
      key: 'ch',
      initial: 'role',
      states: {
        role: {
          on: {
            REPEAT: 'role',
          },
        },
      },
    },
    kv: {
      on: {
        RESET: 'idle',
        DONE: 'complete',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
      on: { RESET: 'idle' },
    },
    paused: {},
  },
};
