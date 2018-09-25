import { helper as buildHelper } from '@ember/component/helper';

//list from http://docs.aws.amazon.com/general/latest/gr/rande.html#sts_region
const REGIONS = [
  'us-east-1',
  'us-east-2',
  'us-west-1',
  'us-west-2',
  'ca-central-1',
  'ap-south-1',
  'ap-northeast-1',
  'ap-northeast-2',
  'ap-southeast-1',
  'ap-southeast-2',
  'eu-central-1',
  'eu-west-1',
  'eu-west-2',
  'sa-east-1',
];

export function regions() {
  return REGIONS.slice(0);
}

export default buildHelper(regions);
