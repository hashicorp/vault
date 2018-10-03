CRDB
====

`crdb` is a wrapper around the logic for issuing SQL transactions which performs
retries (as required by CockroachDB).

Note that unfortunately there is no generic way of extracting a pg error code;
the library has to recognize driver-dependent error types. We currently support
`github.com/lib/pq` and `github.com/jackc/pgx`.
