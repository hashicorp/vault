import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { getOwner } from '@ember/application';
import { pluralize, singularize } from 'ember-inflector';
export default class PkiController extends Controller {
  @service router;
  @service secretMountPath;
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
  getPath(item) {
    const plural = pluralize(item);
    const singular = singularize(item);
    switch (item) {
      case 'pki':
        return 'overview';
      case 'details':
        return null;
      case 'index':
        return null;
      case plural:
        return `${item}.index`;
      case singular:
        return `${plural}.${item}.details`;
      default:
        break;
    }
  }

  generateBreadcrumbs(item) {
    const { currentRoute } = this.router;
    const backend = this.secretMountPath.currentPath;
    const trimmedRoute = currentRoute.name
      .replace('vault.cluster.secrets.backend.pki', 'pki')
      .replace(`${item}.details`, `${item}`)
      .split('.');
    const breadcrumbs = trimmedRoute.map((route) => {
      const label = route === 'pki' ? backend : route;
      return { label, path: this.getPath(route) };
    });
    return breadcrumbs.filter((b) => b.path !== null);
  }
}
