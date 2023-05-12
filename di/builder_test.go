package di_test

import (
	"errors"
	"testing"

	"github.com/kukymbr/core2go/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build_WhenValid_ExpectNoError(t *testing.T) {
	builder := &di.Builder{}

	items := map[string]string{
		"testname1": "testval1",
		"testname2": "testval2",
		"testname3": "testval3",
	}

	for name, val := range items {
		v := val
		err := builder.Add(
			di.Def{
				Name: name,
				Build: func(ctn *di.Container) (any, error) {
					return v, nil
				},
				Validate: func(ctn *di.Container) (err error) {
					return nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		assert.NoError(t, err)
	}

	container, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, container)
	require.Equal(t, 3, container.Len())

	for name, expected := range items {
		val, err := container.SafeGet(name)
		assert.NoError(t, err)
		assert.Equal(t, expected, val)
	}

	err = container.Close()
	assert.NoError(t, err)
}

func TestBuilder_Build_WhenError_ExpectError(t *testing.T) {
	builder := &di.Builder{}

	err := builder.Add(
		di.Def{
			Name: "testname1",
			Build: func(ctn *di.Container) (any, error) {
				return "testname2", nil
			},
		},
	)
	require.NoError(t, err)

	tests := []di.Def{
		{Name: "testname1"},
		{
			Name: "testname2",
			Build: func(ctn *di.Container) (obj any, err error) {
				return "testval2", nil
			},
			Validate: func(ctn *di.Container) (err error) {
				return errors.New("failed to validate")
			},
		},
	}

	for i, test := range tests {
		err = builder.Add(test)
		assert.Error(t, err, i)
	}

	err = builder.Add(di.Def{
		Name: "testname3",
		Build: func(ctn *di.Container) (obj any, err error) {
			return "testval3", errors.New("failed to build")
		},
	})
	require.NoError(t, err)

	container, err := builder.Build()
	assert.Error(t, err)
	assert.Nil(t, container)

	err = builder.Add(di.Def{Name: "testname4"})
	require.NoError(t, err)

	container, err = builder.Build()
	assert.Error(t, err)
	assert.Nil(t, container)
}
