package srun

import "errors"

var (
	ErrResultCannotFound = errors.New("result cannot found from response")
	ErrAcidCannotFound   = errors.New("acid not found")
	ErrEnvCannotFound    = errors.New("enc not found")
)
