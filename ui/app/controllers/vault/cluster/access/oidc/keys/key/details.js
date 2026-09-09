/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default class OidcKeyDetailsController extends Controller {
  @service api;
  @service router;
  @service flashMessages;

  rotateKey = task(
    waitFor(async () => {
      try {
        const { name, verification_ttl } = this.model.key;
        await this.api.identity.oidcRotateKey(name, { verification_ttl });
        this.flashMessages.success(`Success: ${name} connection was rotated.`);
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(message);
      }
    })
  );

  @action
  async delete() {
    try {
      await this.api.identity.oidcDeleteKey(this.model.key.name);
      this.flashMessages.success('Key deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.keys');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(message);
    }
  }
}
