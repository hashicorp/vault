import Service, { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { decodeString } from 'core/utils/b64';

export default class VersionService extends Service {
  @service store;
  @tracked messages = [];

  get bannerMessages() {
    return this.messages.filter((message) => message.type === 'banner');
  }

  get modalMessages() {
    return this.messages.filter((message) => message.type === 'modal');
  }

  @task
  *getMessages(ns) {
    try {
      const url = '/v1/sys/internal/ui/unauthenticated-messages';
      const opts = {
        method: 'GET',
      };
      if (ns) {
        opts.headers['X-Vault-Namespace'] = ns;
      }
      const result = yield fetch(url, opts);
      const body = yield result.json();
      if (body.data) {
        if (body.data?.keys && Array.isArray(body.data.keys)) {
          this.messages = body.data.keys.map((key) => {
            const data = {
              id: key,
              linkTitle: body.data.key_info.link?.title,
              linkHref: body.data.key_info.link?.href,
              ...body.data.key_info[key],
            };
            data.message = decodeString(data.message);
            return data;
          });
        }
      }
      return body;
    } catch (e) {
      return null;
    }
  }

  fetchMessages(ns) {
    return this.getMessages.perform(ns);
  }
}
