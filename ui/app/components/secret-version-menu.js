import { maybeQueryRecord } from 'vault/macros/maybe-query-record';
import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { alias, or } from '@ember/object/computed';

export default Component.extend({
  tagName: '',
  store: service(),
  version: null,
  useDefaultTrigger: false,

  deleteVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      let [backend, id] = JSON.parse(context.version.id);
      return {
        id: `${backend}/delete/${id}`,
      };
    },
    'version.id'
  ),
  canDeleteVersion: alias('deleteVersionPath.canUpdate'),
  destroyVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      let [backend, id] = JSON.parse(context.version.id);
      return {
        id: `${backend}/destroy/${id}`,
      };
    },
    'version.id'
  ),
  canDestroyVersion: alias('destroyVersionPath.canUpdate'),
  undeleteVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      let [backend, id] = JSON.parse(context.version.id);
      return {
        id: `${backend}/undelete/${id}`,
      };
    },
    'version.id'
  ),
  canUndeleteVersion: alias('undeleteVersionPath.canUpdate'),

  isFetchingVersionCapabilities: or(
    'deleteVersionPath.isPending',
    'destroyVersionPath.isPending',
    'undeleteVersionPath.isPending'
  ),
  actions: {
    deleteVersion(deleteType = 'destroy') {
      return this.store
        .adapterFor('secret-v2-version')
        .v2DeleteOperation(this.store, this.version.id, deleteType);
    },
  },
});
