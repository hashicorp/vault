import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { next } from '@ember/runloop';

/**
 * @module StatusMenu
 * StatusMenu component is the drop down menu on the main navigation.
 *
 * @example
 * ```js
 * <StatusMenu @label='user' @onLinkClick={{action Nav.closeDrawer}}/>
 * ```
 * @param {string} [ariaLabel] - aria label for the status icon.
 * @param {string} [label] - label for the status menu.
 * @param {string} [type] - determines where the component is being used. e.g. replication, auth, etc.
 * @param {function} [onLinkClick] - function to handle click on the nested links under content.
 *
 */

export default class StatusMenu extends Component {
  @service currentCluster;
  @service auth;
  @service media;
  @service router;

  get type() {
    return this.args.type || 'cluster';
  }

  get glyphName() {
    return this.type === 'user' ? 'user' : 'circle-dot';
  }

  @action
  onLinkClick(dropdown) {
    if (dropdown) {
      // strange issue where closing dropdown triggers full transition which redirects to auth screen in production builds
      // closing dropdown in next tick of run loop fixes it
      next(() => dropdown.actions.close());
    }
    this.args.onLinkClick();
  }
}
