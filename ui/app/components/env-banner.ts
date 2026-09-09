/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import ENV from 'vault/config/environment';

export default class EnvBannerComponent extends Component {
  get isDevelopment(): boolean {
    return ENV.environment === 'development';
  }

  get vaultBranch(): string | null | undefined {
    return ENV.APP.VAULT_BRANCH;
  }

  get vaultCommit(): string | null | undefined {
    return ENV.APP.VAULT_COMMIT;
  }

  get hasBranchAndCommit(): boolean {
    return !!(this.vaultBranch && this.vaultCommit);
  }
}
