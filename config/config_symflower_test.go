package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSymflowerGetInstance(t *testing.T) {
	type testCase struct {
		Name string

		ExpectedConfig *Config
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualConfig := GetInstance()

			assert.Equal(t, tc.ExpectedConfig, actualConfig)
		})
	}

	validate(t, &testCase{
		ExpectedConfig: &Config{Server: nil, Client: nil},
	})
	validate(t, &testCase{
		ExpectedConfig: &Config{Server: nil, Client: nil},
	})
}
