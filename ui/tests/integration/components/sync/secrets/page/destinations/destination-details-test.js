/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { syncDestinations } from 'vault/helpers/sync-destinations';
import { toLabel } from 'vault/helpers/to-label';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

const SYNC_DESTINATIONS = syncDestinations();
module(
  'Integration | Component | sync | Secrets::Page::Destinations::Destination::Details',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'sync');
    setupMirage(hooks);

    hooks.beforeEach(function () {
      this.store = this.owner.lookup('service:store');

      this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

      this.renderFormComponent = () => {
        return render(
          hbs` <Secrets::Page::Destinations::Destination::Details @destination={{this.model}} />`,
          { owner: this.engine }
        );
      };
    });

    test('it renders toolbar with actions', async function (assert) {
      assert.expect(3);
      const type = SYNC_DESTINATIONS[0].type;
      const data = this.server.create('sync-destination', type);

      const id = `${type}/${data.name}`;
      data.id = id;
      this.store.pushPayload(`sync/destinations/${type}`, {
        modelName: `sync/destinations/${type}`,
        ...data,
      });
      this.model = this.store.peekRecord(`sync/destinations/${type}`, id);

      await this.renderFormComponent();

      assert.dom(PAGE.toolbar('Delete destination')).exists();
      assert.dom(PAGE.toolbar('Sync new secret')).exists();
      assert.dom(PAGE.toolbar('Edit destination')).exists();
    });

    // module runs for each destination type
    for (const destination of SYNC_DESTINATIONS) {
      const { type } = destination;
      module(`destination: ${type}`, function (hooks) {
        hooks.beforeEach(function () {
          const data = this.server.create('sync-destination', type);

          const id = `${type}/${data.name}`;
          data.id = id;
          this.store.pushPayload(`sync/destinations/${type}`, {
            modelName: `sync/destinations/${type}`,
            ...data,
          });
          this.model = this.store.peekRecord(`sync/destinations/${type}`, id);
          const { maskedParams } = this.model;
          this.maskedAttrs = this.model.formFields.filter((attr) => maskedParams.includes(attr.name));
          this.unmaskedAttrs = this.model.formFields.filter((attr) => !maskedParams.includes(attr.name));
        });

        test('it renders destination details with connection_details', async function (assert) {
          assert.expect(this.model.formFields.length);

          await this.renderFormComponent();

          // these values are returned by the API masked: '*****'
          this.maskedAttrs.forEach((attr) => {
            const label = attr.options?.label || toLabel([attr.name]);
            assert.dom(PAGE.infoRowValue(label)).hasText('Destination credentials added');
          });

          // assert the remaining model attributes render
          this.unmaskedAttrs.forEach(({ name, options }) => {
            const label = options.label || toLabel([name]);
            const value = Array.isArray(this.model[name]) ? this.model[name].join(',') : this.model[name];
            assert.dom(PAGE.infoRowValue(label)).hasText(value);
          });
        });

        test('it renders destination details without connection_details', async function (assert) {
          assert.expect(this.maskedAttrs.length + 3);

          this.maskedAttrs.forEach((attr) => {
            // these values are undefined when environment variables are set
            this.model[attr.name] = undefined;
          });

          await this.renderFormComponent();

          assert.dom(PAGE.title).hasTextContaining(this.model.name);
          assert.dom(PAGE.icon(this.model.icon)).exists();
          assert.dom(PAGE.infoRowValue('Name')).hasText(this.model.name);

          this.maskedAttrs.forEach((attr) => {
            const label = attr.options?.label || toLabel([attr.name]);
            assert.dom(PAGE.infoRowValue(label)).hasText('Using environment variable');
          });
        });
      });
    }
  }
);
