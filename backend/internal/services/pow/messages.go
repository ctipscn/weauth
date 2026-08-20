package pow

import "errors"

var (
	ErrChallengeNotFound = errors.New("challenge not found")
	ErrChallengeExpired  = errors.New("challenge expired")
	ErrChallengeSolved   = errors.New("challenge already solved")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenUsed         = errors.New("token already used")
	ErrTokenExpired      = errors.New("token expired")
)