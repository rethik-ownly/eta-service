#!/usr/bin/env python3
"""
Load sublocality ETA estimates from a CSV/TSV file into MongoDB via the eta-service HTTP API.

Each CSV row is processed independently:
  1. POST /api/v1/eta/sublocality  (insert doc with this row's meal filled, others zero)
  2. On HTTP 409 (doc exists) → PATCH /api/v1/eta/sublocality/:sublocalityId for that meal

Usage:
  python3 scripts/load_sublocality_estimates.py \\
    --file sublocality_data.csv \\
    --base-url http://localhost:8080 \\
    --default-city-id 572ca7ff116b5db3057bd814 \\
    --default-zone-id Z_BLR

Run once per environment — point --base-url at local, staging, or production.

Expected CSV columns (tab or comma separated):
  sublocalityId, cityId, zoneId, dayType, mealType,
  cat_seconds, cat_sample_count, fm_seconds, fm_sample_count

dayType accepts one or more values separated by comma, pipe, or semicolon
(e.g. weekday, weekend or MONDAY,TUESDAY,FRIDAY).
"""

from __future__ import annotations

import argparse
import csv
import json
import ssl
import sys
import urllib.error
import urllib.parse
import urllib.request
from typing import Any

MEAL_KEYS = ("breakfast", "lunch", "snacks", "dinner", "latenight")

MEAL_TYPE_TO_KEY = {
    "breakfast": "breakfast",
    "lunch": "lunch",
    "snacks": "snacks",
    "snack": "snacks",
    "dinner": "dinner",
    "latenight": "latenight",
    "late_night": "latenight",
    "late night": "latenight",
}

MEAL_TYPE_TO_API = {
    "breakfast": "BREAKFAST",
    "lunch": "LUNCH",
    "snacks": "SNACK",
    "snack": "SNACK",
    "dinner": "DINNER",
    "latenight": "LATE_NIGHT",
    "late_night": "LATE_NIGHT",
    "late night": "LATE_NIGHT",
}

DAY_TYPE_TO_API = {
    "weekday": "WEEKDAY",
    "weekend": "WEEKEND",
    "monday": "MONDAY",
    "tuesday": "TUESDAY",
    "wednesday": "WEDNESDAY",
    "thursday": "THURSDAY",
    "friday": "FRIDAY",
    "saturday": "SATURDAY",
    "sunday": "SUNDAY",
}

EXPECTED_COLUMNS = (
    "sublocalityId",
    "cityId",
    "zoneId",
    "dayType",
    "mealType",
    "cat_seconds",
    "cat_sample_count",
    "fm_seconds",
    "fm_sample_count",
)


def parse_float(value: str, default: float = 0.0) -> float:
    value = (value or "").strip()
    if not value:
        return default
    return float(value)


def parse_int(value: str, default: int = 0) -> int:
    value = (value or "").strip()
    if not value:
        return default
    return int(float(value))


def zero_meal() -> dict[str, Any]:
    zero = {"seconds": 0, "sampleCount": 0}
    return {
        "cat": dict(zero),
        "fm": dict(zero),
    }


def meal_from_row(row: dict[str, str]) -> dict[str, Any]:
    return {
        "cat": {
            "seconds": parse_float(row["cat_seconds"]),
            "sampleCount": parse_int(row["cat_sample_count"]),
        },
        "fm": {
            "seconds": parse_float(row["fm_seconds"]),
            "sampleCount": parse_int(row["fm_sample_count"]),
        },
    }


def normalize_meal_type(raw: str) -> tuple[str, str]:
    key = MEAL_TYPE_TO_KEY.get(raw.strip().lower())
    if not key:
        raise ValueError(f"unknown mealType: {raw!r}")
    return key, MEAL_TYPE_TO_API[key]


def normalize_day_types(raw: str) -> list[str]:
    """Parse one or more day types from a CSV cell."""
    parts = [part.strip() for part in raw.replace("|", ",").replace(";", ",").split(",")]
    day_types: list[str] = []
    seen: set[str] = set()

    for part in parts:
        if not part:
            continue
        normalized = DAY_TYPE_TO_API.get(part.lower())
        if not normalized:
            raise ValueError(f"unknown dayType: {part!r}")
        if normalized not in seen:
            seen.add(normalized)
            day_types.append(normalized)

    if not day_types:
        raise ValueError("dayType is empty")

    return day_types


def patch_day(day_types: list[str]) -> str:
    """Pick a day value for PATCH; must be a member of the document's dayType array."""
    return day_types[0]


def build_insert_payload(
    row: dict[str, str],
    meal_key: str,
    day_types: list[str],
    city_id: str,
    zone_id: str,
) -> dict[str, Any]:
    meals = {key: zero_meal() for key in MEAL_KEYS}
    meals[meal_key] = meal_from_row(row)

    return {
        "sublocalityId": row["sublocalityId"].strip(),
        "dayType": day_types,
        "cityId": city_id,
        "zoneId": zone_id,
        **meals,
    }


def build_patch_payload(row: dict[str, str], meal_api: str, day_types: list[str]) -> dict[str, Any]:
    payload: dict[str, Any] = {
        "day": patch_day(day_types),
        "mealType": meal_api,
        "catSeconds": parse_float(row["cat_seconds"]),
        "catSampleCount": parse_int(row["cat_sample_count"]),
        "fmSeconds": parse_float(row["fm_seconds"]),
        "fmSampleCount": parse_int(row["fm_sample_count"]),
    }

    if len(day_types) > 1:
        payload["days"] = day_types

    return payload


def detect_dialect(path: str) -> csv.Dialect:
    with open(path, newline="", encoding="utf-8-sig") as f:
        sample = f.read(4096)
        f.seek(0)
        try:
            return csv.Sniffer().sniff(sample, delimiters="\t,")
        except csv.Error:
            return csv.excel_tab


def read_rows(path: str) -> list[dict[str, str]]:
    dialect = detect_dialect(path)
    with open(path, newline="", encoding="utf-8-sig") as f:
        reader = csv.DictReader(f, dialect=dialect)
        if not reader.fieldnames:
            raise ValueError("CSV has no header row")

        fieldnames = [name.strip() for name in reader.fieldnames]
        reader.fieldnames = fieldnames

        missing = [col for col in EXPECTED_COLUMNS if col not in fieldnames]
        if missing:
            raise ValueError(f"CSV missing columns: {', '.join(missing)}")

        return list(reader)


def ssl_context(insecure: bool) -> ssl.SSLContext | None:
    if not insecure:
        return None
    context = ssl.create_default_context()
    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
    return context


def http_request(
    method: str,
    url: str,
    body: dict[str, Any] | None = None,
    insecure: bool = False,
) -> tuple[int, str]:
    data = None
    headers = {"Accept": "application/json"}
    if body is not None:
        data = json.dumps(body).encode("utf-8")
        headers["Content-Type"] = "application/json"

    request = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(request, timeout=60, context=ssl_context(insecure)) as response:
            return response.status, response.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        return exc.code, exc.read().decode("utf-8", errors="replace")


def process_row(
    row: dict[str, str],
    base_url: str,
    default_city_id: str,
    default_zone_id: str,
    line_no: int,
    insecure: bool = False,
) -> tuple[str, str]:
    sublocality_id = row["sublocalityId"].strip()
    if not sublocality_id:
        return "failed", "missing sublocalityId"

    city_id = row["cityId"].strip() or default_city_id
    zone_id = row["zoneId"].strip() or default_zone_id
    if not city_id:
        return "failed", "missing cityId (and no --default-city-id)"
    if not zone_id:
        return "failed", "missing zoneId (and no --default-zone-id)"

    try:
        meal_key, meal_api = normalize_meal_type(row["mealType"])
        day_types = normalize_day_types(row["dayType"])
    except ValueError as exc:
        return "failed", str(exc)

    insert_url = f"{base_url.rstrip('/')}/api/v1/eta/sublocality"
    insert_body = build_insert_payload(row, meal_key, day_types, city_id, zone_id)
    status, response_body = http_request("POST", insert_url, insert_body, insecure=insecure)

    if status == 201:
        return "inserted", ""

    if status == 409:
        patch_url = f"{base_url.rstrip('/')}/api/v1/eta/sublocality/{urllib.parse.quote(sublocality_id, safe='')}"
        patch_body = build_patch_payload(row, meal_api, day_types)
        patch_status, patch_body_text = http_request("PATCH", patch_url, patch_body, insecure=insecure)
        if patch_status == 200:
            return "updated", ""
        return "failed", f"PATCH returned HTTP {patch_status}: {patch_body_text}"

    return "failed", f"POST returned HTTP {status}: {response_body}"


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Load sublocality ETA estimates from CSV via eta-service HTTP API.",
    )
    parser.add_argument("--file", required=True, help="Path to sublocality CSV (tab or comma separated)")
    parser.add_argument("--base-url", required=True, help="eta-service base URL, e.g. http://localhost:8080")
    parser.add_argument("--default-city-id", default="", help="Fallback when cityId column is blank")
    parser.add_argument("--default-zone-id", default="", help="Fallback when zoneId column is blank")
    parser.add_argument("--fail-fast", action="store_true", help="Stop on first failed row")
    parser.add_argument(
        "--insecure",
        action="store_true",
        help="Skip TLS certificate verification (staging self-signed certs)",
    )
    args = parser.parse_args()

    if args.insecure:
        print("warning: TLS certificate verification disabled (--insecure)", file=sys.stderr)

    try:
        rows = read_rows(args.file)
    except (OSError, ValueError) as exc:
        print(f"error reading CSV: {exc}", file=sys.stderr)
        return 1

    counts = {"inserted": 0, "updated": 0, "failed": 0}

    for index, row in enumerate(rows, start=2):
        status, detail = process_row(
            row,
            args.base_url,
            args.default_city_id,
            args.default_zone_id,
            index,
            insecure=args.insecure,
        )
        counts[status] += 1

        sublocality_id = row.get("sublocalityId", "").strip()
        meal_type = row.get("mealType", "").strip()
        if status == "failed":
            print(f"row {index} {sublocality_id} {meal_type}: FAILED — {detail}", file=sys.stderr)
            if args.fail_fast:
                break
        else:
            print(f"row {index} {sublocality_id} {meal_type}: {status}")

    print(
        f"\ndone: inserted={counts['inserted']} updated={counts['updated']} failed={counts['failed']}",
        file=sys.stderr,
    )
    return 1 if counts["failed"] else 0


if __name__ == "__main__":
    sys.exit(main())
