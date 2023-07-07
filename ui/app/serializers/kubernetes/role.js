/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

export default class KubernetesRoleSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  attrs = {
    backend: { serialize: false },
  };
}
