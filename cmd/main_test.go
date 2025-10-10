package main

import (
	"errors"
	"testing"

	"github.com/device-management-toolkit/mps-router/internal/db"
	itest "github.com/device-management-toolkit/mps-router/internal/test"
)

// fakeManager implements db.Manager minimal surface via internal/test mocks
// to observe Health() and Query() behaviors.

type fakeServerStart struct {
	called bool
	addr   string
	target string
	err    error
}

func (f *fakeServerStart) start(_ db.Manager, addr, target string) error {
	f.called = true
	f.addr = addr
	f.target = target
	return f.err
}

// Shorthand aliases for readability
type (
	mongoMgr = itest.MockNOSQLDBManager
	pgMgr    = itest.MockSQLDBManager
)

func TestIsMongoConnectionString(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"mongodb://localhost:27017", true},
		{"mongodb+srv://atlas.test", true},
		{"postgres://localhost:5432", false},
	}

	for _, c := range cases {
		got := isMongoConnectionString(c.input)
		if got != c.want {
			t.Errorf("isMongoConnectionString(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestRun_RequiresConnectionString(t *testing.T) {
	code := run(
		nil,
		func(s string) string { return "" },
		func(_ db.Manager, _, _ string) error { return nil },
		func(s string) db.Manager { return &mongoMgr{} },
		func(s string) db.Manager { return &pgMgr{} },
	)
	if code == 0 {
		t.Fatalf("expected non-zero exit code when MPS_CONNECTION_STRING missing")
	}
}

func TestRun_HealthSuccessAndFailure(t *testing.T) {
	// success path
	getenv := func(key string) string {
		switch key {
		case "MPS_CONNECTION_STRING":
			return "postgres://test"
		default:
			return ""
		}
	}
	start := &fakeServerStart{}
	code := run(
		[]string{"-health"},
		getenv,
		func(m db.Manager, a, tg string) error { return start.start(m, a, tg) },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
	)
	if code != 0 {
		t.Fatalf("health success expected 0, got %d", code)
	}
	if start.called {
		t.Fatalf("server should not start during health check")
	}

	// failure path
	code = run(
		[]string{"-health"},
		getenv,
		func(m db.Manager, a, tg string) error { return start.start(m, a, tg) },
		func(s string) db.Manager { return &pgMgr{HealthResult: false} },
		func(s string) db.Manager { return &pgMgr{HealthResult: false} },
	)
	if code == 0 {
		t.Fatalf("health failure expected non-zero, got %d", code)
	}
}

func TestRun_EnvDefaultsAndOverrides(t *testing.T) {
	// Defaults when missing
	getenvDefaults := func(key string) string {
		if key == "MPS_CONNECTION_STRING" {
			return "postgres://test"
		}
		return ""
	}
	server := &fakeServerStart{}
	code := run(
		nil,
		getenvDefaults,
		func(m db.Manager, a, tg string) error { return server.start(m, a, tg) },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
	)
	if code != 0 {
		t.Fatalf("expected success, got %d", code)
	}
	if !server.called || server.addr != ":8003" || server.target != "mps:3000" {
		t.Fatalf("defaults not applied: got addr=%q target=%q", server.addr, server.target)
	}

	// Overrides
	getenvOverrides := func(key string) string {
		switch key {
		case "MPS_CONNECTION_STRING":
			return "postgres://test"
		case "PORT":
			return "9000"
		case "MPS_HOST":
			return "example.local"
		case "MPS_PORT":
			return "1234"
		default:
			return ""
		}
	}
	server2 := &fakeServerStart{}
	code = run(
		nil,
		getenvOverrides,
		func(m db.Manager, a, tg string) error { return server2.start(m, a, tg) },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
	)
	if code != 0 {
		t.Fatalf("expected success with overrides, got %d", code)
	}
	if server2.addr != ":9000" || server2.target != "example.local:1234" {
		t.Fatalf("overrides not applied: addr=%q target=%q", server2.addr, server2.target)
	}
}

func TestRun_DBSelection(t *testing.T) {
	// capture which constructor was called
	var used []string
	getenv := func(k string) string {
		if k == "MPS_CONNECTION_STRING" {
			return "mongodb://atlas"
		}
		return ""
	}
	// track calls
	newMongo := func(s string) db.Manager {
		used = append(used, "mongo")
		return &mongoMgr{HealthResult: true}
	}
	newPg := func(s string) db.Manager {
		used = append(used, "pg")
		return &pgMgr{HealthResult: true}
	}
	server := &fakeServerStart{}
	code := run(nil, getenv, func(m db.Manager, a, tg string) error { return server.start(m, a, tg) }, newMongo, newPg)
	if code != 0 {
		t.Fatalf("expected success, got %d", code)
	}
	if len(used) != 1 || used[0] != "mongo" {
		t.Fatalf("expected mongo constructor used, got %v", used)
	}
}

func TestRun_ServerErrorPropagates(t *testing.T) {
	expected := errors.New("boom")
	getenv := func(k string) string {
		if k == "MPS_CONNECTION_STRING" {
			return "postgres://test"
		}
		return ""
	}
	server := &fakeServerStart{err: expected}
	code := run(
		nil,
		getenv,
		func(m db.Manager, a, tg string) error { return server.start(m, a, tg) },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
	)
	if code == 0 {
		t.Fatalf("expected non-zero when server returns error")
	}
}

func TestRun_FlagParseError(t *testing.T) {
	getenv := func(k string) string { return "postgres://test" }
	code := run(
		[]string{"-health=maybe"}, // invalid bool value triggers parse error
		getenv,
		func(m db.Manager, a, tg string) error { return nil },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
		func(s string) db.Manager { return &pgMgr{HealthResult: true} },
	)
	if code == 0 {
		t.Fatalf("expected non-zero exit on flag parse error")
	}
}

func TestStartServerReal_InvalidAddr(t *testing.T) {
	// Provide an invalid TCP address to force net.Listen to fail immediately
	err := startServerReal(&pgMgr{HealthResult: true}, "badaddr", "mps:3000")
	if err == nil {
		t.Fatalf("expected error from startServerReal with invalid addr")
	}
}
