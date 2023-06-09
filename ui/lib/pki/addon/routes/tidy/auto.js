import Route from '@ember/routing/route';

export default class PkiTidyAutoRoute extends Route {
  model() {
    const { autoTidyConfig } = this.modelFor('tidy');
    return autoTidyConfig;
  }
}
