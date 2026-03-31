/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

//list from http://docs.aws.amazon.com/general/latest/gr/rande.html#sts_region
const REGIONS = [
  {
    group: 'Africa',
    options: [{ value: 'af-south-1', displayName: 'Cape Town (af-south-1)' }],
  },
  {
    group: 'Asia Pacific',
    options: [
      { value: 'ap-east-1', displayName: 'Hong Kong (ap-east-1)' },
      { value: 'ap-east-2', displayName: 'Taipei (ap-east-2)' },
      { value: 'ap-northeast-1', displayName: 'Tokyo (ap-northeast-1)' },
      { value: 'ap-northeast-2', displayName: 'Seoul (ap-northeast-2)' },
      { value: 'ap-northeast-3', displayName: 'Osaka (ap-northeast-3)' },
      { value: 'ap-south-1', displayName: 'Mumbai (ap-south-1)' },
      { value: 'ap-south-2', displayName: 'Hyderabad (ap-south-2)' },
      { value: 'ap-southeast-1', displayName: 'Singapore (ap-southeast-1)' },
      { value: 'ap-southeast-2', displayName: 'Sydney (ap-southeast-2)' },
      { value: 'ap-southeast-3', displayName: 'Jakarta (ap-southeast-3)' },
      { value: 'ap-southeast-4', displayName: 'Melbourne (ap-southeast-4)' },
      { value: 'ap-southeast-5', displayName: 'Malaysia (ap-southeast-5)' },
      { value: 'ap-southeast-6', displayName: 'New Zealand (ap-southeast-6)' },
      { value: 'ap-southeast-7', displayName: 'Thailand (ap-southeast-7)' },
    ],
  },
  {
    group: 'Canada',
    options: [
      { value: 'ca-central-1', displayName: 'Central (ca-central-1)' },
      { value: 'ca-west-1', displayName: 'Calgary (ca-west-1)' },
    ],
  },
  {
    group: 'Europe',
    options: [
      { value: 'eu-central-1', displayName: 'Frankfurt (eu-central-1)' },
      { value: 'eu-central-2', displayName: 'Zurich (eu-central-2)' },
      { value: 'eu-north-1', displayName: 'Stockholm (eu-north-1)' },
      { value: 'eu-south-1', displayName: 'Milan (eu-south-1)' },
      { value: 'eu-south-2', displayName: 'Spain (eu-south-2)' },
      { value: 'eu-west-1', displayName: 'Ireland (eu-west-1)' },
      { value: 'eu-west-2', displayName: 'London (eu-west-2)' },
      { value: 'eu-west-3', displayName: 'Paris (eu-west-3)' },
    ],
  },
  {
    group: 'Israel',
    options: [{ value: 'il-central-1', displayName: 'Tel Aviv (il-central-1)' }],
  },
  {
    group: 'Mexico',
    options: [{ value: 'mx-central-1', displayName: 'Central (mx-central-1)' }],
  },
  {
    group: 'Middle East',
    options: [
      { value: 'me-central-1', displayName: 'UAE (me-central-1)' },
      { value: 'me-south-1', displayName: 'Bahrain (me-south-1)' },
    ],
  },
  {
    group: 'South America',
    options: [{ value: 'sa-east-1', displayName: 'São Paulo (sa-east-1)' }],
  },
  {
    group: 'US East',
    options: [
      { value: 'us-east-1', displayName: 'N. Virginia (us-east-1)' },
      { value: 'us-east-2', displayName: 'Ohio (us-east-2)' },
    ],
  },
  {
    group: 'US West',
    options: [
      { value: 'us-west-1', displayName: 'N. California (us-west-1)' },
      { value: 'us-west-2', displayName: 'Oregon (us-west-2)' },
    ],
  },
];

export function regions() {
  return REGIONS.map((regionGroup) => ({
    ...regionGroup,
    options: [...regionGroup.options],
  }));
}

export default buildHelper(regions);
