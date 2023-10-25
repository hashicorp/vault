/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type SyncDestinationModel from 'vault/models/sync/destination';
interface Args {
  destination: SyncDestinationModel;
}

export default class DestinationDetailsPage extends Component<Args> {
  credentialValue = (value: string) => {
    // if this value is empty, a destination uses globally set environment variables
    if (value) return 'Destination credentials added';
    return 'Using environment variable';
  };
}
