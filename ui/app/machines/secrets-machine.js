export default {
  key: 'secrets',
  initial: 'idle',
  states: {
    idle: {
      on: {
        CONTINUE: {
          aws: {
            cond: type => type === 'aws',
          },
          cubbyhole: {
            cond: type => type === 'cubbyhole',
          },
          kv: {
            cond: type => type === 'kv',
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
            CONTINUE: 'display',
          },
        },
        display: {
          REPEAT: 'role',
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
