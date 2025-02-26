package reflect

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func toTerraform5ValueErrorDiag(err error, path *tftypes.AttributePath) diag.AttributeErrorDiagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to convert into a Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
	)
}

func toTerraformValueErrorDiag(err error, path *tftypes.AttributePath) diag.AttributeErrorDiagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to convert the Attribute value into a Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
	)
}

func validateValueErrorDiag(err error, path *tftypes.AttributePath) diag.AttributeErrorDiagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to validate the Terraform value type. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
	)
}

func valueFromTerraformErrorDiag(err error, path *tftypes.AttributePath) diag.AttributeErrorDiagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Value Conversion Error",
		"An unexpected error was encountered trying to convert the Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
	)
}

type DiagIntoIncompatibleType struct {
	Val        tftypes.Value
	TargetType reflect.Type
	AttrPath   *tftypes.AttributePath
	Err        error
}

func (d DiagIntoIncompatibleType) Severity() diag.Severity {
	return diag.SeverityError
}

func (d DiagIntoIncompatibleType) Summary() string {
	return "Value Conversion Error"
}

func (d DiagIntoIncompatibleType) Detail() string {
	return fmt.Sprintf("An unexpected error was encountered trying to convert %T into %s. This is always an error in the provider. Please report the following to the provider developer:\n\n%s", d.Val, d.TargetType, d.Err.Error())
}

func (d DiagIntoIncompatibleType) Equal(o diag.Diagnostic) bool {
	od, ok := o.(DiagIntoIncompatibleType)
	if !ok {
		return false
	}
	if !d.Val.Equal(od.Val) {
		return false
	}
	if d.TargetType != od.TargetType {
		return false
	}
	if !d.AttrPath.Equal(od.AttrPath) {
		return false
	}
	if d.Err.Error() != od.Err.Error() {
		return false
	}
	return true
}

func (d DiagIntoIncompatibleType) Path() *tftypes.AttributePath {
	return d.AttrPath
}

type DiagNewAttributeValueIntoWrongType struct {
	ValType    reflect.Type
	TargetType reflect.Type
	SchemaType attr.Type
	AttrPath   *tftypes.AttributePath
}

func (d DiagNewAttributeValueIntoWrongType) Severity() diag.Severity {
	return diag.SeverityError
}

func (d DiagNewAttributeValueIntoWrongType) Summary() string {
	return "Value Conversion Error"
}

func (d DiagNewAttributeValueIntoWrongType) Detail() string {
	return fmt.Sprintf("An unexpected error was encountered trying to convert into a Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\nCannot use attr.Value %s, only %s is supported because %T is the type in the schema", d.TargetType, d.ValType, d.SchemaType)
}

func (d DiagNewAttributeValueIntoWrongType) Equal(o diag.Diagnostic) bool {
	od, ok := o.(DiagNewAttributeValueIntoWrongType)
	if !ok {
		return false
	}
	if d.ValType != od.ValType {
		return false
	}
	if d.TargetType != od.TargetType {
		return false
	}
	if !d.SchemaType.Equal(od.SchemaType) {
		return false
	}
	if !d.AttrPath.Equal(od.AttrPath) {
		return false
	}
	return true
}

func (d DiagNewAttributeValueIntoWrongType) Path() *tftypes.AttributePath {
	return d.AttrPath
}
