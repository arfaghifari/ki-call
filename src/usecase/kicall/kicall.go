package kicall

import (
	"context"
	"fmt"
	"log"
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
	KiCall(host, svc, method string, req map[string]interface{}) (map[string]interface{}, error)
}

func NewUsecase() Usecase {
	return &usecase{}
}

type usecase struct{}

func (u *usecase) GetListService() ([]string, error) {
	// cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	cli := reflect.TypeOf(myClient.ClientKitex)

	var fields []string
	for i := 0; i < cli.NumField(); i++ {
		fields = append(fields, cli.Field(i).Name)
	}

	return fields, nil
}

func (u *usecase) GetListMethod(svc string) ([]string, error) {
	// cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	cli := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !cli.IsValid() {
		return []string{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	cliType := cli.Type()
	var methods []string
	for i := 0; i < cliType.NumMethod(); i++ {
		methods = append(methods, cliType.Method(i).Name)
	}

	return methods, nil
}

func (u *usecase) GetRequestMethod(svc, method string, noEmpty bool) (map[string]interface{}, error) {
	// cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	cli := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !cli.IsValid() {
		return map[string]interface{}{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	cliType := cli.Type()
	mthd, found := cliType.MethodByName(method)
	if !found {
		return nil, fmt.Errorf(constant.ErrMethodNotFound, method)
	}

	inp := mthd.Type.In(1)
	log.Println(inp)
	req := reflect.New(inp.Elem()).Elem()

	return converter.TransformStructToMapJson(req, noEmpty)
}

func (u *usecase) KiCall(host, svc, method string, req map[string]interface{}) (map[string]interface{}, error) {
	var errKitex error
	myClient.ClientKitex.RegisterAllClient(host)
	cli := reflect.ValueOf(myClient.ClientKitex).FieldByName(svc)
	if !cli.IsValid() {
		return map[string]interface{}{}, fmt.Errorf(constant.ErrServiceNotFound, svc)
	}
	cliType := cli.Type()
	// cli2, _ := kitexClient.NewClient("test", client.WithHostPorts(host))
	mthd2 := cli.MethodByName(method)
	// cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	mthd, found := cliType.MethodByName(method)
	if !found {
		return nil, fmt.Errorf(constant.ErrMethodNotFound, method)
	}

	inp := mthd.Type.In(1)
	reqs := reflect.New(inp.Elem()).Elem()

	reqs2, _ := converter.TransformMapJsonToStruct(req, reqs)

	vp := reflect.New(reqs2.Type())
	vp.Elem().Set(reqs2)

	outFunc := mthd2.Call([]reflect.Value{
		reflect.ValueOf(context.Background()),
		vp,
	})

	if !outFunc[1].IsNil() {
		errKitex = errors.New(outFunc[1].MethodByName("Error").Call([]reflect.Value{})[0].String())
		if !outFunc[0].IsNil() {
			resp, _ := converter.TransformStructToMapJson(outFunc[0].Elem(), false)
			return resp, errKitex
		}
		return nil, errKitex
	}
	return converter.TransformStructToMapJson(outFunc[0].Elem(), false)

}
