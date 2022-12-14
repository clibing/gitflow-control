package internal

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRg(t *testing.T) {
	configV1 := &Config{
		Mode:  "auto",
		Issue: &Issue{},
	}
	messagev1 := `chore(pom): add pom dep version

add pom dep version

Signed-off-by: clibing <wmsjhappy@gmail.com>`

	ccm(messagev1, configV1)

	configV2 := &Config{
		Mode: "first",
		Issue: &Issue{
			LeftMarker:  "[",
			RightMarker: "]",
		},
	}
	messagev2 := `[wback-11]

chore(pom): add pom dep version

add pom dep version

Signed-off-by: clibing <wmsjhappy@gmail.com>`

	ccm(messagev2, configV2)
	ccm(messagev2, configV1)
}

func ccm(message string, config *Config) error {
	rg := commitMessageCheckPatternV1
	if config.Mode == "first" {
		rg = fmt.Sprintf(commitMessageCheckPatternV2, config.Issue.LeftMarker, config.Issue.RightMarker)
	}
	// 增加 commit-msg hook时使用
	reg := regexp.MustCompile(rg)

	msgs := reg.FindStringSubmatch(message)
	if config.Mode == "first" {
		if len(msgs) != 4 {
			return fmt.Errorf(commitMessageCheckFailedMsgV2, config.Issue.LeftMarker, config.Issue.RightMarker)
		}
	} else {
		if len(msgs) != 4 {
			return fmt.Errorf(commitMessageCheckFailedMsgV1)
		}
	}

	return nil
}

func TestProjectName(t *testing.T) {
	v, e := GetProjectName()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(v)
}
