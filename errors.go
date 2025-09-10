package main

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMissingKey      = errors.New("missing key")
	ErrInvalidScheme   = errors.New("invalid scheme")
	ErrSecRefTooShort  = errors.New("secret reference too short")
	ErrSecRefTooLong   = errors.New("secret reference too long")
	ErrVaultNotFound   = errors.New("vault not found")
	ErrItemNotFound    = errors.New("item not found")
	ErrEmptyVault      = errors.New("empty vault")
	ErrSectionNotFound = errors.New("section not found")
)
