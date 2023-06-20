/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

export default class DashboardSecretsEnginesCardComponent extends Component {
  get filteredSecretsEngines() {
    const filteredEngines = this.args.secretsEngines.filter(
      (secretEngine) => secretEngine.shouldIncludeInList
    );

    return filteredEngines.slice(0, 5);
  }
}
