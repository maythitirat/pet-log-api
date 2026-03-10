package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Validate validates a struct based on validate tags
// Returns a map of field names to error messages
func Validate(s interface{}) map[string]string {
	errors := make(map[string]string)
	
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	
	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return errors
	}
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		
		// Get JSON field name for error messages
		jsonName := field.Tag.Get("json")
		if jsonName == "" {
			jsonName = strings.ToLower(field.Name)
		} else {
			// Handle json tags like "name,omitempty"
			jsonName = strings.Split(jsonName, ",")[0]
		}
		
		// Parse validation rules
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			err := validateRule(value, rule)
			if err != nil {
				errors[jsonName] = err.Error()
				break // One error per field is enough
			}
		}
	}
	
	return errors
}

// validateRule validates a single rule against a value
func validateRule(v reflect.Value, rule string) error {
	// Handle pointer types - check if nil
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// Only required validation applies to nil pointers
			if rule == "required" {
				return fmt.Errorf("this field is required")
			}
			return nil // Skip other validations for nil pointers
		}
		v = v.Elem()
	}
	
	// Skip validation for omitempty if value is zero
	if strings.HasPrefix(rule, "omitempty") {
		if isZero(v) {
			return nil
		}
		// Continue with other rules if not empty
		return nil
	}
	
	switch {
	case rule == "required":
		if isZero(v) {
			return fmt.Errorf("this field is required")
		}
		
	case strings.HasPrefix(rule, "min="):
		minStr := strings.TrimPrefix(rule, "min=")
		var min int
		fmt.Sscanf(minStr, "%d", &min)
		
		switch v.Kind() {
		case reflect.String:
			if len(v.String()) < min {
				return fmt.Errorf("must be at least %d characters", min)
			}
		case reflect.Int, reflect.Int64:
			if v.Int() < int64(min) {
				return fmt.Errorf("must be at least %d", min)
			}
		case reflect.Float64:
			if v.Float() < float64(min) {
				return fmt.Errorf("must be at least %d", min)
			}
		}
		
	case strings.HasPrefix(rule, "max="):
		maxStr := strings.TrimPrefix(rule, "max=")
		var max int
		fmt.Sscanf(maxStr, "%d", &max)
		
		switch v.Kind() {
		case reflect.String:
			if len(v.String()) > max {
				return fmt.Errorf("must be at most %d characters", max)
			}
		case reflect.Int, reflect.Int64:
			if v.Int() > int64(max) {
				return fmt.Errorf("must be at most %d", max)
			}
		case reflect.Float64:
			if v.Float() > float64(max) {
				return fmt.Errorf("must be at most %d", max)
			}
		}
		
	case rule == "email":
		if v.Kind() == reflect.String {
			email := v.String()
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				return fmt.Errorf("must be a valid email address")
			}
		}
		
	case strings.HasPrefix(rule, "gt="):
		gtStr := strings.TrimPrefix(rule, "gt=")
		var gt float64
		fmt.Sscanf(gtStr, "%f", &gt)
		
		switch v.Kind() {
		case reflect.Int, reflect.Int64:
			if float64(v.Int()) <= gt {
				return fmt.Errorf("must be greater than %v", gt)
			}
		case reflect.Float64:
			if v.Float() <= gt {
				return fmt.Errorf("must be greater than %v", gt)
			}
		}
	}
	
	return nil
}

// isZero checks if a value is the zero value for its type
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	default:
		return false
	}
}
