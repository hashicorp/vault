/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  const checklistState = {};

  server.get('sys/config/ui/checklist-state', () => ({
    data: JSON.parse(JSON.stringify(checklistState)),
  }));

  server.post('sys/config/ui/checklist-state', (schema, request) => {
    const body = JSON.parse(request.requestBody);
    const update = body?.checklist_state ?? {};

    // deep-merge the update: preserve existing steps in each checklist
    for (const [checklistId, steps] of Object.entries(update)) {
      checklistState[checklistId] = { ...(checklistState[checklistId] ?? {}), ...steps };
    }

    return { data: JSON.parse(JSON.stringify(checklistState)) };
  });
}
