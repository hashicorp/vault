/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { fillIn, click } from '@ember/test-helpers';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

export const mountBackend = async (type, path) => {
  await click(MOUNT_BACKEND_FORM.mountType(type));
  if (path) {
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.saveButton);
  } else {
    // save with default path
    await click(GENERAL.saveButton);
  }
};
