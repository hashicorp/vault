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
        CONTINUE: {
          approle: {
            cond: type => type === 'approle',
          },
          aws: {
            cond: type => type === 'aws',
          },
          azure: {
            cond: type => type === 'azure',
          },
          github: {
            cond: type => type === 'github',
          },
          gcp: {
            cond: type => type === 'gcp',
          },
          kubernetes: {
            cond: type => type === 'kubernetes',
          },
          ldap: {
            cond: type => type === 'ldap',
          },
          okta: {
            cond: type => type === 'okta',
          },
          radius: {
            cond: type => type === 'radius',
          },
          cert: {
            cond: type => type === 'cert',
          },
          userpass: {
            cond: type => type === 'userpass',
          },
        },
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
        EDIT: 'edit',
        DETAILS: 'details',
      },
    },
    edit: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/auth-edit' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
      on: {
        CONTINUE: 'details',
      },
    },
    details: {
      onEntry: [
        { type: 'render', level: 'step', component: 'wizard/auth-details' },
        { type: 'render', level: 'feature', component: 'wizard/mounts-wizard' },
      ],
    },
    approle: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/approle-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    aws: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/aws-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    azure: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/azure-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    github: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/github-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    gcp: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/gcp-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    kubernetes: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/kubernetes-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    ldap: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/ldap-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    okta: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/okta-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    radius: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/radius-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    cert: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/cert-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    userpass: {
      onEntry: [
        { type: 'render', level: 'details', component: 'wizard/userpass-method' },
        { type: 'continueFeature' },
      ],
      on: {
        CONTINUE: 'enable',
      },
    },
    complete: {
      onEntry: ['completeFeature'],
    },
  },
};
