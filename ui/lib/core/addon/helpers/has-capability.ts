/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * helper to get capabilities from map for given item by id
 * useful in list views when capabilities are fetched for each item
 * provide the full map of capabilities along with the id of the item and capabilities to be evaluated (read, update etc)
 * multiple capabilities will return true if any of the capabilities are true by default
 * to check if all provided capability types are true, set the all argument to true
 * usage example: {{has-capability this.capabilities "delete" "update" pathKey="customMessages" params=message all=true}}
 */

import Helper from '@ember/component/helper';
import { service } from '@ember/service';
import { capitalize } from '@ember/string';

import type { Capabilities, CapabilitiesMap, CapabilityTypes } from 'vault/app-types';
import type CapabilitiesService from 'vault/services/capabilities';
import type { PATH_MAP } from 'vault/utils/constants/capabilities';

export default class HasCapabilityHelper extends Helper {
  @service declare readonly capabilities: CapabilitiesService;

  compute<T>(
    [capabilitiesMap, ...types]: [CapabilitiesMap, ...CapabilityTypes[]],
    { pathKey, params, all = false }: { pathKey: keyof typeof PATH_MAP; params?: T; all?: boolean }
  ) {
    // since the default is to return true we need to validate the inputs here more thoroughly
    // this is to help devs from thinking that the checks are working properly when the capability lookup is actually failing
    if (!capabilitiesMap) {
      throw new Error('First positional argument must be the capabilities map.');
    }
    if (!types || !types.length) {
      throw new Error('At least one capability type is required as a positional argument.');
    }
    const acceptedTypes = ['read', 'update', 'delete', 'list', 'create', 'patch', 'sudo'];
    const invalidTypes = types.filter((type) => !acceptedTypes.includes(type));
    if (invalidTypes.length) {
      throw new Error(
        `Invalid capability types: ${invalidTypes.join(', ')}. Accepted types are: ${acceptedTypes.join(
          ', '
        )}.`
      );
    }
    if (!pathKey) {
      throw new Error('pathKey is a required named arg for path lookup in capabilities map');
    }

    const path = this.capabilities.pathFor(pathKey, params);
    const capabilities = capabilitiesMap[path];

    if (capabilities) {
      const method = all ? 'every' : 'some';

      return types[method]((type) => {
        const key = `can${capitalize(type)}` as keyof Capabilities;
        return capabilities[key];
      });
    }

    // similar to the Capabilities service, default to allow and the API will deny if needed
    return true;
  }
}
