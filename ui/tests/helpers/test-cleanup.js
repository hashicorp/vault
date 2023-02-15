export async function clearRecord(store, modelType, id) {
  await store
    .findRecord(modelType, id)
    .then((model) => {
      deleteModelRecord(model);
    })
    // swallow error
    .catch(() => {});
}

export async function clearRecordsFromStore(store, modelType) {
  const records = store.peekAll(modelType);
  await records.forEach((model) => deleteModelRecord(model).catch(() => {}));
}

const deleteModelRecord = async (model) => {
  await model.destroyRecord();
};
