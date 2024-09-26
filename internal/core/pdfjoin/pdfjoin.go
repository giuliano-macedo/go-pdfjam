package pdfjoin

import (
	"fmt"
	"os"
	"os/exec"
)

func Join(inFiles []string, outFile string) error {
	pdfJamPath, found := os.LookupEnv("PDFJAM_PATH")
	if !found {
		pdfJamPath = "pdfjam"
	}

	cmdArgs := make([]string, 0, len(inFiles)+10)
	cmdArgs = append(cmdArgs, "--no-tidy", "--outfile", outFile)
	cmdArgs = append(cmdArgs, inFiles...)

	cmd := exec.Command(pdfJamPath, cmdArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("failed executing command %s: '%s'", err, string(out))
	}

	return err
}
