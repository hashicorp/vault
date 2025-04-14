/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

/**
 * ARG TODO
 */

export default class SecretsMountDetails extends Route {
  @service store;

  async model(params) {
    const response = await this.store.query('secret-engine', {
      path: params.mount_name,
    });
    return response[0];

    // if (secretEngineModel.isV2KV) {
    //   const canRead = await this.store
    //     .findRecord('capabilities', `${secretEngineModel.id}/config`)
    //     .then((response) => response.canRead);
    //   // only set these config params if they can read the config endpoint.
    //   if (canRead) {
    //     // design wants specific default to show that can't be set in the model
    //     secretEngineModel.casRequired = secretEngineModel.casRequired
    //       ? secretEngineModel.casRequired
    //       : 'False';
    //     secretEngineModel.deleteVersionAfter = secretEngineModel.deleteVersionAfter
    //       ? secretEngineModel.deleteVersionAfter
    //       : 'Never delete';
    //   } else {
    //     // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
    //     secretEngineModel.set('casRequired', null);
    //     secretEngineModel.set('deleteVersionAfter', null);
    //     secretEngineModel.set('maxVersions', null);
    //   }
    // }

    // return { secretEngineModel };
  }
}
