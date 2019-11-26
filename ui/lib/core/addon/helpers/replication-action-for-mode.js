import { helper as buildHelper } from '@ember/component/helper';
import { get } from '@ember/object';
const ACTIONS = {
  performance: {
    primary: ['disable', 'demote', 'recover', 'reindex'],
    secondary: ['disable', 'promote', 'update-primary', 'recover', 'reindex'],
    bootstrapping: ['disable', 'recover', 'reindex'],
  },
  dr: {
    primary: ['disable', 'recover', 'reindex', 'demote'],
    secondary: ['promote'],
    bootstrapping: ['disable', 'recover', 'reindex'],
  },
};

export function replicationActionForMode([replicationMode, clusterMode] /*, hash*/) {
  return get(ACTIONS, `${replicationMode}.${clusterMode}`);
}

export default buildHelper(replicationActionForMode);
