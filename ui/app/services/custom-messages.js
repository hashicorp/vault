import Service, { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
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
    if (this.messages?.errors) return [];
    return this.messages?.filter((message) => message.type === 'banner');
  }

  get modalMessages() {
    if (this.messages?.errors) return [];
    return this.messages?.filter((message) => message.type === 'modal');
  }

  @task
  *getMessages(ns) {
    try {
      const url = '/v1/sys/internal/ui/unauthenticated-messages';
      const opts = {
        method: 'GET',
        headers: {},
      };
      if (ns) {
        opts.headers['X-Vault-Namespace'] = ns;
      }
      const result = yield fetch(url, opts);
      const body = yield result.json();
      const serializer = this.store.serializerFor('config-ui/message');
      this.messages = serializer.mapPayload(body);
    } catch (e) {
      return e;
    }
  }

  fetchMessages(ns) {
    return this.getMessages.perform(ns);
  }
}
