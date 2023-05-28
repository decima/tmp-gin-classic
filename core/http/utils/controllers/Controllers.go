package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"sigomid/core/http/utils/configuration"
)

const ANY_METHOD = "ANY"

const PUBLIC_GROUP = "PUBLIC"

const routeTag = "route"
const groupTag = "group"
const methodTag = "method"

var allControllers []*ControllerInterface

func Register(controllerInterface ControllerInterface) {
	fmt.Println("Registering controller")
	allControllers = append(allControllers, &controllerInterface)

}

func extractRoutesFromController(controllerInterface *ControllerInterface) *ControllerDeclaration {

	var actions = []*declaredAction{}

	var controller *ControllerDeclaration

	interfaceValue := reflect.ValueOf(*controllerInterface)
	realValue := reflect.Indirect(interfaceValue)
	realType := realValue.Type()

	for i := 0; i < realType.NumField(); i++ {

		if (realType.Field(i).Type.Kind() == reflect.TypeOf(&DefaultController{}).Kind()) {
			controller = createControllerDeclaration(realType.Field(i))
			continue
		} else if _, ok := realType.Field(i).Tag.Lookup(routeTag); !ok {
			continue
		}
		actions = append(actions, extractDeclaredAction(realValue, realType, i))

	}
	for _, action := range actions {
		action.Route = controller.prefixRoute + action.Route
		if action.Method == "" {
			action.Method = controller.defaultMethod
		}
		if action.Group == "" {
			action.Group = controller.defaultGroup
		}
	}
	controller.Actions = actions
	return controller
}

func createControllerDeclaration(s reflect.StructField) *ControllerDeclaration {
	c := &ControllerDeclaration{
		prefixRoute:   "",
		defaultMethod: ANY_METHOD,
		defaultGroup:  PUBLIC_GROUP,
	}
	if value, ok := s.Tag.Lookup(routeTag); ok && value != "" {
		c.prefixRoute = value
	}
	if value, ok := s.Tag.Lookup(groupTag); ok && value != "" {
		c.defaultGroup = value
	}
	if value, ok := s.Tag.Lookup(methodTag); ok && value != "" {
		c.defaultMethod = value
	}

	return c
}

func extractDeclaredAction(v reflect.Value, t reflect.Type, fieldId int) *declaredAction {
	currentRoute := t.Field(fieldId).Tag.Get(routeTag)
	requiredGroup := t.Field(fieldId).Tag.Get(groupTag)
	method := t.Field(fieldId).Tag.Get(methodTag)
	return &declaredAction{
		Route:  currentRoute,
		Method: method,
		Group:  requiredGroup,
		Do:     v.FieldByName(t.Field(fieldId).Name).Interface().(gin.HandlerFunc),
	}
}

func GetComputed(configuration configuration.ServerConfiguration) []*ControllerDeclaration {
	computedControllers := []*ControllerDeclaration{}
	for _, controller := range allControllers {
		(*controller).SetConfiguration(configuration)
		computedControllers = append(computedControllers, extractRoutesFromController(controller))

	}

	return computedControllers
}

type ControllerInterface interface {
	SetConfiguration(s configuration.ServerConfiguration)
}

type Action gin.HandlerFunc

type declaredAction struct {
	Route  string
	Method string
	Group  string
	Do     gin.HandlerFunc
}

type DefaultController struct {
	ServerConfiguration configuration.ServerConfiguration
}

func (d DefaultController) SetConfiguration(s configuration.ServerConfiguration) {
	d.ServerConfiguration = s
}

type ControllerDeclaration struct {
	Actions       []*declaredAction
	prefixRoute   string
	defaultGroup  string
	defaultMethod string
}
