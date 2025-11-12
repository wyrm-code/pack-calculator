# Pack Calculator (Golang)

A tiny HTTP API (with a minimal UI) that calculates how many **whole packs** to ship for an order.

Rules implemented:

1. Only whole packs may be sent (no splitting).
2. Among all options, return the solution with the **least total items shipped** (>= ordered).
3. If multiple solutions ship the same least number of items, choose the one with the **fewest number of packs**.

This is equivalent to a coin‑change problem with a **lexicographic objective**: minimize `(total_items, number_of_packs)`.
The algorithm below is deterministic and works with *any* positive pack sizes. Pack sizes are fully configurable at runtime.

---

## Quick Start

### Local (no Docker)

```bash
# build & test
make test
make run
# server starts at :8080
```

### With Docker

```bash
# build the container
make docker-build
# run the container
make docker-run
# stop container
make docker-stop
```

Server will listen on `:8080` by default. Override with `PORT` env var.

---

## API

### `POST /api/calc`

**Body**
```json
{
  "items": 12001,
  "sizes": [250, 500, 1000, 2000, 5000]
}
```

If `sizes` is omitted, the server uses the defaults from the environment variable `PACK_SIZES` (comma-separated integers, e.g. `250,500,1000,2000,5000`).

**Response**
```json
{
  "itemsOrdered": 12001,
  "totalItems": 12250,
  "packs": [
    {"size": 5000, "qty": 2},
    {"size": 2000, "qty": 1},
    {"size": 250,  "qty": 1}
  ]
}
```

### `GET /api/calc?items=263&sizes=250,500,1000`
URL form for quick testing. `sizes` optional (comma list).

### `GET /health`
Simple health check.

---

## UI

Open `http://localhost:8080/` after starting the server.
- Enter pack sizes (comma‑separated) and `items`.
- Click **Calculate** to see the resulting pack breakdown.

---

## Configuration

- `PORT` – server port (default `8080`)
- `PACK_SIZES` – default sizes, e.g. `250,500,1000,2000,5000`

---

## Implementation Notes

- **Algorithm**: Unbounded dynamic programming up to `N + (maxSize - 1)`.
  - We compute the best exact totals `S` (number of packs only) for all `S` in `[0, limit]`.
  - Then we choose the best `S >= N` by minimizing `(S, packs)`.
  - Backpointers reconstruct the quantities for each pack size.
- Time complexity `O(limit * numSizes)`, where `limit = N + maxSize - 1`.
- Works with any positive sizes, in any order, with duplicates removed.

---

## Tests

Run `make test` to execute unit tests that cover the examples from the assignment and a few edge cases.

---

## Project Layout

```
.
├── Dockerfile
├── Makefile
├── README.md
├── go.mod
├── cmd/server/main.go
├── internal/solver/solver.go
├── internal/solver/solver_test.go
└── web
    └── index.html
```

---

## License

MIT (do whatever you want; attribution appreciated).
