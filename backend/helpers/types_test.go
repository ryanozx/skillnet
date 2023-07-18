package helpers

import (
	"reflect"
	"testing"
)

func TestParseNullableUint(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    *NullableUint
		wantErr bool
	}{
		{
			"Parse Nullable int OK",
			args{
				str: "1",
			},
			&NullableUint{
				value:    1,
				hasValue: true,
			},
			false,
		},
		{
			"Parse empty string",
			args{
				str: "",
			},
			&NullableUint{
				hasValue: false,
			},
			false,
		},
		{
			"Parse non-numeric",
			args{
				str: "abc",
			},
			&NullableUint{},
			true,
		},
		{
			"Parse negative value",
			args{
				str: "-1",
			},
			&NullableUint{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNullableUint(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNullableUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNullableUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableUint_IsNull(t *testing.T) {
	type fields struct {
		value    uint
		hasValue bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"Is null",
			fields{
				hasValue: false,
			},
			true,
		},
		{
			"Is not null",
			fields{
				value:    1,
				hasValue: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &NullableUint{
				value:    tt.fields.value,
				hasValue: tt.fields.hasValue,
			}
			if got := v.IsNull(); got != tt.want {
				t.Errorf("NullableUint.IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableUint_GetValue(t *testing.T) {
	type fields struct {
		value    uint
		hasValue bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint
		wantErr bool
	}{
		{
			"Non-null value",
			fields{
				value:    1,
				hasValue: true,
			},
			1,
			false,
		},
		{
			"Null value",
			fields{
				hasValue: false,
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &NullableUint{
				value:    tt.fields.value,
				hasValue: tt.fields.hasValue,
			}
			got, err := v.GetValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullableUint.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NullableUint.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
