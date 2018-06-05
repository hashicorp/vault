import Ember from 'ember';

const { get, computed } = Ember;

const MODEL_TYPES = {
  'ssh-sign': {
    model: 'ssh-sign',
  },
  'ssh-creds': {
    model: 'ssh-otp-credential',
    title: 'Generate SSH Credentials',
    generatedAttr: 'key',
  },
  'aws-creds': {
    model: 'iam-credential',
    title: 'Generate IAM Credentials',
    generateWithoutInput: true,
    backIsListLink: true,
  },
  'aws-sts': {
    model: 'iam-credential',
    title: 'Generate IAM Credentials with STS',
    generatedAttr: 'accessKey',
  },
  'pki-issue': {
    model: 'pki-certificate',
    title: 'Issue Certificate',
    generatedAttr: 'certificate',
  },
  'pki-sign': {
    model: 'pki-certificate-sign',
    title: 'Sign Certificate',
    generatedAttr: 'certificate',
  },
};

export default Ember.Component.extend({
  store: Ember.inject.service(),
  routing: Ember.inject.service('-routing'),
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
    this.maybeGenerate();
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
    if (this.get('action') === 'sts') {
      attrs.withSTS = true;
    }
    const newModel = this.get('store').createRecord(modelType, attrs);
    this.set('model', newModel);
  },

  /*
   *
   * @function maybeGenerate
   *
   * This method is called on `init`. If there is no input requried (as is the case for AWS IAM creds)
   * then the `create` action is triggered right away.
   *
   */
  maybeGenerate() {
    if (this.get('backend.type') !== 'aws' || this.get('action') === 'sts') {
      return;
    }
    // for normal IAM creds - there's no input, so just generate right away
    this.send('create');
  },

  actions: {
    create() {
      this.set('loading', true);
      this.model.save().finally(() => {
        this.set('loading', false);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        Ember.set(this.get('model'), attr, JSON.parse(val));
      }
    },

    newModel() {
      this.createOrReplaceModel();
    },
  },
});
