package linux

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	clabnodes "github.com/srl-labs/containerlab/nodes"
	clabtypes "github.com/srl-labs/containerlab/types"
)

func TestLinuxLinkApplyMode(t *testing.T) {
	if got := (&linux{}).LinkApplyMode(context.Background()); got != clabnodes.LinkApplyModeLive {
		t.Fatalf("LinkApplyMode() = %q, want %q", got, clabnodes.LinkApplyModeLive)
	}
}

func newTestNode(t *testing.T, cfg *clabtypes.NodeConfig) *linux {
	t.Helper()

	cfg.LabDir = filepath.Join(t.TempDir(), cfg.ShortName)
	if cfg.Sysctls == nil {
		cfg.Sysctls = map[string]string{}
	}

	n := new(linux)
	if err := n.Init(cfg); err != nil {
		t.Fatalf("Init: %v", err)
	}

	return n
}

func TestStartupConfigBind(t *testing.T) {
	tmpl := filepath.Join(t.TempDir(), "app.json.tmpl")
	if err := os.WriteFile(tmpl, []byte(
		`{"name":"{{ .ShortName }}","token":"{{ index .Config.Vars "token" }}"}`,
	), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	n := newTestNode(t, &clabtypes.NodeConfig{
		ShortName:     "app1",
		StartupConfig: tmpl,
		Config: &clabtypes.ConfigDispatcher{
			Vars: map[string]any{"token": "lab-token"},
		},
		KindSpecificConfig: &KindSpecificConfig{
			StartupConfigPath: "/config/app.json",
		},
	})

	wantBind := filepath.Join(n.Cfg.LabDir, "app.json") + ":/config/app.json:ro"
	found := false
	for _, b := range n.Cfg.Binds {
		if b == wantBind {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Binds = %v, want %q", n.Cfg.Binds, wantBind)
	}

	if err := n.PreDeploy(context.Background(), &clabnodes.PreDeployParams{}); err != nil {
		t.Fatalf("PreDeploy: %v", err)
	}

	got, err := os.ReadFile(n.Cfg.ResStartupConfig)
	if err != nil {
		t.Fatalf("read rendered config: %v", err)
	}
	want := `{"name":"app1","token":"lab-token"}`
	if string(got) != want {
		t.Fatalf("rendered config = %q, want %q", got, want)
	}
}

func TestStartupConfigPathNoOp(t *testing.T) {
	n := newTestNode(t, &clabtypes.NodeConfig{ShortName: "app1"})

	if len(n.Cfg.Binds) != 0 {
		t.Fatalf("Binds = %v, want none", n.Cfg.Binds)
	}
	if err := n.PreDeploy(context.Background(), &clabnodes.PreDeployParams{}); err != nil {
		t.Fatalf("PreDeploy: %v", err)
	}
	if _, err := os.Stat(filepath.Join(n.Cfg.LabDir, "app.json")); !os.IsNotExist(err) {
		t.Fatalf("unexpected lab file: %v", err)
	}
}

func TestStartupConfigSuppress(t *testing.T) {
	tmpl := filepath.Join(t.TempDir(), "app.json.tmpl")
	if err := os.WriteFile(tmpl, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	n := newTestNode(t, &clabtypes.NodeConfig{
		ShortName:             "app1",
		StartupConfig:         tmpl,
		SuppressStartupConfig: true,
		KindSpecificConfig: &KindSpecificConfig{
			StartupConfigPath: "/config/app.json",
		},
	})

	if len(n.Cfg.Binds) != 0 {
		t.Fatalf("Binds = %v, want none when suppress-startup-config is set", n.Cfg.Binds)
	}
	if err := n.PreDeploy(context.Background(), &clabnodes.PreDeployParams{}); err != nil {
		t.Fatalf("PreDeploy: %v", err)
	}
	if _, err := os.Stat(filepath.Join(n.Cfg.LabDir, "app.json")); !os.IsNotExist(err) {
		t.Fatalf("unexpected rendered file: %v", err)
	}
}

func TestStartupConfigPathInvalid(t *testing.T) {
	cases := []string{
		"config/app.json",
		"/config/app.json/",
		"/config:app.json",
	}

	for _, path := range cases {
		cfg := &clabtypes.NodeConfig{
			ShortName:     "app1",
			LabDir:        filepath.Join(t.TempDir(), "app1"),
			Sysctls:       map[string]string{},
			StartupConfig: "/tmp/ignored.cfg",
			KindSpecificConfig: &KindSpecificConfig{
				StartupConfigPath: path,
			},
		}
		n := new(linux)
		err := n.Init(cfg)
		if err == nil {
			t.Fatalf("Init(%q) = nil, want error", path)
		}
		if !strings.Contains(err.Error(), "startup-config-path") {
			t.Fatalf("Init(%q) error = %v, want startup-config-path mention", path, err)
		}
	}
}

func TestStartupConfigPathRequiresStartupConfig(t *testing.T) {
	cfg := &clabtypes.NodeConfig{
		ShortName: "app1",
		LabDir:    filepath.Join(t.TempDir(), "app1"),
		Sysctls:   map[string]string{},
		KindSpecificConfig: &KindSpecificConfig{
			StartupConfigPath: "/config/app.json",
		},
	}
	n := new(linux)
	err := n.Init(cfg)
	if err == nil {
		t.Fatal("Init = nil, want error when startup-config is missing")
	}
	if !strings.Contains(err.Error(), "requires startup-config") {
		t.Fatalf("Init error = %v, want requires startup-config", err)
	}
}
