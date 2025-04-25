/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { TrackedObject } from 'tracked-built-ins';
import { decodeString } from 'core/utils/b64';

export default class CustomMessagesService extends Service {
  @service api;
  @service namespace;
  @service auth;
  @tracked messages = [];
  @tracked showMessageModal = true;
  bannerState = new TrackedObject();

  constructor() {
    super(...arguments);
    this.fetchMessages();
  }

  get bannerMessages() {
    return this.messages?.filter((message) => message?.type === 'banner') || [];
  }

  get modalMessages() {
    return this.messages?.filter((message) => message?.type === 'modal') || [];
  }

  async fetchMessages() {
    try {
      const type = this.auth.currentToken ? 'Authenticated' : 'Unauthenticated';
      const method = `internalUiRead${type}ActiveCustomMessages`;
      const { keys = [], keyInfo } = await this.api.sys[method]();

      this.messages = keys.map((key) => {
        const data = keyInfo[key];
        return {
          id: key,
          ...data,
          message: data.message ? decodeString(data.message) : data.message,
        };
      });

      this.bannerMessages?.forEach((bm) => (this.bannerState[bm.id] = true));
    } catch (e) {
      this.clearCustomMessages();
    }
  }

  clearCustomMessages() {
    this.messages = [];
  }

  @action
  onBannerDismiss(id) {
    this.bannerState[id] = false;
  }
}
