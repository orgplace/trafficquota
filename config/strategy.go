//go:generate go-enum -f=$GOFILE --marshal

package config

// Strategy is an enumeration of token bucket strategies.
// ENUM(
// time-slice
// timestamp
// )
type Strategy int32
