# CLAUDE.md

## Rules

- Always invoke commands via `mise run <task>` rather than calling `go`, `pnpm`, etc. directly. This includes routine operations like `go mod tidy` (`mise run backend-mod-tidy`), `go test`, `pnpm install`, `pnpm build`, etc.
- `openapi.yaml` is the single source of truth for the API. After editing it, run `mise run generate` and commit the generated files alongside the schema change.
- `manifests/` is never committed in a built state. The release CI builds `install.yaml` from kustomize and attaches it as a GitHub Release asset. Do not run `kustomize build` and commit the output.
- Ingress/PVC and other environment-specific resources are managed in a separate manifest repository, not here.
