package kicall

import (
	"context"
	"reflect"

	"errors"

	kitexClient "github.com/arfaghifari/ki-call/kitex_gen/merchantvoucher/merchantvoucher"
	"github.com/arfaghifari/ki-call/src/converter"
)

type Usecase interface {
	GetListMethod() ([]string, error)
	GetRequestMethod(method string, noEmpty bool) (map[string]interface{}, error)
	KiCall(host string, method string, req map[string]interface{}) (map[string]interface{}, error)
}

func NewUsecase() Usecase {
	return &usecase{}
}

type usecase struct{}

func (u *usecase) GetListMethod() ([]string, error) {
	cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	var methods []string
	for i := 0; i < cli.NumMethod(); i++ {
		methods = append(methods, cli.Method(i).Name)
	}

	return methods, nil
}

func (u *usecase) GetRequestMethod(method string, noEmpty bool) (map[string]interface{}, error) {
	cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	mthd, found := cli.MethodByName(method)
	if !found {
		return nil, errors.New("method not exist")
	}

	inp := mthd.Type.In(1)
	req := reflect.New(inp.Elem()).Elem()

	return converter.TransformStructToMapJson(req, noEmpty)
}

func (u *usecase) KiCall(host string, method string, req map[string]interface{}) (map[string]interface{}, error) {
	cli2, _ := kitexClient.NewClient("test")
	mthd2 := reflect.ValueOf(cli2).MethodByName(method)

	cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	mthd, found := cli.MethodByName(method)
	if !found {
		return nil, errors.New("method not exist")
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
		return nil, nil
	}

	return converter.TransformStructToMapJson(outFunc[0].Elem(), false)

}
