package app

import (
	"net/url"
	"path"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

const diagnosticPolicy = "scrubbed-v1"
const diagnosticInputLimit = 64 * 1024
const diagnosticMessageLimit = 2048
const diagnosticStackLimit = 8192
const diagnosticStackLines = 40
const diagnosticPolicyDescription = "scrubbed-v1: diagnostic records, not incidents or people. Messages are limited to 2048 UTF-8 bytes; stacks to 8192 bytes and 40 lines; each field processes at most 64 KiB. Known keyed credentials, tokens, cookies, emails, URL queries/paths and home paths are redacted. Scrubbing is imperfect; review synthetic diagnostic text before any real collection. Unknown attributes and ordinary informational log bodies remain omitted. Detail expires after 30 days; summaries retain no diagnostic free text."

var diagnosticANSI = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
var diagnosticURL = regexp.MustCompile(`(?i)\b(?:https?|file)://[^\s<>"')]+`)
var diagnosticEmail = regexp.MustCompile(`(?i)[a-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-z0-9](?:[a-z0-9.-]*[a-z0-9])?\.[a-z]{2,}`)
var diagnosticHeader = regexp.MustCompile(`(?im)(\b(?:authorization|proxy-authorization|cookie|set-cookie)\s*[:=]\s*)[^\n]*`)
var diagnosticBearer = regexp.MustCompile(`(?i)\b(Bearer|Basic)\s+[a-z0-9._~+/=-]+`)
var diagnosticKey = regexp.MustCompile(`(?i)((?:["']?)(?:password|passwd|pwd|passphrase|authorization|proxy-authorization|authentication|cookie|set-cookie|secret|client[_-]?secret|api[_-]?key|access[_-]?token|refresh[_-]?token|auth[_-]?token|token|session[_-]?(?:id|key|token)|connection[_-]?string)(?:["']?)\s*[:=]\s*)(?:"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|[^\s,;)}\]]+)`)
var diagnosticKnownToken = regexp.MustCompile(`\b(?:sk-[A-Za-z0-9_-]{12,}|gh[pousr]_[A-Za-z0-9]{12,}|github_pat_[A-Za-z0-9_]{12,}|AKIA[A-Z0-9]{16}|eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+)\b`)
var diagnosticHome = regexp.MustCompile(`(?i)(?:/Users/|/home/|[a-z]:\\Users\\)[^\s<>"')]+`)
var diagnosticFile = regexp.MustCompile(`(?:[a-zA-Z]:[\\/]|/)[^\s<>"')]+`)
var diagnosticSource = regexp.MustCompile(`(?i)^[a-z0-9_.@+~-]+\.(?:go|cs|fs|vb|js|jsx|mjs|cjs|ts|tsx|py|rb|java|kt|scala|rs|c|cc|cpp|h|php|swift|dart)(?::(?:line|[0-9]+)(?::[0-9]+)?)?:?$`)

func diagnosticUTF8Limit(value string, limit int) (string, bool) {
	if len(value) <= limit {
		return value, false
	}
	end := limit
	for end > 0 && !utf8.RuneStart(value[end]) {
		end--
	}
	return value[:end], true
}
func diagnosticSourceName(value string) string {
	value = strings.ReplaceAll(value, "\\", "/")
	name := path.Base(value)
	if diagnosticSource.MatchString(name) {
		return name
	}
	return "[redacted-path]"
}
func scrubDiagnostic(value string, stack bool) (string, bool) {
	value, clipped := diagnosticUTF8Limit(value, diagnosticInputLimit)
	value = strings.ToValidUTF8(value, "�")
	value = strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\r", "\n")
	value = diagnosticANSI.ReplaceAllString(value, "")
	value = strings.Map(func(r rune) rune {
		if r == '\n' {
			return r
		}
		if r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) || unicode.In(r, unicode.Cf) {
			return -1
		}
		return r
	}, value)
	value = diagnosticHeader.ReplaceAllString(value, "${1}[redacted]")
	value = diagnosticBearer.ReplaceAllString(value, "${1} [redacted]")
	value = diagnosticKey.ReplaceAllString(value, "${1}[redacted]")
	value = diagnosticKnownToken.ReplaceAllString(value, "[redacted-token]")
	value = diagnosticURL.ReplaceAllStringFunc(value, func(raw string) string {
		u, err := url.Parse(raw)
		if err != nil {
			return "[redacted-url]"
		}
		if u.Scheme == "file" {
			if stack {
				return diagnosticSourceName(u.Path)
			}
			return "[redacted-path]"
		}
		suffix := ""
		if u.Path != "" && u.Path != "/" {
			suffix = "/[redacted-path]"
			if stack {
				suffix = "/" + diagnosticSourceName(u.Path)
			}
		}
		return u.Scheme + "://" + u.Host + suffix
	})
	value = diagnosticEmail.ReplaceAllString(value, "[redacted-email]")
	value = diagnosticHome.ReplaceAllStringFunc(value, func(raw string) string {
		if stack {
			return diagnosticSourceName(raw)
		}
		return "[redacted-path]"
	})
	if stack {
		// Absolute filesystem paths are reduced to source basenames. URL paths have
		// already been projected above; do not consume their scheme/host here.
		value = diagnosticFile.ReplaceAllStringFunc(value, func(raw string) string {
			if strings.HasPrefix(raw, "//") {
				return raw
			}
			if strings.HasPrefix(raw, "/[redacted-") {
				return raw
			}
			return diagnosticSourceName(raw)
		})
		lines := strings.Split(value, "\n")
		if len(lines) > diagnosticStackLines {
			value = strings.Join(lines[:diagnosticStackLines], "\n")
			clipped = true
		}
	}
	limit := diagnosticMessageLimit
	if stack {
		limit = diagnosticStackLimit
	}
	value = strings.TrimSpace(value)
	value, cut := diagnosticUTF8Limit(value, limit)
	return value, clipped || cut
}
func retainDiagnosticText(x *Evidence, key, value string, stack bool) {
	if value == "" {
		return
	}
	text, truncated := scrubDiagnostic(value, stack)
	x.Attrs["diagnostic.policy"] = diagnosticPolicy
	if text != "" {
		x.Attrs[key] = text
	}
	if truncated {
		x.Attrs["diagnostic.truncated"] = "true"
	}
}
func projectDiagnosticAttributes(x *Evidence, a pcommon.Map) {
	if x.Attrs["exception.type"] != "" {
		x.Attrs["diagnostic.policy"] = diagnosticPolicy
	}
	for _, key := range []string{"exception.message", "exception.stacktrace"} {
		value, ok := a.Get(key)
		if ok && value.Type() == pcommon.ValueTypeStr {
			retainDiagnosticText(x, key, value.Str(), key == "exception.stacktrace")
		}
	}
}

// comparisonAtReceiptPolicy prevents a projection upgrade from turning an exact
// replay into a conflict. It never backfills old rows or weakens comparisons for
// fields already retained by that row's policy.
func comparisonAtReceiptPolicy(incoming, previous Evidence) Evidence {
	if previous.Attrs["diagnostic.policy"] != "" || incoming.Attrs["diagnostic.policy"] != diagnosticPolicy {
		return incoming
	}
	attrs := make(map[string]string, len(incoming.Attrs))
	for key, value := range incoming.Attrs {
		attrs[key] = value
	}
	for _, key := range []string{"exception.message", "exception.stacktrace", "log.message", "diagnostic.policy", "diagnostic.truncated"} {
		if _, existed := previous.Attrs[key]; !existed {
			delete(attrs, key)
		}
	}
	incoming.Attrs = attrs
	if len(attrs) == 0 && previous.Attrs == nil {
		incoming.Attrs = nil
	}
	// The former span projection did not retain Error status. Log severity was
	// retained, so it must still participate in legacy-log conflict detection.
	if incoming.Signal == "traces" {
		incoming.Severity = previous.Severity
	}
	return incoming
}

// isErrorDiagnostic accepts standard exception text even when the producer has
// no exception.type. Metric attributes remain receipt metadata, not errors.
func isErrorDiagnostic(e Evidence) bool {
	if e.Signal == "metrics" {
		return false
	}
	if e.Severity >= 17 || e.Attrs["exception.type"] != "" {
		return true
	}
	if e.Signal == "logs" || e.Signal == "traces" || e.Signal == "span-events" {
		return e.Attrs["exception.message"] != "" || e.Attrs["exception.stacktrace"] != ""
	}
	return false
}
