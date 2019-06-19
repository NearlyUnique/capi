package run

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

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
	filename := "some-file.json"
	_ = os.Remove(filename)
	err = ioutil.WriteFile(filename, []byte(`{"apis":[{"Name":"any-name"}]}`), 0666)
	require.NoError(t, err)
	defer func() { _ = os.Remove(filename) }()

	loader := NewConfigLoader(noHome, ioutil.ReadFile)
	loader.RegisterFileExtension(".json", JSONFormatReader)

	cfg, err := loader.Load("some-file")

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "any-name", cfg.APIs[0].Name)
}

func Test_ConfigLoader_list_will_list_all_filenames_with_registered_extns(t *testing.T) {
	const noHome = ""
	noFileReader := func(filename string) ([]byte, error) { return nil, nil }

	pwd := os.TempDir()
	err := os.Chdir(pwd)
	require.NoError(t, err)
	searchPrefix := randomString(10, 11)
	expected := []string{
		searchPrefix + randomString(3, 20),
		searchPrefix + randomString(3, 20),
		searchPrefix + randomString(3, 20)}
	fakeFiles := []string{
		expected[0] + ".ext1",
		expected[1] + ".ext1",
		expected[2] + ".ext2",
		searchPrefix + ".ignore"}

	clean := createTestFilesAndCleanUp(t, fakeFiles)
	defer clean()

	loader := NewConfigLoader(noHome, noFileReader)
	loader.RegisterFileExtension(".ext1", JSONFormatReader)
	loader.RegisterFileExtension(".ext2", JSONFormatReader)

	options := loader.List(searchPrefix)

	require.Equal(t, 3, len(options))
	require.Contains(t, options, expected[0])
	require.Contains(t, options, expected[1])
	require.Contains(t, options, expected[2])
}

func Test_what_happens_when_one_extn_ends_with_another(t *testing.T) {
	const noHome = ""
	noFileReader := func(filename string) ([]byte, error) { return nil, nil }

	pwd := os.TempDir()
	err := os.Chdir(pwd)
	require.NoError(t, err)
	searchPrefix := randomString(10, 11)
	expected := []string{
		searchPrefix + randomString(3, 20),
		searchPrefix + randomString(3, 20),
	}
	fakeFiles := []string{
		expected[0] + ".json",
		expected[1] + ".some-long-extension-name.json",
	}

	clean := createTestFilesAndCleanUp(t, fakeFiles)
	defer clean()

	loader := NewConfigLoader(noHome, noFileReader)
	loader.RegisterFileExtension(".json", JSONFormatReader)
	loader.RegisterFileExtension(".some-long-extension-name.json", JSONFormatReader)

	options := loader.List(searchPrefix)

	require.Equal(t, 2, len(options))
	require.Contains(t, options, expected[0])
	require.Contains(t, options, expected[1])
}

func createTestFilesAndCleanUp(t *testing.T, filenames []string) func() {
	t.Helper()
	for _, fn := range filenames {
		_ = os.Remove(fn)
		err := ioutil.WriteFile(fn, []byte(`any`), 0666)
		require.NoError(t, err)
	}

	return func() {
		for _, fn := range filenames {
			_ = os.Remove(fn)
		}
	}
}

func Test_what_happens_when_the_same_filename_is_used_for_two_extns(t *testing.T) {
	t.Skip()
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
	t.Run("search in order filename locally then in homefolder", func(t *testing.T) {
		expectedSearch := []string{
			"api-name.json",
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

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(min, max int32) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"_-~+="
	r := rand.Int31n(int32(math.Abs(float64(max - min))))
	return stringWithCharset(min+r, charset)
}

func stringWithCharset(length int32, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
