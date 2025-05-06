/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import yaml from 'js-yaml';

export function yamlify(target) {
  return yaml.dump(target);
}

export default buildHelper(yamlify);
