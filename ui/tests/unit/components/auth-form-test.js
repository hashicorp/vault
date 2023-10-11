/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { settled } from '@ember/test-helpers';

module('Unit | Component | auth-form', function (hooks) {
  setupTest(hooks);

  test('it should use token for oidc and jwt auth method types when processing form submit', async function (assert) {
    assert.expect(4);

    const component = this.owner.lookup('component:auth-form');
    component.reopen({
      methods: [], // eslint-disable-line
      // eslint-disable-next-line
      authenticate: {
        unlinked() {
          return {
            perform(type, data) {
              assert.deepEqual(
                type,
                'token',
                `Token type correctly passed to authenticate method for ${component.providerName}`
              );
              assert.deepEqual(
                data,
                { token: component.token },
                `Token passed to authenticate method for ${component.providerName}`
              );
            },
          };
        },
      },
    });

    const event = new Event('submit');

    for (const type of ['oidc', 'jwt']) {
      component.set('selectedAuth', type);
      await settled();
      await component.actions.doSubmit.apply(component, [undefined, event, 'foo-bar']);
    }
  });
});
