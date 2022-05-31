/**
 * @module StatusMenu
 * StatusMenu component is used as a landing page with a box horizontally and center aligned on the page. It's used as the login landing page.
 *
 *
 * @example
 * ```js
 * <StatusMenu @type='user' @onLinkClick={{action Nav.closeDrawer}}/>
 * ```
 * @param {string} [ariaLabel] - aria label for the status icon.
 * @param {string} [label] - label for the status menu.
 * @param {string} [type] - determines where the component is being used. e.g. replication, auth, etc.
 * @param {function} [onLinkClick] - function to handle click on the nested links under content.
 *
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class StatusMenu extends Component {
  @service currentCluster; // was current-cluster
  @service auth;
  @service media;

  get type() {
    return this.args.type || 'cluster';
  }

  get glyphName() {
    return this.type === 'user' ? 'user' : 'circle-dot';
  }

  @action
  handleClick(d) {
    this.args.onLinkClick;
    d.actions.close();
  }
}
