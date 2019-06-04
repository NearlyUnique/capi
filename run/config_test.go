package run

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func Test_can_use_ioutil_ReadFile_as_reader_to_read_from_disk(t *testing.T) {
	const noHome = ""
	pwd := os.TempDir()
	err := os.Chdir(pwd)
	require.NoError(t, err)
	filename := "apiset.json"
	_ = os.Remove(filename)
	err = ioutil.WriteFile(filename, []byte(`{"apis":[{"Name":"any-name"}]}`), 0666)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	loader := NewConfigLoader(noHome, ioutil.ReadFile)
	loader.RegisterFileExtension(".json", JSONFormatReader)

	cfg, err := loader.Load("")

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "any-name", cfg.APIs[0].Name)
}

func Test_config_can_be_loaded_using_json_reader(t *testing.T) {
	const noFilename = ""
	const noHome = ""
	t.Run("formats must be registered", func(t *testing.T) {
		noFileReader := func(filename string) ([]byte, error) { return nil, nil }

		loader := NewConfigLoader(noHome, noFileReader)
		cfg, err := loader.LoadRaw(noFilename)

		assert.Error(t, err)
		_, ok := err.(builder.InvalidOperation)
		assert.True(t, ok)
		assert.Nil(t, cfg)
	})
	t.Run("when no name given, defaults to apiset.json", func(t *testing.T) {
		fakeReader := func(filename string) ([]byte, error) {
			assert.Equal(t, "apiset.json", filename)
			return []byte(`some-content`), nil
		}

		loader := NewConfigLoader(noHome, fakeReader)
		loader.RegisterFileExtension(".json", JSONFormatReader)
		cfg, err := loader.LoadRaw(noFilename)

		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "some-content", string(cfg.Data))
	})
	t.Run("search in order apiset, filename locally then in homefolder", func(t *testing.T) {
		expectedSearch := []string{
			"apiset.json",
			"api-name.json",
			path.Join("home-dir", "apiset.json"),
			path.Join("home-dir", "api-name.json"),
		}
		fakeReader := func(filename string) ([]byte, error) {
			assert.Equal(t, expectedSearch[0], filename)
			expectedSearch = expectedSearch[1:]
			return nil, xerrors.New("not found")
		}

		loader := NewConfigLoader("home-dir", fakeReader)
		loader.RegisterFileExtension(".json", func([]byte) (set *builder.APISet, e error) {
			t.Fail()
			return nil, nil
		})
		_, err := loader.LoadRaw("api-name")
		require.Error(t, err)

		assert.Empty(t, expectedSearch)
	})
}
