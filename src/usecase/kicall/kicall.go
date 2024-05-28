package kicall

import (
	"context"
	"reflect"
	"strings"

	"errors"

	kitexClient "github.com/arfaghifari/ki-call/kitex_gen/merchantvoucher/merchantvoucher"
)

type Usecase interface {
	GetListMethod() ([]string, error)
	GetRequestMethod(method string) (map[string]interface{}, error)
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

func (u *usecase) GetRequestMethod(method string) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	mthd, found := cli.MethodByName(method)
	if !found {
		return nil, errors.New("method not exist")
	}

	inp := mthd.Type.In(1)
	req := reflect.New(inp.Elem()).Elem()
	reqT := req.Type()

	for i := 0; i < reqT.NumField(); i++ {
		f := reqT.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		res[strings.Split(v, ",")[0]] = req.FieldByName(f.Name).Interface()
	}

	return res, nil
}

func (u *usecase) KiCall(host string, method string, req map[string]interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	cli2, _ := kitexClient.NewClient("test")
	mthd2 := reflect.ValueOf(cli2).MethodByName(method)

	cli := reflect.TypeOf((*kitexClient.Client)(nil)).Elem()
	mthd, found := cli.MethodByName(method)
	if !found {
		return nil, errors.New("method not exist")
	}

	inp := mthd.Type.In(1)
	reqs := reflect.New(inp.Elem()).Elem()
	reqT := reqs.Type()

	for i := 0; i < reqT.NumField(); i++ {
		f := reqT.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		jsonName := strings.Split(v, ",")[0]
		val, ok := req[jsonName]
		if !ok || val == nil {
			continue
		}
		reqs.FieldByName(f.Name).Set(reflect.ValueOf(val).Convert(f.Type))
	}

	vp := reflect.New(reqT)
	vp.Elem().Set(reqs)

	outFunc := mthd2.Call([]reflect.Value{
		reflect.ValueOf(context.Background()),
		vp,
	})

	out := mthd.Type.Out(0)
	resp := reflect.New(out.Elem()).Elem()
	resT := resp.Type()

	for i := 0; i < resT.NumField(); i++ {
		f := resT.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok || f.Name != "Status" {
			continue
		}

		res[strings.Split(v, ",")[0]] = outFunc[0].Elem().FieldByName(f.Name).Interface()
	}

	return res, nil
}
