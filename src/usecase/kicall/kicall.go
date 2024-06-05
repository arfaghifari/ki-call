package kicall

import (
	"context"
	"fmt"
	"reflect"

	"errors"

	myClient "github.com/arfaghifari/ki-call/src/client"
	"github.com/arfaghifari/ki-call/src/constant"
	"github.com/arfaghifari/ki-call/src/converter"
)

type Usecase interface {
	GetListService() ([]string, error)
	GetListMethod(svc string) ([]string, error)
	GetRequestMethod(svc, method string, noEmpty bool) (map[string]interface{}, error)
	GetResponseMethod(svc, method string, noEmpty bool) (map[string]interface{}, error)
	KiCall(host, svc, method string, req map[string]interface{}) (map[string]interface{}, error)
}

func NewUsecase() Usecase {
	return &usecase{}
}

type usecase struct{}

func (u *usecase) GetListService() ([]string, error) {
	cli := reflect.TypeOf(myClient.ClientKitex)

	var fields []string
	for i := 0; i < cli.NumField(); i++ {
		fields = append(fields, cli.Field(i).Name)
	}

	return fields, nil
}

func (u *usecase) GetListMethod(svc string) ([]string, error) {
	service := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !service.IsValid() {
		return []string{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	serviceType := service.Type()
	var methods []string
	for i := 0; i < serviceType.NumMethod(); i++ {
		methods = append(methods, serviceType.Method(i).Name)
	}

	return methods, nil
}

func (u *usecase) GetRequestMethod(svc, method string, noEmpty bool) (map[string]interface{}, error) {
	service := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !service.IsValid() {
		return map[string]interface{}{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	serviceType := service.Type()
	mthd, found := serviceType.MethodByName(method)
	if !found {
		return nil, fmt.Errorf(constant.ErrMethodNotFound, method)
	}

	inp := mthd.Type.In(1)
	req := reflect.New(inp.Elem()).Elem().Interface()

	return converter.TransformStructToMapJson(req, noEmpty)
}

func (u *usecase) GetResponseMethod(svc, method string, noEmpty bool) (map[string]interface{}, error) {
	service := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !service.IsValid() {
		return map[string]interface{}{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	serviceType := service.Type()
	mthd, found := serviceType.MethodByName(method)
	if !found {
		return nil, fmt.Errorf(constant.ErrMethodNotFound, method)
	}

	inp := mthd.Type.Out(0)
	req := reflect.New(inp.Elem()).Elem().Interface()

	return converter.TransformStructToMapJson(req, noEmpty)
}

func (u *usecase) KiCall(host, svc, method string, req map[string]interface{}) (map[string]interface{}, error) {
	var errKitex error
	myClient.ClientKitex.RegisterAllClient(host)
	service := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !service.IsValid() {
		return map[string]interface{}{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	serviceType := service.Type()
	mthd2 := service.MethodByName(method)
	mthd, found := serviceType.MethodByName(method)
	if !found {
		return nil, fmt.Errorf(constant.ErrMethodNotFound, method)
	}

	inp := mthd.Type.In(1)
	methodRequest, _ := converter.TransformMapJsonToStruct(req, inp)

	outFunc := mthd2.Call([]reflect.Value{
		reflect.ValueOf(context.Background()),
		reflect.ValueOf(methodRequest),
	})

	if !outFunc[1].IsNil() {
		errKitex = errors.New(outFunc[1].MethodByName("Error").Call([]reflect.Value{})[0].String())
		if !outFunc[0].IsNil() {
			resp, _ := converter.TransformStructToMapJson(outFunc[0].Interface(), false)
			return resp, errKitex
		}
		return nil, errKitex
	}
	return converter.TransformStructToMapJson(outFunc[0].Interface(), false)

}
