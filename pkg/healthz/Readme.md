# Healthz Server

A lightweight, dependency-free health check server implemented using a raw TCP listener.  
It is designed for high-performance environments such as Kubernetes, Docker, and load balancers
where liveness/readiness checks must be fast, simple, and extremely reliable.

---

## 🚀 Features

- Ultra-lightweight (no HTTP framework required)
- Supports multiple health check functions
- Catches panics (prevents healthz server from crashing if a check fails)
- Gracefully shuts down on SIGINT / SIGTERM
- Fully compatible with Kubernetes readiness/liveness probes
- Uses direct TCP for maximum performance

---

## 🧠 Why TCP Healthz Instead of Full HTTP Server?

Using `net.ListenTCP` instead of `http.Server` has several advantages:

### ✔ Zero overhead

No routing, no middleware, no allocations → fastest possible health check.

### ✔ Consistent & deterministic

No goroutine explosions, no HTTP keep-alive side effects.

### ✔ Safe in container environments

Check functions can safely panic without killing the server  
(thanks to `runCatchPanic()`).

### ✔ Perfect for K8s probes

Container orchestrators only care about:

- "Is port open?"
- "Do you return 200 or 503?"

Nothing more.

---

### Usage

package main

import (
"your-module/healthz"
)

func main() {
healthz.RunServer()
}

Healthz server will run on port 9999.
