package install_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestInstallCleanMachineJourney(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("the initial supported platform is Linux amd64")
	}
	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "built-pathframe")
	build := exec.Command("go", "build", "-trimpath", "-ldflags", "-X github.com/0xkhdr/pathframe/internal/app.Version=install-test", "-o", source, "./cmd/pathframe")
	build.Dir = repository
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %s: %v", output, err)
	}
	target := filepath.Join(dir, "bin", "pathframe")
	run(t, filepath.Join(repository, "scripts", "install.sh"), "install", source, target)
	command := exec.Command(target)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "Pathframe") || !strings.Contains(string(output), "Configured: false") {
		t.Fatalf("clean-machine orientation = %s, %v", output, err)
	}
}

func TestDownloadInstallVerifiesRelease(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("the initial supported platform is Linux amd64")
	}
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	dir := t.TempDir()
	release := filepath.Join(dir, "release")
	if err := os.Mkdir(release, 0o755); err != nil {
		t.Fatal(err)
	}
	executable(t, release, "pathframe", "#!/bin/sh\necho downloaded\n")
	archive := filepath.Join(release, "pathframe_linux_amd64.tar.gz")
	command := exec.Command("tar", "-czf", archive, "pathframe")
	command.Dir = release
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("archive: %s: %v", output, err)
	}
	sum := exec.Command("sha256sum", "pathframe_linux_amd64.tar.gz")
	sum.Dir = release
	checksum, err := sum.Output()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive+".sha256", checksum, 0o644); err != nil {
		t.Fatal(err)
	}
	fakebin := filepath.Join(dir, "fakebin")
	if err := os.Mkdir(fakebin, 0o755); err != nil {
		t.Fatal(err)
	}
	executable(t, fakebin, "curl", "#!/bin/sh\nwhile [ \"$1\" != -o ]; do shift; done\nout=$2\nshift 2\ncp \"$PATHFRAME_TEST_RELEASE/${1##*/}\" \"$out\"\n")
	targetDir := filepath.Join(dir, "bin")
	command = exec.Command(filepath.Join(repository, "scripts", "install.sh"))
	environment := append(os.Environ(), "PATH="+fakebin+":"+os.Getenv("PATH"), "PATHFRAME_TEST_RELEASE="+release, "PATHFRAME_INSTALL_DIR="+targetDir)
	command.Env = environment
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("download install: %s: %v", output, err)
	}
	assertOutput(t, filepath.Join(targetDir, "pathframe"), "downloaded\n")
	if err := os.WriteFile(archive+".sha256", []byte(strings.Repeat("0", 64)+"  pathframe_linux_amd64.tar.gz\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	command = exec.Command(filepath.Join(repository, "scripts", "install.sh"))
	command.Env = environment
	if err := command.Run(); err == nil {
		t.Fatal("installer accepted a release with a mismatched checksum")
	}
	assertOutput(t, filepath.Join(targetDir, "pathframe"), "downloaded\n")
}

func TestInstallUpdateFailureAndUninstall(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the initial installer supports the claimed Linux platform")
	}
	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(repository, "scripts", "install.sh")
	dir := t.TempDir()
	target := filepath.Join(dir, "bin", "pathframe")
	v1, v2 := executable(t, dir, "v1", "#!/bin/sh\necho v1\n"), executable(t, dir, "v2", "#!/bin/sh\necho v2\n")

	run(t, script, "install", v1, target)
	assertOutput(t, target, "v1\n")
	run(t, script, "install", v2, target)
	assertOutput(t, target, "v2\n")

	missing := filepath.Join(dir, "missing")
	command := exec.Command(script, "install", missing, target)
	if err := command.Run(); err == nil {
		t.Fatal("failed update unexpectedly succeeded")
	}
	assertOutput(t, target, "v2\n")

	run(t, script, "uninstall", target)
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("installed binary remains after uninstall: %v", err)
	}
}

func TestInstallRejectsSymlinkTarget(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the initial installer supports the claimed Linux platform")
	}
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	script := filepath.Join(repository, "scripts", "install.sh")
	dir := t.TempDir()
	source := executable(t, dir, "source", "#!/bin/sh\necho safe\n")
	outside := filepath.Join(dir, "outside")
	if err := os.WriteFile(outside, []byte("unchanged"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "pathframe")
	if err := os.Symlink(outside, target); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command(script, "install", source, target).Run(); err == nil {
		t.Fatal("installer accepted symlink target")
	}
	data, _ := os.ReadFile(outside)
	if string(data) != "unchanged" {
		t.Fatalf("symlink target was changed: %q", data)
	}
}

func TestInstallInterruptedUpdatePreservesPrevious(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the initial installer supports the claimed Linux platform")
	}
	repository, _ := filepath.Abs(filepath.Join("..", ".."))
	script := filepath.Join(repository, "scripts", "install.sh")
	dir := t.TempDir()
	target := filepath.Join(dir, "bin", "pathframe")
	v1, v2 := executable(t, dir, "v1", "#!/bin/sh\necho v1\n"), executable(t, dir, "v2", "#!/bin/sh\necho v2\n")
	run(t, script, "install", v1, target)

	fakebin := filepath.Join(dir, "fakebin")
	if err := os.Mkdir(fakebin, 0o755); err != nil {
		t.Fatal(err)
	}
	executable(t, fakebin, "cp", "#!/bin/sh\ntouch \"$PATHFRAME_CP_STARTED\"\nsleep 60\n")
	started := filepath.Join(dir, "copy-started")
	command := exec.Command(script, "install", v2, target)
	command.Env = append(os.Environ(), "PATH="+fakebin+":/usr/bin:/bin", "PATHFRAME_CP_STARTED="+started)
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(started); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-command.Process.Pid, syscall.SIGKILL)
			t.Fatal("installer did not reach staged copy")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := syscall.Kill(-command.Process.Pid, syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	_ = command.Wait()
	assertOutput(t, target, "v1\n")
	if leftovers, _ := filepath.Glob(filepath.Join(filepath.Dir(target), ".pathframe-install.*")); len(leftovers) != 0 {
		t.Fatalf("interrupted install left staging files: %v", leftovers)
	}
}

func executable(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func run(t *testing.T, script string, args ...string) {
	t.Helper()
	if output, err := exec.Command(script, args...).CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %s: %v", script, args, output, err)
	}
}

func assertOutput(t *testing.T, binary, want string) {
	t.Helper()
	output, err := exec.Command(binary).CombinedOutput()
	if err != nil || string(output) != want {
		t.Fatalf("%s output = %q, %v", binary, output, err)
	}
}
