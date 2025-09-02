/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Store from '@ember-data/store';
import { AdapterRegistry } from 'ember-data/adapter';

interface ExportDataQuery {
  format?: string;
  start_time?: string;
  end_time?: string;
  namespace?: string;
}

export default interface ActivityAdapter extends AdapterRegistry {
  exportData(query?: ExportDataQuery): Promise<Blob>;
}
