import RoleEdit from './role-edit';
import Ember from 'ember';

const { get, set } = Ember;
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default RoleEdit.extend({
  useARN: false,
  init() {
    this._super(...arguments);
    const arn = get(this, 'model.arn');
    if (arn) {
      set(this, 'useARN', true);
    }
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      const modelId = this.get('model.id');
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && Ember.isBlank(modelId)) {
        return;
      }
      // clear the policy or arn before save depending on "useARN"
      if (get(this, 'useARN')) {
        set(this, 'model.policy', '');
      } else {
        set(this, 'model.arn', '');
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.get('model'), attr, val);
      }
    },
  },
});
