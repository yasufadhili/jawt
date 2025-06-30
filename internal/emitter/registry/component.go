package registry

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/ast"
	"strconv"
	"strings"
)

// ComponentType represents the type of component for categorisation
type ComponentType int

const (
	ComponentTypeLayout ComponentType = iota
	ComponentTypeText
	ComponentTypeInput
	ComponentTypeNavigation
	ComponentTypeMedia
	ComponentTypeList
	ComponentTypeContainer
)

// AttributeType defines the type of attribute for validation and processing
type AttributeType int

const (
	AttributeTypeString AttributeType = iota
	AttributeTypeBool
	AttributeTypeInt
	AttributeTypeFloat
	AttributeTypeURL
	AttributeTypeEvent
)

// AttributeDefinition defines metadata for component attributes
type AttributeDefinition struct {
	Name         string
	Type         AttributeType
	HTMLAttr     string // Maps to HTML attribute name
	Required     bool
	DefaultValue interface{}
	Description  string
}

// ComponentDefinition holds all information about a built-in component
type ComponentDefinition struct {
	Name           string
	HTMLTag        string
	Type           ComponentType
	DefaultClasses []string                        // Default style classes
	Attributes     map[string]*AttributeDefinition // Supported attributes
	SelfClosing    bool                            // Whether the HTML tag is self-closing
	Description    string
}

// ComponentRegistry manages all built-in components
type ComponentRegistry struct {
	components map[string]*ComponentDefinition
}

// NewComponentRegistry creates and initialises the component registry
func NewComponentRegistry() *ComponentRegistry {
	registry := &ComponentRegistry{
		components: make(map[string]*ComponentDefinition),
	}
	registry.initialiseBuiltInComponents()
	return registry
}

// GetComponent retrieves a component definition by name
func (cr *ComponentRegistry) GetComponent(name string) (*ComponentDefinition, bool) {
	comp, exists := cr.components[name]
	return comp, exists
}

// RegisterComponent adds a new component to the registry
func (cr *ComponentRegistry) RegisterComponent(comp *ComponentDefinition) {
	cr.components[comp.Name] = comp
}

// GetAllComponents returns all registered components
func (cr *ComponentRegistry) GetAllComponents() map[string]*ComponentDefinition {
	return cr.components
}

// GetComponentsByType returns components of a specific type
func (cr *ComponentRegistry) GetComponentsByType(compType ComponentType) []*ComponentDefinition {
	var result []*ComponentDefinition
	for _, comp := range cr.components {
		if comp.Type == compType {
			result = append(result, comp)
		}
	}
	return result
}

// initialiseBuiltInComponents populates the registry with built-in components
func (cr *ComponentRegistry) initialiseBuiltInComponents() {
	cr.initialiseLayoutComponents()
	cr.initialiseTextComponents()
	cr.initialiseInputComponents()
	cr.initialiseNavigationComponents()
	cr.initialiseListComponents()
}

func (cr *ComponentRegistry) initialiseLayoutComponents() {
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Container",
		HTMLTag:        "div",
		Type:           ComponentTypeContainer,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Just an HTML div",
	})
}

func (cr *ComponentRegistry) initialiseSemanticLayoutComponents() {
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Main",
		HTMLTag:        "main",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Main content landmark element",
	})
}

func (cr *ComponentRegistry) initialiseTextComponents() {
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Text",
		HTMLTag:        "p",
		Type:           ComponentTypeText,
		DefaultClasses: []string{},
		Attributes: map[string]*AttributeDefinition{
			"variant": {
				Name:        "variant",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Text variant (h1, h2, h3, h4, h5, h6, body, caption, small)",
			},
		},
		Description: "Just basic text",
	})
}

func (cr *ComponentRegistry) initialiseInputComponents() {}

func (cr *ComponentRegistry) initialiseNavigationComponents() {}

func (cr *ComponentRegistry) initialiseListComponents() {}

func (cr *ComponentRegistry) initialiseMediaComponents() {}

// ComponentProcessor handles component processing with the registry
type ComponentProcessor struct {
	registry *ComponentRegistry
}

// NewComponentProcessor creates a new component processor
func NewComponentProcessor() *ComponentProcessor {
	return &ComponentProcessor{
		registry: NewComponentRegistry(),
	}
}

// GetHTMLTag returns the HTML tag for a component, considering special attributes
func (cp *ComponentProcessor) GetHTMLTag(componentName string, properties []*ast.PropertyNode) string {
	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		return "div" // Default fallback
	}

	// Handle special cases where attributes affect the tag
	switch componentName {
	case "List":
		if cp.getBoolProperty(properties, "ordered") {
			return "ol"
		}
		return "ul"
	case "Heading":
		if level := cp.getIntProperty(properties, "level"); level >= 1 && level <= 6 {
			return fmt.Sprintf("h%d", level)
		}
		return "h1" // Default to h1
	default:
		return comp.HTMLTag
	}
}

// BuildClasses combines default classes with style property and variant-specific classes
func (cp *ComponentProcessor) BuildClasses(componentName string, properties []*ast.PropertyNode) []string {
	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		return []string{}
	}

	var classes []string

	// Add default classes
	classes = append(classes, comp.DefaultClasses...)

	// Add variant-specific classes
	classes = append(classes, cp.getVariantClasses(componentName, properties)...)

	// Add style property classes (custom Tailwind classes)
	if styleClasses := cp.getStyleClasses(properties); len(styleClasses) > 0 {
		classes = append(classes, styleClasses...)
	}

	return cp.deduplicateClasses(classes)
}

// BuildAttributes builds HTML attributes from component properties
func (cp *ComponentProcessor) BuildAttributes(componentName string, properties []*ast.PropertyNode) map[string]string {
	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		return map[string]string{}
	}

	attributes := make(map[string]string)

	// Build classes
	if classes := cp.BuildClasses(componentName, properties); len(classes) > 0 {
		attributes["class"] = strings.Join(classes, " ")
	}

	// Process-component-specific attributes
	for _, prop := range properties {
		if prop.Name == "text" {
			continue
		}
		if attrDef, exists := comp.Attributes[prop.Name]; exists {
			switch attrDef.HTMLAttr {
			case "class":
				// Already handled in BuildClasses
				continue
			case "tag":
				// Special case - affects the tag itself, not an attribute
				continue
			default:
				if value := cp.formatPropertyValue(prop, attrDef.Type); value != "" {
					attributes[attrDef.HTMLAttr] = value
				}
			}
		}
	}

	return attributes
}

// IsSelfClosing returns whether a component should be self-closing
func (cp *ComponentProcessor) IsSelfClosing(componentName string) bool {
	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		return false
	}
	return comp.SelfClosing
}

// ValidateComponent validates a component and its properties
func (cp *ComponentProcessor) ValidateComponent(componentName string, properties []*ast.PropertyNode) []string {
	var errors []string

	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		errors = append(errors, fmt.Sprintf("Unknown component: %s", componentName))
		return errors
	}

	// Check required attributes
	propertyMap := make(map[string]*ast.PropertyNode)
	for _, prop := range properties {
		propertyMap[prop.Name] = prop
	}

	for attrName, attrDef := range comp.Attributes {
		if attrDef.Required {
			if _, exists := propertyMap[attrName]; !exists {
				errors = append(errors, fmt.Sprintf("Required attribute '%s' missing for component '%s'", attrName, componentName))
			}
		}
	}

	// Validate attribute types
	for _, prop := range properties {
		if attrDef, exists := comp.Attributes[prop.Name]; exists {
			if !cp.validatePropertyType(prop, attrDef.Type) {
				errors = append(errors, fmt.Sprintf("Invalid type for attribute '%s' in component '%s'", prop.Name, componentName))
			}
		}
	}

	return errors
}

// Helper methods

func (cp *ComponentProcessor) GetTextContent(properties []*ast.PropertyNode) string {
	for _, prop := range properties {
		if prop.Name == "text" {
			if literal, ok := prop.Value.(*ast.LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

func (cp *ComponentProcessor) getVariantClasses(componentName string, properties []*ast.PropertyNode) []string {
	variant := cp.getStringProperty(properties, "variant")
	if variant == "" {
		return []string{}
	}

	switch componentName {
	case "Text":
		return cp.getTextVariantClasses(variant)
	default:
		return []string{}
	}
}

func (cp *ComponentProcessor) getTextVariantClasses(variant string) []string {
	switch variant {
	case "h1":
		return []string{"text-4xl", "font-bold"}
	case "h2":
		return []string{"text-3xl", "font-bold"}
	case "h3":
		return []string{"text-2xl", "font-semibold"}
	case "h4":
		return []string{"text-xl", "font-semibold"}
	case "h5":
		return []string{"text-lg", "font-medium"}
	case "h6":
		return []string{"text-base", "font-medium"}
	case "body":
		return []string{"text-base"}
	case "caption":
		return []string{"text-sm", "text-gray-600"}
	case "small":
		return []string{"text-xs"}
	default:
		return []string{}
	}
}

func (cp *ComponentProcessor) getStyleClasses(properties []*ast.PropertyNode) []string {
	for _, prop := range properties {
		if prop.Name == "style" {
			if literal, ok := prop.Value.(*ast.LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return strings.Fields(str)
				}
			}
		}
	}
	return []string{}
}

func (cp *ComponentProcessor) getStringProperty(properties []*ast.PropertyNode, name string) string {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*ast.LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

func (cp *ComponentProcessor) getBoolProperty(properties []*ast.PropertyNode, name string) bool {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*ast.LiteralNode); ok {
				if b, ok := literal.Value.(bool); ok {
					return b
				}
			}
		}
	}
	return false
}

func (cp *ComponentProcessor) getIntProperty(properties []*ast.PropertyNode, name string) int {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*ast.LiteralNode); ok {
				if i, ok := literal.Value.(int); ok {
					return i
				}
			}
		}
	}
	return 0
}

func (cp *ComponentProcessor) formatPropertyValue(prop *ast.PropertyNode, attrType AttributeType) string {
	if literal, ok := prop.Value.(*ast.LiteralNode); ok {
		switch attrType {
		case AttributeTypeString, AttributeTypeURL:
			if str, ok := literal.Value.(string); ok {
				return str
			}
		case AttributeTypeBool:
			if b, ok := literal.Value.(bool); ok {
				if b {
					return prop.Name // For boolean attributes like "required"
				}
				return ""
			}
		case AttributeTypeInt:
			if i, ok := literal.Value.(int); ok {
				return strconv.Itoa(i)
			}
		case AttributeTypeFloat:
			if f, ok := literal.Value.(float64); ok {
				return fmt.Sprintf("%.2f", f)
			}
		case AttributeTypeEvent:
			if str, ok := literal.Value.(string); ok {
				return str
			}
		}
	}
	return ""
}

func (cp *ComponentProcessor) validatePropertyType(prop *ast.PropertyNode, expectedType AttributeType) bool {
	if literal, ok := prop.Value.(*ast.LiteralNode); ok {
		switch expectedType {
		case AttributeTypeString, AttributeTypeURL, AttributeTypeEvent:
			_, ok := literal.Value.(string)
			return ok
		case AttributeTypeBool:
			_, ok := literal.Value.(bool)
			return ok
		case AttributeTypeInt:
			_, ok := literal.Value.(int)
			return ok
		case AttributeTypeFloat:
			_, ok := literal.Value.(float64)
			return ok
		}
	}
	return false
}

func (cp *ComponentProcessor) deduplicateClasses(classes []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, class := range classes {
		if class != "" && !seen[class] {
			seen[class] = true
			result = append(result, class)
		}
	}

	return result
}
