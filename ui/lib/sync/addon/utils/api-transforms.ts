/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findDestination } from 'core/helpers/sync-destinations';

import type { DestinationType, ListDestination } from 'vault/sync';
import type { SystemListSyncDestinationsResponse } from '@hashicorp/vault-client-typescript';

// transforms the systemListSyncDestinations response to a flat array
// destination name and type are the only properties returned from the request
// to satisfy the list views, combine this data with static properties icon and typeDisplayName
// the flat array is then filtered by name and type if filter values are provided
export const listDestinationsTransform = (
  response: SystemListSyncDestinationsResponse,
  nameFilter?: string,
  typeFilter?: string
) => {
  const { keyInfo } = response;
  const destinations: ListDestination[] = [];
  // build ListDestination objects from keyInfo
  for (const key in keyInfo) {
    // iterate through each type's destination names
    const names = (keyInfo as Record<string, string[]>)[key];
    // remove trailing slash from key
    const type = key.replace(/\/$/, '') as DestinationType;

    names?.forEach((name: string) => {
      const id = `${type}/${name}`;
      const { icon, name: typeDisplayName } = findDestination(type);
      // create object with destination's id and attributes
      destinations.push({ id, name, type, icon, typeDisplayName });
    });
  }

  // optionally filter by name and type
  let filteredDestinations = [...destinations];

  const filter = (key: 'type' | 'name', value: string) => {
    return filteredDestinations.filter((destination: ListDestination) => {
      return destination[key].toLowerCase().includes(value.toLowerCase());
    });
  };
  if (typeFilter) {
    filteredDestinations = filter('type', typeFilter);
  }
  if (nameFilter) {
    filteredDestinations = filter('name', nameFilter);
  }

  return filteredDestinations;
};
