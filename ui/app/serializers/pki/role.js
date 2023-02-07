/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class PkiRoleSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
  };
}
