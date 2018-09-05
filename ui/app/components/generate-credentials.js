import Ember from 'ember';

const { get, set, computed, Component, inject } = Ember;

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
  wizard: inject.service(),
  store: inject.service(),
  routing: inject.service('-routing'),
  // set on the component
  backend: null,
  action: null,
  role: null,

  model: null,
  loading: false,
  emptyData: '{\n}',

  modelForType() {
    const type = this.get('options');
    if (type) {
      return type.model;
    }
    // if we don't have a mode for that type then redirect them back to the backend list
    const router = this.get('routing.router');
    router.transitionTo.call(router, 'vault.cluster.secrets.backend.list-root', this.get('model.backend'));
  },

  options: computed('action', 'backend.type', function() {
    const action = this.get('action') || 'creds';
    return MODEL_TYPES[`${this.get('backend.type')}-${action}`];
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
        this.get('backend.type')
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
    const roleModel = this.get('role');
    if (!modelType) {
      return;
    }
    if (model) {
      model.unloadRecord();
    }
    const attrs = {
      role: roleModel,
      id: `${get(roleModel, 'backend')}-${get(roleModel, 'name')}`,
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
              this.get('backend.type')
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
