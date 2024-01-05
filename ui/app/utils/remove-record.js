/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Unlinks a record from all its relationships and unloads it from
// the store.
export default function removeRecord(store, record) {
  const id = record.path || record.id;
  if (id) {
    // Collect relationship property names and types
    const relationshipMeta = [];
    record.eachRelationship((key, { kind }) => {
      relationshipMeta.push({ key, kind });
    });

    // Push an update to this record with the relationships nulled out.
    // This unlinks the relationship from the models that aren't about to
    // be unloaded.
    store.push({
      data: {
        id,
        type: record.constructor.modelName,
        relationships: relationshipMeta.reduce((hash, rel) => {
          hash[rel.key] = { data: rel.kind === 'hasMany' ? [] : null };
          return hash;
        }, {}),
      },
    });
  }

  // Now that the record has no attachments, it can be safely unloaded
  // from the store.
  store.unloadRecord(record);
}
