package peanut

import "testing"

func TestSupportedKinds(t *testing.T) {
	for k := range supportedKind {
		// SQLiteWriter's lookup table should have entries for all supported types.
		if _, ok := kindToDBType[k]; !ok {
			t.Fail()
		}
	}
}
