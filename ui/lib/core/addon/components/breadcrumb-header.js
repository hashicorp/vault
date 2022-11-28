import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 * @module BreadcrumbHeader
 * BreadcrumbHeader components are used to display a header with breadcrumb links and an optional title below
 *
 * @example
 * ```js
 * <BreadcrumbHeader @breadcrumbs={{breadcrumbs}} @pageTitle="View key" @icon="certificate" />
 * ```
 * @param {array} breadcrumbs - array of objects with a label and path to display as breadcrumbs
 * @param {string} [pageTitle] - optional title to display in the header
 * @param {string} [icon] - icon name that displays to the left of the page title
 */

export default class BreadcrumbHeader extends Component {
  @service router;

  get localName() {
    return this.router.currentRoute.localName;
  }
}
