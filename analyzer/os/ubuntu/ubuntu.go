package ubuntu

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	aos "github.com/aquasecurity/fanal/analyzer/os"
	"github.com/aquasecurity/fanal/types"
	"github.com/aquasecurity/fanal/utils"
)

func init() {
	analyzer.RegisterAnalyzer(&ubuntuOSAnalyzer{})
}

var requiredFiles = []string{"etc/lsb-release"}

type ubuntuOSAnalyzer struct{}

func (a ubuntuOSAnalyzer) Analyze(content []byte) (analyzer.AnalyzeReturn, error) {
	isUbuntu := false
	scanner := bufio.NewScanner(bytes.NewBuffer(content))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "DISTRIB_ID=Ubuntu" {
			isUbuntu = true
			continue
		}

		if isUbuntu && strings.HasPrefix(line, "DISTRIB_RELEASE=") {
			return analyzer.AnalyzeReturn{
				OS: types.OS{
					Family: aos.Ubuntu,
					Name:   strings.TrimSpace(line[16:]),
				},
			}, nil
		}
	}
	return analyzer.AnalyzeReturn{}, xerrors.Errorf("ubuntu: %w", aos.AnalyzeOSError)
}

func (a ubuntuOSAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return utils.StringInSlice(filePath, requiredFiles)
}

func (a ubuntuOSAnalyzer) Name() string {
	return aos.Ubuntu
}
