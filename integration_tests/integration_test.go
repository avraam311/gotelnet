package integrationtests

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const binName = "gotelnet"

func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/gotelnet")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\\n%s", err, string(out))
	}
	return binPath
}

func runCmd(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Start()
	time.Sleep(2 * time.Second) 
	cmd.Process.Kill()
	time.Sleep(1 * time.Second)
	return strings.TrimRight(out.String(), "\\n")
}

func TestConnection(t *testing.T) {
	t.Parallel()
	bin := buildBinary(t)
	out := runCmd(t, bin, "smtp.freesmtpservers.com", "25", "--timeout", "5")
	expected := "connected to server successfully"
	if !strings.Contains(out, expected) {
		t.Errorf("unexpected output:\\nGot:\\n%s\\nWant contains:\\n%s", out, expected)
	}
}
