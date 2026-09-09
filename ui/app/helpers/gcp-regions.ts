/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

interface GcpRegionOption {
  value: string;
  displayName: string;
}

interface GcpRegionGroup {
  group: string;
  options: GcpRegionOption[];
}

//list from https://docs.cloud.google.com/secret-manager/docs/locations
const REGIONS: GcpRegionGroup[] = [
  {
    group: 'Africa',
    options: [{ value: 'africa-south1', displayName: 'Johannesburg (africa-south1)' }],
  },
  {
    group: 'Asia Pacific',
    options: [
      { value: 'asia-east1', displayName: 'Taiwan (asia-east1)' },
      { value: 'asia-east2', displayName: 'Hong Kong (asia-east2)' },
      { value: 'asia-northeast1', displayName: 'Tokyo (asia-northeast1)' },
      { value: 'asia-northeast2', displayName: 'Osaka (asia-northeast2)' },
      { value: 'asia-northeast3', displayName: 'Seoul (asia-northeast3)' },
      { value: 'asia-south1', displayName: 'Mumbai (asia-south1)' },
      { value: 'asia-south2', displayName: 'Delhi (asia-south2)' },
      { value: 'asia-southeast1', displayName: 'Singapore (asia-southeast1)' },
      { value: 'asia-southeast2', displayName: 'Jakarta (asia-southeast2)' },
      { value: 'asia-southeast3', displayName: 'Bangkok (asia-southeast3)' },
      { value: 'australia-southeast1', displayName: 'Sydney (australia-southeast1)' },
      { value: 'australia-southeast2', displayName: 'Melbourne (australia-southeast2)' },
    ],
  },
  {
    group: 'Europe',
    options: [
      { value: 'europe-central2', displayName: 'Warsaw (europe-central2)' },
      { value: 'europe-north1', displayName: 'Finland (europe-north1)' },
      { value: 'europe-north2', displayName: 'Stockholm (europe-north2)' },
      { value: 'europe-southwest1', displayName: 'Madrid (europe-southwest1)' },
      { value: 'europe-west1', displayName: 'Belgium (europe-west1)' },
      { value: 'europe-west2', displayName: 'London (europe-west2)' },
      { value: 'europe-west3', displayName: 'Frankfurt (europe-west3)' },
      { value: 'europe-west4', displayName: 'Netherlands (europe-west4)' },
      { value: 'europe-west6', displayName: 'Zurich (europe-west6)' },
      { value: 'europe-west8', displayName: 'Milan (europe-west8)' },
      { value: 'europe-west9', displayName: 'Paris (europe-west9)' },
      { value: 'europe-west10', displayName: 'Berlin (europe-west10)' },
      { value: 'europe-west12', displayName: 'Turin (europe-west12)' },
    ],
  },
  {
    group: 'Middle East',
    options: [
      { value: 'me-central1', displayName: 'Doha (me-central1)' },
      { value: 'me-central2', displayName: 'Dammam (me-central2)' },
      { value: 'me-west1', displayName: 'Tel Aviv (me-west1)' },
    ],
  },
  {
    group: 'North America',
    options: [
      { value: 'northamerica-northeast1', displayName: 'Montréal (northamerica-northeast1)' },
      { value: 'northamerica-northeast2', displayName: 'Toronto (northamerica-northeast2)' },
      { value: 'northamerica-south1', displayName: 'Mexico (northamerica-south1)' },
      { value: 'us-central1', displayName: 'Iowa (us-central1)' },
      { value: 'us-east1', displayName: 'South Carolina (us-east1)' },
      { value: 'us-east4', displayName: 'Northern Virginia (us-east4)' },
      { value: 'us-east5', displayName: 'Columbus (us-east5)' },
      { value: 'us-south1', displayName: 'Dallas (us-south1)' },
      { value: 'us-west1', displayName: 'Oregon (us-west1)' },
      { value: 'us-west2', displayName: 'Los Angeles (us-west2)' },
      { value: 'us-west3', displayName: 'Salt Lake City (us-west3)' },
      { value: 'us-west4', displayName: 'Las Vegas (us-west4)' },
    ],
  },
  {
    group: 'South America',
    options: [
      { value: 'southamerica-east1', displayName: 'São Paulo (southamerica-east1)' },
      { value: 'southamerica-west1', displayName: 'Santiago (southamerica-west1)' },
    ],
  },
];

export function gcpRegions(): GcpRegionGroup[] {
  return REGIONS.map((regionGroup) => ({
    ...regionGroup,
    options: [...regionGroup.options],
  }));
}

export default buildHelper(gcpRegions);
