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
 * usage example: {{has-capability this.capabilities "delete" "update" id=item.id all=true}}
 */

import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';

import type { Capabilities, CapabilitiesMap, CapabilityTypes } from 'vault/app-types';

type NamedArgs = {
  id: string;
  all?: boolean;
};

export function hasCapability(
  capabilitiesMap?: CapabilitiesMap,
  types?: CapabilityTypes[],
  id?: string,
  all?: boolean
) {
  if (capabilitiesMap && Array.isArray(types)) {
    const path = Object.keys(capabilitiesMap).find((key) => (id ? key.includes(id) : false));

    if (path) {
      const capabilities = capabilitiesMap[path];

      if (capabilities) {
        const method = all ? 'every' : 'some';

        return types[method]((type) => {
          const key = `can${capitalize(type)}` as keyof Capabilities;
          // default to true if type provided is not valid - edit rather than update for example
          return !(key in capabilities) ? true : capabilities[key];
        });
      }
    }
  }
  // similar to the Capabilities service, default to allow and the API will deny if needed
  return true;
}

export default buildHelper(function ([capabilitiesMap, ...types], { id, all = false }: NamedArgs) {
  return hasCapability(capabilitiesMap as CapabilitiesMap, types as CapabilityTypes[], id, all);
});
