/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { capitalize } from '@ember/string';
import { singularize } from 'ember-inflector';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type ControlGroupService from 'vault/vault/services/control-group';
import type { ControlGroupErrorLog } from 'vault/vault/services/control-group';
import type { ApiParsedError } from 'vault/api';
import type GeneratedItemForm from 'vault/forms/generated-item';
import type { GeneratedItemUrls } from 'vault/utils/generated-item-urls';

interface Args {
  form: GeneratedItemForm;
  mode: 'create' | 'edit' | 'show';
  itemType: string;
  mountPath: string;
  urls: GeneratedItemUrls;
}

/**
 * @module GeneratedItem
 * Create, edit and show views for auth method items whose fields come from OpenAPI (e.g. userpass users, LDAP groups).
 *
 * @example
 * <GeneratedItem @form={{this.model}} @urls={{this.urls}} @mountPath="userpass" @mode="show" @itemType="user" />
 */
export default class GeneratedItem extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly controlGroup: ControlGroupService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked errorMessage = '';
  @tracked controlGroupError: ControlGroupErrorLog | null = null;

  get title() {
    const { form, mode, itemType } = this.args;
    if (mode === 'show') {
      return form.id;
    }
    const title = `${capitalize(mode)} ${singularize(itemType)}`;
    return mode === 'edit' ? `${title} ${form.id}` : title;
  }

  // A write to a path behind a Control Group needs approval; save the token so a retry after approval unwraps it.
  handleControlGroupError(response: ApiParsedError['response']) {
    if (!response?.isControlGroupError) {
      return false;
    }
    this.controlGroup.saveTokenFromError(response);
    this.controlGroupError = this.controlGroup.logFromError(response);
    return true;
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      this.errorMessage = '';
      this.controlGroupError = null;
      const { form, mode, urls, itemType } = this.args;
      const id = form.id as string;
      const { data } = form.toJSON();
      try {
        if (mode === 'create') {
          await this.api.request.post(urls.create(id), data);
        } else {
          await this.api.request.put(urls.update(id), data);
        }
      } catch (error) {
        const { message, response } = await this.api.parseError(error);
        if (!this.handleControlGroupError(response)) {
          this.errorMessage = message;
        }
        return;
      }
      this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
      this.flashMessages.success(`Successfully saved ${itemType} ${id}.`);
    })
  );

  @action
  async deleteItem() {
    const { form, urls, itemType } = this.args;
    const id = form.id as string;
    this.controlGroupError = null;
    try {
      await this.api.request.delete(urls.delete(id));
    } catch (error) {
      const { message, response } = await this.api.parseError(error);
      if (!this.handleControlGroupError(response)) {
        this.flashMessages.danger(`There was an error deleting this ${itemType} ${message}`);
      }
      return;
    }
    this.router.transitionTo('vault.cluster.access.method.item.list').followRedirects();
    this.flashMessages.success(`Successfully deleted ${itemType} ${id}.`);
  }
}
