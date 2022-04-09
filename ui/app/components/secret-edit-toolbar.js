/**
 * @module SecretEditToolbar
 * SecretEditToolbar component is the toolbar component displaying the JSON toggle and the actions like delete in the show mode.
 *
 * @example
 * ```js
 * <SecretEditToolbar
 * @mode={{mode}}
 * @model={{this.model}}
 * @isV2={{isV2}}
 * @isWriteWithoutRead={{isWriteWithoutRead}}
 * @secretDataIsAdvanced={{secretDataIsAdvanced}}
 * @showAdvancedMode={{showAdvancedMode}}
 * @modelForData={{this.modelForData}}
 * @canUpdateSecretData={{canUpdateSecretData}}
 * @codemirrorString={{codemirrorString}}
 * @wrappedData={{wrappedData}}
 * @editActions={{hash
    toggleAdvanced=(action "toggleAdvanced")
    refresh=(action "refresh")
  }}
 * />
 * ```
 
 * @param {string} mode - show, create, edit. The view.
 * @param {object} model - the model passed from the parent secret-edit
 * @param {boolean} isV2 - KV type
 * @param {boolean} isWriteWithoutRead - boolean describing permissions
 * @param {boolean} secretDataIsAdvanced - used to determine if show JSON toggle
 * @param {boolean} showAdvacnedMode - used for JSON toggle
 * @param {object} modelForData - a modified version of the model with secret data
 * @param {boolean} canUpdateSecretData - permissions that show the create new version button or not.
 * @param {string} codemirrorString - used to copy the JSON
 * @param {object} wrappedData - when copy the data it's the token of the secret returned.
 * @param {object} editActions - actions passed from parent to child
 */
/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { not } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class SecretEditToolbar extends Component {
  @service store;
  @service flashMessages;

  @tracked wrappedData = null;
  @tracked isWrapping = false;
  @not('wrappedData') showWrapButton;

  @action
  clearWrappedData() {
    this.wrappedData = null;
  }

  @action
  handleCopyError() {
    this.flashMessages.danger('Could Not Copy Wrapped Data');
    this.send('clearWrappedData');
  }

  @action
  handleCopySuccess() {
    this.flashMessages.success('Copied Wrapped Data!');
    this.send('clearWrappedData');
  }

  @action
  handleWrapClick() {
    this.isWrapping = true;
    if (this.args.isV2) {
      this.store
        .adapterFor('secret-v2-version')
        .queryRecord(this.args.modelForData.id, { wrapTTL: 1800 })
        .then((resp) => {
          this.wrappedData = resp.wrap_info.token;
          this.flashMessages.success('Secret Successfully Wrapped!');
        })
        .catch(() => {
          this.flashMessages.danger('Could Not Wrap Secret');
        })
        .finally(() => {
          this.isWrapping = false;
        });
    } else {
      this.store
        .adapterFor('secret')
        .queryRecord(null, null, {
          backend: this.args.model.backend,
          id: this.args.modelForData.id,
          wrapTTL: 1800,
        })
        .then((resp) => {
          this.wrappedData = resp.wrap_info.token;
          this.flashMessages.success('Secret Successfully Wrapped!');
        })
        .catch(() => {
          this.flashMessages.danger('Could Not Wrap Secret');
        })
        .finally(() => {
          this.isWrapping = false;
        });
    }
  }
}
