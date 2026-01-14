package terraform

import (
	"testing"
)

func TestParser_ParsePlan(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		output   string
		expected *PlanInfo
	}{
		{
			name:   "normal plan output",
			output: "Plan: 3 to add, 1 to change, 2 to destroy.",
			expected: &PlanInfo{
				ToAdd:     3,
				ToChange:  1,
				ToDestroy: 2,
			},
		},
		{
			name:   "no changes",
			output: "Plan: 0 to add, 0 to change, 0 to destroy.",
			expected: &PlanInfo{
				ToAdd:     0,
				ToChange:  0,
				ToDestroy: 0,
			},
		},
		{
			name:   "no plan line",
			output: "Some other output",
			expected: &PlanInfo{
				ToAdd:     0,
				ToChange:  0,
				ToDestroy: 0,
			},
		},
		{
			name:   "plan in multiline output",
			output: "Refreshing state...\nPlan: 5 to add, 2 to change, 1 to destroy.\nApply complete!",
			expected: &PlanInfo{
				ToAdd:     5,
				ToChange:  2,
				ToDestroy: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ParsePlan(tt.output)
			if result.ToAdd != tt.expected.ToAdd {
				t.Errorf("ToAdd = %d, want %d", result.ToAdd, tt.expected.ToAdd)
			}
			if result.ToChange != tt.expected.ToChange {
				t.Errorf("ToChange = %d, want %d", result.ToChange, tt.expected.ToChange)
			}
			if result.ToDestroy != tt.expected.ToDestroy {
				t.Errorf("ToDestroy = %d, want %d", result.ToDestroy, tt.expected.ToDestroy)
			}
		})
	}
}

func TestParser_ParseTfstate(t *testing.T) {
	parser := NewParser()

	tfstate := []byte(`{
		"version": 4,
		"resources": [
			{
				"type": "aws_instance",
				"name": "web",
				"instances": [
					{
						"attributes": {
							"id": "i-1234567890abcdef0",
							"arn": "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"
						}
					}
				]
			}
		]
	}`)

	attrs := parser.ParseTfstate(tfstate)

	if len(attrs) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(attrs))
	}

	if attrs["aws_instance.web.id"] != "i-1234567890abcdef0" {
		t.Errorf("wrong id: %s", attrs["aws_instance.web.id"])
	}

	if attrs["aws_instance.web.arn"] != "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0" {
		t.Errorf("wrong arn: %s", attrs["aws_instance.web.arn"])
	}
}

func TestParser_ParseTfstate_Empty(t *testing.T) {
	parser := NewParser()

	// 无 resources
	tfstate := []byte(`{"version": 4}`)
	attrs := parser.ParseTfstate(tfstate)
	if len(attrs) != 0 {
		t.Errorf("expected 0 attributes, got %d", len(attrs))
	}

	// 空 instances
	tfstate2 := []byte(`{"resources": [{"type": "test", "name": "test", "instances": []}]}`)
	attrs2 := parser.ParseTfstate(tfstate2)
	if len(attrs2) != 0 {
		t.Errorf("expected 0 attributes, got %d", len(attrs2))
	}
}

func TestParser_ParseTfstate_NoIdOrArn(t *testing.T) {
	parser := NewParser()

	tfstate := []byte(`{
		"resources": [
			{
				"type": "null_resource",
				"name": "test",
				"instances": [{"attributes": {"triggers": {}}}]
			}
		]
	}`)

	attrs := parser.ParseTfstate(tfstate)
	if len(attrs) != 0 {
		t.Errorf("expected 0 attributes for resource without id/arn, got %d", len(attrs))
	}
}

func TestParser_ExtractAttributes(t *testing.T) {
	parser := NewParser()

	tfstate := []byte(`{
		"resources": [
			{
				"instances": [
					{
						"attributes": {
							"id": "test-id",
							"name": "test-name"
						}
					}
				]
			}
		]
	}`)

	paths := map[string]string{
		"id":   "resources.0.instances.0.attributes.id",
		"name": "resources.0.instances.0.attributes.name",
	}

	result := parser.ExtractAttributes(tfstate, paths)

	if result["id"] != "test-id" {
		t.Errorf("wrong id: %s", result["id"])
	}
	if result["name"] != "test-name" {
		t.Errorf("wrong name: %s", result["name"])
	}
}

func TestParser_ExtractAttributes_MissingPath(t *testing.T) {
	parser := NewParser()

	tfstate := []byte(`{"resources": []}`)
	paths := map[string]string{
		"missing": "nonexistent.path",
	}

	result := parser.ExtractAttributes(tfstate, paths)
	if _, exists := result["missing"]; exists {
		t.Error("missing path should not be in result")
	}
}

func TestParser_ParseTfstateJSON(t *testing.T) {
	parser := NewParser()

	tfstate := []byte(`{
		"version": 4,
		"serial": 10,
		"lineage": "test-lineage",
		"resources": [
			{
				"mode": "managed",
				"type": "aws_instance",
				"name": "web",
				"provider": "provider.aws",
				"instances": [
					{
						"schema_version": 1,
						"attributes": {
							"id": "i-12345"
						}
					}
				]
			}
		]
	}`)

	state, err := parser.ParseTfstateJSON(tfstate)
	if err != nil {
		t.Fatalf("ParseTfstateJSON should succeed: %v", err)
	}

	if state.Version != 4 {
		t.Errorf("version should be 4, got %d", state.Version)
	}
	if state.Serial != 10 {
		t.Errorf("serial should be 10, got %d", state.Serial)
	}
	if state.Lineage != "test-lineage" {
		t.Errorf("lineage wrong: %s", state.Lineage)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("should have 1 resource, got %d", len(state.Resources))
	}
	if state.Resources[0].Type != "aws_instance" {
		t.Errorf("resource type wrong: %s", state.Resources[0].Type)
	}
}

func TestParser_ParseTfstateJSON_Invalid(t *testing.T) {
	parser := NewParser()

	_, err := parser.ParseTfstateJSON([]byte("invalid json"))
	if err == nil {
		t.Error("invalid JSON should return error")
	}
}

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("NewParser should not return nil")
	}
	if parser.planRegex == nil {
		t.Fatal("planRegex should be initialized")
	}
}

func TestParser_ParseJSONOutput(t *testing.T) {
	parser := NewParser()

	// Sample JSON output from terraform -json
	output := `{"@level":"info","@message":"Terraform 1.14.2","@module":"terraform.ui","@timestamp":"2025-12-27T00:13:17.758668+08:00","terraform":"1.14.2","type":"version","ui":"1.2"}
{"@level":"info","@message":"aws_instance.test_vm: Refreshing state...","@module":"terraform.ui","@timestamp":"2025-12-27T00:13:21.955742+08:00","hook":{"resource":{"addr":"aws_instance.test_vm","resource_type":"aws_instance","resource_name":"test_vm"}},"type":"refresh_start"}
{"@level":"info","@message":"Plan: 2 to add, 1 to change, 3 to destroy.","@module":"terraform.ui","@timestamp":"2025-12-27T00:13:44.414609+08:00","changes":{"add":2,"change":1,"import":0,"remove":3,"operation":"plan"},"type":"change_summary"}`

	result := parser.ParseJSONOutput(output)

	if result.Version != "1.14.2" {
		t.Errorf("Version should be '1.14.2', got '%s'", result.Version)
	}
	if result.Changes == nil {
		t.Fatal("Changes should not be nil")
	}
	if result.Changes.Add != 2 {
		t.Errorf("Add should be 2, got %d", result.Changes.Add)
	}
	if result.Changes.Change != 1 {
		t.Errorf("Change should be 1, got %d", result.Changes.Change)
	}
	if result.Changes.Remove != 3 {
		t.Errorf("Remove should be 3, got %d", result.Changes.Remove)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
	if len(result.Messages) != 3 {
		t.Errorf("Should have 3 messages, got %d", len(result.Messages))
	}
}

func TestParser_ParseJSONOutput_WithErrors(t *testing.T) {
	parser := NewParser()

	output := `{"@level":"info","@message":"Terraform 1.14.2","type":"version","terraform":"1.14.2"}
{"@level":"error","@message":"Error: connection refused","type":"diagnostic","diagnostic":{"severity":"error","summary":"connection refused","address":"aws_instance.test"}}`

	result := parser.ParseJSONOutput(output)

	if result.Success {
		t.Error("Success should be false when there are errors")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("Should have 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Severity != "error" {
		t.Errorf("Severity should be 'error', got '%s'", result.Errors[0].Severity)
	}
	if result.Errors[0].Summary != "connection refused" {
		t.Errorf("Summary wrong: %s", result.Errors[0].Summary)
	}
}

func TestParser_ParsePlan_JSON(t *testing.T) {
	parser := NewParser()

	// JSON format
	output := `{"@level":"info","@message":"Plan: 5 to add, 2 to change, 1 to destroy.","type":"change_summary","changes":{"add":5,"change":2,"remove":1,"operation":"plan"}}`

	info := parser.ParsePlan(output)

	if info.ToAdd != 5 {
		t.Errorf("ToAdd should be 5, got %d", info.ToAdd)
	}
	if info.ToChange != 2 {
		t.Errorf("ToChange should be 2, got %d", info.ToChange)
	}
	if info.ToDestroy != 1 {
		t.Errorf("ToDestroy should be 1, got %d", info.ToDestroy)
	}
}
