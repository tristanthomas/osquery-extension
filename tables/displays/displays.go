package displays

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kolide/osquery-go/plugin/table"
	"github.com/pkg/errors"
)

type DisplayInformation struct {
	InternalDisplay   int
	TotalDisplayCount int
}

func DisplayInformationColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("internal_display"),
		table.TextColumn("total_display_count"),
	}
}

// getSystemReport returns standard output of system_profiler
func getSystemReport(arg ...string) ([]byte, error) {
	cmd := exec.Command("/usr/sbin/system_profiler", arg...)
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "calling /usr/sbin/system_profiler to get system report")
	}
	return out, nil
}

func DisplayInformationGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	displayReport, err := getSystemReport("SPDisplaysDataType")
	if err != nil {
		return nil, err
	}
	displayReportString := string(displayReport)
	info := DisplayInformation{
		InternalDisplay:   strings.Count(displayReportString, "Connection Type: Internal"),
		TotalDisplayCount: strings.Count(displayReportString, "Connection Type:"),
	}
	results := make([]map[string]string, 0)
	results = append(results, map[string]string{
		"internal_display":    strconv.Itoa(info.InternalDisplay),
		"total_display_count": strconv.Itoa(info.TotalDisplayCount),
	})
	return results, nil
}
