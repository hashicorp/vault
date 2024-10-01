/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  mountType: (name) => `[data-test-mount-type="${name}"]`,
  submit: '[data-test-mount-submit]',
};

const assertFields = (assert, fields, customSelectors = {}) => {
  fields.forEach((param) => {
    if (Object.keys(customSelectors).includes(param)) {
      assert.dom(customSelectors[param]).exists();
    } else {
      assert.dom(GENERAL.inputByAttr(param)).exists();
    }
  });
};
export default (test) => {
  test('it renders mount fields', async function (assert) {
    await click(SELECTORS.mountType(this.type));
    await click(GENERAL.toggleGroup('Method Options'));
    assertFields(assert, this.mountFields, this.customSelectors);
  });

  test('it renders tune fields', async function (assert) {
    // enable auth method to check tune fields
    await click(SELECTORS.mountType(this.type));
    await fillIn(GENERAL.inputByAttr('path'), this.path);
    await click(SELECTORS.submit);
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${this.path}/configuration`,
      `${this.type}: it mounts navigates to tune form`
    );

    assertFields(assert, this.tuneFields, this.customSelectors);

    for (const toggle in this.tuneToggles) {
      const fields = this.tuneToggles[toggle];
      await click(GENERAL.toggleGroup(toggle));
      assertFields(assert, fields, this.customSelectors);
    }
  });
};
