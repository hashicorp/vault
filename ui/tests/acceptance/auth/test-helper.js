/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

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
    await click(MOUNT_BACKEND_FORM.mountType(this.type));
    await click(GENERAL.toggleGroup('Method Options'));
    assertFields(assert, this.mountFields, this.customSelectors);
  });

  test('it renders tune fields', async function (assert) {
    // enable auth method to check tune fields
    await mountBackend(this.type, this.path);
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
