package runner

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/result"
	fileutil "github.com/zan8in/pins/file"
)

type FileType = uint8

const (
	FILE_TXT = iota
	FILE_JSON
	FILE_CSV
	NOT_FOUND
)

func FileExt(filename string) FileType {
	ext := path.Ext(filename)
	switch ext {
	case ".txt":
		return FILE_TXT
	case ".json":
		return FILE_JSON
	case ".csv":
		return FILE_CSV
	default:
		return NOT_FOUND
	}
}

func WriteOutput(results *result.Result) {
	if !results.HasResult() {
		return
	}

	var (
		file     *os.File
		output   string
		err      error
		fileType uint8
		csvutil  *csv.Writer
	)

	output = fmt.Sprintf("output-%d.csv", time.Now().UnixMilli())

	fileType = FileExt(output)

	outputFolder := filepath.Dir(output)
	if fileutil.FolderExists(outputFolder) {
		mkdirErr := os.MkdirAll(outputFolder, 0700)
		if mkdirErr != nil {
			gologger.Error().Msgf("Could not create output folder %s: %s\n", outputFolder, mkdirErr)
			return
		}
	}

	file, err = os.Create(output)
	if err != nil {
		gologger.Error().Msgf("Could not create file %s: %s\n", output, err)
		return
	}
	defer file.Close()

	if fileType == FILE_CSV {
		csvutil = csv.NewWriter(file)
		file.WriteString("\xEF\xBB\xBF")
	}

	csvutil.Write([]string{"Query conditions: ", results.GetQuery()})
	csvutil.Flush()

	for rs := range results.GetResult() {
		csvutil.Write(rs)
		csvutil.Flush()
	}

}
