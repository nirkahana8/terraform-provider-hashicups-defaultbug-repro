# Repro: object `Default` on a Computed `SingleNestedAttribute` has no effect when the object's children are also `computed`

Minimal, self-contained reproduction for
[hashicorp/terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework).

## The bug

A `SingleNestedAttribute` that is `Computed` and carries an object-level
`Default` (`objectdefault.StaticValue(...)`) loses that default for any child
attribute that is itself `Computed` and has **no default of its own**.

During `PlanResourceChange`:

1. `TransformDefaults` applies the object default → `locations = {is_any = true}`.
2. `MarkComputedNilsAsUnknown` then walks each attribute independently. It skips
   the `locations` object (it has an `ObjectDefaultValue`), but re-marks the child
   `locations.is_any` unknown — it is `Computed`, null in config, and has no
   default of its own — discarding the value the object default just set.

This contradicts the maintainer behavior table in
[#726](https://github.com/hashicorp/terraform-plugin-framework/issues/726)
(row: Default-on-nested = Yes, Default-on-child = No, config = null nested
attribute → "single nested attribute default").

See [ISSUE.md](./ISSUE.md) for the full write-up.

## Shape (mirrors a real provider)

```
scope     (Optional+Computed)
  users     (Optional+Computed) { is_any bool }          -- set in config
  locations (Optional+Computed, Default {is_any=true})   -- omitted in config
    is_any  (bool, Optional+Computed, no default of its own)
```

## Reproduce (acceptance test)

```console
$ TF_ACC=1 go test ./internal/provider/ -run TestNestedObjectDefault_childClobbered -v
```

Fails with:

```
Error: Provider returned invalid result object after apply
After the apply operation, the provider still indicated an unknown value
for defaultbug_thing.test.scope.locations.is_any. All values must be known after apply.
```

## Reproduce (manual `terraform plan`)

```console
$ go build -o terraform-provider-defaultbug .
$ cat > examples/.terraformrc <<EOF
provider_installation {
  dev_overrides { "registry.terraform.io/hashicorp/defaultbug" = "$(pwd)" }
  direct {}
}
EOF
$ cd examples && TF_CLI_CONFIG_FILE=.terraformrc terraform plan
```

Shows (object default applied, but the child was re-marked unknown):

```hcl
scope = {
  locations = {
      is_any = (known after apply)   # expected true
    }
  users = {
      is_any = true
    }
}
```

## Versions

- `github.com/hashicorp/terraform-plugin-framework v1.19.0`
- `github.com/hashicorp/terraform-plugin-go v0.31.0`
- `github.com/hashicorp/terraform-plugin-testing v1.16.0`
- Terraform v1.5.7
