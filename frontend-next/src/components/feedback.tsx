import type { ReactElement } from 'react';

type FeedbackProps = {
  success?: string | null;
  error?: string | null;
};

export function Feedback({ success, error }: FeedbackProps): ReactElement {
  return (
    <>
      {success ? (
        <div className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-success">
          {success}
        </div>
      ) : null}
      {error ? (
        <div className="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      ) : null}
    </>
  );
}
