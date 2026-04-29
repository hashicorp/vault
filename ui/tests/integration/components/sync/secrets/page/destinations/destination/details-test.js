/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { syncDestinations, findDestination } from 'vault/helpers/sync-destinations';
import { toLabel } from 'vault/helpers/to-label';
import { setupDataStubs } from 'vault/tests/helpers/sync/setup-hooks';
import { DestinationType, CLOUD_DESTINATION_TYPES } from 'sync/utils/constants';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SYNC_DESTINATIONS = syncDestinations();
module(
  'Integration | Component | sync | Secrets::Page::Destinations::Destination::Details',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'sync');
    setupMirage(hooks);
    setupDataStubs(hooks);

    hooks.beforeEach(function () {
      this.renderComponent = () => {
        return render(
          hbs` <Secrets::Page::Destinations::Destination::Details @destination={{this.destination}} @capabilities={{this.capabilities}} />`,
          { owner: this.engine }
        );
      };
    });

    test('it renders toolbar with actions', async function (assert) {
      assert.expect(3);

      await this.renderComponent();

      assert.dom(PAGE.toolbar('Delete destination')).exists();
      assert.dom(PAGE.toolbar('Sync secrets')).exists();
      assert.dom(PAGE.toolbar('Edit destination')).exists();
    });

    // module runs for each destination type
    for (const destination of SYNC_DESTINATIONS) {
      const { type } = destination;
      module(`destination: ${type}`, function (hooks) {
        hooks.beforeEach(function () {
          this.setupStubsForType(type);

          const { name, connection_details, options } = this.destination;
          this.details = { name, ...connection_details, ...options };
          this.fields = Object.keys(this.details).reduce((arr, key) => {
            const noCustomTags = !CLOUD_DESTINATION_TYPES.includes(type) && key === 'custom_tags';
            return noCustomTags ? arr : [...arr, key];
          }, []);

          const { maskedParams } = findDestination(type);
          this.maskedParams = maskedParams;

          this.getLabel = (field) => {
            const customLabel = {
              granularity_level: 'Secret sync granularity',
              access_key_id: 'Access key ID',
              role_arn: 'Role ARN',
              external_id: 'External ID',
              key_vault_uri: 'Key Vault URI',
              client_id: 'Client ID',
              tenant_id: 'Tenant ID',
              project_id: 'Project ID',
              credentials: 'JSON credentials',
              team_id: 'Team ID',
            }[field];

            return customLabel || toLabel([field]);
          };
        });

        test('it renders destination details with connection_details and options', async function (assert) {
          assert.expect(this.fields.length);

          await this.renderComponent();

          this.fields.forEach((field) => {
            if (this.maskedParams.includes(field)) {
              // these values are returned by the API masked: '*****'
              const label = this.getLabel(field);
              assert.dom(GENERAL.infoRowValue(label)).hasText('Destination credentials added');
            } else {
              // assert the remaining model attributes render
              const fieldValue = this.details[field];
              let label, value;
              if (field === 'custom_tags') {
                [label] = Object.keys(fieldValue);
                [value] = Object.values(fieldValue);
              } else {
                label = this.getLabel(field);
                value = Array.isArray(fieldValue) ? fieldValue.join(',') : fieldValue;
              }
              assert.dom(GENERAL.infoRowValue(label)).hasText(value);
            }
          });
        });

        test('it renders destination details without connection_details or options', async function (assert) {
          // Filter maskedParams to only include fields that are actually displayed on the details page
          // identity_token_audience and identity_token_key are masked but not displayed
          const displayedMaskedParams = this.maskedParams.filter((param) => {
            return !['identity_token_audience', 'identity_token_key'].includes(param);
          });

          assert.expect(displayedMaskedParams.length + 4);

          displayedMaskedParams.forEach((param) => {
            // these values are undefined when environment variables are set
            this.destination.connection_details[param] = undefined;
          });

          // assert custom tags section header does not render
          this.destination.options.custom_tags = undefined;

          await this.renderComponent();

          assert
            .dom(PAGE.destinations.details.sectionHeader)
            .doesNotExist('does not render Custom tags header');
          assert.dom(GENERAL.hdsPageHeaderTitle).hasTextContaining(this.destination.name);
          assert.dom(GENERAL.icon(findDestination(destination.type).icon)).exists();
          assert.dom(GENERAL.infoRowValue('Name')).hasText(this.destination.name);

          displayedMaskedParams.forEach((param) => {
            const label = this.getLabel(param);
            assert.dom(GENERAL.infoRowValue(label)).hasText('Using environment variable');
          });
        });
      });
    }

    // WIF-specific tests for cloud destinations
    module('WIF credential type display', function (hooks) {
      hooks.beforeEach(function () {
        // AWS auth setup helpers
        this.setupAwsWifAuth = () => {
          this.destination.connection_details.identity_token_audience = '*****';
          this.destination.connection_details.identity_token_ttl = 7200;
          this.destination.connection_details.role_arn = 'arn:aws:iam::123456789012:role/test-role';
          delete this.destination.connection_details.access_key_id;
          delete this.destination.connection_details.secret_access_key;
        };

        this.setupAwsAccountAuth = () => {
          this.destination.connection_details.access_key_id = '*****';
          this.destination.connection_details.secret_access_key = '*****';
        };

        // Azure auth setup helpers
        this.setupAzureWifAuth = () => {
          this.destination.connection_details.identity_token_audience = '*****';
          this.destination.connection_details.identity_token_ttl = 3600;
          delete this.destination.connection_details.client_secret;
        };

        this.setupAzureAccountAuth = () => {
          this.destination.connection_details.client_secret = '*****';
        };

        // GCP auth setup helpers
        this.setupGcpWifAuth = () => {
          this.destination.connection_details.identity_token_audience = '*****';
          this.destination.connection_details.identity_token_ttl = 3600;
          this.destination.connection_details.service_account_email = 'test@project.iam.gserviceaccount.com';
          delete this.destination.connection_details.credentials;
        };

        this.setupGcpAccountAuth = () => {
          this.destination.connection_details.credentials = '*****';
        };
      });

      test('aws-sm: it displays IAM credential type for account-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.AwsSm);
        this.setupAwsAccountAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert.dom(GENERAL.infoRowValue('Credential type')).hasText('IAM', 'Shows IAM credential type');
      });

      test('aws-sm: it displays WIF credential type for WIF-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.AwsSm);
        this.setupAwsWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert.dom(GENERAL.infoRowValue('Credential type')).hasText('WIF', 'Shows WIF credential type');
      });

      test('aws-sm: it formats identity_token_ttl as duration', async function (assert) {
        this.setupStubsForType(DestinationType.AwsSm);
        this.setupAwsWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert
          .dom(GENERAL.infoRowLabel('Identity token time to live'))
          .exists('Identity token TTL label is displayed');
        assert
          .dom(GENERAL.infoRowValue('Identity token time to live'))
          .hasText('2 hours', 'TTL is formatted as duration (7200s = 2 hours)');
      });

      test('azure-kv: it displays Client secret credential type for account-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.AzureKv);
        this.setupAzureAccountAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert
          .dom(GENERAL.infoRowValue('Credential type'))
          .hasText('Client secret', 'Shows Client secret credential type');
      });

      test('azure-kv: it displays WIF credential type for WIF-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.AzureKv);
        this.setupAzureWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert.dom(GENERAL.infoRowValue('Credential type')).hasText('WIF', 'Shows WIF credential type');
      });

      test('gcp-sm: it displays JSON credential type for account-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.GcpSm);
        this.setupGcpAccountAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert.dom(GENERAL.infoRowValue('Credential type')).hasText('JSON', 'Shows JSON credential type');
      });

      test('gcp-sm: it displays WIF credential type for WIF-based auth', async function (assert) {
        this.setupStubsForType(DestinationType.GcpSm);
        this.setupGcpWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Credential type')).exists('Credential type label is displayed');
        assert.dom(GENERAL.infoRowValue('Credential type')).hasText('WIF', 'Shows WIF credential type');
      });

      test('gcp-sm: it displays service_account_email for WIF auth', async function (assert) {
        this.setupStubsForType(DestinationType.GcpSm);
        this.setupGcpWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert
          .dom(GENERAL.infoRowLabel('Service account email'))
          .exists('Service account email label is displayed');
        assert
          .dom(GENERAL.infoRowValue('Service account email'))
          .hasText('test@project.iam.gserviceaccount.com', 'Shows service account email');
      });

      test('aws-sm: it does not display account fields when WIF is configured', async function (assert) {
        this.setupStubsForType(DestinationType.AwsSm);
        this.setupAwsWifAuth();
        assert.expect(2);

        await this.renderComponent();

        assert.dom(GENERAL.infoRowLabel('Access key ID')).doesNotExist('Access key ID is not displayed');
        assert
          .dom(GENERAL.infoRowLabel('Secret access key'))
          .doesNotExist('Secret access key is not displayed');
      });

      test('aws-sm: it does not display WIF fields when account credentials are configured', async function (assert) {
        this.setupStubsForType(DestinationType.AwsSm);
        this.setupAwsAccountAuth();
        assert.expect(2);

        await this.renderComponent();

        assert
          .dom(GENERAL.infoRowLabel('Identity token audience'))
          .doesNotExist('Identity token audience is not displayed');
        assert
          .dom(GENERAL.infoRowLabel('Identity token time to live'))
          .doesNotExist('Identity token TTL is not displayed');
      });
    });
  }
);
