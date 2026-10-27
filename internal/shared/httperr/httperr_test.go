package httperr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/apperr"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

// Reads the codes out of the source rather than listing them here, so a code added later is
// checked without anyone updating this test. A code with no status reaches the client as a 500,
// which reads as "the server is broken" for what is usually the caller's own mistake.
func TestEveryDeclaredCodeMapsToAStatus(t *testing.T) {
	declared := declaredCodes(t)

	if len(declared) < 4 {
		t.Fatalf("found only %d codes; the discovery is probably broken", len(declared))
	}
	for code, file := range declared {
		if _, ok := codeToStatus[code]; !ok {
			t.Errorf("%s (%s) has no row in codeToStatus", code, file)
		}
	}
}

func declaredCodes(t *testing.T) map[apperr.Code]string {
	t.Helper()

	files, err := filepath.Glob("../../modules/*/*/codes.go")
	if err != nil {
		t.Fatalf("globbing codes.go: %v", err)
	}

	codes := map[apperr.Code]string{}
	for _, file := range files {
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", file, err)
		}

		ast.Inspect(parsed, func(n ast.Node) bool {
			spec, ok := n.(*ast.ValueSpec)
			if !ok || len(spec.Values) == 0 {
				return true
			}
			sel, ok := spec.Type.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Code" {
				return true
			}
			for _, value := range spec.Values {
				lit, ok := value.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				raw, err := strconv.Unquote(lit.Value)
				if err != nil {
					continue
				}
				codes[apperr.Code(raw)] = file
			}
			return true
		})
	}
	return codes
}

// Each domain error maps to the status its caller expects. A wrapped error maps the same way.
func TestDomainErrorsMapToTheirStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "invalid email", err: &user.InvalidEmailError{Email: "not-an-email"}, want: http.StatusUnprocessableEntity},
		{name: "weak password", err: &user.WeakPasswordError{MinLength: 3, MaxBytes: 72}, want: http.StatusUnprocessableEntity},
		{name: "email taken", err: &user.EmailTakenError{Email: "ana@example.com"}, want: http.StatusConflict},
		{name: "invalid credentials", err: &user.InvalidCredentialsError{}, want: http.StatusUnauthorized},
		{name: "wrapped email taken", err: fmt.Errorf("registering: %w", &user.EmailTakenError{Email: "ana@example.com"}), want: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToHumaError(context.Background(), discard, tt.err)

			status, ok := got.(huma.StatusError)
			if !ok {
				t.Fatalf("got %T, want a huma.StatusError", got)
			}
			if status.GetStatus() != tt.want {
				t.Errorf("status = %d, want %d", status.GetStatus(), tt.want)
			}
		})
	}
}

// An unrecognised failure stays opaque: a generic 500, with no message from inside, and it is logged
// once so someone can find it.
func TestUnmappedErrorsStayInternalAndAreLoggedOnce(t *testing.T) {
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))

	got := ToHumaError(context.Background(), log, errors.New("pq: connection refused at 10.0.0.5"))

	status, ok := got.(huma.StatusError)
	if !ok {
		t.Fatalf("got %T, want a huma.StatusError", got)
	}
	if status.GetStatus() != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", status.GetStatus())
	}
	if strings.Contains(got.Error(), "10.0.0.5") {
		t.Errorf("client message leaked the internal detail: %q", got.Error())
	}
	if n := strings.Count(logs.String(), "level=ERROR"); n != 1 {
		t.Errorf("logged %d error lines, want exactly 1:\n%s", n, logs.String())
	}
}

// A mapped domain error is an answer, not a failure: it must not be logged as an error.
func TestMappedErrorsAreNotLoggedAsErrors(t *testing.T) {
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))

	ToHumaError(context.Background(), log, &user.EmailTakenError{Email: "ana@example.com"})

	if logs.Len() != 0 {
		t.Errorf("a mapped domain error was logged:\n%s", logs.String())
	}
}
