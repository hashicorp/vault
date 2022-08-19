import Controller from '@ember/controller';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default class OidcKeyDetailsController extends Controller {
  @service router;
  @service flashMessages;

  @task
  @waitFor
  *rotateKey() {
    const adapter = this.store.adapterFor('oidc/key');
    yield adapter
      .rotate(this.model.name, this.model.verificationTTL)
      .then(() => {
        this.flashMessages.success(`Success: ${this.model.name} connection was rotated`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      });
  }
  @action
  async delete() {
    try {
      await this.model.destroyRecord();
      this.flashMessages.success('Key deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.keys');
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }
}
