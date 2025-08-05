/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import camelizeKeys from 'vault/utils/camelize-object-keys';

// creates destination and association model for use in sync integration tests
// ensure that setupMirage is used prior to setupModels since this.server is used
export function setupDataStubs(hooks) {
  hooks.beforeEach(function () {
    // most tests are good with the default data generated here
    // allow for this to be overridden to test different types
    this.setupStubsForType = (destType) => {
      const {
        id, // eslint-disable-line no-unused-vars
        name,
        type,
        granularity,
        secret_name_template,
        custom_tags,
        purge_initiated_at,
        purge_error,
        ...connection_details
      } = this.server.create('sync-destination', destType);

      this.destination = {
        name,
        type,
        connectionDetails: camelizeKeys(connection_details),
        options: {
          granularityLevel: granularity,
          secretNameTemplate: secret_name_template,
          customTags: custom_tags,
        },
        purgeInitiatedAt: purge_initiated_at,
        purgeError: purge_error,
      };

      this.destinations = [this.destination];
      this.destinations.meta = {
        filteredTotal: this.destinations.length,
        currentPage: 1,
        pageSize: 5,
      };

      const association = this.server.create('sync-association', {
        type: this.destination.type,
        name: this.destination.name,
        mount: 'kv',
        secret_name: 'my-secret',
        sync_status: 'SYNCED',
        updated_at: '2023-09-20T10:51:53.961861096', // removed tz offset so time is consistently displayed
      });
      this.association = {
        ...camelizeKeys(association),
        destinationType: this.destination.type,
        destinationName: this.destination.name,
      };
      this.associations = [this.association];
      this.associations.meta = {
        filteredTotal: this.associations.length,
        currentPage: 1,
        pageSize: 5,
      };

      const capabilitiesService = this.owner.lookup('service:capabilities');
      const paths = [
        capabilitiesService.pathFor('syncDestination', this.destination),
        capabilitiesService.pathFor('syncSetAssociation', this.destination),
        capabilitiesService.pathFor('syncRemoveAssociation', this.destination),
      ];
      this.capabilities = paths.reduce((obj, path) => {
        obj[path] = { canRead: true, canCreate: true, canUpdate: true, canDelete: true };
        return obj;
      }, {});
    };

    this.setupStubsForType('aws-sm');
  });
}
