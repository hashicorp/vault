/* eslint-disable no-undef */
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module DiffVersionSelector
 * DiffVersionSelector component is a specific component that has a toolbar for selecting KV 2 versions and showing a diff between those versions.
 *
 * @example
 * ```js
 * <DiffVersionSelector @model={model}/>
 * ```
 * @param {object} model - model of formed from secret-v2-version
 */

export default class DiffVersionSelector extends Component {
  @tracked leftSideVersionDataSelected = null;
  @tracked leftSideVersionSelected = null;
  @tracked rightSideVersionDataSelected = null;
  @tracked rightSideVersionSelected = null;
  @tracked statesMatch = false;
  @tracked visualDiff = null;
  @service store;

  adapter = this.store.adapterFor('secret-v2-version');

  constructor() {
    super(...arguments);
    this.createVisualDiff();
  }

  get leftSideDataInit() {
    // return secretData from hitting the get secret endpoint
    let string = `["${this.args.model.engineId}", "${this.args.model.id}", "${this.args.model.currentVersion}"]`;
    return this.adapter
      .querySecretDataByVersion(string)
      .then(response => response.data) // using ember promise helpers to await in the hbs file
      .catch(() => null);
  }
  get rightSideDataInit() {
    // return secretData from hitting the get secret endpoint
    let string = `["${this.args.model.engineId}", "${this.args.model.id}", "${this.rightSideVersionInit}"]`;
    return this.adapter
      .querySecretDataByVersion(string)
      .then(response => response.data) // using ember promise helpers to await in the hbs file\
      .catch(() => null);
  }
  get rightSideVersionInit() {
    // initial value of right side version is one less than the current version
    return this.args.model.currentVersion === 1 ? 0 : this.args.model.currentVersion - 1;
  }

  async createVisualDiff() {
    let diffpatcher = jsondiffpatch.create({});
    let leftSideVersionData = this.leftSideVersionDataSelected || (await this.leftSideDataInit);
    let rightSideVersionData = this.rightSideVersionDataSelected || (await this.rightSideDataInit);
    let delta = diffpatcher.diff(rightSideVersionData, leftSideVersionData);
    if (delta === undefined) {
      this.statesMatch = true;
      this.visualDiff = JSON.stringify(leftSideVersionData, undefined, 2); // params: value, replacer (all properties included), space (white space and indentation, line break, etc.)
    } else {
      this.statesMatch = false;
      this.visualDiff = jsondiffpatch.formatters.html.format(delta, rightSideVersionData);
    }
  }

  @action
  async selectVersion(selectedVersion, actions, side) {
    let string = `["${this.args.model.engineId}", "${this.args.model.id}", "${selectedVersion}"]`;
    let secretData = await this.adapter.querySecretDataByVersion(string);
    if (side === 'left') {
      this.leftSideVersionDataSelected = secretData.data;
      this.leftSideVersionSelected = selectedVersion;
    }
    if (side === 'right') {
      this.rightSideVersionDataSelected = secretData.data;
      this.rightSideVersionSelected = selectedVersion;
    }
    await this.createVisualDiff();
    // close dropdown menu.
    actions.close();
  }
}
