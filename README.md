# Kaiju Registry API

A registry of known kaiju, their classifications, and reported sightings across the globe.

Built as a sample API for demonstrating [Kong Konnect](https://konghq.com) platform capabilities.

## The Registry

The Kaiju Registry tracks giant monsters discovered worldwide. Each kaiju entry includes:

- **Species classification** — from Mega Primates to Leviathans
- **Biometrics** — height in meters, mass in metric tonnes
- **Threat level** — `alpha`, `beta`, `gamma`, or `omega`
- **Status** — `active`, `dormant`, `neutralized`, or `unknown`
- **Sighting history** — timestamped, geolocated reports with confirmation status

### Sample Entries

| ID | Name | Species | Height | Threat | Status |
|----|------|---------|--------|--------|--------|
| k-001 | Goraxus | Mega Primate | 102m | omega | active |
| k-002 | Tidestrider | Leviathan | 78m | gamma | dormant |

## API Endpoints

```
GET /kaiju                        List all kaiju (filterable by threat_level)
GET /kaiju/{kaijuId}              Get a single kaiju record
GET /kaiju/{kaijuId}/sightings    List reported sightings for a kaiju
```

All list endpoints support pagination via `page` and `page_size` query parameters.

## Quick Start

Point your client at the base URL and start querying:

```bash
# List all omega-level threats
curl "https://api.kaiju-registry.example.com/v1/kaiju?threat_level=omega"

# Get details on Goraxus
curl "https://api.kaiju-registry.example.com/v1/kaiju/k-001"

# Check recent sightings
curl "https://api.kaiju-registry.example.com/v1/kaiju/k-001/sightings"
```

## Environments

| Environment | Base URL |
|-------------|----------|
| Production | `https://api.kaiju-registry.example.com/v1` |
| Staging | `https://staging-api.kaiju-registry.example.com/v1` |

## Spec

The full OpenAPI 3.0.3 specification lives in [`openapi.yaml`](openapi.yaml).

## License

[Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0)
