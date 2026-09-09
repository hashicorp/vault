/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { SECRET_TYPE_CONFIGS } from 'sync/utils/secret-type-config';

import type { SecretType } from 'vault/sync';

interface Args {
  type?: SecretType | null;
}

export default class SecretTypeBadge extends Component<Args> {
  get config() {
    return this.args.type ? SECRET_TYPE_CONFIGS[this.args.type] : null;
  }

  get text() {
    return this.config?.accessorType || 'Engine mount path';
  }

  get icon() {
    return this.config?.icon || 'lock';
  }
}
