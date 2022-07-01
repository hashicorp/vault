/**
 * @module SplashPage
 * SplashPage component is used as a landing page with a box horizontally and center aligned on the page. It's used as the login landing page.
 *
 *
 * @example
 * ```js
 * <SplashPage >
 * content here
 * </SplashPage
 * ```
 * @param {boolean} [hasAltContent] - boolean to bypass the UiWizard.
 * @param {boolean} [showTruncatedNavBar = true] - boolean to hide or show the navBar. By default this is true.
 *
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

export default class SplashPage extends Component {
  @service version;
  @service auth;
  @service store;

  get showTruncatedNavBar() {
    // default is true unless showTruncatedNavBar is defined as false
    return this.args.showTruncatedNavBar === false ? false : true;
  }

  get activeCluster() {
    return this.store.peekRecord('cluster', this.auth.activeCluster);
  }
}
