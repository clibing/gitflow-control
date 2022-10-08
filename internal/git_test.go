package internal

import "testing"

func TestRg(t *testing.T) {
	configV1 := &Config{
		Issue: &Issue{
			FirstEnable: false,
		},
	}
	messagev1 := `chore(pom): add pom dep version

add pom dep version

Signed-off-by: clibing <wmsjhappy@gmail.com>`

	CheckCommitMessage(messagev1, configV1)

	configV2 := &Config{
		Issue: &Issue{
			FirstEnable: true,
			LeftMarker:  "[",
			RightMarker: "]",
		},
	}
	messagev2 := `[wback-11]chore(pom): add pom dep version

add pom dep version

Signed-off-by: clibing <wmsjhappy@gmail.com>`

	CheckCommitMessage(messagev2, configV2)
	CheckCommitMessage(messagev2, configV1)
}
