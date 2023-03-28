/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
export default Controller.extend({
  namespaceService: service('namespace'),
});
