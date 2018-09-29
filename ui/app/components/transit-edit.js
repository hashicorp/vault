import { inject as service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import $ from 'jquery';
import Component from '@ember/component';
import { set, get } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend(FocusOnInsertMixin, {
  router: service(),
  wizard: service(),
  mode: null,
  onDataChange() {},
  onRefresh() {},
  key: null,
  requestInFlight: or('key.isLoading', 'key.isReloading', 'key.isSaving'),

  init() {
    this._super(...arguments);
  },

  didInsertElement() {
    this._super(...arguments);
    $(document).on('keyup.keyEdit', this.onEscape.bind(this));
  },

  willDestroyElement() {
    this._super(...arguments);
    const key = this.get('key');
    if (get(key, 'isError')) {
      key.rollbackAttributes();
    }
    $(document).off('keyup.keyEdit');
  },

  transitionToRoute() {
    this.get('router').transitionTo(...arguments);
  },

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
      if (!get(key, 'isError')) {
        if (this.get('wizard.featureState') === 'secret') {
          this.get('wizard').transitionFeatureMachine('secret', 'CONTINUE');
        } else {
          if (this.get('wizard.featureState') === 'encryption') {
            this.get('wizard').transitionFeatureMachine('encryption', 'CONTINUE', 'transit');
          }
        }
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
      if (type === 'create' && isBlank(keyId)) {
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
      this.get('onRefresh')();
    },

    deleteKey() {
      this.persistKey('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },
  },
});
