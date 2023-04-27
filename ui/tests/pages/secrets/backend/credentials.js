/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable } from 'ember-cli-page-object';

export const Base = {
  visit: visitable('/vault/secrets/:backend/credentials/:id'),
  visitRoot: visitable('/vault/secrets/:backend/credentials'),
};

export default create(Base);
