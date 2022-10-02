# argocd-terraform-plugin
![Pipeline](https://github.com/KazanExpress/argocd-terraform-plugin/workflows/Pipeline/badge.svg)
![Code Scanning](https://github.com/KazanExpress/argocd-terraform-plugin/workflows/Code%20Scanning/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/KazanExpress/argocd-terraform-plugin)](https://goreportcard.com/report/github.com/KazanExpress/argocd-terraform-plugin)
![Downloads](https://img.shields.io/github/downloads/IBM/argocd-terraform-plugin/total?logo=github)
[![codecov](https://codecov.io/gh/KazanExpress/argocd-terraform-plugin/branch/main/graph/badge.svg?token=6Xr7V8AMTE)](https://codecov.io/gh/KazanExpress/argocd-terraform-plugin)

<img src="https://github.com/KazanExpress/argocd-terraform-plugin/raw/main/assets/argo_vault_logo.png" width="300">

An Argo CD plugin to retrieve secrets from various Secret Management tools (HashiCorp Vault, IBM Cloud Secrets Manager, AWS Secrets Manager, etc.) and inject them into Kubernetes resources

### Why use this plugin?
This plugin is aimed at helping to solve the issue of secret management with GitOps and Argo CD. We wanted to find a simple way to utilize Secret Management tools without having to rely on an operator or custom resource definition. This plugin can be used not just for secrets but also for deployments, configMaps or any other Kubernetes resource.

## Documentation
You can our full set of documentation at https://argocd-terraform-plugin.readthedocs.io/

## Contributing
Interested in contributing? Please read our contributing documentation [here](./CONTRIBUTING.md) to get started!

## Blogs
- [Solving ArgoCD Secret Management with the argocd-terraform-plugin](https://itnext.io/argocd-secret-management-with-argocd-terraform-plugin-539f104aff05)
- [Introducing argocd-terraform-plugin v1.0!](https://itnext.io/introducing-argocd-terraform-plugin-v1-0-708433294b2d)
- [How to Use HashiCorp Vault and Argo CD for GitOps on OpenShift](https://cloud.redhat.com/blog/how-to-use-hashicorp-vault-and-argo-cd-for-gitops-on-openshift)

## Presentations
- [Shh, It’s a Secret: Managing Your Secrets in a GitOps Way - Jake Wernette & Josh Kayani, IBM](https://youtu.be/7L6nSuKbC2c)
