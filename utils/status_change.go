package utils

// `StatusChange` is a struct with two fields, `To` and `Model`. The `To` field is a string, and the `Model` field is an
// interface{}.
// @property {string} To - The status to change to.
// @property Model - The model that is being changed.
type StatusChange struct {
	To    string
	Model interface{}
}
