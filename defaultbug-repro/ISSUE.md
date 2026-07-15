# Object `Default` on a Computed `SingleNestedAttribute` is discarded for computed children (clobbered by `MarkComputedNilsAsUnknown`)

### Module version

```
github.com/hashicorp/terraform-plugin-framework v1.19.0
github.com/hashicorp/terraform-plugin-go        v0.31.0
Terraform                                        v1.5.7
```

### Relevant provider source code

A grandparent `scope` (Computed) containing a `locations` object (Computed) that
carries a classic object-level `Default`. `locations`' child `is_any` is Computed
with **no default of its own**:

```go
"scope": schema.SingleNestedAttribute{
    Optional: true,
    Computed: true,
    Attributes: map[string]schema.Attribute{
        "users": schema.SingleNestedAttribute{ /* is_any bool, no default */ },
        "locations": schema.SingleNestedAttribute{
            Optional: true,
            Computed: true,
            Default: objectdefault.StaticValue(types.ObjectValueMust(
                map[string]attr.Type{
                    "is_any":       types.BoolType,
                    "location_ids": types.ListType{ElemType: types.StringType},
                },
                map[string]attr.Value{
                    "is_any":       types.BoolValue(true),
                    "location_ids": types.ListNull(types.StringType),
                },
            )),
            Attributes: map[string]schema.Attribute{
                "is_any":       schema.BoolAttribute{Optional: true, Computed: true}, // no Default of its own
                "location_ids": schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType}, // no Default
            },
        },
    },
},
```

Config sets `scope` (via `users`) but omits `locations`:

```hcl
resource "defaultbug_thing" "test" {
  scope = { users = { is_any = true } }
}
```

### Expected behavior

Per the maintainer behavior table in #726 — row **Default-on-nested = Yes,
Default-on-child = No, config = null nested attribute → "single nested attribute
default"** — the plan should apply the object default:

```hcl
scope = {
  locations = { is_any = true, location_ids = null }
  users     = { is_any = true, user_ids = (known after apply) }
}
```

### Actual behavior

`terraform plan` (the object default is applied, then every computed child is re-marked unknown):

```hcl
scope = {
  locations = {
    is_any       = (known after apply)
    location_ids = (known after apply)
  }
  users = { is_any = true, user_ids = (known after apply) }
}
```

`terraform apply` fails:

```
Error: Provider returned invalid result object after apply
After the apply operation, the provider still indicated an unknown value
for defaultbug_thing.test.scope.locations.is_any. All values must be known
after apply.
```

The same happens if `locations = {}` is written explicitly (present but with a
null child): the object default's child is still clobbered.

### Root cause

In `PlanResourceChange` (`internal/fwserver/server_planresourcechange.go`) the order is:

1. `data.TransformDefaults(...)` — applies the object default. Trace confirms:
   `setting attribute scope.locations to default value: {"is_any":true}` (`fwschemadata/data_default.go`).
2. `MarkComputedNilsAsUnknown(...)` — walks **each attribute independently**. It
   correctly returns early for the `locations` object because it has an
   `ObjectDefaultValue`:

   ```go
   case fwschema.AttributeWithObjectDefaultValue:
       if a.ObjectDefaultValue() != nil {
           return val, nil // object node kept
       }
   ```

   …but it then visits the child `scope.locations.is_any`, which is `Computed`,
   null in config, and has no default *of its own*, so it falls through and is
   re-marked unknown:

   ```go
   return tftypes.NewValue(newValueType, tftypes.UnknownValue), nil
   ```

   This discards the value the parent object default just populated.

So an object `Default` is effectively ineffective for **any** nested object whose
children are `Computed` and lack their own per-attribute defaults. It only appears
to work when every child *also* carries its own default (the case in #726's
workaround) — then `MarkComputedNilsAsUnknown` skips those children too.

### Suggested direction

When `MarkComputedNilsAsUnknown` encounters a computed, null-in-config attribute,
it should not re-mark it unknown if an **ancestor** attribute already supplied a
non-null value via `TransformDefaults` (i.e. respect a value set by an ancestor
object default, not only a per-attribute default). Alternatively, document this
sharp edge on the [Default](https://developer.hashicorp.com/terraform/plugin/framework/resources/default)
page and in #726's table, since the `Yes | No | null` row does not hold when the
children are `Computed`.

### References

- Related: #726 (Defaults don't work in nested attributes) — the maintainer behavior table this contradicts.
- Related: #777 (Consider Object Nested Attribute Default Implementation) — enhancement, different ask.
- Minimal repro: https://github.com/nirkahana8/terraform-provider-hashicups-defaultbug-repro/tree/defaultbug-repro/defaultbug-repro
