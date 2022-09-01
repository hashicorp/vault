import Component from '@glimmer/component';

/**
 * @module SecretListHeader
 * SecretListHeader component is breadcrumb, title with icon and menu with tabs component.
 *
 * @example
 * ```js
 * <SecretListHeader
   @model={{this.model}}
   @backendCrumb={{hash
    label=this.model.id
    text=this.model.id
    path="vault.cluster.secrets.backend.list-root"
    model=this.model.id
   }}
  />
 * ```
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 * @param {string} [baseKey] - Provided for navigation on the breadcrumbs.
 * @param {object} [backendCrumb] - Includes label, text, path and model ID.
 * @param {string} [baseKey] - Passthrough variable to the component KeyValueHeader.
 * @param {booelan} [isEngine=false] - If it's an engine we need to use the Ember Engines link-to-external and not the LinkTo component.

 */

export default class SecretListHeader extends Component {
  get isEngine() {
    return this.args.isEngine || false;
  }

  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
}
