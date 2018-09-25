import { inject as service } from '@ember/service';
import { computed, set } from '@ember/object';
import Component from '@ember/component';

const MODEL_TYPES = {
  'ssh-sign': {
    model: 'ssh-sign',
  },
  'ssh-creds': {
    model: 'ssh-otp-credential',
    title: 'Generate SSH Credentials',
  },
  'aws-creds': {
    model: 'aws-credential',
    title: 'Generate AWS Credentials',
    backIsListLink: true,
  },
  'pki-issue': {
    model: 'pki-certificate',
    title: 'Issue Certificate',
  },
  'pki-sign': {
    model: 'pki-certificate-sign',
    title: 'Sign Certificate',
  },
};

export default Component.extend({
  wizard: service(),
  store: service(),
  router: service(),
  // set on the component
  backendType: null,
  backendPath: null,
  roleName: null,
  action: null,

  model: null,
  loading: false,
  emptyData: '{\n}',

  modelForType() {
    const type = this.get('options');
    if (type) {
      return type.model;
    }
    // if we don't have a mode for that type then redirect them back to the backend list
    this.get('router').transitionTo('vault.cluster.secrets.backend.list-root', this.get('backendPath'));
  },

  options: computed('action', 'backendType', function() {
    const action = this.get('action') || 'creds';
    return MODEL_TYPES[`${this.get('backendType')}-${action}`];
  }),

  init() {
    this._super(...arguments);
    this.createOrReplaceModel();
  },

  didReceiveAttrs() {
    if (this.get('wizard.featureState') === 'displayRole') {
      this.get('wizard').transitionFeatureMachine(
        this.get('wizard.featureState'),
        'CONTINUE',
        this.get('backendType')
      );
    }
  },

  willDestroy() {
    this.get('model').unloadRecord();
    this._super(...arguments);
  },

  createOrReplaceModel() {
    const modelType = this.modelForType();
    const model = this.get('model');
    const roleName = this.get('roleName');
    const backendPath = this.get('backendPath');
    if (!modelType) {
      return;
    }
    if (model) {
      model.unloadRecord();
    }
    const attrs = {
      role: {
        backend: backendPath,
        name: roleName,
      },
      id: `${backendPath}-${roleName}`,
    };
    const newModel = this.get('store').createRecord(modelType, attrs);
    this.set('model', newModel);
  },

  actions: {
    create() {
      let model = this.get('model');
      this.set('loading', true);
      this.model
        .save()
        .catch(() => {
          if (this.get('wizard.featureState') === 'credentials') {
            this.get('wizard').transitionFeatureMachine(
              this.get('wizard.featureState'),
              'ERROR',
              this.get('backendType')
            );
          }
        })
        .finally(() => {
          model.set('hasGenerated', true);
          this.set('loading', false);
        });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.get('model'), attr, JSON.parse(val));
      }
    },

    newModel() {
      this.createOrReplaceModel();
    },
  },
});
