package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
    "time"

    "os/exec"

    "github.com/cucumber/godog"
    "github.com/cucumber/godog/gherkin"
)

var InterceptedStdout bytes.Buffer

func executableInstalledOnLocalMachine(name string) error {
    _, err := exec.LookPath(name)
    if err != nil {
        return fmt.Errorf("cannot find %s on local machine: %v", name, err)
    }

    // return godog.ErrPending
    return nil
}

func iRunCmd(desc string, input *gherkin.DocString) error {
    InterceptedStdout.Reset()
    cmd, err := runCmd(input.Content, &InterceptedStdout)
    if err != nil {
        return fmt.Errorf("failed to run %s: %v", cmd, err)
    }

    return nil
}

func iRunCmdStr(desc, input string) error {
    InterceptedStdout.Reset()
    cmd, err := runCmd(input, &InterceptedStdout)
    if err != nil {
        return fmt.Errorf("failed to run %s: %v", cmd, err)
    }

    return nil
}

func iCanRunCmd(desc string, input *gherkin.DocString) error {
    return iRunCmd(desc, input)
}

func iCanRunCmdStr(desc, input string) error {
    return iRunCmdStr(desc, input)
}

func runCmd(input string, out io.Writer) (string, error) {
    var cmd *exec.Cmd
    cmd = exec.Command("bash", "-c", input)
    // switch isPipe {
    // case true: 
    //     cmd = exec.Command("bash", "-c", input)
    // case false:
    //     args := strings.Fields(input)
    //     cmd = exec.Command(args[0], args[1:]...)
    // }
    cmd.Stdout = out
    if err := cmd.Run(); err != nil {
        return cmd.String(), err
    }
    return cmd.String(), nil
}

func outputPrints(expectedDocString *gherkin.DocString) error {
    var actual = InterceptedStdout.String()
    var expected = expectedDocString.Content
    if expected != actual {
        return fmt.Errorf("\nexpected:\n%s\nactual:\n%s", expected, actual)
    }
    return nil
}

func outputContainsStringInLines(match string) error {
    lines := strings.Split(InterceptedStdout.String(),"\n")

    for _, line := range lines {
        if len(line) == 0 {
            continue
        }
        if !strings.Contains(line, match) {
            return fmt.Errorf("expect output to contain %s\nactual:\n%s", match, line)
        }
    }
    return nil
}

func createK8sResourceWithManifest(resource string, manifest *gherkin.DocString) error {
    return iRunCmdStr("", fmt.Sprintf("cat <<EOF | kubectl apply -f -\n%s\nEOF", manifest.Content))
}

func waitForResourceCondition(resource, match string, input *gherkin.DocString) error {
    err := retry(160, 5*time.Second, func() (err error) {
        iRunCmdStr("", input.Content)
        err = outputContainsStringInLines(match)
        return
    })
    if err != nil {
        return err
    }
    return nil
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
    for i := 0; ; i++ {
        err = callback()
        if err == nil {
            return
        }

        if i >= (attempts - 1) {
            break
        }

        time.Sleep(sleep)

        // fmt.Printf("retrying after error:%v\n", err)
    }
    return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func GeneralContext(s *godog.Suite) {
    s.Step(`^output prints:$`, outputPrints)
    s.Step(`^"([^"]*)" installed on local machine$`, executableInstalledOnLocalMachine)
    s.Step(`^I can run command to "([^"]*)":$`, iCanRunCmd)
    s.Step(`^I run command to "([^"]*)":$`, iRunCmd)
    s.Step(`^output contains "([^"]*)"$`, outputContainsStringInLines)
    s.Step(`^I create "([^"]*)" with manifest:$`, createK8sResourceWithManifest)
    s.Step(`^I wait until "([^"]*)" "([^"]*)":$`, waitForResourceCondition)
}
