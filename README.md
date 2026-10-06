# Kubernetes Config Reloader Operator 🚀

[![Go Version](https://img.shields.io/github/go-mod/go-version/Techbrainzz/kubernetes-config-reloader)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

An enterprise-grade, declarative Kubernetes custom controller built in Go using `controller-runtime` and Operator SDK. It solves the **runtime config-reload gap** left behind by GitOps continuous delivery tools (like Argo CD) and package managers (like Helm).

---

## 💡 The Architectural Problem Solved

While tools like Argo CD excel at keeping cluster state synchronized with Git, they stop at the API server boundary:
* **The Blind Spot:** When a ConfigMap is updated via GitOps, the data in `etcd` changes successfully, but **running application pods continue using stale configuration values loaded into memory** at startup. 
* **The Traditional Solution:** Requires manual rolling bounces, custom file-watching sidecars, or bloated CI/CD pipelines.

**Config Reloader Operator** bridges this runtime gap automatically by introducing a declarative controller that watches custom configuration resources and actively orchestrates safe, zero-downtime updates.

---

## ✨ Key Features

1. **Declarative Custom Resource (`ConfigSync`):** Define configuration maps and target workloads using high-level custom resources.
2. **Automated Rolling Restarts via Hashing:** Computes deterministic cryptographic **SHA-256 hashes** of configuration blocks and injects them into deployment pod template annotations (`config.infra.io/config-hash`), forcing Kubernetes to seamlessly roll out updates only when data actually changes.
3. **Safe Garbage Collection:** Utilizes Kubernetes Owner References so downstream ConfigMaps are automatically cleaned up if the parent custom resource is deleted.
4. **Production Observability:** Reports real-time status subresources (`Synced`/`Failed`) and exposes custom Prometheus metrics for reconciliation tracking.

---

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Framework:** Kubebuilder / Controller-Runtime (`controller-runtime v0.21.0`)
* **Environment:** Kubernetes, Kind, Custom Resource Definitions (CRDs)

---

## 🚀 Quick Navigation

* 📖 **[Getting Started Guide (For Beginners)](GETTING_STARTED.md)** - Step-by-step instructions on setting up a local cluster and running the operator.
* 🏗️ **[Architectural Build Log](ARCHITECTURE.md)** - Deep-dive technical breakdown of every phase of engineering.

---

## 📂 Project Structure

```text
├── api/v1alpha1/           # Custom Resource Definition (CRD) Go types & schemas
├── cmd/main.go             # Operator entrypoint and manager initialization
├── config/                 # Kubebuilder manifests (CRDs, RBAC, samples)
├── internal/controller/    # Reconciliation loop, hashing logic, and Prometheus metrics
├── README.md               # Main documentation
├── ARCHITECTURE.md         # Engineering build log
└── GETTING_STARTED.md      # Beginner execution runbook