/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { generateCsv, type CsvColumn } from 'vault/utils/generate-csv';

interface TestRow {
  name?: string;
  status?: string;
  notes?: string;
}

interface GroupedRow {
  name: string;
  id: string;
  age: string;
  type: string;
}

module('Unit | Util | generate-csv', function () {
  test('it builds CSV with headers and escaped values', function (assert) {
    const columns: CsvColumn<TestRow>[] = [
      { header: 'Name', key: 'name' },
      { header: 'Status', key: 'status' },
      { header: 'Notes', key: 'notes' },
    ];

    const rows: TestRow[] = [
      { name: 'build, bot', status: 'Enabled', notes: 'says "hello"' },
      { name: 'deploy-bot', status: 'Disabled', notes: 'line 1\nline 2' },
    ];

    const csv = generateCsv({ rows, columns });

    assert.strictEqual(
      csv,
      'Name,Status,Notes\n"build, bot",Enabled,"says ""hello"""\ndeploy-bot,Disabled,"line 1\nline 2"',
      'CSV includes escaped commas, quotes, and newlines'
    );
  });

  test('it renders headers when rows are empty', function (assert) {
    const columns: CsvColumn<TestRow>[] = [
      { header: 'Name', key: 'name' },
      { header: 'Status', key: 'status' },
    ];

    const csv = generateCsv({ rows: [], columns });

    assert.strictEqual(csv, 'Name,Status', 'CSV contains only headers with no data rows');
  });

  test('it converts nullish values to empty strings', function (assert) {
    const columns: CsvColumn<TestRow>[] = [
      { header: 'Name', key: 'name' },
      { header: 'Status', key: 'status' },
      { header: 'Notes', key: 'notes' },
    ];

    const csv = generateCsv({
      rows: [{ name: undefined, status: undefined, notes: undefined }],
      columns,
    });

    assert.strictEqual(csv, 'Name,Status,Notes\n,,', 'Nullish values are emitted as empty CSV cells');
  });

  test('it supports callback-based column values', function (assert) {
    const columns: CsvColumn<TestRow>[] = [
      { header: 'Upper name', value: ({ row }) => row.name?.toUpperCase() },
    ];

    const csv = generateCsv({ rows: [{ name: 'build-bot' }], columns });

    assert.strictEqual(csv, 'Upper name\nBUILD-BOT', 'Callback columns remain supported for computed values');
  });

  test('it preserves grouped row ordering for parent and child rows', function (assert) {
    const columns: CsvColumn<GroupedRow>[] = [
      { header: 'Name', key: 'name' },
      { header: 'ID', key: 'id' },
      { header: 'Age', key: 'age' },
      { header: 'Type', key: 'type' },
    ];

    const rows: GroupedRow[] = [
      {
        name: 'Pokemon',
        id: 'pokemon-12345',
        age: '10',
        type: 'Electric',
      },
      {
        name: '',
        id: 'pichu-2222',
        age: '10',
        type: 'Electric',
      },
      {
        name: '',
        id: 'raichu-67890',
        age: '20',
        type: 'Electric',
      },
      {
        name: '',
        id: 'pikachu-9033',
        age: '30',
        type: 'Electric',
      },
    ];

    const csv = generateCsv({ rows, columns });

    assert.strictEqual(
      csv,
      'Name,ID,Age,Type\nPokemon,pokemon-12345,10,Electric\n,pichu-2222,10,Electric\n,raichu-67890,20,Electric\n,pikachu-9033,30,Electric',
      'CSV preserves the flattened parent row followed by its child alias row'
    );
  });
});
