/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import Controller from '@ember/controller';

export default class NamespaceCreateController extends Controller {
  @service api;
  @service flashMessages;
  @service('namespace') namespaceService;
  @service router;

  @tracked modelValidations = null;
  @tracked errorMessage = null;

  save = task(
    { drop: true },
    waitFor(async (event) => {
      event.preventDefault();
      this.errorMessage = null;
      const { isValid, state, data } = this.model.toJSON();
      // refresh validations on every submit so a corrected field clears its previous error
      this.modelValidations = state;
      if (!isValid) {
        return;
      }
      try {
        await this.api.sys.systemWriteNamespacesPath(data.path, {});
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorMessage = message;
        return;
      }
      this.flashMessages.success('Saved!');
      // fetch new namespaces for the namespace picker
      this.namespaceService.findNamespacesForUser.perform();
      this.router.transitionTo('vault.cluster.access.namespaces.index');
    })
  );
}
