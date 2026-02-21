package commands

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/trey.copeland/bootdev_blogag/internal/database"
)

type mockConfig struct {
	currentUser string
	setUserErr  error
}

func (m *mockConfig) SetUser(currentUserName string) error {
	if m.setUserErr != nil {
		return m.setUserErr
	}
	m.currentUser = currentUserName
	return nil
}

func (m *mockConfig) CurrentUserName() string {
	return m.currentUser
}

type mockQueries struct {
	getUserFn    func(ctx context.Context, name string) (database.User, error)
	createUserFn func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	clearUsersFn func(ctx context.Context) error
	getUsersFn   func(ctx context.Context) ([]string, error)
}

func (m *mockQueries) GetUser(ctx context.Context, name string) (database.User, error) {
	if m.getUserFn == nil {
		return database.User{}, errors.New("get user not mocked")
	}
	return m.getUserFn(ctx, name)
}

func (m *mockQueries) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	if m.createUserFn == nil {
		return database.User{}, errors.New("create user not mocked")
	}
	return m.createUserFn(ctx, arg)
}

func (m *mockQueries) ClearUsers(ctx context.Context) error {
	if m.clearUsersFn == nil {
		return nil
	}
	return m.clearUsersFn(ctx)
}

func (m *mockQueries) GetUsers(ctx context.Context) ([]string, error) {
	if m.getUsersFn == nil {
		return nil, nil
	}
	return m.getUsersFn(ctx)
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()
	originalStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}
	os.Stdout = writePipe

	run()

	_ = writePipe.Close()
	os.Stdout = originalStdout

	outputBytes, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	_ = readPipe.Close()
	return string(outputBytes)
}

func TestRunUnknownCommand(t *testing.T) {
	appCmds := New()
	RegisterDefault(appCmds)

	state := State{Config: &mockConfig{}, Queries: &mockQueries{}}
	err := appCmds.Run(&state, Command{Name: "poop"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "Command not registered") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHelpCommandPrintsCommands(t *testing.T) {
	appCmds := New()
	RegisterDefault(appCmds)
	state := State{Config: &mockConfig{}, Queries: &mockQueries{}}

	output := captureStdout(t, func() {
		err := appCmds.Run(&state, Command{Name: "help"})
		if err != nil {
			t.Fatalf("run help: %v", err)
		}
	})

	for _, commandName := range []string{"help", "login", "register", "reset", "users"} {
		if !strings.Contains(output, commandName) {
			t.Fatalf("help output missing command %q: %s", commandName, output)
		}
	}
}

func TestLoginReturnsUserNotExistsForNoRows(t *testing.T) {
	appCmds := New()
	RegisterDefault(appCmds)

	state := State{
		Config: &mockConfig{},
		Queries: &mockQueries{
			getUserFn: func(ctx context.Context, name string) (database.User, error) {
				return database.User{}, sql.ErrNoRows
			},
		},
	}

	err := appCmds.Run(&state, Command{Name: "login", Args: []string{"trey"}})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "User doesn't exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterCreatesUserWhenNoRows(t *testing.T) {
	appCmds := New()
	RegisterDefault(appCmds)

	config := &mockConfig{}
	queries := &mockQueries{
		getUserFn: func(ctx context.Context, name string) (database.User, error) {
			return database.User{}, sql.ErrNoRows
		},
		createUserFn: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
			return database.User{
				ID:        arg.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      arg.Name,
			}, nil
		},
	}

	state := State{Config: config, Queries: queries}
	err := appCmds.Run(&state, Command{Name: "register", Args: []string{"trey"}})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if config.CurrentUserName() != "trey" {
		t.Fatalf("expected current user to be set, got %q", config.CurrentUserName())
	}
}

func TestRegisterReturnsErrorOnUnexpectedLookupFailure(t *testing.T) {
	appCmds := New()
	RegisterDefault(appCmds)

	state := State{
		Config: &mockConfig{},
		Queries: &mockQueries{
			getUserFn: func(ctx context.Context, name string) (database.User, error) {
				return database.User{}, errors.New("database unavailable")
			},
			createUserFn: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
				return database.User{ID: uuid.New()}, nil
			},
		},
	}

	err := appCmds.Run(&state, Command{Name: "register", Args: []string{"trey"}})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Error checking existing user") {
		t.Fatalf("unexpected error: %v", err)
	}
}
