# RakshaGrid — System Design & Architecture

**Offline-resilient hyperlocal emergency response mesh, designed to scale from a single campus pilot to a multi-tenant campus-safety SaaS.**

---

## 1. Requirements

### 1.1 Functional Requirements
- A verified user can raise an SOS with one tap, carrying live geolocation and incident type.
- The system geo-matches the SOS to the nearest available *verified* responder (warden, security staff, trained volunteer).
- Responders receive real-time alerts and can accept/escalate within an SLA window.
- Every incident produces an immutable, timestamped audit trail (raised → matched → acked → resolved).
- Institute admins get a dashboard for live incidents, response-time analytics, and compliance reporting (UGC/AICTE anti-ragging & women-safety mandates).
- The client app continues to function (raise + relay SOS) on local Wi-Fi/LAN even when the internet is down.
- New campuses (tenants) can be onboarded without code changes — configuration-driven provisioning.

### 1.2 Non-Functional Requirements
| Attribute | Target | Notes |
|---|---|---|
| Dispatch match latency | < 1s p95 | Time from SOS received to responder-list computed |
| End-to-end alert delivery | < 3s p95 | Includes WebSocket/push fan-out |
| Availability (data plane) | 99.9% | Degrades to offline-LAN mode rather than hard-failing |
| Offline capability | Core requirement | PWA must raise/relay SOS over LAN with no WAN internet |
| Data durability | No incident record loss | Audit log is the compliance backbone — durability > raw speed here |
| Multi-tenancy | Campus-level isolation | Data residency & privacy per institute |

### 1.3 Constraints
- Small founding team (1–4 people) — favors a **modular monolith** over premature microservices.
- Must be demonstrably buildable with a known stack (Go, Redis, PostgreSQL, React) rather than exotic tech that can't be defended under questioning.
- Hackathon-stage budget — architecture must have a credible path from free-tier/single-VM to real multi-region scale, not require it on day one.

---

## 2. High-Level Design

### 2.1 Component Architecture
See `01-hld-component-architecture.mermaid`.

- **Client layer**: User PWA, Responder PWA, Admin Dashboard — all offline-first via service workers, syncing when connectivity returns.
- **Edge**: A single Go-based API Gateway handles auth (JWT), rate limiting, and request routing. This is the only public-facing entry point.
- **Application services** (shipped as one deployable Go binary at launch, split later — see §5):
  - **Auth & Identity** — verifies users/responders, issues JWTs, manages institute-level roles.
  - **SOS Dispatch** — the core: takes an SOS, queries Redis for nearest verified responders, orchestrates matching + escalation.
  - **Responder Registry** — tracks responder availability/presence via periodic heartbeat into Redis.
  - **Notification/Realtime Gateway** — WebSocket fan-out to responders and live status back to the user, with push notification fallback when WS isn't available.
  - **Audit & Compliance** — writes the immutable incident record and generates compliance reports.
  - **Tenant Management** — onboards new campuses, manages billing hooks and per-tenant configuration.
- **Data layer**:
  - **Redis** — geospatial sorted sets (`GEOADD`/`GEOSEARCH`) for responder location + presence pub/sub. Chosen for sub-millisecond in-memory geo queries over PostGIS, which is better suited to heavier analytical spatial queries we don't need on the hot path.
  - **PostgreSQL** — source of truth for users, responders, incidents, tenants.
  - **Object storage (S3-compatible)** — incident evidence/attachments.
  - **Event queue** — durable log of incident/notification events, enabling replay and analytics without coupling services directly.
- **Observability**: Prometheus + Grafana across all services from day one — response-time SLAs are the product's core promise, so they must be measured, not assumed.

### 2.2 Core Data Flow — Raising and Resolving an SOS
See `02-sos-dispatch-sequence.mermaid`.

1. User raises SOS with live geolocation → API Gateway → Dispatch Service.
2. Dispatch queries Redis `GEOSEARCH` for verified, on-duty responders within a configurable radius, ranked by distance and availability.
3. Notification Gateway pushes the alert to the top-N responders via WebSocket (push notification fallback if WS unavailable).
4. First responder to ACK within the SLA window (e.g. 60s) is assigned; user gets a live status update.
5. If no ACK within the SLA window, the system auto-escalates to the next-nearest responder and notifies the admin.
6. On resolution, the full timeline is written to the Audit Service as an immutable record.

### 2.3 Offline-Resilience Design (the actual differentiator)
- The PWA ships with a service worker that caches the app shell and queues SOS requests locally.
- On a campus LAN, a **local relay node** (a lightweight instance of the Dispatch + Notification services, deployable on a single on-prem box or Raspberry-Pi-class device) can serve SOS matching *without* internet, syncing the incident log to the cloud once connectivity returns.
- This means the product's core promise — "works when the network doesn't" — is an architectural property, not a UI trick.

---

## 3. Data Model (core entities)

```
Tenant(id, name, region, plan_tier, created_at)
User(id, tenant_id, name, verified_flag, role[student|staff], contact_info)
Responder(id, tenant_id, user_id, verification_status, responder_type[warden|security|volunteer])
ResponderPresence(responder_id, lat, lng, status[available|busy|offline], last_heartbeat)   -- lives in Redis
Incident(id, tenant_id, raised_by, type, lat, lng, status, created_at, resolved_at)
IncidentEvent(id, incident_id, event_type[raised|matched|acked|escalated|resolved], actor_id, timestamp)  -- immutable, append-only
ComplianceReport(id, tenant_id, period, generated_at, metrics_json)
```

- `IncidentEvent` is intentionally append-only — it's the audit trail the compliance story depends on, so it is never updated or deleted, only appended to.
- Multi-tenancy is enforced with `tenant_id` on every row at the smaller-tenant tier (shared schema); larger tenants can be promoted to a dedicated schema/database for stricter data-residency or contractual requirements (see §5).

---

## 4. API Contract (representative, not exhaustive)

| Method & Path | Purpose |
|---|---|
| `POST /v1/sos` | Raise an SOS `{geo, type}` — auth required |
| `GET /v1/incidents/{id}` | Poll incident status (fallback if WS unavailable) |
| `POST /v1/incidents/{id}/ack` | Responder accepts an assigned incident |
| `POST /v1/incidents/{id}/resolve` | Mark incident resolved, triggers audit write |
| `WS /v1/stream` | Realtime channel: SOS alerts to responders, status updates to users |
| `GET /v1/admin/incidents` | Admin dashboard feed — filterable by date/type/status |
| `GET /v1/admin/compliance-report` | Auto-generated UGC/AICTE-aligned report |
| `POST /v1/admin/responders` | Onboard/verify a new responder |
| `POST /v1/tenants` | Provision a new campus (internal/ops use) |

---

## 5. Scaling This Into a Business

### 5.1 Multi-Tenant SaaS Model
See `03-scale-out-deployment.mermaid`.

- **Control plane** (single region): tenant registry, billing/provisioning, Route 53 latency + health-check based routing to the nearest healthy data-plane region.
- **Data plane** (per region): ALB → Auto Scaling Group running the Go services → Redis cluster + RDS Postgres, with a warm-standby region for cross-region failover — this directly reuses the warm-standby / Route 53 failover pattern from the AWS Multi-Region DR project, giving a concrete, previously-validated failure story to defend under questioning.
- **Tenant isolation tiers**:
  - *Trial/small campus*: shared Postgres schema, `tenant_id`-scoped rows — fast onboarding, low cost.
  - *Large/enterprise campus*: isolated schema or dedicated database — needed once a tenant has data-residency or contractual isolation requirements.

### 5.2 Growth Path (what to revisit as it scales)
| Stage | Trigger | Change |
|---|---|---|
| 1 campus (pilot) | Launch | Single-region, modular monolith, shared Postgres |
| ~5–20 campuses | Response-time SLA under load / multi-institute demand | Split Dispatch + Notification into independently scaled services; introduce Redis Cluster sharding |
| ~20+ campuses | Cross-region customers, compliance-heavy tenants | Multi-region active/warm-standby (as designed above); promote large tenants to isolated schemas |
| Heavy analytics demand | Institutes asking for safety-trend insights (a sellable upsell) | Introduce durable event streaming (Kafka/SQS) so audit events can be replayed into an analytics warehouse without touching the hot path |
| National scale | Cross-border data-residency law differences | Region-pinned data planes per compliance jurisdiction |

### 5.3 Business Model Angles This Architecture Enables
- **Per-campus SaaS subscription**, tiered by student count / responder count.
- **Compliance-reporting as a paid add-on** — auto-generated UGC/AICTE-aligned reports are a genuine budget line item for institutes today done manually.
- **Anonymized safety-analytics** sold back to institutes/city safety bodies (aggregate incident-density heatmaps, response-time benchmarking) — the append-only `IncidentEvent` log is the raw asset this is built from.

---

## 6. Trade-off Analysis

| Decision | Chosen | Rejected alternative | Why |
|---|---|---|---|
| Geo-matching store | Redis geospatial | PostGIS | Sub-ms in-memory lookups matter more on the hot path than rich spatial queries; PostGIS is a good fit later for analytics, not for real-time matching |
| Service topology at launch | Modular monolith | Microservices from day one | Small team, needs to ship and defend a working demo fast; premature service boundaries cost more than they save at this scale |
| Tenant isolation | Shared schema (small) → isolated schema (large) | Isolated schema for everyone from day one | Isolated-per-tenant infra is operationally expensive before there's paying-tenant volume to justify it |
| Realtime transport | WebSocket + push fallback | Pure polling | Sub-second alert delivery is the product's core SLA; polling can't hit that reliably |
| Multi-region strategy | Warm standby + async replication | Active-active multi-region | Active-active adds conflict-resolution complexity that isn't justified until there's real cross-region traffic; warm-standby is a known, previously-validated pattern for this team |

---

## 7. What This Design Deliberately Does *Not* Solve Yet
- Indoor positioning still depends on a Wi-Fi/BLE-beacon fallback layer that requires physical beacon deployment per building — a hardware dependency, not purely software.
- Responder verification is a manual onboarding workflow with each institute's security office; this system does not (and should not) automate identity verification of responders.
- Cold-start: the safety guarantee only holds once a campus has a minimum density of onboarded, verified responders — this is a go-to-market sequencing problem as much as an engineering one.
