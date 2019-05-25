package log

import "testing"

func TestNewZeroLogger(t *testing.T) {
	l := NewZeroLogger()
	l.Info().
		Str("name", "Hung Le").
		Str("email", "hunglm@vzota.com.vn").
		Msg("Profile")
}
