package run

import (
	"os"
	"path"
	"sort"
	"strings"

	"github.com/NearlyUnique/capi/builder"
	"golang.org/x/xerrors"
)

type (
	//FileReader interface to arbitrary file reader
	FileReader func(filename string) ([]byte, error)
	// FormatReader can convert bytes to an APISet
	FormatReader func(content []byte) (*builder.APISet, error)
	// ConfigLoader find a config and load it
	ConfigLoader struct {
		formats map[string]FormatReader
		home    string
		reader  FileReader
	}
	rawConfig struct {
		Format string
		Data   []byte
	}
)

// NewConfigLoader
func NewConfigLoader(home string, reader FileReader) ConfigLoader {
	loader := ConfigLoader{
		reader:  reader,
		formats: make(map[string]FormatReader),
		home:    home,
	}
	return loader
}

func (loader ConfigLoader) RegisterFileExtension(extn string, reader FormatReader) {
	loader.formats[extn] = reader
}

func (loader ConfigLoader) Load(filename string) (*builder.APISet, error) {
	raw, err := loader.LoadRaw(filename)
	if err != nil {
		return nil, err
	}

	return loader.formats[raw.Format](raw.Data)
}

// LoadRaw using filename, with extn json|xml|yaml or blank
// if blank, then look for file called apiset with the same file extns
func (loader ConfigLoader) LoadRaw(filename string) (*rawConfig, error) {
	if len(loader.formats) == 0 {
		return nil, builder.InvalidOperation("no formats registered")
	}
	targets := []string{
		filename,
		path.Join(loader.home, filename),
	}
	var extn string
	var buf []byte
	for _, fname := range targets {
		extn, buf = loader.tryOpen(fname)
		if extn != "" && buf != nil {
			break
		}
	}
	if buf == nil {
		return nil, xerrors.New("no config found")
	}

	return &rawConfig{Format: extn, Data: buf}, nil
}

func (loader ConfigLoader) List(search string) []string {
	var keys []string
	for k := range loader.formats {
		keys = append(keys, k)
	}
	root, err := os.Getwd()
	if err != nil {
		return nil
	}
	list, err := osReadDir(root, search, keys)
	if err != nil {
		return []string{"error", err.Error()}
	}
	return list
}

func osReadDir(root, search string, extns []string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	defer func() { _ = f.Close() }()
	fileInfo, err := f.Readdir(-1)
	if err != nil {
		return files, err
	}

	sort.SliceStable(extns, sortByLongestFirst(extns))

	for _, file := range fileInfo {
		if file.IsDir() {
			continue
		}
		_, name := path.Split(file.Name())
		if strings.HasPrefix(name, search) {
			for _, extn := range extns {
				if strings.HasSuffix(name, extn) {
					files = append(files, name[:len(name)-len(extn)])
					break
				}
			}
		}
	}

	return files, nil
}

func (loader ConfigLoader) tryOpen(filename string) (string, []byte) {
	if filename == "" {
		return "", nil
	}
	for extn := range loader.formats {
		buf, err := loader.reader(filename + extn)
		if err == nil {
			return extn, buf
		}
	}
	return "", nil
}

func indexOrEmpty(args []string, i int) string {
	if i < 0 || i >= len(args) {
		return ""
	}
	return args[i]
}

func sortByLongestFirst(extns []string) func(i, j int) bool {
	return func(i, j int) bool { return len(extns[i]) > len(extns[j]) }
}
