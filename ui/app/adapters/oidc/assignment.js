/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import NamedPathAdapter from '../named-path';

export default class OidcAssignmentAdapter extends NamedPathAdapter {
  pathForType() {
    return 'identity/oidc/assignment';
  }
}
