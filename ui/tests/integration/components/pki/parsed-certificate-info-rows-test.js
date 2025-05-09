/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | parsed-certificate-info-rows', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  test('it renders nothing if no valid attributes passed', async function (assert) {
    this.set('parsedCertificate', {
      foo: '',
      common_name: 'not-shown',
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });

    assert.dom(this.element).hasText('');
  });

  test('it renders only valid attributes with values', async function (assert) {
    this.set('parsedCertificate', {
      common_name: 'not-shown',
      use_pss: false,
      alt_names: ['something', 'here'],
      ttl: undefined,
    });
    await render(hbs`<ParsedCertificateInfoRows @model={{this.parsedCertificate}} />`, {
      owner: this.engine,
    });
    assert.dom('[data-test-component="info-table-row"]').exists({ count: 2 }, 'renders 2 rows');

    assert.dom(GENERAL.infoRowValue('Common name')).doesNotExist('common name is never rendered');
    assert.dom('[data-test-row-value="Subject Alternative Names (SANs)"]').hasText('something,here');
    assert.dom(GENERAL.infoRowValue('Use PSS')).hasText('No', 'Booleans are rendered');
    assert.dom(GENERAL.infoRowValue('ttl')).doesNotExist('ttl is not rendered because value undefined');
    assert
      .dom('[data-test-parsing-error-alert-banner]')
      .doesNotExist('does not render parsing error info banner');
  });
});
