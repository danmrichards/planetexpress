package event

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromValues(t *testing.T) {
	tests := []struct {
		name      string
		values    map[string]interface{}
		expPkgEvt *PackageEvent
		expErr    error
	}{
		{
			name: "valid",
			values: map[string]interface{}{
				"type":         "package_allocate",
				"package_id":   "bar",
				"package_size": "123",
			},
			expPkgEvt: &PackageEvent{
				EventID:     "test",
				Typ:         PackageAllocate,
				PackageID:   "bar",
				PackageSize: 123,
			},
		},
		{
			name: "missing type",
			values: map[string]interface{}{
				"package_id":   "bar",
				"package_size": "123",
			},
			expErr: MissingValueError("type"),
		},
		{
			name: "missing package_id",
			values: map[string]interface{}{
				"type":         "package_allocate",
				"package_size": "123",
			},
			expErr: MissingValueError("package_id"),
		},
		{
			name: "missing package_size",
			values: map[string]interface{}{
				"type":       "package_allocate",
				"package_id": "bar",
			},
			expErr: MissingValueError("package_size"),
		},
		{
			name: "invalid type",
			values: map[string]interface{}{
				"type":         123,
				"package_id":   "bar",
				"package_size": "123",
			},
			expErr: ValueTypeError{
				val: "type",
				v:   123,
			},
		},
		{
			name: "invalid package_id",
			values: map[string]interface{}{
				"type":         "package_allocate",
				"package_id":   123,
				"package_size": "123",
			},
			expErr: ValueTypeError{
				val: "package_id",
				v:   123,
			},
		},
		{
			name: "invalid package_size",
			values: map[string]interface{}{
				"type":         "package_allocate",
				"package_id":   "bar",
				"package_size": 1.23,
			},
			expErr: ValueTypeError{
				val: "package_size",
				v:   1.23,
			},
		},
		{
			name: "parse package_size",
			values: map[string]interface{}{
				"type":         "package_allocate",
				"package_id":   "bar",
				"package_size": "baz",
			},
			expErr: &strconv.NumError{
				Func: "Atoi",
				Num:  "baz",
				Err:  strconv.ErrSyntax,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgEvt, err := FromValues("test", tt.values)
			if tt.expErr != nil {
				require.EqualError(t, err, tt.expErr.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expPkgEvt, pkgEvt)
		})
	}
}

func TestPackageEvent_Values(t *testing.T) {
	pkgEvt := PackageEvent{
		EventID:     "foo",
		Typ:         PackageAllocate,
		PackageID:   "bar",
		PackageSize: 123,
	}

	require.Equal(
		t,
		pkgEvt.Values(),
		map[string]interface{}{
			"type":         "package_allocate",
			"package_id":   "bar",
			"package_size": 123,
		},
	)
}
