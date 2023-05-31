package runner

import (
	"fmt"
	"log"
	"sync"
	"bufio"
	"os"
	"strings"
	"sort"
	"path/filepath"
	"github.com/zisadlier/linecount/internal/pkg/types"
	"github.com/zisadlier/linecount/internal/pkg/utils"
)

var version = "v1.0.0"

type Runner struct {
	Options *types.Options // Execution options
	FilesToCount map[string]int // Mapping of filepath to line count
	Mutex sync.RWMutex // Write mutex for FilesToCount
	TotalLines int // Total lines of counted files
}

type File struct {
	Name string
	Lines int
}

func New(options *types.Options) (*Runner, error) {

	if errs := options.Validate(); len(errs) > 0 {
		return nil, &types.OptionsValidationErrorGroup{ Errors: errs }
	}

	if options.Version {
		fmt.Printf("Version: %s\n", version)
		return nil, nil
	}

	runner := &Runner{
		Options: options,
		FilesToCount: make(map[string]int),
		Mutex: sync.RWMutex{},
		TotalLines: 0,
	}

	return runner, nil
}

func (r *Runner) ProcessFile(path string, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	file, fileErr := os.Open(path)
	if fileErr != nil {
		log.Println(fmt.Sprintf("Could not open %s due to %s, skipping", path, fileErr))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fileLines := 0
	for scanner.Scan() {
		text := scanner.Text()
		regex := utils.ShouldLineBeCounted(text, r.Options.SkipRegex, r.Options.IncludeRegex)
		whiteSpace := r.Options.CountWhiteSpace || strings.TrimSpace(text) != ""
		if whiteSpace && regex {
			fileLines++
		}
	}

	r.Mutex.Lock()
	r.FilesToCount[path] = fileLines
	r.TotalLines += fileLines
	r.Mutex.Unlock()
}

func (r *Runner) GetFiles() []*File {
	// Get all paths and sort them
	names := make([]string, len(r.FilesToCount))
	i := 0
	for k := range r.FilesToCount {
		names[i] = k
		i++
	}
	sort.Sort(sort.StringSlice(names))

	// Create slice of file data
	files := make([]*File, len(names))
	i = 0
	for _, n := range names {
		files[i] = &File{
			Name: n,
			Lines: r.FilesToCount[n],
		}
		i++
	}

	return files
}

func (r *Runner) Run() error {
	var err error = nil
	if r.Options.HasFile() {
		r.ProcessFile(r.Options.File, nil)
	} else if r.Options.HasDirectory()  {
		var wg sync.WaitGroup
		err = filepath.Walk(r.Options.Directory, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if err != nil {
				return err
			}

			if r.Options.ShouldFileBeProcessed(path) {
				wg.Add(1)
				go r.ProcessFile(path, &wg)
			}

			return nil
		})
		wg.Wait()
	}

	return err
}