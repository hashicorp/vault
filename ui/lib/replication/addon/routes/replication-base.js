import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModelRouteMixin from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModelRouteMixin, {
  store: service(),
  version: service(),
  rm: service('replication-mode'),
  modelPath: 'model.config', // TODO (unload mixin): when removing mixin, remove prepended 'model'
  replicationMode: alias('rm.mode'),
});
