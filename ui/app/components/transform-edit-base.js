import { inject as service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import Component from '@ember/component';
import { set, get } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export const addToList = (list, itemToAdd) => {
  if (!list || !Array.isArray(list)) return list;
  list.push(itemToAdd);
  return list.uniq();
};

export const removeFromList = (list, itemToRemove) => {
  if (!list) return list;
  const index = list.indexOf(itemToRemove);
  if (index < 0) return list;
  const newList = list.removeAt(index, 1);
  return newList.uniq();
};

export default Component.extend(FocusOnInsertMixin, {
  store: service(),
  flashMessages: service(),
  router: service(),

  mode: null,
  onDataChange() {},
  onRefresh() {},
  model: null,
  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  init() {
    this._super(...arguments);
    this.set('backendType', 'transform');
  },

  willDestroyElement() {
    this._super(...arguments);
    if (this.model && this.model.isError) {
      this.model.rollbackAttributes();
    }
  },

  transitionToRoute() {
    this.get('router').transitionTo(...arguments);
  },

  modelPrefixFromType(modelType) {
    let modelPrefix = '';
    if (modelType && modelType.startsWith('transform/')) {
      modelPrefix = `${modelType.replace('transform/', '')}/`;
    }
    return modelPrefix;
  },
  persist(method, successCallback) {
    const model = get(this, 'model');
    return model[method]()
      .then(() => {
        successCallback(model);
      })
      .catch(e => {
        // console.log('e', e);
        model.set('displayErrors', e.errors);
      });
  },

  applyChanges(type, callback = () => {}) {
    const modelId = this.get('model.id') || this.get('model.name'); // transform comes in as model.name
    const modelPrefix = this.modelPrefixFromType(this.get('model.constructor.modelName'));
    // prevent from submitting if there's no key
    // maybe do something fancier later
    if (type === 'create' && isBlank(modelId)) {
      return;
    }

    this.persist('save', () => {
      callback();
      this.transitionToRoute(SHOW_ROUTE, `${modelPrefix}${modelId}`);
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges(type);
    },

    setValue(key, event) {
      set(get(this, 'model'), key, event.target.checked);
    },

    refresh() {
      this.get('onRefresh')();
    },

    delete() {
      this.persist('destroyRecord', () => {
        this.onDataChange();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },
  },
});
