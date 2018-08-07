import Ember from 'ember';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';

const { get, set, computed } = Ember;
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Ember.Component.extend(FocusOnInsertMixin, {
  mode: null,
  emptyData: '{\n}',
  onDataChange: () => {},
  refresh: 'refresh',
  model: null,
  routing: Ember.inject.service('-routing'),
  requestInFlight: computed.or('model.isLoading', 'model.isReloading', 'model.isSaving'),
  willDestroyElement() {
    const model = this.get('model');
    if (get(model, 'isError')) {
      model.rollbackAttributes();
    }
  },

  transitionToRoute() {
    const router = this.get('routing.router');
    router.transitionTo.apply(router, arguments);
  },

  bindKeys: Ember.on('didInsertElement', function() {
    Ember.$(document).on('keyup.keyEdit', this.onEscape.bind(this));
  }),

  unbindKeys: Ember.on('willDestroyElement', function() {
    Ember.$(document).off('keyup.keyEdit');
  }),

  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.get('mode') !== 'show') {
      return;
    }
    this.transitionToRoute(LIST_ROOT_ROUTE);
  },

  hasDataChanges() {
    get(this, 'onDataChange')(get(this, 'model.hasDirtyAttributes'));
  },

  persist(method, successCallback) {
    const model = get(this, 'model');
    return model[method]().then(() => {
      if (!Ember.get(model, 'isError')) {
        successCallback(model);
      }
    });
  },

  actions: {
    handleKeyDown(_, e) {
      e.stopPropagation();
      if (!(e.keyCode === keys.ENTER && e.metaKey)) {
        return;
      }
      let $form = this.$('form');
      if ($form.length) {
        $form.submit();
      }
      $form = null;
    },

    createOrUpdate(type, event) {
      event.preventDefault();

      const modelId = this.get('model.id');
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && Ember.isBlank(modelId)) {
        return;
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    },

    handleChange() {
      this.hasDataChanges();
    },

    setValue(key, event) {
      set(get(this, 'model'), key, event.target.checked);
    },

    refresh() {
      this.sendAction('refresh');
    },

    delete() {
      this.persist('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.get('model'), attr, JSON.parse(val));
      }
    },
  },
});
