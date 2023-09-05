/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';

module('Integration | Component | ldap | Page::Library::Details::Accounts', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.store = this.owner.lookup('service:store');

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      ...this.server.create('ldap-library', { name: 'test-library' }),
    });
    this.model = this.store.peekRecord('ldap/library', 'test-library');
    this.statuses = [
      {
        account: 'foo.bar',
        available: false,
        library: 'test-library',
        borrower_client_token: '123',
        borrower_entity_id: '456',
      },
      { account: 'bar.baz', available: true, library: 'test-library' },
    ];
    this.renderComponent = () => {
      return render(
        hbs`
          <div id="modal-wormhole"></div>
          <Page::Library::Details::Accounts @library={{this.model}} @statuses={{this.statuses}} />
        `,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should render account cards', async function (assert) {
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await this.renderComponent();

    assert.dom('[data-test-account-name="foo.bar"]').hasText('foo.bar', 'Account name renders');
    assert
      .dom('[data-test-account-status="foo.bar"]')
      .hasText('Unavailable', 'Correct badge renders for checked out account');
    assert
      .dom('[data-test-account-status="bar.baz"]')
      .hasText('Available', 'Correct badge renders for available account');

    await click('[data-test-check-out]');
    await fillIn('[data-test-ttl-value="TTL"]', 4);
    await click('[data-test-check-out="save"]');

    const didTransition = transitionStub.calledWith(
      'vault.cluster.secrets.backend.ldap.libraries.library.check-out',
      { queryParams: { ttl: '4h' } }
    );
    assert.true(didTransition, 'Transitions to check out route on action click');

    assert.dom('[data-test-checked-out-card]').exists('Accounts checked out card renders');

    assert
      .dom('[data-test-cli-command]')
      .hasText('vault lease renew ad/library/test-library/check-out/:lease_id', 'Renew cli command renders');
    assert.dom(`[data-test-cli-command-copy]`).exists('Renew cli command copy button renders');
  });
});
