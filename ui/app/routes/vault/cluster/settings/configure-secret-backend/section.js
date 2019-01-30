import { set } from '@ember/object';
import Route from '@ember/routing/route';
import DS from 'ember-data';

const SECTIONS_FOR_TYPE = {
  pki: ['cert', 'urls', 'crl', 'tidy'],
};
export default Route.extend({
  fetchModel() {
    const { section_name: sectionName } = this.paramsFor(this.routeName);
    const backendModel = this.modelFor('vault.cluster.settings.configure-secret-backend');
    const modelType = `${backendModel.get('type')}-config`;
    return this.store
      .queryRecord(modelType, {
        backend: backendModel.id,
        section: sectionName,
      })
      .then(model => {
        model.set('backendType', backendModel.get('type'));
        model.set('section', sectionName);
        return model;
      });
  },

  model(params) {
    const { section_name: sectionName } = params;
    const backendModel = this.modelFor('vault.cluster.settings.configure-secret-backend');
    const sections = SECTIONS_FOR_TYPE[backendModel.get('type')];
    const hasSection = sections.includes(sectionName);
    if (!backendModel || !hasSection) {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return this.fetchModel();
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('onRefresh', () => this.fetchModel());
  },
});
