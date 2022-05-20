terraform {
  required_version = ">= 1.0.3"
}

provider "google" {
  credentials = var.credentials
  project = var.project
  region  = var.region
  zone = var.zone
}

resource "google_compute_network" "vpc" {
  name                    = "${var.name}-vpc"
  auto_create_subnetworks = "false"
}

# Subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "${var.name}-subnet"
  region        = var.region
  network       = google_compute_network.vpc.name
  ip_cidr_range = "10.255.0.0/16"

  secondary_ip_range {
    range_name    = "cluster-secondary-ip-range"
    ip_cidr_range = "10.0.0.0/12"
  }

  secondary_ip_range {
    range_name    = "services-secondary-ip-range"
    ip_cidr_range = "10.64.0.0/12"
  }
}

resource "google_container_cluster" "gitpod-cluster" {
  name     = var.name
  location = "${var.zone}" == "" ? "${var.region}" : "${var.region}-${var.zone}"

  cluster_autoscaling {
    enabled = true

    resource_limits {
        resource_type = "cpu"
        minimum       = 2
        maximum       = 8
    }

    resource_limits {
      resource_type = "memory"
      minimum       = 4
      maximum       = 64
    }
  }

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
  release_channel {
    channel = "UNSPECIFIED"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = "cluster-secondary-ip-range"
    services_secondary_range_name = "services-secondary-ip-range"
  }

  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  addons_config {
    http_load_balancing {
      disabled = false
    }

    horizontal_pod_autoscaling {
      disabled = false
    }


    # only available in beta
    # dns_cache_config {
    #   enabled = true
    # }
  }

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name


}

resource "google_container_node_pool" "services" {
  name       = "${var.name}-services"
  location   = google_container_cluster.gitpod-cluster.location
  cluster    = google_container_cluster.gitpod-cluster.name
  version    = var.kubernetes_version // kubernetes version
  initial_node_count = 1
  max_pods_per_node = 110

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    labels = {
      "gitpod.io/workload_meta" =true
      "gitpod.io/workload_ide" = true
    }

    preemptible  = var.pre-emptible
    image_type = "UBUNTU_CONTAINERD"
    disk_type = "pd-ssd"
    disk_size_gb  = var.disk_size_gb
    machine_type = var.machine_type
    tags         = ["gke-node", "${var.project}-gke"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  autoscaling {
    min_node_count = var.min_count
    max_node_count = var.max_count
  }


  management {
    auto_repair = true
    auto_upgrade = false
  }
}

resource "google_container_node_pool" "workspaces" {
  name       = "${var.name}-workspaces"
  location     = google_container_cluster.gitpod-cluster.location
  cluster    = google_container_cluster.gitpod-cluster.name
  version    = var.kubernetes_version // kubernetes version
  initial_node_count = 1
  max_pods_per_node = 110

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    labels = {
      "gitpod.io/workload_metal" =true
      "gitpod.io/workload_ide" = true
      "gitpod.io/workload_workspace_services" = true
      "gitpod.io/workload_workspace_regular" = true
      "gitpod.io/workload_workspace_headless" = true
    }

    preemptible  = var.pre-emptible
    image_type = "UBUNTU_CONTAINERD"
    disk_type = "pd-ssd"
    disk_size_gb  = var.disk_size_gb
    machine_type = var.machine_type
    tags         = ["gke-node", "${var.project}-gke"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  autoscaling {
    min_node_count = var.min_count
    max_node_count = var.max_count
  }

  management {
    auto_repair = true
    auto_upgrade = false
  }
}

module "gke_auth" {
  depends_on = [google_container_node_pool.workspaces]

  source = "terraform-google-modules/kubernetes-engine/google//modules/auth"

  project_id   = var.project
  location     = google_container_cluster.gitpod-cluster.location
  cluster_name = var.name
}

resource "local_file" "kubeconfig" {
  filename = var.kubeconfig
  content = module.gke_auth.kubeconfig_raw
}
