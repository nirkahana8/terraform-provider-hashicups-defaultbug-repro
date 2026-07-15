# Object `Default` on a Computed `SingleNestedAttribute` is discarded for computed children (clobbered by `MarkComputedNilsAsUnknown`)

### Module version

```
github.com/hashicorp/terraform-plugin-framework v1.19.0
github.com/hashicorp/terraform-plugin-go        v0.31.0
Terraform                                        v1.5.7
```

### Relevant provider source code

A `Computed` `SingleNestedAttribute` with an object-level `Default`, whose child is `Computed` with **no** default of its own:

```go
"nested": schema.SingleNestedAttribute{
    Optional: true,
    Computed: true,
    Default: objectdefault.StaticValue(types.ObjectValueMust(
        map[string]attr.Type{"child": types.StringType},
        map[string]attr.Value{"child": types.StringValue("default-value")},
    )),
    Attributes: map[string]schema.Attribute{
        "child": schema.StringAttribute{
            Optional: true,
            Computed: true,
            // no Default
        },
    },
},
```

Config omits the block entirely (null nested attribute):

```hcl
resource "defaultbug_thing" "test" {}
```

### Expected behavior

Per the maintainer behavior table in #726 — row **Default-on-nested = Yes, Default-on-child = No, config = null nested attribute → "single nested attribute default"** — the plan should be:

```hcl
+ nested = {
    + child = "default-value"
  }
```

and apply should succeed.

### Actual behavior

`terraform plan` shows the child as unknown despite the object default:

```hcl
+ nested = {
    + child = (known after apply)
  }
```

and `terraform apply` fails:

```
Error: Provider returned invalid result object after apply
After the apply operation, the provider still indicated an unknown value
for defaultbug_thing.test.nested.child. All values must be known after apply.
```

### Root cause

In `PlanResourceChange` (`internal/fwserver/server_planresourcechange.go`) the ordering is:

1. `data.TransformDefaults(...)` — applies the object default. Trace confirms:
   `setting attribute scope.locations to default value: {"is_any":true,...}` (`fwschemadata/data_default.go`).
2. `MarkComputedNilsAsUnknown(...)` — walks **each attribute independently**. For the `nested` object it correctly returns early because it has an `ObjectDefaultValue`:

   ```go
   case fwschema.AttributeWithObjectDefaultValue:
       if a.ObjectDefaultValue() != nil {
           return val, nil // object node skipped — keeps the default
       }
   ```

   …but it then visits the **child** `nested.child`, which is `Computed`, null in config, and has no default *of its own*, so it falls through the switch and is re-marked unknown:

   ```go
   return tftypes.NewValue(newValueType, tftypes.UnknownValue), nil
   ```

   This discards the value the parent object default just populated.

So an object `Default` is effectively ineffective for **any** nested object whose children are `Computed` and lack their own per-attribute defaults. It only appears to work when every child *also* carries its own default (the case in #726's workaround), because then `MarkComputedNilsAsUnknown` skips those children too.

### Suggested direction

When `MarkComputedNilsAsUnknown` encounters a computed, null-in-config attribute, it should not re-mark it unknown if an **ancestor** attribute supplied a (non-null) default value for it during `TransformDefaults` — i.e. respect a value already set by an ancestor object default, not only a per-attribute default. (Alternatively, document this sharp edge on the [Default](https://developer.hashicorp.com/terraform/plugin/framework/resources/default) page and in #726's table, since the `Yes | No | null` row does not hold when children are `Computed`.)

### References

- Related: #726 (Defaults don't work in nested attributes) — the maintainer behavior table this contradicts.
- Related: #777 (Consider Object Nested Attribute Default Implementation) — enhancement, different ask.
- Minimal repro: https://github.com/nirkahana8/terraform-provider-hashicups-defaultbug-repro/tree/defaultbug-repro/defaultbug-repro
```
