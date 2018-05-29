import Ember from 'ember';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';

const { get, set, computed } = Ember;
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Ember.Component.extend(FocusOnInsertMixin, {
  mode: null,
  onDataChange: null,
  refresh: 'refresh',
  key: null,
  routing: Ember.inject.service('-routing'),
  requestInFlight: computed.or('key.isLoading', 'key.isReloading', 'key.isSaving'),
  willDestroyElement() {
    const key = this.get('key');
    if (get(key, 'isError')) {
      key.rollbackAttributes();
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
    get(this, 'onDataChange')(get(this, 'key.hasDirtyAttributes'));
  },

  persistKey(method, successCallback) {
    const key = get(this, 'key');
    return key[method]().then(() => {
      if (!Ember.get(key, 'isError')) {
        successCallback(key);
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

    createOrUpdateKey(type, event) {
      event.preventDefault();

      const keyId = this.get('key.id');
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && Ember.isBlank(keyId)) {
        return;
      }

      this.persistKey(
        'save',
        () => {
          this.hasDataChanges();
          this.transitionToRoute(SHOW_ROUTE, keyId);
        },
        type === 'create'
      );
    },

    handleChange() {
      this.hasDataChanges();
    },

    setValueOnKey(key, event) {
      set(get(this, 'key'), key, event.target.checked);
    },

    derivedChange(val) {
      get(this, 'key').setDerived(val);
    },

    convergentEncryptionChange(val) {
      get(this, 'key').setConvergentEncryption(val);
    },

    refresh() {
      this.sendAction('refresh');
    },

    deleteKey() {
      this.persistKey('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },
  },
});
