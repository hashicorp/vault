/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';
import ENV from 'vault/config/environment';

export default helper(function ([path]) {
  return path.replace(/^~\//, `${ENV.rootURL}images/`);
});
