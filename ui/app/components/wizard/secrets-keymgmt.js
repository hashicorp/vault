/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

export default class WizardSecretsKeymgmtComponent extends Component {
  get headerText() {
    return {
      provider: 'Creating a provider',
      displayProvider: 'Distributing a key',
      distribute: 'Creating a key',
    }[this.args.featureState];
  }

  get body() {
    return {
      provider: 'This process connects an external provider to Vault. You will need its credentials.',
      displayProvider: 'A key can now be created and distributed to this destination.',
      distribute: 'This process creates a key and distributes it to your provider.',
    }[this.args.featureState];
  }

  get instructions() {
    return {
      provider: 'Enter your provider details and click “Create provider“.',
      displayProvider: 'Click “Distribute key” in the toolbar.',
      distribute: 'Enter your key details and click “Distribute key”.',
    }[this.args.featureState];
  }
}
