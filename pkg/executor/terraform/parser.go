package terraform

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// TerraformMessage represents a single JSON message from terraform output.
type TerraformMessage struct {
	Level     string `json:"@level"`
	Message   string `json:"@message"`
	Module    string `json:"@module"`
	Timestamp string `json:"@timestamp"`
	Type      string `json:"type"`

	// For version message
	Terraform string `json:"terraform,omitempty"`
	UI        string `json:"ui,omitempty"`

	// For change_summary
	Changes *ChangeSummary `json:"changes,omitempty"`

	// For diagnostic (errors)
	Diagnostic *Diagnostic `json:"diagnostic,omitempty"`

	// For hooks (refresh_start, apply_start, etc)
	Hook *HookInfo `json:"hook,omitempty"`
}

// ChangeSummary contains plan change statistics.
type ChangeSummary struct {
	Add       int    `json:"add"`
	Change    int    `json:"change"`
	Import    int    `json:"import"`
	Remove    int    `json:"remove"`
	Operation string `json:"operation"`
}

// Diagnostic contains error/warning information.
type Diagnostic struct {
	Severity string       `json:"severity"`
	Summary  string       `json:"summary"`
	Detail   string       `json:"detail"`
	Address  string       `json:"address"`
	Range    *SourceRange `json:"range,omitempty"`
}

// SourceRange points to source code location.
type SourceRange struct {
	Filename string `json:"filename"`
	Start    struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"start"`
}

// HookInfo contains resource operation information.
type HookInfo struct {
	Resource *ResourceInfo `json:"resource,omitempty"`
	Action   string        `json:"action,omitempty"`
	IDKey    string        `json:"id_key,omitempty"`
	IDValue  string        `json:"id_value,omitempty"`
}

// ResourceInfo contains resource details.
type ResourceInfo struct {
	Addr         string `json:"addr"`
	Module       string `json:"module"`
	Resource     string `json:"resource"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
}

// PlanInfo plan result summary.
type PlanInfo struct {
	ToAdd     int
	ToChange  int
	ToDestroy int
}

// ParseResult contains parsed terraform output.
type ParseResult struct {
	Messages []TerraformMessage
	Changes  *ChangeSummary
	Errors   []Diagnostic
	Version  string
	Success  bool
}

// Parser parses terraform output.
type Parser struct {
	planRegex *regexp.Regexp
}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{
		planRegex: regexp.MustCompile(`Plan:\s*(\d+)\s*to add,\s*(\d+)\s*to change,\s*(\d+)\s*to destroy`),
	}
}

// ParseJSONOutput parses terraform JSON output (one JSON per line).
func (p *Parser) ParseJSONOutput(output string) *ParseResult {
	result := &ParseResult{
		Messages: []TerraformMessage{},
		Errors:   []Diagnostic{},
		Success:  true,
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var msg TerraformMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		result.Messages = append(result.Messages, msg)

		switch msg.Type {
		case "version":
			result.Version = msg.Terraform
		case "change_summary":
			if msg.Changes != nil {
				result.Changes = msg.Changes
			}
		case "diagnostic":
			if msg.Diagnostic != nil {
				result.Errors = append(result.Errors, *msg.Diagnostic)
				if msg.Diagnostic.Severity == "error" {
					result.Success = false
				}
			}
		}
	}

	return result
}

// ParsePlan parses plan output (legacy text format).
func (p *Parser) ParsePlan(output string) *PlanInfo {
	info := &PlanInfo{}

	// First try JSON format
	result := p.ParseJSONOutput(output)
	if result.Changes != nil {
		info.ToAdd = result.Changes.Add
		info.ToChange = result.Changes.Change
		info.ToDestroy = result.Changes.Remove
		return info
	}

	// Fallback to regex for text format
	matches := p.planRegex.FindStringSubmatch(output)
	if len(matches) == 4 {
		info.ToAdd, _ = strconv.Atoi(matches[1])
		info.ToChange, _ = strconv.Atoi(matches[2])
		info.ToDestroy, _ = strconv.Atoi(matches[3])
	}

	return info
}

// ParseTfstate extracts attributes from tfstate.
func (p *Parser) ParseTfstate(data []byte) map[string]string {
	attrs := make(map[string]string)

	resources := gjson.GetBytes(data, "resources")
	if !resources.Exists() {
		return attrs
	}

	resources.ForEach(func(_, resource gjson.Result) bool {
		resourceType := resource.Get("type").String()
		name := resource.Get("name").String()

		instances := resource.Get("instances")
		instances.ForEach(func(_, instance gjson.Result) bool {
			attributes := instance.Get("attributes")
			if !attributes.Exists() {
				return true
			}

			if id := attributes.Get("id"); id.Exists() {
				key := resourceType + "." + name + ".id"
				attrs[key] = id.String()
			}
			if arn := attributes.Get("arn"); arn.Exists() {
				key := resourceType + "." + name + ".arn"
				attrs[key] = arn.String()
			}

			return true
		})
		return true
	})

	return attrs
}

// ParseTfstateJSON parses tfstate to structured data.
func (p *Parser) ParseTfstateJSON(data []byte) (*TfState, error) {
	var state TfState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// TfState represents terraform.tfstate structure.
type TfState struct {
	Version   int               `json:"version"`
	Serial    int               `json:"serial"`
	Lineage   string            `json:"lineage"`
	Resources []TfStateResource `json:"resources"`
}

// TfStateResource represents a resource in tfstate.
type TfStateResource struct {
	Mode      string            `json:"mode"`
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Provider  string            `json:"provider"`
	Instances []TfStateInstance `json:"instances"`
}

// TfStateInstance represents a resource instance.
type TfStateInstance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// ExtractAttributes extracts attributes by path config.
func (p *Parser) ExtractAttributes(data []byte, paths map[string]string) map[string]string {
	result := make(map[string]string)

	for name, path := range paths {
		value := gjson.GetBytes(data, path)
		if value.Exists() {
			result[name] = value.String()
		}
	}

	return result
}
