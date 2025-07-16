/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
  States for disaster recovery and performance
  (More states might be added once this is hooked up to the backend)
*/
let REPLICATION_ENABLED_STATE = /*#__PURE__*/function (REPLICATION_ENABLED_STATE) {
  REPLICATION_ENABLED_STATE["PRIMARY"] = "primary";
  REPLICATION_ENABLED_STATE["SECONDARY"] = "secondary";
  REPLICATION_ENABLED_STATE["BOOTSTRAPPING"] = "bootstrapping";
  return REPLICATION_ENABLED_STATE;
}({});
const REPLICATION_DISABLED_STATE = 'disabled';

export { REPLICATION_DISABLED_STATE, REPLICATION_ENABLED_STATE };
//# sourceMappingURL=index.js.map
