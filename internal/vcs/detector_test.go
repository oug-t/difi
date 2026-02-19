package vcs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectVCS_GitPriority(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create both .git and .hg directories
	gitDir := filepath.Join(tempDir, ".git")
	hgDir := filepath.Join(tempDir, ".hg")

	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	if err := os.Mkdir(hgDir, 0755); err != nil {
		t.Fatalf("Failed to create .hg dir: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test detection - should prefer Git over Mercurial
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(GitVCS{}) {
		t.Errorf("Expected GitVCS, got %T", vcs)
	}
}

func TestDetectVCS_GitOnly(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create only .git directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test detection
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(GitVCS{}) {
		t.Errorf("Expected GitVCS, got %T", vcs)
	}
}

func TestDetectVCS_HgOnly(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create only .hg directory
	hgDir := filepath.Join(tempDir, ".hg")
	if err := os.Mkdir(hgDir, 0755); err != nil {
		t.Fatalf("Failed to create .hg dir: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test detection
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(HgVCS{}) {
		t.Errorf("Expected HgVCS, got %T", vcs)
	}
}

func TestDetectVCS_NoVCS(t *testing.T) {
	// Create temporary directory structure without any VCS
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test detection - should default to Git
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(GitVCS{}) {
		t.Errorf("Expected GitVCS (default), got %T", vcs)
	}
}

func TestDetectVCS_NestedDirectories(t *testing.T) {
	// Create temporary directory structure with nested directories
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir", "nested")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Create .git in parent directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Change to nested directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to nested dir: %v", err)
	}

	// Test detection - should find Git repository in parent
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(GitVCS{}) {
		t.Errorf("Expected GitVCS (from parent), got %T", vcs)
	}
}

func TestDetectVCS_GitInParentHgInChild(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create .git in parent directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Create .hg in child directory
	hgDir := filepath.Join(subDir, ".hg")
	if err := os.Mkdir(hgDir, 0755); err != nil {
		t.Fatalf("Failed to create .hg dir: %v", err)
	}

	// Change to child directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to subdir: %v", err)
	}

	// Test detection - should prefer Git (due to Git-first priority) over closer .hg
	// This tests our Git-first priority behavior
	vcs := DetectVCS()
	if reflect.TypeOf(vcs) != reflect.TypeOf(GitVCS{}) {
		t.Errorf("Expected GitVCS (priority), got %T", vcs)
	}
}

func TestVCSInterface_GitVCS(t *testing.T) {
	var vcs VCS = GitVCS{}

	// Test that GitVCS implements VCS interface
	// This is mainly a compile-time check, but we can test some basic calls
	_ = vcs.GetCurrentBranch()
	_ = vcs.GetRepoName()

	// Test that the interface methods exist and can be called
	// (actual functionality would require a git repo, so we just test the interface)
	files, _ := vcs.ListChangedFiles("main")
	if files == nil {
		files = []string{} // Just to use the variable
	}
}

func TestVCSInterface_HgVCS(t *testing.T) {
	var vcs VCS = HgVCS{}

	// Test that HgVCS implements VCS interface
	// This is mainly a compile-time check, but we can test some basic calls
	_ = vcs.GetCurrentBranch()
	_ = vcs.GetRepoName()

	// Test that the interface methods exist and can be called
	// (actual functionality would require an hg repo, so we just test the interface)
	files, _ := vcs.ListChangedFiles("default")
	if files == nil {
		files = []string{} // Just to use the variable
	}
}

func TestDetectVCS_ErrorHandling(t *testing.T) {
	// Test behavior when os.Getwd() might fail
	// We can't easily simulate os.Getwd() failure, but we can test
	// the default behavior when no VCS is found

	// This is more of a documentation test showing expected behavior
	vcs := DetectVCS()

	// Should always return something that implements VCS interface
	if vcs == nil {
		t.Error("DetectVCS() returned nil, should always return a VCS implementation")
	}

	// The default should be GitVCS when no specific VCS is detected
	if _, ok := vcs.(GitVCS); !ok {
		// Allow either GitVCS or HgVCS, depending on the environment
		if _, ok := vcs.(HgVCS); !ok {
			t.Errorf("DetectVCS() returned unexpected type %T", vcs)
		}
	}
}