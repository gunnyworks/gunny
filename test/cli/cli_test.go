// All tests in this package assume that the gunny executable has been built
// and is available to run from build/gunny in the root of the repository.
//
// Multiple CLI test cases can be housed in a single directory. Each test case
// consists of multiple files of the same file name, but different file
// extensions. For example:
//
//   - test-case1.in     - CLI args for test-case1 (one arg per line)
//   - test-case1.stdin  - Optional data to pipe in via stdin (if missing, no data will be piped in)
//   - test-case1.yaml   - Optional pipeline configuration for test-case1
//   - test-case1.json   - Optional pipeline configuration for test-case1
//   - test-case1.stdout - Expected stdout from running Gunny with the given inputs for test-case1
//   - test-case1.out    - Optional expected output file content for test-case1
//
// NOTE: One newline is trimmed off of expected stdout content. If newlines are
// to be expected in test case output, be sure to cater for this.
package cli_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type cliTestCase struct {
	name                      string
	args                      []string
	stdinData                 string
	expectedStdout            string
	outputFile                string
	expectedOutputFileContent string
}

func TestWithNoPipelineConfig(t *testing.T) {
	testCases, err := loadTestCases("./data/nopipelineconfig/")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stdout, stderr, err := runGunny(testCase.args, testCase.stdinData)
			require.NoError(t, err, stderr)
			require.Equal(t, testCase.expectedStdout, stdout)
		})
	}
}

func TestWithPipelineConfig(t *testing.T) {
	testCases, err := loadTestCases("./data/pipelineconfig/")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tempDir := t.TempDir()
			args := append(testCase.args, "--output-base-path", tempDir)
			stdout, stderr, err := runGunny(args, testCase.stdinData)
			require.NoError(t, err, stderr)
			require.Equal(t, testCase.expectedStdout, stdout)

			if len(testCase.outputFile) > 0 {
				// We expect the output file to be in the temp dir we created.
				outputFileBytes, err := os.ReadFile(filepath.Join(tempDir, testCase.outputFile))
				require.NoError(t, err)
				require.Equal(t, testCase.expectedOutputFileContent, string(outputFileBytes))
			}
		})
	}
}

func runGunny(args []string, stdinData string) (string, string, error) {
	cmd := exec.Command("../../build/gunny", args...)
	if len(stdinData) > 0 {
		cmd.Stdin = strings.NewReader(stdinData)
	}
	var stderrBytes bytes.Buffer
	cmd.Stderr = &stderrBytes
	output, err := cmd.Output()
	if err != nil {
		return "", stderrBytes.String(), fmt.Errorf("failed to run Gunny: %w", err)
	}
	return string(output), stderrBytes.String(), nil
}

func loadTestCases(path string) ([]*cliTestCase, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	testCases := make([]*cliTestCase, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".in" {
			continue
		}

		inputData, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		args := strings.Split(strings.TrimSpace(string(inputData)), "\n")

		filename := strings.TrimSuffix(file.Name(), ext)
		stdinDataFile := fmt.Sprintf("%s.stdin", filename)
		expectedStdoutFile := fmt.Sprintf("%s.stdout", filename)
		expectedOutputFile := fmt.Sprintf("%s.out", filename)

		stdinData := ""
		stdinDataBytes, err := os.ReadFile(filepath.Join(path, stdinDataFile))
		if err == nil {
			// We have data to pipe into stdin.
			stdinData = string(stdinDataBytes)
		}

		expectedStdout, err := os.ReadFile(filepath.Join(path, expectedStdoutFile))
		if err != nil {
			return nil, err
		}

		outputFile := ""
		expectedOutputFileContent := ""
		outputFileDataBytes, err := os.ReadFile(filepath.Join(path, expectedOutputFile))
		if err == nil {
			// We expect an output file.
			outputFile = expectedOutputFile
			expectedOutputFileContent = string(outputFileDataBytes)
		}

		testCases = append(testCases, &cliTestCase{
			name:                      filename,
			args:                      args,
			stdinData:                 stdinData,
			expectedStdout:            strings.TrimSuffix(string(expectedStdout), "\n"),
			outputFile:                outputFile,
			expectedOutputFileContent: expectedOutputFileContent,
		})
	}
	return testCases, nil
}
