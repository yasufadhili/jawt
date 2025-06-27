package compiler

import (
	"fmt"
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
	DefaultClasses []string                        // Default Tailwind classes
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
	registry.initializeBuiltInComponents()
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

// initializeBuiltInComponents populates the registry with built-in components
func (cr *ComponentRegistry) initializeBuiltInComponents() {
	// Layout Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Row",
		HTMLTag:        "div",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{"flex", "flex-row"},
		Attributes: map[string]*AttributeDefinition{
			"justify": {
				Name:        "justify",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Justify content alignment (start, center, end, between, around, evenly)",
			},
			"align": {
				Name:        "align",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Align items (start, center, end, stretch, baseline)",
			},
			"gap": {
				Name:        "gap",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Gap between items (0, 1, 2, 4, 8, etc.)",
			},
		},
		Description: "Horizontal flex container",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Column",
		HTMLTag:        "div",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{"flex", "flex-col"},
		Attributes: map[string]*AttributeDefinition{
			"justify": {
				Name:        "justify",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Justify content alignment",
			},
			"align": {
				Name:        "align",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Align items",
			},
			"gap": {
				Name:        "gap",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Gap between items",
			},
		},
		Description: "Vertical flex container",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Container",
		HTMLTag:        "div",
		Type:           ComponentTypeContainer,
		DefaultClasses: []string{"container", "mx-auto"},
		Attributes: map[string]*AttributeDefinition{
			"padding": {
				Name:        "padding",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Padding (p-0, p-4, px-4, py-2, etc.)",
			},
			"maxWidth": {
				Name:        "maxWidth",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Maximum width (sm, md, lg, xl, 2xl, etc.)",
			},
		},
		Description: "Responsive container with centered content",
	})

	// Text Components
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
			"color": {
				Name:        "color",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Text color",
			},
			"align": {
				Name:        "align",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Text alignment (left, center, right, justify)",
			},
			"weight": {
				Name:        "weight",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Font weight (thin, light, normal, medium, semibold, bold, extrabold, black)",
			},
		},
		Description: "Text element with typography variants",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Heading",
		HTMLTag:        "h1",
		Type:           ComponentTypeText,
		DefaultClasses: []string{"text-2xl", "font-bold"},
		Attributes: map[string]*AttributeDefinition{
			"level": {
				Name:        "level",
				Type:        AttributeTypeInt,
				HTMLAttr:    "tag",
				Description: "Heading level (1-6)",
			},
			"size": {
				Name:        "size",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Text size (xs, sm, base, lg, xl, 2xl, etc.)",
			},
		},
		Description: "Heading element with configurable level",
	})

	// Input Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Input",
		HTMLTag:        "input",
		Type:           ComponentTypeInput,
		DefaultClasses: []string{"border", "rounded", "px-3", "py-2"},
		SelfClosing:    true,
		Attributes: map[string]*AttributeDefinition{
			"type": {
				Name:         "type",
				Type:         AttributeTypeString,
				HTMLAttr:     "type",
				DefaultValue: "text",
				Description:  "Input type (text, email, password, number, etc.)",
			},
			"placeholder": {
				Name:        "placeholder",
				Type:        AttributeTypeString,
				HTMLAttr:    "placeholder",
				Description: "Placeholder text",
			},
			"required": {
				Name:        "required",
				Type:        AttributeTypeBool,
				HTMLAttr:    "required",
				Description: "Whether the input is required",
			},
			"disabled": {
				Name:        "disabled",
				Type:        AttributeTypeBool,
				HTMLAttr:    "disabled",
				Description: "Whether the input is disabled",
			},
			"onChange": {
				Name:        "onChange",
				Type:        AttributeTypeEvent,
				HTMLAttr:    "onchange",
				Description: "Change event handler",
			},
		},
		Description: "Input field with various types",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Button",
		HTMLTag:        "button",
		Type:           ComponentTypeInput,
		DefaultClasses: []string{"bg-blue-500", "text-white", "px-4", "py-2", "rounded", "hover:bg-blue-600"},
		Attributes: map[string]*AttributeDefinition{
			"variant": {
				Name:        "variant",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Button variant (primary, secondary, outline, ghost, danger)",
			},
			"size": {
				Name:        "size",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Button size (xs, sm, base, lg, xl)",
			},
			"disabled": {
				Name:        "disabled",
				Type:        AttributeTypeBool,
				HTMLAttr:    "disabled",
				Description: "Whether the button is disabled",
			},
			"onClick": {
				Name:        "onClick",
				Type:        AttributeTypeEvent,
				HTMLAttr:    "onclick",
				Description: "Click event handler",
			},
			"type": {
				Name:         "type",
				Type:         AttributeTypeString,
				HTMLAttr:     "type",
				DefaultValue: "button",
				Description:  "Button type (button, submit, reset)",
			},
		},
		Description: "Interactive button element",
	})

	// Navigation Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Link",
		HTMLTag:        "a",
		Type:           ComponentTypeNavigation,
		DefaultClasses: []string{"text-blue-600", "hover:text-blue-800", "underline"},
		Attributes: map[string]*AttributeDefinition{
			"href": {
				Name:        "href",
				Type:        AttributeTypeURL,
				HTMLAttr:    "href",
				Required:    true,
				Description: "Link destination URL",
			},
			"target": {
				Name:        "target",
				Type:        AttributeTypeString,
				HTMLAttr:    "target",
				Description: "Link target (_blank, _self, _parent, _top)",
			},
			"rel": {
				Name:        "rel",
				Type:        AttributeTypeString,
				HTMLAttr:    "rel",
				Description: "Link relationship",
			},
		},
		Description: "Navigation link element",
	})

	// List Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "List",
		HTMLTag:        "ul",
		Type:           ComponentTypeList,
		DefaultClasses: []string{"list-disc", "list-inside"},
		Attributes: map[string]*AttributeDefinition{
			"ordered": {
				Name:        "ordered",
				Type:        AttributeTypeBool,
				HTMLAttr:    "tag",
				Description: "Whether the list is ordered (ol) or unordered (ul)",
			},
			"listStyle": {
				Name:        "listStyle",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "List style (disc, decimal, none, etc.)",
			},
		},
		Description: "List container (ordered or unordered)",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "ListItem",
		HTMLTag:        "li",
		Type:           ComponentTypeList,
		DefaultClasses: []string{},
		Attributes: map[string]*AttributeDefinition{
			"value": {
				Name:        "value",
				Type:        AttributeTypeString,
				HTMLAttr:    "value",
				Description: "List item value (for ordered lists)",
			},
		},
		Description: "List item element",
	})

	// Media Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Image",
		HTMLTag:        "img",
		Type:           ComponentTypeMedia,
		DefaultClasses: []string{},
		SelfClosing:    true,
		Attributes: map[string]*AttributeDefinition{
			"src": {
				Name:        "src",
				Type:        AttributeTypeURL,
				HTMLAttr:    "src",
				Required:    true,
				Description: "Image source URL",
			},
			"alt": {
				Name:        "alt",
				Type:        AttributeTypeString,
				HTMLAttr:    "alt",
				Required:    true,
				Description: "Alternative text for accessibility",
			},
			"width": {
				Name:        "width",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Image width",
			},
			"height": {
				Name:        "height",
				Type:        AttributeTypeString,
				HTMLAttr:    "class",
				Description: "Image height",
			},
		},
		Description: "Image element",
	})

	// Semantic Layout Components
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Header",
		HTMLTag:        "header",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Header landmark element",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Footer",
		HTMLTag:        "footer",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Footer landmark element",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Nav",
		HTMLTag:        "nav",
		Type:           ComponentTypeNavigation,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Navigation landmark element",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Main",
		HTMLTag:        "main",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Main content landmark element",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Section",
		HTMLTag:        "section",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Section element for thematic grouping",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Article",
		HTMLTag:        "article",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Article element for standalone content",
	})

	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Aside",
		HTMLTag:        "aside",
		Type:           ComponentTypeLayout,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Aside element for sidebar content",
	})

	// Generic container
	cr.RegisterComponent(&ComponentDefinition{
		Name:           "Div",
		HTMLTag:        "div",
		Type:           ComponentTypeContainer,
		DefaultClasses: []string{},
		Attributes:     map[string]*AttributeDefinition{},
		Description:    "Generic container div element",
	})
}

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
func (cp *ComponentProcessor) GetHTMLTag(componentName string, properties []*PropertyNode) string {
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
func (cp *ComponentProcessor) BuildClasses(componentName string, properties []*PropertyNode) []string {
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
func (cp *ComponentProcessor) BuildAttributes(componentName string, properties []*PropertyNode) map[string]string {
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
func (cp *ComponentProcessor) ValidateComponent(componentName string, properties []*PropertyNode) []string {
	var errors []string

	comp, exists := cp.registry.GetComponent(componentName)
	if !exists {
		errors = append(errors, fmt.Sprintf("Unknown component: %s", componentName))
		return errors
	}

	// Check required attributes
	propertyMap := make(map[string]*PropertyNode)
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

func (cp *ComponentProcessor) GetTextContent(properties []*PropertyNode) string {
	for _, prop := range properties {
		if prop.Name == "text" {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

func (cp *ComponentProcessor) getVariantClasses(componentName string, properties []*PropertyNode) []string {
	variant := cp.getStringProperty(properties, "variant")
	if variant == "" {
		return []string{}
	}

	switch componentName {
	case "Text":
		return cp.getTextVariantClasses(variant)
	case "Button":
		return cp.getButtonVariantClasses(variant)
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

func (cp *ComponentProcessor) getButtonVariantClasses(variant string) []string {
	switch variant {
	case "primary":
		return []string{"bg-blue-500", "text-white", "hover:bg-blue-600"}
	case "secondary":
		return []string{"bg-gray-500", "text-white", "hover:bg-gray-600"}
	case "outline":
		return []string{"border", "border-blue-500", "text-blue-500", "hover:bg-blue-50"}
	case "ghost":
		return []string{"text-blue-500", "hover:bg-blue-50"}
	case "danger":
		return []string{"bg-red-500", "text-white", "hover:bg-red-600"}
	default:
		return []string{}
	}
}

func (cp *ComponentProcessor) getStyleClasses(properties []*PropertyNode) []string {
	for _, prop := range properties {
		if prop.Name == "style" {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return strings.Fields(str)
				}
			}
		}
	}
	return []string{}
}

func (cp *ComponentProcessor) getStringProperty(properties []*PropertyNode, name string) string {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

func (cp *ComponentProcessor) getBoolProperty(properties []*PropertyNode, name string) bool {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if b, ok := literal.Value.(bool); ok {
					return b
				}
			}
		}
	}
	return false
}

func (cp *ComponentProcessor) getIntProperty(properties []*PropertyNode, name string) int {
	for _, prop := range properties {
		if prop.Name == name {
			if literal, ok := prop.Value.(*LiteralNode); ok {
				if i, ok := literal.Value.(int); ok {
					return i
				}
			}
		}
	}
	return 0
}

func (cp *ComponentProcessor) formatPropertyValue(prop *PropertyNode, attrType AttributeType) string {
	if literal, ok := prop.Value.(*LiteralNode); ok {
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

func (cp *ComponentProcessor) validatePropertyType(prop *PropertyNode, expectedType AttributeType) bool {
	if literal, ok := prop.Value.(*LiteralNode); ok {
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
