/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { inject as controller } from '@ember/controller';

export default class SidebarNavComponent extends Component {
  @service currentCluster;
  @service console;
  @controller('vault.cluster') clusterController;
}
