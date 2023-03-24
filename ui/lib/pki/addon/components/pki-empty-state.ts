/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

interface Args {
  roles: Array<string> | null;
  certificates: Array<string> | null;
}

export default class PkiEmptyState extends Component<Args> {
  get message() {
    const defaultMessage = "This PKI mount hasn't yet been configured with a certificate issuer.";
    let cliMessage = '';

    if (this.args.roles?.length)
      cliMessage =
        'There are existing roles. Use the CLI to perform any operations with them until an issuer is configured.';
    if (this.args.certificates?.length)
      cliMessage =
        'There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured.';
    if (this.args.roles?.length && this.args.certificates?.length)
      cliMessage =
        'There are existing roles and certificates. Use the CLI to perform any operations with them until an issuer is configured.';

    return `${defaultMessage} ${cliMessage}`;
  }
}
