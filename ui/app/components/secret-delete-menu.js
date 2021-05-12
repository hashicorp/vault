import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { alias } from '@ember/object/computed';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';

export default class SecretDeleteMenu extends Component {
  @service store;
  @service router;

  @tracked showDeleteModal = false;

  @maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.args.model) {
        return;
      }
      let backend = context.args.model.backend;
      let id = context.args.model.id;
      let path = context.args.isV2 ? `${backend}/data/${id}` : `${backend}/${id}`;
      return {
        id: path,
      };
    },
    'isV2',
    'model',
    'model.id',
    'mode'
  )
  updatePath;
  @alias('updatePath.canDelete') canDelete;
  @alias('updatePath.canUpdate') canUpdate;

  @maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.args || !context.modelForData || !context.modelForData.id) return;
      let [backend, id] = JSON.parse(context.args.modelForData.id);
      return {
        id: `${backend}/delete/${id}`,
      };
    },
    'model.id'
  )
  deleteVersionPath;
  @alias('deleteVersionPath.canUpdate') canDeleteAnyVersion;

  @maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.args || !context.args.modelForData || !context.args.modelForData.id) return;
      let [backend, id] = JSON.parse(context.args.modelForData.id);
      return {
        id: `${backend}/undelete/${id}`,
      };
    },
    'model.id'
  )
  undeleteVersionPath;
  @alias('undeleteVersionPath.canUpdate') canUndeleteVersion;

  @maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.args || !context.args.modelForData || !context.args.modelForData.id) return;
      let [backend, id] = JSON.parse(context.args.modelForData.id);
      return {
        id: `${backend}/destroy/${id}`,
      };
    },
    'model.id'
  )
  destroyVersionPath;
  @alias('destroyVersionPath.canUpdate') canDestroyVersion;

  @maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.args.model || !context.args.model.engine || !context.args.model.id) return;
      let backend = context.args.model.engine.id;
      let id = context.args.model.id;
      return {
        id: `${backend}/metadata/${id}`,
      };
    },
    'model',
    'model.id',
    'mode'
  )
  v2UpdatePath;
  @alias('v2UpdatePath.canDelete') canDestroyAllVersions;

  get isLatestVersion() {
    let { model } = this.args;
    if (!model) return false;
    let latestVersion = model.currentVersion;
    let selectedVersion = model.selectedVersion.version;
    if (latestVersion !== selectedVersion) {
      return false;
    }
    return true;
  }

  @action toggleDeleteModal() {
    this.showDeleteModal = !this.showDeleteModal;
  }

  @action
  handleDelete(deleteType) {
    // deleteType should be 'delete', 'destroy', 'undelete', 'delete-latest-version', 'destroy-all-versions'
    if (!deleteType) {
      return;
    }
    if (deleteType === 'destroy-all-versions' || deleteType === 'v1') {
      let { id } = this.args.model;
      this.args.model.destroyRecord().then(() => {
        this.args.navToNearestAncestor.perform(id);
      });
    } else {
      return this.store
        .adapterFor('secret-v2-version')
        .v2DeleteOperation(this.store, this.args.modelForData.id, deleteType)
        .then(() => {
          location.reload();
        });
    }
  }
}
