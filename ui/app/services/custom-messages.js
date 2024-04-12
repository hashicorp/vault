/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { TrackedObject } from 'tracked-built-ins';

export default class CustomMessagesService extends Service {
  @service store;
  @service namespace;
  @service auth;
  @tracked messages = [];
  @tracked showMessageModal = true;
  bannerState = new TrackedObject();

  constructor() {
    super(...arguments);
    this.fetchMessages(this.namespace.path);
  }

  get bannerMessages() {
    if (!this.messages || !this.messages.length) return [];
    return this.messages?.filter((message) => message?.type === 'banner');
  }

  get modalMessages() {
    if (!this.messages || !this.messages.length) return [];
    return this.messages?.filter((message) => message?.type === 'modal');
  }

  async fetchMessages(ns) {
    try {
      const url = this.auth.currentToken
        ? '/v1/sys/internal/ui/authenticated-messages'
        : '/v1/sys/internal/ui/unauthenticated-messages';
      const opts = {
        method: 'GET',
        headers: {},
      };
      if (this.auth.currentToken) opts.headers['X-Vault-Token'] = this.auth.currentToken;
      if (ns) opts.headers['X-Vault-Namespace'] = ns;
      const result = await fetch(url, opts);
      const body = await result.json();
      if (body.errors) return (this.messages = []);
      const serializer = this.store.serializerFor('config-ui/message');
      this.messages = serializer.mapPayload(body);
      this.bannerMessages?.forEach((bm) => (this.bannerState[bm.id] = true));
    } catch (e) {
      return e;
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
