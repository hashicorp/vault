/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';

/* sample response

{
  "request_id": "d81bba81-e8a1-0ee9-240e-a77d36e3e08f",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "cluster_id": "ab7d4191-d1a3-b4d6-6297-5a41af6154ae",
    "known_secondaries": [
      "test"
    ],
    "last_performance_wal": 72,
    "last_reindex_epoch": "1588281113",
    "last_wal": 73,
    "merkle_root": "c8d258d376f01d98156f74e8d8f82ea2aca8dc4a",
    "mode": "primary",
    "primary_cluster_addr": "",
    "reindex_building_progress": 26838,
    "reindex_building_total": 305443,
    "reindex_in_progress": true,
    "reindex_stage": "building",
    "state": "running"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}


*/

export default Model.extend({
  status: attr('object'),
});
