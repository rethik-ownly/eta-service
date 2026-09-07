# Scripts

## load_restaurant_estimates.py

Loads restaurant ETA estimates from a CSV/TSV file into MongoDB by calling the eta-service HTTP API. Run **once per environment** — point `--base-url` at local, staging, or production.

### Prerequisites

- eta-service running and connected to the target MongoDB
- Python 3 (stdlib only — no pip install)

### CSV format

Tab- or comma-separated with header row:

```
sublocalityId	cityId	zoneId	dayType	mealType	restaurantId	rat_seconds	rat_sample_count	kpt_seconds	kpt_sample_count	pickup_seconds	pickup_sample_count	delay_dispatch_seconds	delay_dispatch_sample_count
```

- `dayType`: one or more values separated by comma, pipe, or semicolon
  - Single: `weekday` → `["WEEKDAY"]`
  - Multiple: `MONDAY,TUESDAY,WEDNESDAY,THURSDAY,FRIDAY` → `["MONDAY","TUESDAY",...]`
  - Mixed: `weekday|weekend` → `["WEEKDAY","WEEKEND"]`
- `mealType`: `Breakfast`, `Lunch`, `Snacks`, `Dinner`, `Latenight` (normalized to API enums)
- Blank `cityId` / `zoneId` → use `--default-city-id` / `--default-zone-id`
- Missing numeric sample counts default to `0`

### Usage

```bash
# Local
python scripts/load_restaurant_estimates.py \
  --file data.csv \
  --base-url http://localhost:8080 \
  --default-city-id 1 \
  --default-zone-id Z_BLR

# Staging (self-signed TLS: add --insecure)
python scripts/load_restaurant_estimates.py \
  --file data.csv \
  --base-url https://eta-service.staging.example.com \
  --default-city-id 1 \
  --default-zone-id Z_BLR \
  --insecure

# Production (run manually when ready)
python scripts/load_restaurant_estimates.py \
  --file data.csv \
  --base-url https://eta-service.prod.example.com \
  --default-city-id 1 \
  --default-zone-id Z_BLR
```

### Behaviour

For each CSV row:

1. **POST** `/api/v1/eta/restaurant` — insert with this row's meal filled, other meals zero
2. On **409 Conflict** (doc already exists) → **PATCH** `/api/v1/eta/restaurant/:restaurantId` for that meal
   - If the CSV row has **multiple** day types (e.g. `weekday|weekend`), PATCH also sends `"days"` to update the `dayType` array

Re-running against the same database is safe: existing documents are updated via PATCH.

Use `--fail-fast` to stop on the first failed row.

---

## load_sublocality_estimates.py

Loads sublocality ETA estimates (cat/fm) from a CSV/TSV file via the eta-service HTTP API. Same row-by-row insert-then-patch flow as the restaurant loader.

### CSV format

```
sublocalityId	cityId	zoneId	dayType	mealType	cat_seconds	cat_sample_count	fm_seconds	fm_sample_count
```

- `dayType`: one or more values (`weekday`, `weekday|weekend`, `MONDAY,TUESDAY`, etc.)
- `mealType`: `Breakfast`, `Lunch`, `Snacks`, `Dinner`, `Latenight`

### Usage

```bash
python3 scripts/load_sublocality_estimates.py \
  --file sublocality_data.csv \
  --base-url http://localhost:8080

# Staging (self-signed TLS)
python3 scripts/load_sublocality_estimates.py \
  --file sublocality_data.csv \
  --base-url https://eta-service.staging.example.com \
  --insecure
```

### Behaviour

For each CSV row:

1. **POST** `/api/v1/eta/sublocality` — insert with this row's meal filled, other meals zero
2. On **409 Conflict** → **PATCH** `/api/v1/eta/sublocality/:sublocalityId` for that meal
   - Multiple day types in CSV → PATCH also sends `"days"` to update the `dayType` array

Sublocality IDs with spaces (e.g. `R T Nagar`) are URL-encoded in the PATCH path.
