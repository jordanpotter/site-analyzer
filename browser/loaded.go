package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const isLoadedTemplate = `
var cb = arguments[arguments.length - 1];

{{if .isEmpty -}}

if (document.readyState === 'complete') {
	cb();
} else {
	window.addEventListener('load', function() { cb(); }, {once: true});
}

{{- else -}}

if ({{.isLoadedConditional}}) {
	cb();
}

var observer = new MutationObserver(function(mutations) {
	mutations.forEach(function(mutation) {
		if ({{.isLoadedConditional}}) {
			observer.disconnect();
			cb();
		}
	});
});

var target = document;
var config = {childList: true, subtree: true};
observer.observe(target, config);

{{- end }}
`

type LoadedSpec struct {
	Operand  string       `json:"operand"`
	Elements []string     `json:"elements"`
	Children []LoadedSpec `json:"children"`
}

func ParseLoadedSpec(data []byte) (*LoadedSpec, error) {
	var spec LoadedSpec
	err := json.Unmarshal(data, &spec)
	return &spec, errors.Wrap(err, "failed to unmarshal json")
}

func (spec *LoadedSpec) IsLoadedScript() (string, error) {
	t, err := template.New("isLoadedTemplate").Parse(isLoadedTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	isLoadedConditional, err := spec.isLoadedConditional()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine loaded conditional")
	}

	data := map[string]interface{}{
		"isEmpty":             spec.isEmpty(),
		"isLoadedConditional": isLoadedConditional,
	}

	var str bytes.Buffer
	if err = t.Execute(&str, data); err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return str.String(), nil
}

func (spec *LoadedSpec) isEmpty() bool {
	return spec.Operand == "" && len(spec.Elements) == 0 && len(spec.Children) == 0
}

func (spec *LoadedSpec) isLoadedConditional() (string, error) {
	var conditionals []string

	for _, elem := range spec.Elements {
		elemConditional := fmt.Sprintf("!!document.querySelector('%s')", elem)
		conditionals = append(conditionals, elemConditional)
	}

	for _, child := range spec.Children {
		childConditional, err := child.isLoadedConditional()
		if err != nil {
			return "", errors.Wrap(err, "failed to process child loaded conditional")
		}
		conditionals = append(conditionals, childConditional)
	}

	separator, err := spec.operandSeparator()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine separator")
	}

	masterConditional := fmt.Sprintf("(%s)", strings.Join(conditionals, separator))
	return masterConditional, nil
}

func (spec *LoadedSpec) operandSeparator() (string, error) {
	if len(spec.Elements) < 2 {
		return "", nil
	}

	switch strings.ToLower(spec.Operand) {
	case "and":
		return " && ", nil
	case "or":
		return " || ", nil
	default:
		return "", errors.Errorf("unexpected operand %q", spec.Operand)
	}
}

func (b *Browser) load(url string, spec *LoadedSpec) (time.Duration, error) {
	start := time.Now()

	if err := b.session.Url(url); err != nil {
		return 0, errors.Wrap(err, "failed to set url")
	}

	script, err := spec.IsLoadedScript()
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve script")
	}

	_, err = b.session.ExecuteScriptAsync(script, []interface{}{})
	return time.Since(start), errors.Wrap(err, "failed to execute async script")
}
