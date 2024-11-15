[![](https://godoc.org/github.com/jackc/pgtype?status.svg)](https://godoc.org/github.com/jackc/pgtype)
![CI](https://github.com/jackc/pgtype/workflows/CI/badge.svg)

---

This version is used with pgx `v4`. In pgx `v5` it is part of the https://github.com/jackc/pgx repository. Only bug
fixes will be made to this version. New features will only be considered for the current release.

---

# pgtype

pgtype implements Go types for over 70 PostgreSQL types. pgtype is the type system underlying the
https://github.com/jackc/pgx PostgreSQL driver. These types support the binary format for enhanced performance with pgx.
They also support the database/sql `Scan` and `Value` interfaces and can be used with https://github.com/lib/pq.
