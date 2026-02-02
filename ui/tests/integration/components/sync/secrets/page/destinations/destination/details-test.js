/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { syncDestinations, findDestination } from 'vault/helpers/sync-destinations';
import { toLabel } from 'vault/helpers/to-label';
import { setupDataStubs } from 'vault/tests/helpers/sync/setup-hooks';
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
            const noCustomTags = ['gh', 'vercel-project'].includes(type) && key === 'custom_tags';
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
              assert.dom(PAGE.infoRowValue(label)).hasText('Destination credentials added');
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
              assert.dom(PAGE.infoRowValue(label)).hasText(value);
            }
          });
        });

        test('it renders destination details without connection_details or options', async function (assert) {
          assert.expect(this.maskedParams.length + 4);

          this.maskedParams.forEach((param) => {
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
          assert.dom(PAGE.icon(findDestination(destination.type).icon)).exists();
          assert.dom(PAGE.infoRowValue('Name')).hasText(this.destination.name);

          this.maskedParams.forEach((param) => {
            const label = this.getLabel(param);
            assert.dom(PAGE.infoRowValue(label)).hasText('Using environment variable');
          });
        });
      });
    }
  }
);
