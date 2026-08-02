# Repro: object plan modifier on a generated (CustomType) nested attribute crashes with "Attribute Missing"

Minimal reproduction for
[hashicorp/terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework),
built from an IR via `tfplugingen-framework generate all` so the schema uses the
**generated** strict CustomType (not a hand-written plain object type).

## The bug

Attaching an object plan modifier (here the stock
`objectplanmodifier.UseStateForUnknown()`) to a generated `SingleNestedAttribute`
(backed by a CustomType) crashes `terraform plan` when the attribute is **unknown
in the plan with no prior state to substitute** — i.e. on **create** when it is
omitted (Optional+Computed):

```
Error: Attribute Missing
kind is missing from object
```

The framework materializes a partial object and runs it through the generated
`EngineType.ValueFromObject`, which rejects it.

It does **not** crash when:
- the attribute is present in config (known in the plan),
- on update (prior state substitutes via `UseStateForUnknown`), or
- the modifier is not attached.

The child attribute types (**scalar or object**) are irrelevant.

Related (closed/locked): terraform-plugin-framework #767, #754.

## How it's built

- `provider_code_spec.json` — the IR: a `car` resource with `id`, `name`, `color`,
  and a computed `engine` object with `kind` (string) and `volume` (int64).
- `internal/generated/` — produced by `tfplugingen-framework generate all` (the
  strict CustomType `EngineType` lives here).
- `internal/provider/car_resource.go` — uses the generated `CarResourceSchema`
  and attaches `objectplanmodifier.UseStateForUnknown()` to `engine`.

Regenerate with:

```console
$ tfplugingen-framework generate all --input provider_code_spec.json --output internal/generated
```

## Reproduce (acceptance test)

```console
$ TF_ACC=1 go test ./internal/provider/ -run TestObjectUseStateForUnknown_crash -v
```

Fails with `Error: Attribute Missing / kind is missing from object`.

## Reproduce (manual `terraform plan`)

```console
$ go build -o terraform-provider-usfu .
$ cat > examples/.terraformrc <<EOF
provider_installation {
  dev_overrides { "registry.terraform.io/hashicorp/usfu" = "../" }
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
