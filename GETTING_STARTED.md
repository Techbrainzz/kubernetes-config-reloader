Here is the complete, single copy-pasteable document for your **`GETTING_STARTED.md`** file. You can copy the block below, paste it directly into your `GETTING_STARTED.md` file, and save it:

```markdown
# 🚀 Getting Started Guide: Running the Operator Locally

This guide walks you through setting up a local test environment and running the **Kubernetes Config Reloader Operator** from scratch. Perfect for testing or demonstrating how the operator functions!

---

## 🛠️ Prerequisites
Before starting, ensure you have the following tools installed and running on your workstation:
* **Docker Desktop** (Must be running in the background)
* **Go (Golang)** (Version 1.22 or newer)
* **kubectl** (Kubernetes command-line tool)
* **Kind** (Kubernetes in Docker, for spinning up a local test cluster)

---

## Step 1: Create a Local Kubernetes Cluster
If you don't already have a cluster running, spin up a lightweight test cluster using `kind`:
```bash
kind create cluster --name kind

```

*(Verify your connection by running: `kubectl cluster-info`)*

---

## Step 2: Install the Custom Resource Definition (CRD)

Your operator introduces a new custom type to Kubernetes (`ConfigSync`). Apply the CRD manifest to teach your cluster about this new resource type:

```bash
kubectl apply -f config/crd/bases/config.infra.io_configsyncs.yaml

```

---

## Step 3: Run the Operator Manager

Open your **first terminal window**, navigate to your project directory, and start the controller manager locally:

```bash
go run ./cmd/main.go

```

*Leave this terminal open! Log messages will appear here whenever the operator processes a resource.*

---

## Step 4: Test the Live Rolling Restart

Open a **second terminal window** (keeping your operator running in the first one) and follow these steps to test the automation:

1. **Deploy a test workload (Nginx):**
```bash
kubectl create deployment my-app --image=nginx --replicas=1

```


2. **Apply your ConfigSync custom resource sample:**
```bash
kubectl apply -f config/samples/config_v1alpha1_configsync.yaml

```


3. **Verify the results:**
* **Check Operator Logs (Terminal 1):** You will see log outputs confirming it computed the SHA-256 hash and triggered a rolling restart for `my-app`.
* **Inspect the Deployment:** Check that the dynamic configuration hash annotation was injected into your pod template:
```bash
kubectl get deployment my-app -o=jsonpath='{.spec.template.metadata.annotations}'

```


* **Check the Synced ConfigMap:** Verify that the operator automatically provisioned the target ConfigMap:
```bash
kubectl get configmaps

```





```

```