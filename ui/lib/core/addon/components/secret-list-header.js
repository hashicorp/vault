import Component from '@glimmer/component';
import { optionsForBackend } from 'core/helpers/options-for-backend';
import { engineOptionsForBackend } from 'core/helpers/engine-options-for-backend';

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
 * @param {boolean} [isEngine=false] - Changes link type if the component is being used inside an Ember engine.
 */

export default class SecretListHeader extends Component {
  get isConfigure() {
    return this.args.isConfigure || false;
  }

  get isCertTab() {
    return this.args.isCertTab || false;
  }

  get isEngine() {
    return this.args.isEngine || false;
  }
  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
  get optionsForBackendHelper() {
    return this.isEngine
      ? engineOptionsForBackend([this.args.model.engineType])
      : optionsForBackend([this.args.model.engineType]);
  }
}
