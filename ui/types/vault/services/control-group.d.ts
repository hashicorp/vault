/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

export default class ControlGroupService extends Service {
  tokenForUrl(url: string): string;
}
