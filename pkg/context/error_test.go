package context

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (suite *ErrorSuite) TestContextError() {
	cases := []struct {
		name       string
		ctxFactory contextFactory
		input      error
		expected   error
	}{
		{
			name: "Input error is context.Canceled",
			ctxFactory: func() context.Context {
				return context.Background()
			},
			input:    context.Canceled,
			expected: context.Canceled,
		},
		{
			name: "Input error is a timeout error",
			ctxFactory: func() context.Context {
				return context.Background()
			},
			input:    &timeoutError{true},
			expected: context.DeadlineExceeded,
		},
		{
			name: "Context error is context.DeadlineExceeded",
			ctxFactory: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now())
				defer cancel()
				return ctx
			},
			input:    errors.New("some error"),
			expected: context.DeadlineExceeded,
		},
		{
			name: "Context error is context.Canceled",
			ctxFactory: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			},
			input:    errors.New("some error"),
			expected: context.Canceled,
		},
		{
			name: "No error",
			ctxFactory: func() context.Context {
				return context.Background()
			},
			input:    errors.New("some error"),
			expected: nil,
		},
	}

	for _, tc := range cases {
		suite.T().Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.expected, ContextError(tc.ctxFactory(), tc.input))
		})
	}
}

func (suite *ErrorSuite) TestOutputError() {
	suite.T().Run("returns same error", func(t *testing.T) {
		input := errors.New("some error")
		assert.Equal(t, input, OutputError(context.Background(), input))
	})
	suite.T().Run("returns context error", func(t *testing.T) {
		input := errors.New("some error")
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.NotEqual(t, input, OutputError(ctx, errors.New("some error")))
	})
}

func (suite *ErrorSuite) TestOutputErrorf() {
	suite.T().Run("returns same error formatted without args", func(t *testing.T) {
		err := OutputErrorf(context.Background(), errors.New("some error"), "test")
		assert.EqualError(t, err, "test")
	})
	suite.T().Run("returns same error formatted with args", func(t *testing.T) {
		err := OutputErrorf(context.Background(), errors.New("some error"), "test %s and %s", "this", "that")
		assert.EqualError(t, err, "test this and that")
	})
	suite.T().Run("returns context error", func(t *testing.T) {
		input := errors.New("some error")
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.NotEqual(t, input, OutputErrorf(ctx, errors.New("some error"), "this won't be used"))
	})
}

// Utils

// Implements net.Error interface
type timeoutError struct {
	timeout bool
}

func (e *timeoutError) Error() string {
	return "timeout error"
}

func (e *timeoutError) Timeout() bool {
	return e.timeout
}

func (e *timeoutError) Temporary() bool {
	return e.timeout
}

type contextFactory func() context.Context
