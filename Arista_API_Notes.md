Aaron, here’s the **full Arista API enumeration you asked for**—organized so you can wire it straight into your Go/Wails explorer. I split it into (A) EOS eAPI (CLI command API) and (B) CloudVision APIs (gRPC models + uniform RPCs, with REST mappings), plus pointers to the Wi-Fi/CV-CUE and AGNI surfaces.

---

# A) EOS eAPI (EOS Command API)

**Transport & endpoint**

* Protocol: JSON-RPC over HTTP/HTTPS
* Path: `POST /command-api` on the switch
* Core method: `runCmds` (versioned)
* Auth: local switch credentials / token as configured
* Output: JSON (or text) per command
* **Scope:** Anything you can do in EOS CLI (show/config) is accessible via eAPI. In other words, **the eAPI surface == the EOS CLI command tree for your EOS version**. ([Arista][1], [GitHub][2], [arista.my.site.com][3])

> Practical implication: “Enumerating” eAPI means enumerating EOS CLI commands for your device image. The whitepaper explicitly calls the API *comprehensive* with parity to CLI, and community docs/libraries confirm JSON-RPC semantics. ([Arista][1], [GitHub][2])

**Canonical request shape (example)**

```json
{
  "jsonrpc": "2.0",
  "method": "runCmds",
  "params": {
    "version": 1,
    "cmds": ["enable", "show version"],
    "format": "json"
  },
  "id": "1"
}
```

(Commands, revisions, `text` vs `json` output modes are versioned/negotiated.) ([anta.arista.com][4])

**How to enumerate on a live device**

* Pull CLI help trees (e.g., `show ?`, `show running-config ?`, config mode help) and cache in your catalog. That catalog = your eAPI “endpoint” list since `cmds[]` are the verbs. (Arista doesn’t publish a single static global list because it varies by EOS release, features, and platform.) ([Arista][1])

---

# B) CloudVision APIs (gRPC models with uniform RPCs, REST mappings)

CloudVision’s public API surface is **model-driven gRPC**. Each model (e.g., `inventory.v1`, `event.v1`) exposes a **consistent set of RPCs** (GetOne/GetAll/Subscribe and Set/Delete for writable models). The official, living reference is published by Arista and is the **authoritative enumeration**. ([aristanetworks.github.io][5])

## B1) Uniform RPCs (apply across models)

Every readable model provides:

* `GetOne`, `GetAll` (server stream), `Subscribe` (server stream)

Writable models add:

* `Set`, `Delete` (and `SetSome`/`DeleteSome` for bulk)

These are standardized—your client can treat models uniformly. ([aristanetworks.github.io][6])

## B2) Complete list of **current CloudVision models** (Trunk)

Below is the **full model enumeration** from Arista’s reference (each item links to its model page where you can read fields, services, and messages). Use this list to auto-generate your explorer’s left-nav. ([aristanetworks.github.io][7])

* `action.v1`  ([aristanetworks.github.io][7])
* `alert.v1`  ([aristanetworks.github.io][7])
* `bugexposure.v1`  ([aristanetworks.github.io][7])
* `changecontrol.v1`  ([aristanetworks.github.io][7])
* `configlet.v1`  ([aristanetworks.github.io][7])
* `configstatus.v1`  ([aristanetworks.github.io][7])
* `connectivitymonitor.v1`  ([aristanetworks.github.io][7])
* `dashboard.v1`  ([aristanetworks.github.io][7])
* `endpointlocation.v1`  ([aristanetworks.github.io][7])
* `event.v1`  ([aristanetworks.github.io][7])
* `identityprovider.v1`  ([aristanetworks.github.io][7])
* `imagestatus.v1`  ([aristanetworks.github.io][7])
* `inventory.v1`  *(see sample services below)*  ([aristanetworks.github.io][8])
* `lifecycle.v1`  ([aristanetworks.github.io][7])
* `redirector.v1`  ([aristanetworks.github.io][7])
* `serviceaccount.v1`  ([aristanetworks.github.io][7])
* `softwaremanagement.v1`  ([aristanetworks.github.io][7])
* `studio.v1`  ([aristanetworks.github.io][7])
* `studio_topology.v1`  ([aristanetworks.github.io][7])
* `tag.v2`  ([aristanetworks.github.io][7])
* `workspace.v1`  ([aristanetworks.github.io][7])

> Tip: parse the models index and the per-model pages to build a **machine-readable catalog** in your app. The repo and site are kept up-to-date by Arista. ([GitHub][9], [aristanetworks.github.io][5])

### Example: Services within `inventory.v1` (excerpt)

Services surfaced for this model include (not exhaustive):

* `DeviceService`, `DeviceOnboardingService`, `DeviceOnboardingConfigService`,
  `DeviceDecommissioningService`, `DeviceDecommissioningConfigService`,
  `ProvisionedDeviceService`, etc., each with their `GetOne`/`GetAll`/`Subscribe` (and Set/Delete for config).
  See the model page for the full message types and service RPCs. ([aristanetworks.github.io][8])

## B3) REST mapping

Arista documents CloudVision APIs primarily as gRPC; **REST/HTTP mappings are generated** via standard gRPC-Gateway style. If your deployment exposes REST, paths will map to the proto services/resources (API Gateway), but the canonical reference remains the gRPC model pages above. ([aristanetworks.github.io][5])

---

# C) Additional Arista API Surfaces (adjacent)

* **CloudVision (platform overview & API gateway)** – context on the CV API gateway, OpenConfig usage, and streaming (gRPC) and REST availability. ([Arista][10])
* **Wi-Fi / CV-CUE REST APIs** – separate portal and docs for wireless/CUE (LaunchPad/Wireless Manager/Guest Manager). If you need to enumerate these too, pull their OpenAPI/portal index at build-time. ([apihelp.wifi.arista.com][11], [Arista][12])
* **AGNI (Identity) REST APIs** – separate, fully documented API guide (2025). Include as an optional module if your environment uses AGNI. ([Arista][13])

---

## Machine-readable seed (drop into your repo)

Here’s a compact seed catalog you can hydrate at runtime from the Arista docs:

```json
{
  "eapi": {
    "transport": "jsonrpc-http(s)",
    "endpoint": "/command-api",
    "method": "runCmds",
    "enumerationStrategy": "discover CLI tree for running EOS image; cmds[] mirrors CLI verbs",
    "notes": [
      "Parity with EOS CLI; varies by EOS version/platform/features",
      "Support text or json output per command; revision/version negotiable"
    ]
  },
  "cloudvision": {
    "models": [
      "action.v1","alert.v1","bugexposure.v1","changecontrol.v1",
      "configlet.v1","configstatus.v1","connectivitymonitor.v1","dashboard.v1",
      "endpointlocation.v1","event.v1","identityprovider.v1","imagestatus.v1",
      "inventory.v1","lifecycle.v1","redirector.v1","serviceaccount.v1",
      "softwaremanagement.v1","studio.v1","studio_topology.v1","tag.v2","workspace.v1"
    ],
    "uniformRPCs": ["GetOne","GetAll","Subscribe","Set","Delete","SetSome","DeleteSome"],
    "enumerationStrategy": "crawl model index + per-model pages to list services/messages; cache in local DB",
    "restMapping": "optional gRPC-Gateway/HTTP mapping if enabled by deployment"
  },
  "adjacent": {
    "wifi_cue": "Arista Wi-Fi Services/CV-CUE REST APIs",
    "agni": "AGNI REST APIs"
  }
}
```

---

## Why this is the “full” enumeration (and how your app stays current)

* **eAPI:** by design, *every* CLI command is accessible—so the true, complete list is the **current device CLI** (discover at runtime). ([Arista][1])
* **CloudVision:** Arista’s **Models** index + per-model pages are the **authoritative list**; you can fetch and cache these for offline/GUI use, and refresh on demand. ([aristanetworks.github.io][7])

If you want, I can generate a tiny Go helper that:

1. Scrapes the **CloudVision models index** at startup to refresh the catalog,
2. Probes a device over **eAPI** to learn the CLI tree for the logged-in EOS version, and
3. Emits a merged `api_catalog.json` your explorer can bind to.

[1]: https://www.arista.com/assets/data/pdf/Whitepapers/Arista_eAPI_FINAL.pdf?utm_source=chatgpt.com "Arista eAPI"
[2]: https://github.com/arista-netdevops-community/arista_eos_automation_with_eAPI?utm_source=chatgpt.com "This repo has examples of Arista EOS automation ..."
[3]: https://arista.my.site.com/AristaCommunity/s/question/0D52I00007ERpvrSAD/supported-commands-for-eapi?utm_source=chatgpt.com "Supported Commands for eAPI?"
[4]: https://anta.arista.com/stable/api/commands/?utm_source=chatgpt.com "Commands - Arista Network Test Automation - ANTA"
[5]: https://aristanetworks.github.io/cloudvision-apis "CloudVision APIs | Arista Networks FMP APIs"
[6]: https://aristanetworks.github.io/cloudvision-apis/rpcs "Uniform RPCs | Arista Networks FMP APIs"
[7]: https://aristanetworks.github.io/cloudvision-apis/models "Models | Arista Networks FMP APIs"
[8]: https://aristanetworks.github.io/cloudvision-apis/models/inventory.v1 "inventory.v1 | Arista Networks FMP APIs"
[9]: https://github.com/aristanetworks/cloudvision-apis "GitHub - aristanetworks/cloudvision-apis: gRPC APIs for integrating with Arista's CloudVision platform"
[10]: https://www.arista.com/en/products/eos/eos-cloudvision?utm_source=chatgpt.com "CloudVision"
[11]: https://apihelp.wifi.arista.com/?utm_source=chatgpt.com "Wi-Fi API Portal - Arista"
[12]: https://www.arista.com/assets/data/pdf/Whitepapers/Arista-Wi-Fi-API-App-Note.pdf?utm_source=chatgpt.com "Introduction to Arista CloudVision CUE API"
[13]: https://www.arista.com/assets/data/pdf/user-manual/um-books/AGNI-API-Guide.pdf?utm_source=chatgpt.com "CloudVision AGNI API Guide"


---

Here’s the **end-to-end path** your Go + Wails app will use to issue commands to Arista endpoints—covering EOS eAPI (on switches) and CloudVision (REST / gRPC). I’ll give you the **runtime flow**, **data contracts**, and **production-grade Go snippets** you can drop in.

---

# How the app issues commands

## 0) High-level runtime flow

1. **User picks an Endpoint** in the UI (switch or CloudVision cluster).
2. **UI builds a Request** (method, path/command, optional JSON body).
3. **Wails binding → Go backend** (uiapi): validates request, applies policy.
4. **Client layer** signs/attaches auth headers and calls the endpoint.
5. **Response** is parsed (JSON first, fallback text if needed).
6. **Result** is logged (request/response, timing) and streamed to the UI.

```mermaid
sequenceDiagram
  participant UI as UI (Explorer)
  participant API as Go uiapi (Wails)
  participant POL as Policy
  participant CLI as Client(eAPI/REST/gRPC)
  participant EP as Arista Endpoint

  UI->>API: Run(template/custom, endpointID, payload)
  API->>POL: Enforce(request)
  POL-->>API: allow/deny
  API->>CLI: Do(request)
  CLI->>EP: HTTP/JSON-RPC (eAPI) or REST/gRPC (CV)
  EP-->>CLI: Response (JSON / stream)
  CLI-->>API: Parsed payload + metadata
  API-->>UI: Render Table/JSON/Raw; append to log
```

---

## 1) Endpoint types & auth

* **EOS eAPI (switch)**

  * Transport: **JSON-RPC over HTTPS** → `POST https://<switch>/command-api`
  * Auth: Basic auth or token; some setups use cookie auth after login
* **CloudVision (controller)**

  * REST: `https://<cv>/api/...` with **Bearer token**
  * gRPC: TLS + token (you’ll generate Go stubs from Arista protos)

Your app stores: `Name, URL, Type (eapi|cv), Auth (username/password or token), TLS mode`.

---

## 2) Data contracts (Go)

```go
// internal/core/models.go
package core

import "time"

type EndpointType string
const (
	EndpointEAPI EndpointType = "eapi"
	EndpointCV   EndpointType = "cloudvision"
)

type Endpoint struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      EndpointType `json:"type"`
	URL       string       `json:"url"`
	Username  string       `json:"username,omitempty"`
	Password  string       `json:"password,omitempty"` // store encrypted
	Token     string       `json:"token,omitempty"`    // store encrypted
	Created   time.Time    `json:"created"`
	Tags      []string     `json:"tags"`
	TLSVerify bool         `json:"tlsVerify"`
}

type ExplorerRequest struct {
	EndpointID string                 `json:"endpointId"`
	Method     string                 `json:"method"` // GET/POST/PUT/DELETE (CV REST) or "runCmds" (eAPI)
	Path       string                 `json:"path"`   // e.g. "/command-api" or "/api/resources/..."
	Body       map[string]any         `json:"body,omitempty"`
	TimeoutMs  int                    `json:"timeoutMs,omitempty"`
}

type ExplorerResponse struct {
	Status     int                    `json:"status"`
	Headers    map[string][]string    `json:"headers"`
	JSON       any                    `json:"json,omitempty"`
	Text       string                 `json:"text,omitempty"`
	ElapsedMs  int64                  `json:"elapsedMs"`
	EndpointID string                 `json:"endpointId"`
	LogID      string                 `json:"logId"`
}
```

---

## 3) The eAPI path (switches)

### 3.1 Canonical eAPI call

* Method: **`runCmds`**
* Path: **`/command-api`**
* Body:

```json
{
  "jsonrpc":"2.0",
  "method":"runCmds",
  "params":{
    "version":1,
    "cmds":["enable","show version"],
    "format":"json",
    "autoComplete": true,
    "expandAliases": true
  },
  "id":"1"
}
```

> Notes:
>
> * Prepend `"enable"` to `cmds` if your next command requires privileged mode.
> * Prefer `format:"json"`; fall back to `text` when JSON isn’t available.
> * `autoComplete/expandAliases` make UX forgiving.

### 3.2 Go client (production-ready)

```go
// internal/client/eapi.go
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type EAPIClient struct {
	http *http.Client
}

func NewEAPIClient(tlsVerify bool, timeout time.Duration) *EAPIClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerify}, // set true in prod
	}
	return &EAPIClient{http: &http.Client{Transport: tr, Timeout: timeout}}
}

type runCmdsParams struct {
	Version       int           `json:"version"`
	Cmds          []string      `json:"cmds"`
	Format        string        `json:"format,omitempty"`        // json|text
	AutoComplete  bool          `json:"autoComplete,omitempty"`
	ExpandAliases bool          `json:"expandAliases,omitempty"`
}

type jsonRPCReq struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  runCmdsParams `json:"params"`
	ID      string        `json:"id"`
}

type jsonRPCResp struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  []any            `json:"result,omitempty"`
	Error   *jsonRPCError    `json:"error,omitempty"`
	ID      string           `json:"id"`
}
type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (c *EAPIClient) RunCmds(ctx context.Context, baseURL, user, pass string, body runCmdsParams) (*jsonRPCResp, *http.Response, time.Duration, error) {
	rpc := jsonRPCReq{
		JSONRPC: "2.0",
		Method:  "runCmds",
		Params:  body,
		ID:      "1",
	}
	raw, _ := json.Marshal(rpc)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/command-api", bytes.NewReader(raw))
	if err != nil { return nil, nil, 0, err }
	req.Header.Set("Content-Type", "application/json")
	if user != "" {
		req.SetBasicAuth(user, pass)
	}

	start := time.Now()
	resp, err := c.http.Do(req)
	elapsed := time.Since(start)
	if err != nil { return nil, resp, elapsed, err }
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var out jsonRPCResp
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, resp, elapsed, err
	}
	if out.Error != nil {
		return &out, resp, elapsed, errors.New(out.Error.Message)
	}
	return &out, resp, elapsed, nil
}
```

### 3.3 Example usage (show version)

```go
params := runCmdsParams{
	Version:      1,
	Cmds:         []string{"enable", "show version"},
	Format:       "json",
	AutoComplete: true,
}
rpc, httpResp, dur, err := eapi.RunCmds(ctx, endpoint.URL, endpoint.Username, endpoint.Password, params)
```

---

## 4) The CloudVision path

### Option A: REST (simpler to start)

* **Bearer token** in `Authorization: Bearer <token>`
* Example: `GET https://<cv>/api/resources/inventory/v1/Devices`

```go
// internal/client/cv_rest.go
func (c *HTTPClient) DoREST(ctx context.Context, method, url, bearer string, body any) (*http.Response, time.Duration, error) {
	var rdr io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		rdr = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil { return nil, 0, err }
	if body != nil { req.Header.Set("Content-Type", "application/json") }
	if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }

	start := time.Now()
	resp, err := c.http.Do(req)
	return resp, time.Since(start), err
}
```

### Option B: gRPC (authoritative, streaming)

* Generate stubs from Arista’s protos (models like `inventory.v1`, `event.v1`).
* Use TLS + token credentials; call `GetAll`, `GetOne`, `Subscribe`, etc.
* Stream results → forward to UI incrementally.

---

## 5) UI → Go binding (Wails)

```go
// internal/uiapi/explorer.go
package uiapi

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"arista_engine/internal/client"
	"arista_engine/internal/core"
)

type ExplorerAPI struct {
	Store   Store // interface for endpoint/kvs/log
	EAPI    *client.EAPIClient
	HTTP    *client.HTTPClient
	Policy  Policy // interface: Enforce(req) error
}

func (x *ExplorerAPI) Run(ctx context.Context, req core.ExplorerRequest) (core.ExplorerResponse, error) {
	ep, err := x.Store.GetEndpoint(req.EndpointID)
	if err != nil { return core.ExplorerResponse{}, err }

	// Policy check (deny dangerous ops if configured)
	if err := x.Policy.Enforce(req, ep); err != nil {
		return core.ExplorerResponse{}, err
	}

	timeout := 30 * time.Second
	if req.TimeoutMs > 0 { timeout = time.Duration(req.TimeoutMs) * time.Millisecond }
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	switch ep.Type {
	case core.EndpointEAPI:
		// Expect Path=/command-api; Body to include cmds/version/format
		params := client.RunCmdsParamsFromMap(req.Body) // helper: map -> struct
		rpc, httpResp, dur, err := x.EAPI.RunCmds(ctx, ep.URL, ep.Username, ep.Password, params)
		if err != nil && rpc == nil { return core.ExplorerResponse{}, err }

		out := core.ExplorerResponse{
			Status:     httpResp.StatusCode,
			Headers:    httpResp.Header,
			JSON:       rpc, // result or error already structured
			ElapsedMs:  dur.Milliseconds(),
			EndpointID: ep.ID,
		}
		logID, _ := x.Store.AppendLog(req, out)
		out.LogID = logID
		return out, err

	case core.EndpointCV:
		resp, dur, err := x.HTTP.DoREST(ctx, req.Method, ep.URL+req.Path, ep.Token, req.Body)
		if err != nil { return core.ExplorerResponse{}, err }
		defer resp.Body.Close()

		var parsed any
		_ = json.NewDecoder(resp.Body).Decode(&parsed)

		out := core.ExplorerResponse{
			Status:     resp.StatusCode,
			Headers:    resp.Header,
			JSON:       parsed,
			ElapsedMs:  dur.Milliseconds(),
			EndpointID: ep.ID,
		}
		logID, _ := x.Store.AppendLog(req, out)
		out.LogID = logID
		return out, nil

	default:
		return core.ExplorerResponse{}, errors.New("unsupported endpoint type")
	}
}
```

---

## 6) Command templates (drive the Explorer)

Put a JSON file (or Go map) in your repo and hydrate the left-nav:

```json
{
  "eapi": {
    "showVersion": {
      "title": "Show Version",
      "method": "runCmds",
      "path": "/command-api",
      "body": {"version":1,"format":"json","autoComplete":true,"cmds":["enable","show version"]}
    },
    "showInterfaces": {
      "title": "Show Interfaces",
      "method": "runCmds",
      "path": "/command-api",
      "body": {"version":1,"format":"json","autoComplete":true,"cmds":["enable","show interfaces"]}
    }
  },
  "cloudvision": {
    "devices": {
      "title": "Inventory Devices",
      "method": "GET",
      "path": "/api/resources/inventory/v1/Devices"
    }
  }
}
```

---

## 7) Connection test & enumeration

* **Connection test (switch):** call `/command-api` with `show version` (json) → if 200 + valid JSON, mark **Connected**.
* **Connection test (CV):** GET `/api/resources/inventory/v1/Devices?limit=1` with token.
* **Enumeration:**

  * **eAPI:** call `enable; show ?` (text), `show running-config ?`, enter config modes and introspect `?` output; cache discovered verbs.
  * **CV:** crawl the models/REST index your deployment exposes; cache service list.

---

## 8) Safety/policy (optional but recommended)

Before sending, evaluate a small rule set:

* Deny `"configure terminal"` unless a role allows it.
* Require explicit confirmation for `delete`, `set`, `write memory`, etc.
* Limit batch sizes / rate to avoid control-plane hits.

---

## 9) Logging & export

* Append **compact audit** per call: endpoint, method, path/commands, status, elapsed, response hash, sample.
* Export **JSON/CSV/PDF** from the UI (CSV for tables; JSON for raw; PDF with a clean template).

---

## 10) Error handling & resiliency

* Timeouts per call, exponential backoff on transient 5xx.
* Detect JSON parsing errors → show **Raw** tab automatically.
* For gRPC streams, surface progress indicators and partial rows.

---

### TL;DR

* **EOS switches:** Issue **JSON-RPC `runCmds`** to `/command-api` with `cmds:["enable","show ..."]`, parse JSON first, fallback to text.
* **CloudVision:** Use **REST (Bearer)** to `/api/...` or **gRPC** with generated stubs.
* Everything flows **UI → Wails binding → client**, with policy, logging, and exports around it.

If you want, I can scaffold these files (clients, uiapi, models, a `templates.json`, and a minimal Wails `main.go`) so you can `wails dev` and fire commands at a lab switch immediately.


---


# 🗺️ Updated roadmap (lean & Windows-focused)

## Phase 1 — Core hardening (keep)

* Windows packaging (Wails + **NSIS**), single EXE, auto-update optional later.
* Encrypted local secrets (Windows **DPAPI** via CryptProtectData).
* Structured logs to local files (JSON lines).
* Policy v0: deny “configure terminal”, “write memory”, “delete …” unless explicitly confirmed.
* Connection profiles limited to TLS verify and proxy (WinHTTP/Schannel).

## Phase 2 — Real API enumeration (keep)

* **EOS eAPI crawler** (BFS of `show ?` and safe config scopes like `interface ?`) with cache per device/EOS.
* **CloudVision model crawler** (models index → services → RPCs) to drive autocomplete.
* Diff on recrawl after upgrades.

## Phase 3 — Operator UX (keep)

* Saved queries, workspaces, run-on-multiple-endpoints (if desired).
* JSON→Table views with drag-select columns.
* Response diffs (same query, two runs/endpoints).
* CSV/JSON/PDF export.

## Phase 4 — Automation & workflows (keep, minimal)

* **Jobs Scheduler** (Windows-only): Run saved queries on a schedule; QPS guardrails.
* **Params/Macros**: `${VRF}`, `${INTF}`; matrix runs.
* **Guarded writes**: confirmation modal + pre/post checks.

## Phase 5 — Quality & release (keep)

* Golden test fixtures (sample EOS/CV payloads).
* Contract smoke tests for common EOS versions.
* Windows installer: **NSIS** script, code-signing if available.
* User backup/restore of app data folder.

### Phase 7 — Scale & performance (single-user value?)

**Short answer:** keep a **light subset**; drop the rest.

**Valuable even for one user:**

* **Connection pooling & reuse** (Schannel/WinHTTP under the hood): cuts latency for repeated calls.
* **Bounded worker pool (2–4 workers)**: lets you run a few queries concurrently without hammering control plane.
* **Streaming render** for very large JSON responses: UI stays responsive.

**Not worth it now:**

* Complex backpressure systems, distributed queues, external caches, advanced metrics → **remove**.

So Phase 7 becomes a tiny “Performance hygiene” checklist rather than a big phase.

---

# 🧩 Windows-only specifics (implementation notes)

* **Schannel/TLS** via Go stdlib on Windows; set `MinVersion: TLS12`.
* **Proxy**: respect WinHTTP system proxy; optional per-endpoint override.
* **Secrets**: use **DPAPI** (e.g., github.com/zalando/go-keyring or a thin DPAPI wrapper) to encrypt tokens/passwords.
* **Installer**: NSIS script with Start Menu shortcut, uninstall entry, optional “run at startup”.
* **PDF export**: wkhtmltopdf Windows binary or a pure-Go PDF lib; ship binary in `resources/`.

---

# 📒 What are Playbooks (and how to use them)?

Think “**operational runbooks**” (procedures) the app can execute or guide step-by-step—like lightweight Ansible **for read/verify** tasks (not device config). They’re great for **repeatable checks** and **pre/post-change validation**.

**Examples:**

* **Pre-upgrade health check**

  1. `show version`, `show environment cooling`, `show interfaces counters errors`
  2. Flag anomalies (power supply down, high CRCs)
  3. Export CSV & PDF summary
* **LLDP topology snapshot**

  1. `show lldp neighbors detail` across core devices
  2. Normalize → table → export to CSV/JSON for NetBox import
* **VLAN/Trunk consistency audit**

  1. `show interfaces trunk`, `show vlan`, `show spanning-tree detail`
  2. Detect mismatched VLAN allow lists & STP states
  3. Report deltas vs last baseline

**Format (simple JSON/YAML):**

```yaml
name: "Pre-Upgrade Health Check"
steps:
  - template: eapi_show_version
  - template: eapi_show_environment
  - template: eapi_show_interfaces_counters
post:
  - transform: summarize_errors
  - export: ["csv","pdf"]
```

You can run Playbooks ad-hoc or schedule them in Phase 4 jobs.

---

# 📚 What is a Schema Cookbook?

A **Schema Cookbook** is a collection of **field maps** turning raw JSON into operator-friendly tables—no JMESPath needed. It standardizes how you present common outputs so exports are consistent.

**Examples (excerpt):**

```json
{
  "tables": {
    "interfaces": {
      "source": "result[1].interfaces",
      "columns": [
        {"name":"Name","path":"name"},
        {"name":"Admin","path":"interfaceStatus.adminStatus"},
        {"name":"Oper","path":"interfaceStatus.operStatus"},
        {"name":"Speed","path":"bandwidth"},
        {"name":"Errors","path":"counters.inputErrors"}
      ]
    },
    "lldp_neighbors": {
      "source": "result[1].lldpNeighbors",
      "columns": [
        {"name":"Local Port","path":"port"},
        {"name":"Neighbor","path":"neighborDevice"},
        {"name":"Neighbor Port","path":"neighborPort"},
        {"name":"VLAN","path":"vlan"}
      ]
    }
  }
}
```

**Value:** one place to maintain column names and JSON paths; every table/export stays uniform across devices and time.

---

# 🧩 What is an Extensibility Doc?

It’s a **developer guide** explaining how others (or future you) can add **new command templates, playbooks, transforms, or exporters** without touching core code.

**Suggested sections:**

1. **Templates** — file location (`/configs/templates.json`), schema, examples, how UI auto-loads them.
2. **Playbooks** — YAML/JSON schema, built-in actions (`template`, `transform`, `export`), return contracts.
3. **Transforms** — how to drop a Go function or JS snippet to reshape JSON (e.g., compute error rates).
4. **Exports** — adding a new exporter (CSV/PDF/Markdown); interface and example.
5. **Policy rules** — TOML schema, matchers (`bodyContainsAny`, endpoint type), adding custom matchers.
6. **Enumeration plugins** — how to add new crawlers (e.g., Wi-Fi/CUE or AGNI) guarded behind a feature flag.
7. **Versioning** — how you version template packs/cookbooks so users can update them independently.

**Micro-example (Templates doc):**

```markdown
## Adding a Template
- Edit `configs/templates.json`
- Each entry must include: `title`, `method`, `path`, and optional `body`.
- Example:
{
  "eapi": {
    "show_mlag": {
      "title": "Show MLAG",
      "method": "runCmds",
      "path": "/command-api",
      "body": {"version":1,"format":"json","cmds":["enable","show mlag"]}
    }
  }
}
- Save → App reloads templates automatically.
```

---

# 📦 Final trimmed deliverables checklist

* [ ] Windows-only Wails project + NSIS installer.
* [ ] DPAPI-backed secret store.
* [ ] eAPI + CV clients (REST first; gRPC optional later).
* [ ] Enumeration cache + diff.
* [ ] Templates + **Playbooks** + **Schema Cookbook** folders.
* [ ] Policy v0 denylist.
* [ ] Exports: CSV/JSON/PDF.
* [ ] Minimal perf hygiene (keep connection reuse + small worker pool).

If you want, I’ll generate the folders/files for **Playbooks**, **Schema Cookbook**, and **Templates** with a few concrete examples so you can run a “Pre-Upgrade Health Check” out of the box on Windows.
