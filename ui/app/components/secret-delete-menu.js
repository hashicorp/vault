import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';

export default Component.extend(FocusOnInsertMixin, WithNavToNearestAncestor, {
  tagName: '',
  router: service(),
  store: service(),

  showDeleteModal: false,

  deleteVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.modelForData) return;
      let [backend, id] = JSON.parse(context.modelForData.id);
      return {
        id: `${backend}/delete/${id}`,
      };
    },
    'model.id'
  ),
  canDeleteAnyVersion: alias('deleteVersionPath.canUpdate'),

  undeleteVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.modelForData) return;
      if (!context.modelForData.id) return;
      let [backend, id] = JSON.parse(context.modelForData.id);
      return {
        id: `${backend}/undelete/${id}`,
      };
    },
    'model.id'
  ),
  canUndeleteVersion: alias('undeleteVersionPath.canUpdate'),

  destroyVersionPath: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.modelForData || !context.modelForData.id) return;

      let [backend, id] = JSON.parse(context.modelForData.id);
      return {
        id: `${backend}/destroy/${id}`,
      };
    },
    'model.id'
  ),
  canDestroyVersion: alias('destroyVersionPath.canUpdate'),

  v2UpdatePath: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.model || context.mode === 'create') {
        return;
      }
      let backend = context.get('model.engine.id');
      let id = context.model.id;
      return {
        id: `${backend}/metadata/${id}`,
      };
    },
    'model',
    'model.id',
    'mode'
  ),

  canDestroyAllVersions: alias('v2UpdatePath.canDelete'),

  isLatestVersion: computed('model.{currentVersion,selectedVersion}', function() {
    let { model } = this;
    if (!model) return;
    let latestVersion = model.currentVersion;
    let selectedVersion = model.selectedVersion.version;
    if (latestVersion !== selectedVersion) {
      return false;
    }
    return true;
  }),

  actions: {
    handleDelete(deleteType) {
      // deleteType should be 'delete', 'destroy', 'undelete', 'delete-latest-version', 'destroy-all-versions'
      if (!deleteType) {
        return;
      }
      if (deleteType === 'destroy-all-versions') {
        let { id } = this.model;
        this.model.destroyRecord().then(() => {
          this.navToNearestAncestor.perform(id);
        });
      } else {
        return this.store
          .adapterFor('secret-v2-version')
          .v2DeleteOperation(this.store, this.modelForData.id, deleteType)
          .then(() => {
            location.reload();
          });
      }
    },
  },
});
