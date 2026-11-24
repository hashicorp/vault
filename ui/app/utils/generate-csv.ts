/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

interface CsvColumnBase {
  header: string;
}

interface CsvFieldColumn<Row> extends CsvColumnBase {
  key: keyof Row;
}

interface CsvValueColumn<Row> extends CsvColumnBase {
  value: (args: CsvValueArgs<Row>) => unknown;
}

export type CsvColumn<Row> = CsvFieldColumn<Row> | CsvValueColumn<Row>;

interface CsvValueArgs<Row> {
  row: Row;
  column: CsvColumn<Row>;
}

interface GenerateCsvArgs<Row> {
  rows: Row[];
  columns: CsvColumn<Row>[];
}

/**
 * Builds a CSV document from typed row data and column definitions.
 */
export function generateCsv<Row>({ rows, columns }: GenerateCsvArgs<Row>): string {
  const headerLine = columns.map((column) => escapeForCsv(column.header)).join(',');
  const rowLines = rows.map((row) => {
    const values = columns.map((column) =>
      'key' in column ? row[column.key] : column.value({ row, column })
    );
    return values.map(escapeForCsv).join(',');
  });

  return [headerLine, ...rowLines].join('\n');
}

function escapeForCsv(value: unknown): string {
  const stringValue = String(value ?? '');

  if (/[",\n\r]/.test(stringValue)) {
    return `"${stringValue.replace(/"/g, '""')}"`;
  }

  return stringValue;
}
