/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class KubernetesRoleSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  attrs = {
    backend: { serialize: false },
  };
}
