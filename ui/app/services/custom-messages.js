/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class CustomMessageService extends Service {
  @service store;
  @service namespace;
  @tracked messages = [];
  @tracked showMessageModal = true;

  constructor() {
    super(...arguments);
    this.fetchMessages(this.namespace.path);
  }

  get bannerMessages() {
    return (this.messages || []).filter((message) => message?.type === 'banner');
  }

  get modalMessages() {
    return (this.messages || []).filter((message) => message?.type === 'modal');
  }

  async fetchMessages(ns) {
    try {
      const url = '/v1/sys/internal/ui/unauthenticated-messages';
      const opts = {
        method: 'GET',
        headers: {},
      };
      if (ns) {
        opts.headers['X-Vault-Namespace'] = ns;
      }
      const result = await fetch(url, opts);
      const body = await result.json();
      if (body.errors) return (this.messages = []);
      const serializer = this.store.serializerFor('config-ui/message');
      this.messages = serializer.mapPayload(body);
    } catch (e) {
      return e;
    }
  }
}
