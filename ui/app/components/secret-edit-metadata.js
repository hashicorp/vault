import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

// ARG TODO FILL OUT
/**
 * @module AuthInfo
 *
 * @example
 * ```js
 * <AuthInfo @activeClusterName={{cluster.name}} @onLinkClick={{action "onLinkClick"}} />
 * ```
 *
 * @param {string} activeClusterName - name of the current cluster, passed from the parent.
 * @param {Function} onLinkClick - parent action which determines the behavior on link click
 */
export default class SecretEditMetadata extends Component {
  @service router;
  @service store;

  async save() {
    let model = this.args.model;
    console.log(model, 'MODEL');
    try {
      await model.save();
    } catch (err) {
      // ARG TODO handle error in error object
      console.log(err, 'ERROR');
    }
    this.router.transitionTo('vault.cluster.secrets.backend.metadata', this.args.model.id);
  }

  @action
  onSaveChanges(event) {
    event.preventDefault();
    const changed = this.args.model.hasDirtyAttributes; // ARG TODO when API done double check this is working
    if (changed) {
      this.save();
      return;
    }
    // ARG TODO else validation error?
  }
}
