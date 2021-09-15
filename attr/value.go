package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Value defines an interface for describing data associated with an attribute.
// Values allow provider developers to specify data in a convenient format, and
// have it transparently be converted to formats Terraform understands.
type Value interface {
	// Type returns the Type that created the Value.
	Type(context.Context) Type

	// ToTerraformValue returns the data contained in the Value as
	// a Go type that tftypes.NewValue will accept.
	ToTerraformValue(context.Context) (interface{}, error)

	// Equal must return true if the Value is considered semantically equal
	// to the Value passed as an argument.
	Equal(Value) bool
}

func ValueToTerraform(ctx context.Context, val Value) (tftypes.Value, error) {
	raw, err := val.ToTerraformValue(ctx)
	if err != nil {
		return tftypes.Value{}, err
	}
	err = tftypes.ValidateValue(val.Type(ctx).TerraformType(ctx), raw)
	if err != nil {
		return tftypes.Value{}, err
	}
	return tftypes.NewValue(val.Type(ctx).TerraformType(ctx), raw), nil
}
