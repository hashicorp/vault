/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * Util to add/remove models in a hasMany relationship via the search-select component
 *
 * @example of using the util within an action in a component
 *```js
 * @action
 * async onSearchSelectChange(selectedIds) {
 *    const methods = await this.args.model.mfa_methods;
 *   handleHasManySelection(selectedIds, methods, this.store, 'mfa-method');
 *  }
 *
 * @param selectedIds array of selected options from search-select component
 * @param modelCollection array-like, list of models from the hasMany relationship
 * @param store the store so we can call peekRecord()
 * @param modelRecord string passed to peekRecord
 */

export default function handleHasManySelection(selectedIds, modelCollection, store, modelRecord) {
  // first check for existing models that have been removed from selection
  modelCollection.forEach((model) => {
    if (!selectedIds.includes(model.id)) {
      modelCollection.removeObject(model);
    }
  });
  // now check for selected items that don't exist and add them to the model
  const modelIds = modelCollection.mapBy('id');
  selectedIds.forEach((id) => {
    if (!modelIds.includes(id)) {
      const model = store.peekRecord(modelRecord, id);
      modelCollection.addObject(model);
    }
  });
}
