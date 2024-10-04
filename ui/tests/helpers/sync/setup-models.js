/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// creates destination and association model for use in sync integration tests
// ensure that setupMirage is used prior to setupModels since this.server is used
export function setupModels(hooks) {
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    const destination = this.server.create('sync-destination', 'aws-sm', { name: 'us-west-1' });
    const destinationModelName = 'sync/destinations/aws-sm';
    this.store.pushPayload(destinationModelName, {
      modelName: destinationModelName,
      ...destination,
      id: destination.name,
    });
    this.destination = this.store.peekRecord(destinationModelName, destination.name);
    this.destinations = this.store.peekAll(destinationModelName);
    this.destinations.meta = {
      filteredTotal: this.destinations.length,
      currentPage: 1,
      pageSize: 5,
    };

    const association = this.server.create('sync-association', {
      type: 'aws-sm',
      name: 'us-west-1',
      mount: 'kv',
      secret_name: 'my-secret',
      sync_status: 'SYNCED',
      updated_at: '2023-09-20T10:51:53.961861096', // removed tz offset so time is consistently displayed
    });
    const associationModelName = 'sync/association';
    const associationId = `${association.mount}/${association.secret_name}`;
    this.store.pushPayload(associationModelName, {
      modelName: associationModelName,
      ...association,
      destinationType: 'aws-sm',
      destinationName: 'us-west-1',
      id: associationId,
    });

    this.association = this.store.peekRecord(associationModelName, associationId);
    this.associations = this.store.peekAll(associationModelName);
    this.associations.meta = {
      filteredTotal: this.associations.length,
      currentPage: 1,
      pageSize: 5,
    };
  });
}
