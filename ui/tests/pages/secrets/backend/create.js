/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable } from 'ember-cli-page-object';

export const Base = {
  visit: visitable('/vault/secrets-engines/:backend/create/:id'),
  visitRoot: visitable('/vault/secrets-engines/:backend/create'),
};
export default create(Base);
