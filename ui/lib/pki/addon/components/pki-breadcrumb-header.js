import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { pluralize } from 'ember-inflector';

export default class PkiBreadcrumbHeader extends Component {
  @service router;
  @service secretMountPath;

  get breadcrumbs() {
    const { currentRoute } = this.router;
    const activeRoutes = this.parents(currentRoute); // array of active routes
    return activeRoutes
      .map((route) => {
        // resource is a singular string: 'key', 'issuer', 'role', 'certificate'
        return this.resourceToRoute(this.args.resource, route.localName);
      })
      .filter((b) => b !== null);
  }

  parents(route) {
    // gets parent of each route until base route
    return route.name === 'vault.cluster.secrets.backend.pki'
      ? [route]
      : [...this.parents(route.parent), route];
  }

  resourceToRoute(resource, route) {
    const resourcePlural = pluralize(resource);
    switch (route) {
      case 'pki':
        return { label: this.secretMountPath.currentPath, path: 'overview' };
      case 'details':
        // filter out 'details' routes since we use the resource name for breadcrumb instead
        return null;
      case 'edit':
        return { label: 'edit', path: `${resourcePlural}.${resource}.edit` };
      case resourcePlural:
        return { label: route, path: `${resourcePlural}.index` };
      case resource:
        return { label: route, path: `${resourcePlural}.${resource}.details` };
      default:
        return null;
    }
  }
}
