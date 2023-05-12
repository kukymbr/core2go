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

	builder = &di.Builder{}

	err = builder.Add(di.Def{Name: "testname4"})
	require.NoError(t, err)

	container, err = builder.Build()
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestContainer(t *testing.T) {
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
			},
		)

		require.NoError(t, err)
	}

	container, err := builder.Build()
	require.NoError(t, err)

	assert.Equal(t, 3, container.Len())

	for name := range items {
		assert.True(t, container.Has(name))
	}

	def, err := container.SafeGet("unknown_def")
	assert.Error(t, err)
	assert.Nil(t, def)

	err = container.Close()
	assert.NoError(t, err)
}

func TestContainer_Close(t *testing.T) {
	builder := &di.Builder{}

	err := builder.Add(di.Def{
		Name: "testname1",
		Build: func(ctn *di.Container) (obj any, err error) {
			return "testval1", nil
		},
		Close: func(obj any) (err error) {
			return errors.New("close error")
		},
	})
	require.NoError(t, err)

	container, err := builder.Build()
	require.NoError(t, err)

	err = container.Close()
	assert.Error(t, err)
}

func TestContainer_Lazy(t *testing.T) {
	builder := &di.Builder{}

	err := builder.Add(
		di.Def{
			Name: "testname1",
			Build: func(ctn *di.Container) (obj any, err error) {
				return "testval1", nil
			},
			Close: func(obj any) (err error) {
				return errors.New("close error")
			},
			Lazy: true,
		},
		di.Def{
			Name: "testname2",
			Build: func(ctn *di.Container) (obj any, err error) {
				return "testval2", errors.New("build error")
			},
			Lazy: true,
		},
	)
	require.NoError(t, err)

	container, err := builder.Build()
	require.NoError(t, err)

	err = container.Close()
	assert.NoError(t, err)

	val := container.Get("testname1")
	assert.Equal(t, "testval1", val)

	err = container.Close()
	assert.Error(t, err)

	assert.Panics(t, func() {
		container.Get("testname2")
	})
}
