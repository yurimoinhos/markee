import type { ReactElement } from 'react';

type DataTableProps<T> = {
  columns: { key: string; label: string; render?: (row: T) => ReactElement | string | number | null }[];
  rows: T[];
  emptyText?: string;
};

export function DataTable<T extends Record<string, unknown>>({
  columns,
  rows,
  emptyText = 'Sem dados',
}: DataTableProps<T>): ReactElement {
  return (
    <div className="overflow-x-auto rounded-2xl border border-slate-200 bg-white shadow-sm">
      <table className="min-w-full text-sm">
        <thead className="bg-slate-50 text-left text-slate-700">
          <tr>
            {columns.map((column) => (
              <th key={column.key} className="px-4 py-3 font-semibold">
                {column.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.length === 0 ? (
            <tr>
              <td className="px-4 py-5 text-slate-500" colSpan={columns.length}>
                {emptyText}
              </td>
            </tr>
          ) : (
            rows.map((row, index) => (
              <tr key={`${String(row.id ?? index)}`} className="border-t border-slate-100">
                {columns.map((column) => (
                  <td key={column.key} className="px-4 py-3 align-top text-slate-800">
                    {column.render
                      ? column.render(row)
                      : typeof row[column.key] === 'string' || typeof row[column.key] === 'number'
                        ? String(row[column.key])
                        : '-'}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
