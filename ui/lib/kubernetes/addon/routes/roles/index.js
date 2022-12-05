import FetchConfigRoute from '../fetch-config';
import { hash } from 'rsvp';

export default class KubernetesRolesRoute extends FetchConfigRoute {
  model(params, transition) {
    // filter roles based on pageFilter value
    const { pageFilter } = transition.to.queryParams;
    const roles = this.store
      .query('kubernetes/role', { backend: this.secretMountPath.get() })
      .then((models) =>
        pageFilter
          ? models.filter((model) => model.name.toLowerCase().includes(pageFilter.toLowerCase()))
          : models
      );
    return hash({
      backend: this.modelFor('application'),
      config: this.configModel,
      roles,
    });
  }
}
