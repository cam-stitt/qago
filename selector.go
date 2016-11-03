package qago

import (
	"fmt"
	"reflect"
)

const (
	// Default is the default selector type
	Default string = "default"
	// ByButton selects by button
	ByButton = "by_button"
	// ByClass selects by element class
	ByClass = "by_class"
	// ByID selects by element ID
	ByID = "by_id"
	// ByLabel selects by element label
	ByLabel = "by_label"
	// ByLink selects by element link
	ByLink = "by_link"
	// ByName selects by element name attribute
	ByName = "by_name"
	// ByXPath selects by provided xpath
	ByXPath = "by_xpath"
	// ForAppium selects using appium selector
	ForAppium = "for_appium"
)

const (
	// Find is the default finder
	Find = "find"
	// Multi is for finding multiple elements
	Multi = "multi"
	// First finds the first element
	First = "first"
)

// ExtractMethodArgs will return the arguments that should be passed into
// the selector method.
func ExtractMethodArgs(selectorType string, predicate *Predicate) []reflect.Value {
	args := []reflect.Value{
		reflect.ValueOf(predicate.Selector),
	}
	switch selectorType {
	case ForAppium:
		args = append(args, reflect.ValueOf(predicate.Text))
	}
	return args
}

// FindSelectMethod will return a string representation of the method to run
func FindSelectMethod(selector, selectorType string) string {
	var prefix string
	var suffix string

	switch selectorType {
	case Find:
		prefix = "Find"
	case Multi:
		prefix = "All"
	case First:
		prefix = "First"
	}

	switch selector {
	case Default:
		suffix = ""
	case ByButton:
		suffix = "ByButton"
	case ByClass:
		suffix = "ByClass"
	case ByID:
		suffix = "ByID"
	case ByLabel:
		suffix = "ByLabel"
	case ByLink:
		suffix = "ByLink"
	case ByName:
		suffix = "ByName"
	case ByXPath:
		suffix = "ByXPath"
	case ForAppium:
		suffix = "ForAppium"
	}

	return fmt.Sprintf("%s%s", prefix, suffix)
}
