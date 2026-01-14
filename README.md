
# Prism

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Terraform Compatible](https://img.shields.io/badge/Terraform-Compatible-844FBA?style=flat&logo=terraform)](https://www.terraform.io/)
[![GitHub Stars](https://img.shields.io/github/stars/cylonchau/prism?style=social)](https://github.com/cylonchau/prism/stargazers)


**Prism** is an IaA (Infrastructure as API) platform. that provides unified API interfaces for multi-cloud through EAV (Entity-Attribute-Value) model-based Terraform trasform, Prism provides schema-free, provider-agnostic infrastructure operations across multi-cloud environments.

**Inspired by:** [WeCube Terraform Plugin](https://github.com/WeBankPartners/wecube-plugins-terraform) - evolved with modern API design and flexible data modeling.

## Concept

Prism unifies these differences into standard APIs through an abstraction layer:


```
User API Request
    ↓
Prism (Unified API)
    ↓
Terraform (Multi-Cloud)
    ↓
AWS | GCP | Tencent Cloud | Alibaba Cloud | Kubernetes | Openstack | vSphere ...
```

## Key Innovations
- Platform Engineering - Multi-Cloud.
- Dynamic EAV Model Architecture.
    - Schema-free resource modeling supporting any cloud provider without database migrations.
- Safe
    - Safe operations with built-in validation, state management, and rollback mechanisms.
    - Provider abstraction layer hiding cloud-specific complexities.