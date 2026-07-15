# Repro: object `Default` on a Computed `SingleNestedAttribute` has no effect when the object's children are also `computed`

Minimal, self-contained reproduction for
[hashicorp/terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework).

## The bug

A `SingleNestedAttribute` that is `Computed` and carries an object-level
`Default` (`objectdefault.StaticValue(...)`) loses that default for any child
attribute that is itself `Computed` and has **no default of its own**.

During `PlanResourceChange` the framework:

1. `TransformDefaults` applies the object default → `nested = {child = "default-value"}`.
2. `MarkComputedNilsAsUnknown` then walks each attribute independently. It
   correctly skips the `nested` object (it has an `ObjectDefaultValue`), but it
   still visits the child `nested.child` — which is `Computed`, null in config,
   and has no default of its own — and re-marks it **unknown**, discarding the
   value the object default just set.

This contradicts the maintainer behavior table in
[#726](https://github.com/hashicorp/terraform-plugin-framework/issues/726)
(row: Default-on-nested = Yes, Default-on-child = No, config = null nested
attribute → "single nested attribute default").

## Files

- `internal/provider/nested_default_resource.go` — the resource. Schema is hand-written
  (no codegen): a `Computed` `nested` object with `objectdefault.StaticValue({child="default-value"})`
  whose `child` is `Computed` with **no** default.
- `internal/provider/nested_default_resource_test.go` — an acceptance test that omits
  `nested` and expects `nested.child == "default-value"`.

## Reproduce (acceptance test)

```console
$ TF_ACC=1 go test ./internal/provider/ -run TestNestedObjectDefault_childClobbered -v
```

Fails with:

```
Error: Provider returned invalid result object after apply
After the apply operation, the provider still indicated an unknown value
for defaultbug_thing.test.nested.child. All values must be known after apply.
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

Shows (the object default was applied, but `child` was re-marked unknown):

```hcl
  + resource "defaultbug_thing" "test" {
      + id     = (known after apply)
      + nested = {
          + child = (known after apply)   # expected "default-value"
        }
    }
```

## Versions

- `github.com/hashicorp/terraform-plugin-framework v1.19.0`
- `github.com/hashicorp/terraform-plugin-go v0.31.0`
- `github.com/hashicorp/terraform-plugin-testing v1.16.0`
- Terraform v1.5.7
