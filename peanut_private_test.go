package peanut

import "testing"

func TestSupportedTypes(t *testing.T) {
	for k := range supportedType {
		// SQLiteWriter's lookup table should have entries for all supported types.
		if _, ok := kindToDBType[k]; !ok {
			t.Fail()
		}
	}
}
