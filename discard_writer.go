package peanut

var _ Writer = &DiscardWriter{}

// DiscardWriter is a Writer that does nothing.
type DiscardWriter struct{}

// Write does nothing and returns nil.
func (*DiscardWriter) Write(x interface{}) error { return nil }

// Close does nothing and returns nil.
func (*DiscardWriter) Close() error { return nil }

// Cancel does nothing and returns nil.
func (*DiscardWriter) Cancel() error { return nil }
