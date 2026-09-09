/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { fillIn, click } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

export const mountBackend = async (type, path, isSecret) => {
  await click(GENERAL.cardContainer(type));

  // catalog for secrets requires clicking next to route, auth method will auto-route on card click
  if (isSecret) {
    await click(GENERAL.button('next'));
  }
  if (path) {
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.submitButton);
  } else {
    // save with default path
    await click(GENERAL.submitButton);
  }
};
