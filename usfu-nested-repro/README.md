# Repro: object `UseStateForUnknown` on a generated (CustomType) nested attribute crashes with "Attribute Missing"

Minimal reproduction for
[hashicorp/terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework),
built from an IR via `tfplugingen-framework generate all` so the schema uses the
**generated** strict CustomType (not a hand-written plain object type).

## The bug

Attaching the stock `objectplanmodifier.UseStateForUnknown()` to a generated
`SingleNestedAttribute` (backed by a CustomType) crashes `terraform plan` when:

- the attribute is **omitted** from config (null + `Computed` ⇒ unknown in the plan), and
- its child is **itself an object**.

The framework materializes a partial parent object and runs it through the
generated `NestedType.ValueFromObject`, which rejects it:

```
Error: Attribute Missing
sub is missing from object
```

It does **not** crash when the attribute is fully specified in config.

Related (closed/locked): terraform-plugin-framework #767, #754.

## How it's built

- `provider_code_spec.json` — the IR: a `thing` resource with `id` and a computed
  `nested` object whose child `sub` is itself an object (`{ ref_id }`).
- `internal/generated/` — produced by `tfplugingen-framework generate all` (the
  strict CustomType `NestedType`/`SubType` live here).
- `internal/provider/thing_resource.go` — uses the generated `ThingResourceSchema`
  and attaches `objectplanmodifier.UseStateForUnknown()` to `nested`.

Regenerate with:

```console
$ tfplugingen-framework generate all --input provider_code_spec.json --output internal/generated
```

## Reproduce (acceptance test)

```console
$ TF_ACC=1 go test ./internal/provider/ -run TestObjectUseStateForUnknown_crash -v
```

Fails with `Error: Attribute Missing / sub is missing from object`.

## Reproduce (manual `terraform plan`)

```console
$ go build -o terraform-provider-usfu .
$ cat > examples/.terraformrc <<EOF
provider_installation {
  dev_overrides { "registry.terraform.io/hashicorp/usfu" = "$(pwd)" }
  direct {}
}
EOF
$ cd examples && TF_CLI_CONFIG_FILE=.terraformrc terraform plan
```

## Versions

- `github.com/hashicorp/terraform-plugin-framework v1.19.0`
- `github.com/hashicorp/terraform-plugin-go v0.31.0`
- `github.com/hashicorp/terraform-plugin-testing v1.16.0`
- Terraform v1.5.7
