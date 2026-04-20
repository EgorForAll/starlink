package users

import "fmt"

type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string { return fmt.Sprintf("retryable: %v", e.Err) }
func (e *RetryableError) Unwrap() error { return e.Err }

type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string { return fmt.Sprintf("non-retryable: %v", e.Err) }
func (e *NonRetryableError) Unwrap() error { return e.Err }
