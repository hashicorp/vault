/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import routerLookup from 'core/utils/router-lookup';

import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ManageDropdown
 * Reusable component for displaying the Manage dropdown used in secret engine headers & secret engine mount list.
 *
 * @example
 * In main app page headers and list components — uses the resource getter for the full absolute route
 * <ManageDropdown
 *   @model={{this.backendModel}}
 *   @configRoute={{this.backendConfigurationLink}}
 * />
 *
 * In Ember engine templates (pki, kubernetes, ldap, kmip, kv) — pass the short relative route,
 * since HDS @route resolves relative to the engine's router mount
 * <ManageDropdown
 *   @model={{@model}}
 *   @configRoute="configuration"
 * />
 *
 * With custom menu items (like KV's Generate policy) — icon variant in a Ember engine list
 * <ManageDropdown
 *   @model={{@backendModel}}
 *   @configRoute="configuration"
 *   as |D|
 * >
 *   <D.Interactive @icon="shield-check" {{on "click" openFlyout}}>Generate policy</D.Interactive>
 * </ManageDropdown>
 *
 * @param {SecretsEngineResource} model - The secrets engine resource containing the engine details
 * @param {string} configRoute - Route for the Configure action.
 * @param {string} variant - Set to "icon" for "..." icon button, otherwise shows "Manage" text button (default)
 */

interface Args {
  model: SecretsEngineResource;
  configRoute: string;
  variant?: 'icon';
}

export default class ManageDropdown extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;

  get router(): RouterService {
    return routerLookup(this);
  }

  get isIcon() {
    return this.args.variant === 'icon';
  }

  get configureRouteModel() {
    return this.args.model.id;
  }

  get shouldShowDelete() {
    // Don't show delete for cubbyhole engine
    return this.args.model.type !== 'cubbyhole';
  }

  transitionOrRefresh() {
    const { currentRouteName } = this.router;
    // Call refresh() when currently on the route so data properly refreshes even when in a namespace.
    const method = currentRouteName === 'vault.cluster.secrets.backends' ? 'refresh' : 'transitionTo';
    this.router[method]('vault.cluster.secrets.backends');
  }

  @action
  handleDeleteClick(engine: SecretsEngineResource) {
    this.engineToDisable = engine;
  }

  @action
  handleModalClose() {
    this.engineToDisable = undefined;
  }

  @action
  async handleModalConfirm() {
    if (this.engineToDisable) {
      const { engineType, id, path } = this.engineToDisable;

      try {
        await this.api.sys.mountsDisableSecretsEngine(id);
        this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
        this.transitionOrRefresh();
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.flashMessages.danger(
          `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${message}.`
        );
      } finally {
        this.engineToDisable = undefined;
      }
    }
  }
}
