/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

const assertFields = (assert, fields, customSelectors = {}) => {
  fields.forEach((param) => {
    if (Object.keys(customSelectors).includes(param)) {
      assert.dom(customSelectors[param]).exists();
    } else if (param === 'config.listing_visibility') {
      assert.dom(GENERAL.toggleInput('toggle-config.listing_visibility')).exists();
    } else {
      assert.dom(GENERAL.inputByAttr(param)).exists();
    }
  });
};
export default (test) => {
  test('it renders mount fields', async function (assert) {
    await click(GENERAL.cardContainer(this.type));
    // This is where the "tune" parameters are rendered.
    await click(GENERAL.button('Method Options'));
    assertFields(assert, this.mountFields, this.customSelectors);
  });

  test('it renders tune fields', async function (assert) {
    // enable auth method to check tune fields
    await mountBackend(this.type, this.path);
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${this.path}/configuration`,
      `${this.type}: it mounts and navigates to configuration form`
    );

    assertFields(assert, this.configFields, this.customSelectors);

    for (const toggle in this.configToggles) {
      const fields = this.configToggles[toggle];
      await click(GENERAL.button(toggle));
      assertFields(assert, fields, this.customSelectors);
    }
  });
};
