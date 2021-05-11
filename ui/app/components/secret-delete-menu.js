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

  // requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  // buttonDisabled: or('requestInFlight', 'model.isFolder', 'model.flagsIsInvalid', 'hasLintError', 'error'),

  // modelForData: computed('isV2', 'model', function() {
  //   let { model } = this;
  //   if (!model) return null;
  //   return this.isV2 ? model.belongsTo('selectedVersion').value() : model;
  // }),

  // basicModeDisabled: computed('secretDataIsAdvanced', 'showAdvancedMode', function() {
  //   return this.secretDataIsAdvanced || this.showAdvancedMode === false;
  // }),

  // secretDataAsJSON: computed('secretData', 'secretData.[]', function() {
  //   return this.secretData.toJSON();
  // }),

  // secretDataIsAdvanced: computed('secretData', 'secretData.[]', function() {
  //   return this.secretData.isAdvanced();
  // }),

  // showAdvancedMode: or('secretDataIsAdvanced', 'preferAdvancedEdit'),

  // isWriteWithoutRead: computed('model.failedServerRead', 'modelForData.failedServerRead', 'isV2', function() {
  //   if (!this.model) return;
  //   // if the version couldn't be read from the server
  //   if (this.isV2 && this.modelForData.failedServerRead) {
  //     return true;
  //   }
  //   // if the model couldn't be read from the server
  //   if (!this.isV2 && this.model.failedServerRead) {
  //     return true;
  //   }
  //   return false;
  // }),

  // transitionToRoute() {
  //   return this.router.transitionTo(...arguments);
  // },

  // onEscape(e) {
  //   if (e.keyCode !== keys.ESC || this.mode !== 'show') {
  //     return;
  //   }
  //   const parentKey = this.model.parentKey;
  //   if (parentKey) {
  //     this.transitionToRoute(LIST_ROUTE, parentKey);
  //   } else {
  //     this.transitionToRoute(LIST_ROOT_ROUTE);
  //   }
  // },

  // // successCallback is called in the context of the component
  // persistKey(successCallback) {
  //   let secret = this.model;
  //   let secretData = this.modelForData;
  //   let isV2 = this.isV2;
  //   let key = secretData.get('path') || secret.id;

  //   if (key.startsWith('/')) {
  //     key = key.replace(/^\/+/g, '');
  //     secretData.set(secretData.pathAttr, key);
  //   }

  //   if (this.mode === 'create') {
  //     key = JSON.stringify({
  //       backend: secret.backend,
  //       id: key,
  //     });
  //   }

  //   return secretData
  //     .save()
  //     .then(() => {
  //       if (!secretData.isError) {
  //         if (isV2) {
  //           secret.set('id', key);
  //         }
  //         if (isV2 && Object.keys(secret.changedAttributes()).length) {
  //           // save secret metadata
  //           secret
  //             .save()
  //             .then(() => {
  //               this.saveComplete(successCallback, key);
  //             })
  //             .catch(e => {
  //               this.set(e, e.errors.join(' '));
  //             });
  //         } else {
  //           this.saveComplete(successCallback, key);
  //         }
  //       }
  //     })
  //     .catch(error => {
  //       if (error instanceof ControlGroupError) {
  //         let errorMessage = this.controlGroup.logFromError(error);
  //         this.set('error', errorMessage.content);
  //       }
  //       throw error;
  //     });
  // },
  // saveComplete(callback, key) {
  //   if (this.wizard.featureState === 'secret') {
  //     this.wizard.transitionFeatureMachine('secret', 'CONTINUE');
  //   }
  //   callback(key);
  // },

  // checkRows() {
  //   if (this.secretData.length === 0) {
  //     this.send('addRow');
  //   }
  // },

  actions: {
    //submit on shift + enter
    // handleKeyDown(e) {
    //   e.stopPropagation();
    //   if (!(e.keyCode === keys.ENTER && e.metaKey)) {
    //     return;
    //   }
    //   let $form = this.element.querySelector('form');
    //   if ($form.length) {
    //     $form.submit();
    //   }
    // },

    // handleChange() {
    //   this.set('codemirrorString', this.secretData.toJSONString(true));
    //   set(this.modelForData, 'secretData', this.secretData.toJSON());
    // },

    // handleWrapClick() {
    //   this.set('isWrapping', true);
    //   if (this.isV2) {
    //     this.store
    //       .adapterFor('secret-v2-version')
    //       .queryRecord(this.modelForData.id, { wrapTTL: 1800 })
    //       .then(resp => {
    //         this.set('wrappedData', resp.wrap_info.token);
    //         this.flashMessages.success('Secret Successfully Wrapped!');
    //       })
    //       .catch(() => {
    //         this.flashMessages.danger('Could Not Wrap Secret');
    //       })
    //       .finally(() => {
    //         this.set('isWrapping', false);
    //       });
    //   } else {
    //     this.store
    //       .adapterFor('secret')
    //       .queryRecord(null, null, { backend: this.model.backend, id: this.modelForData.id, wrapTTL: 1800 })
    //       .then(resp => {
    //         this.set('wrappedData', resp.wrap_info.token);
    //         this.flashMessages.success('Secret Successfully Wrapped!');
    //       })
    //       .catch(() => {
    //         this.flashMessages.danger('Could Not Wrap Secret');
    //       })
    //       .finally(() => {
    //         this.set('isWrapping', false);
    //       });
    //   }
    // },

    // clearWrappedData() {
    //   this.set('wrappedData', null);
    // },

    // handleCopySuccess() {
    //   this.flashMessages.success('Copied Wrapped Data!');
    //   this.send('clearWrappedData');
    // },

    // handleCopyError() {
    //   this.flashMessages.danger('Could Not Copy Wrapped Data');
    //   this.send('clearWrappedData');
    // },

    // createOrUpdateKey(type, event) {
    //   event.preventDefault();
    //   const MAXIMUM_VERSIONS = 9999999999999999;
    //   let model = this.modelForData;
    //   let secret = this.model;
    //   // prevent from submitting if there's no key
    //   if (type === 'create' && isBlank(model.path || model.id)) {
    //     this.flashMessages.danger('Please provide a path for the secret');
    //     return;
    //   }
    //   const maxVersions = secret.get('maxVersions');
    //   if (MAXIMUM_VERSIONS < maxVersions) {
    //     this.flashMessages.danger('Max versions is too large');
    //     return;
    //   }

    //   this.persistKey(key => {
    //     let secretKey;
    //     try {
    //       secretKey = JSON.parse(key).id;
    //     } catch (error) {
    //       secretKey = key;
    //     }
    //     this.transitionToRoute(SHOW_ROUTE, secretKey);
    //   });
    // },

    // deleteKey() {
    //   let { id } = this.model;
    //   this.model.destroyRecord().then(() => {
    //     this.navToNearestAncestor.perform(id);
    //   });
    // },

    // refresh() {
    //   this.onRefresh();
    // },

    // addRow() {
    //   const data = this.secretData;
    //   if (isNone(data.findBy('name', ''))) {
    //     data.pushObject({ name: '', value: '' });
    //     this.send('handleChange');
    //   }
    //   this.checkRows();
    // },

    // deleteRow(name) {
    //   const data = this.secretData;
    //   const item = data.findBy('name', name);
    //   if (isBlank(item.name)) {
    //     return;
    //   }
    //   data.removeObject(item);
    //   this.checkRows();
    //   this.send('handleChange');
    // },

    // toggleAdvanced(bool) {
    //   this.onToggleAdvancedEdit(bool);
    // },

    // codemirrorUpdated(val, codemirror) {
    //   this.set('error', null);
    //   codemirror.performLint();
    //   const noErrors = codemirror.state.lint.marked.length === 0;
    //   if (noErrors) {
    //     try {
    //       this.secretData.fromJSONString(val);
    //       set(this.modelForData, 'secretData', this.secretData.toJSON());
    //     } catch (e) {
    //       this.set('error', e.message);
    //     }
    //   }
    //   this.set('hasLintError', !noErrors);
    //   this.set('codemirrorString', val);
    // },

    // formatJSON() {
    //   this.set('codemirrorString', this.secretData.toJSONString(true));
    // },

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
