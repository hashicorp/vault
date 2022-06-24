import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class OidcConfigureRoute extends Route {
  @service router;

  header = null;

  beforeModel(transition) {
    // set correct header state based on child route
    // when no clients have been created display button as call to action to create
    // list views share the same header with tabs as resource links
    // the remaining routes are responsible for their own header
    const routeName = transition.to.name;
    if (routeName.includes('oidc.index')) {
      this.header = 'cta';
    } else {
      const isList = ['clients', 'assignments', 'keys', 'scopes', 'providers'].find((resource) => {
        return routeName.includes(`${resource}.index`);
      });
      this.header = isList ? 'list' : null;
    }
  }

  setupController(controller) {
    controller.header = this.header;
  }
}
