import { inject as service } from '@ember/service';
import Component from '@ember/component';

/**
 * @module GeneratedItemList
 * The `GeneratedItemList` component lists generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * ```js
 * <GeneratedItemList @model={{model}} @itemType={{itemType/>
 * ```
 *
 * @property model=null {DS.Model} - The corresponding item model that is being configured.
 * @property itemType {String} - the type of item displayed
 *
 */

export default Component.extend({
  model: null,
  itemType: null,
  router: service(),
  store: service(),
  actions: {
    refreshItemList() {
      let route = Ember.getOwner(this).lookup(`route:${this.router.currentRouteName}`);
      this.store.clearAllDatasets();
      route.refresh();
    },
  },
});
