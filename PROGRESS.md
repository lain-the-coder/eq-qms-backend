# PROGRESS — EA QMS Change Control Backend (Go)

**Scope of this file:** what is built, what is next, decisions made in working sessions
that are not recorded in any guardrail document, and open flags. Nothing else — the six
guardrail docs carry the substance and are always attached.

- **Repo:** `github.com/lain-the-coder/ea-qms-backend`
- **Last checkpoint:** 1 — repo scaffold + `001_users.sql`
- **Next task:** checkpoint 2 — `002_change_controls.sql` (not started)

---

## Completed

### ✅ Checkpoint 1 — Scaffold + users migration

**Repo**

- `go mod init github.com/lain-the-coder/ea-qms-backend`
- `sql/schema/` and `sql/queries/` created
- `.gitignore` contains `.env`; `.env` filled; `.env.example` committed with empty values
- Keys in both: `DB_URL`, `PLATFORM`, `JWT_SECRET`
- Local database `ea_qms` created via psql
- Placeholder `main.go` — `ServeMux`, `WelcomeHome` handler, port `:1304`

**`sql/schema/001_users.sql`** — applied, up → `\d` → down → `\dt` → up cycle clean.

Verified against DB Design §3.1 / §5.1 / §6.1 by reading `\d users`:

| Check                                                                  | Result |
| ---------------------------------------------------------------------- | ------ |
| 8 columns, §3.1 order and names                                        | ✅     |
| `TIMESTAMPTZ` on `created_on`, `updated_on`                            | ✅     |
| 4 schema-level defaults: `gen_random_uuid()`, `true`, `now()`, `now()` | ✅     |
| `ck_users_role` — four values, ASCII                                   | ✅     |
| `uq_users_email` **functional** unique index on `lower(email)`         | ✅     |
| `idx_users_role_active` composite `(role, is_active)`                  | ✅     |

---

## Next

### ⬜ Checkpoint 2 — `002_change_controls.sql`

Sources: DB Design **§3.2** (columns, six field groups) · **§4.1** (five FKs, all
`RESTRICT`) · **§5.1 #4–#9** (indexes) · **§6.1** (13 CHECKs) · **§6.2** (7 defaults) ·
**§8.1** (sequence + generated `cc_id`).

Expect **50 physical columns** (see flag #6), 13 CHECKs, 5 FKs, 1 UNIQUE constraint,
6 `CREATE INDEX` statements, 1 sequence.

### ⬜ Remaining migrations

`003_file_attachments` · `004_audit_logs` · `005_esignatures` · `006_refresh_tokens` ·
seed (4 users, one per role, hashed with the app's own argon2id — never pasted hashes,
gated on `PLATFORM=dev`, DB Design §7)

### ⬜ Then

sqlc setup (`emit_pointers_for_null_types: true`) → build in API Endpoint Plan order,
starting `POST /api/login`.

---

## Session decisions not in any guardrail doc

These were settled in working sessions and are binding. They exist nowhere else.

| #   | Decision                                                                                                    | Rationale                                                                                                                                                                                                                                  |
| --- | ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | **Module path `github.com/lain-the-coder/ea-qms-backend`** — not `-cc-backend`                              | Future QMS modules (Deviation, CAPA) live under the same module rather than forcing a second repo                                                                                                                                          |
| 2   | **Constraint naming: `ck_cc_*` and `idx_cc_*`**, per DB §5.1/§6.1 verbatim — _not_ §1.3's long-form example | §5.1/§6.1 are definitions and are cross-referenced by name elsewhere (§8.2 cites `idx_cc_owner_state`); §1.3 is a convention statement with one stale example. Also keeps names clear of Postgres's 63-byte identifier truncation          |
| 3   | **Two naming exceptions kept verbatim:** `uq_change_controls_cc_id` and `ck_cc_post_impl_issues`            | Spelled that way in §3.2/§5.1 and §6.1 respectively — do not "regularize" them                                                                                                                                                             |
| 4   | **FK constraints use the long form** — `fk_change_controls_change_owner_id`                                 | §4 lists all eleven FKs but never names the constraints, so §1.3 stands unopposed for this object type. Decision #2 does not extend to FKs                                                                                                 |
| 5   | **PostgreSQL 14.23 accepted** (doc §1.2 specifies 15+)                                                      | Every feature the schema needs was traced and predates 14: `gen_random_uuid()` core (13), `GENERATED ALWAYS AS ... STORED` (12), `ON CONFLICT DO UPDATE` (9.5), `SELECT ... FOR UPDATE`, functional/composite indexes. Nothing requires 15 |
| 6   | **`log.Fatal(server.ListenAndServe())`**, not a bare call                                                   | A discarded error means a bind failure exits silently with status 0 and no message                                                                                                                                                         |
| 7   | **Goose run as a global CLI from `sql/schema`**                                                             | Matches prior boot.dev workflow; keeps migration files free of Go wiring                                                                                                                                                                   |

---

## Open flags

| #   | Flag                                                                                                                                                                                                                                                                                                                                                                                                            | Status                                                     |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------- |
| 1   | **`change_controls` column count contradiction.** §3.2 and the §3 Summary state **48 columns**, but §3.2's own parenthetical reads "id + cc_number + 48 of the 50 BRD fields" = **50**. Counted the field-group tables: 48 BRD-derived (50 minus #24 and #34, which are `file_attachments` rows) + `id` + `cc_number` = **50 physical columns**. The doc's "48" counts BRD fields and mislabels them as columns | **Resolved in practice: build 50.** Doc correction pending |
| 2   | **`change_controls` DEFAULT count.** §6.4 says 8; §6.2 enumerates 7. The likely 8th is `cc_id`, which §6.2 explicitly states has _no_ DEFAULT (it's a generation expression)                                                                                                                                                                                                                                    | **Resolved in practice: 7.** Doc correction pending        |
| 3   | **`updated_at` vs `updated_on` on `refresh_tokens`.** Blueprint §7's code sample uses `updated_at`; DB Design §3.6 and CONTEXT_HANDOFF both say `updated_on`. Per precedence the DB doc wins — the Blueprint snippet is stale. Affects migration 006 and the refresh handler                                                                                                                                    | Logged, resolution clear                                   |
| 4   | **En-dash in HTML prototype `<option value="...">`.** A frontend built from the prototypes verbatim fails `ck_cc_requires_testing` on every submit. Frontend must normalize at the API boundary, or the prototypes get fixed                                                                                                                                                                                    | Open risk, pre-existing (DB §6.5)                          |
| 5   | **BRD §13.1 deferral note** for the three descoped password flows                                                                                                                                                                                                                                                                                                                                               | Lain to add on next BRD touch                              |
| 6   | **Production version parity.** Dev is on PostgreSQL 14.23; if production runs 15/16 there's a major-version gap. No feature dependency, so not a blocker — belongs in deployment notes                                                                                                                                                                                                                          | Noted                                                      |
| 7   | **The two `.docx` guardrail files are stored as plain text** despite the extension. Read them directly; do not attempt to unzip                                                                                                                                                                                                                                                                                 | Environmental note                                         |

---

## Environment & commands

WSL Ubuntu 22.04 · VS Code Remote-WSL · Go 1.25.x · PostgreSQL 14.23 · DBeaver for
browsing — but **read `\d` in psql first**; the GUI normalizes and can hide whether an
index is on `email` or `lower(email)`.

```bash
# migrations — run from sql/schema
goose postgres "postgres://postgres:PASS@localhost:5432/ea_qms?sslmode=disable" up
goose postgres "postgres://postgres:PASS@localhost:5432/ea_qms?sslmode=disable" down

# psql
psql "postgres://postgres:PASS@localhost:5432/ea_qms?sslmode=disable"
#   \l  databases    \dt  tables    \d <table>  full detail    \q  quit
```

**Every migration gets the up → `\d` → down → `\dt` → up cycle before it counts as done.**
If down doesn't cleanly reverse up, the migration has a bug — find it now, locally
(Blueprint §13).

Two observations from checkpoint 1 worth carrying forward:

- Postgres **rewrites `IN (...)` as `= ANY (ARRAY[...])`** in the catalog, so `\d` will
  never read back character-for-character as written. Normalization, not drift — expect
  it on all 13 CHECKs in 002.
- `TIMESTAMPTZ` displays as `timestamp with time zone`; that's the canonical spelling.
