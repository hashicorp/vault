/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable } from 'ember-cli-page-object';

export const Base = {
  visit: visitable('/vault/secrets/:backend/show/:id'),
  visitRoot: visitable('/vault/secrets/:backend/show'),
};
export default create(Base);
