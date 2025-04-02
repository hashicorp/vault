/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';
import type CapabilitiesModel from 'vault/models/capabilities';

type PkiConfigAcmeModel = Model & {
  enabled: boolean;
  defaultDirectoryPolicy: string;
  allowedRoles: string[];
  allowRoleExtKeyUsage: boolean;
  allowedIssuers: string[];
  eabPolicy: string;
  dnsResolver: string;
  maxTtl: string;
  acmePath: CapabilitiesModel;
  get canSet(): boolean;
};

export default PkiConfigAcmeModel;
