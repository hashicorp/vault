/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { action } from '@ember/object';
import { ClientFilterTypes } from 'core/utils/client-count-utils';

export default class ClientsClientListPageComponent extends ActivityComponent {
  @action
  handleFilter(filters: Record<ClientFilterTypes, string>) {
    const { nsLabel, mountPath, mountType } = filters;
    this.args.onFilterChange({ nsLabel, mountPath, mountType });
  }
}
