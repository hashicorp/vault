/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import IdentityManager from 'vault/utils/identity-manager';
// to more closely match the Vault backend this will return UUIDs as identifiers for records in mirage
export default IdentityManager;
