import Mixin from '@ember/object/mixin';
import { get } from '@ember/object';

export default Mixin.create({
  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  },

  setupController(controller, resolvedModel) {
    let { pageFilter } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.setProperties({
      filter: pageFilter || '',
      page: get(resolvedModel || {}, 'meta.currentPage') || 1,
    });
  },

  resetController(controller, isExiting) {
    this._super(...arguments);
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('filter', null);
    }
  },
  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
