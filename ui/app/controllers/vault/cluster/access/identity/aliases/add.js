/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ROUTES } from 'vault/utils/routes';
import CreateController from '../create';

export default CreateController.extend({
  showRoute: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_ALIASES_SHOW,
});
