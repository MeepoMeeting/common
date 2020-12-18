package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type function struct {
	funcType reflect.Value
	in       []reflect.Type
	out      []reflect.Type
}

type MethodsMap struct {
	funcMap map[string]function
	imp     interface{} // imp
}

func (methodsMap *MethodsMap) Init(imp interface{}) {
	impType := reflect.TypeOf(imp)
	methodsMap.imp = imp
	methodsMap.funcMap = make(map[string]function)
	for i := 0; i < impType.NumMethod(); i++ {
		method := impType.Method(i)
		// 此处函数参数必须为 this, context,req*,
		if method.Type.NumIn() != 3 {
			continue
		}
		// 返回值必须为rsp*,err
		if method.Type.NumOut() != 2 {
			continue
		}
		f := function{
			in:  []reflect.Type{},
			out: []reflect.Type{},
		}
		for j := 0; j < method.Type.NumIn(); j++ {
			f.in = append(f.in, method.Type.In(j))
		}
		for j := 0; j < method.Type.NumOut(); j++ {
			f.out = append(f.out, method.Type.Out(j))
		}
		f.funcType = method.Func
		methodsMap.funcMap[method.Name] = f
	}
}

func (methodsMap *MethodsMap) Call(methodName, req string) (string, error) {
	method, ok := methodsMap.funcMap[methodName]
	if !ok {
		return "", fmt.Errorf("%s,MethodNotFound", methodName)
	}
	reqValue := reflect.New(method.in[2].Elem()).Interface()
	err := json.Unmarshal([]byte(req), &reqValue)
	if err != nil {
		return "", fmt.Errorf("req type error:%v", req)
	}
	values := []reflect.Value{reflect.ValueOf(methodsMap.imp), reflect.ValueOf(context.Background()), reflect.ValueOf(reqValue)}

	out := method.funcType.Call(values)
	if len(out) != 2 {
		return "", fmt.Errorf("values returned num error")
	}
	jsonstr, err := json.Marshal(out[0].Interface())
	if err != nil {
		return "", fmt.Errorf("rsp type error:%v", err)
	}
	// 保留原始error
	if out[1].Interface() == nil {
		return string(jsonstr), nil
	}
	callErr := out[1].Interface().(error)
	return string(jsonstr), callErr
}

func (methodsMap *MethodsMap) RegisteHttpRouter(router *httprouter.Router, basePath string) {
	for k, _ := range methodsMap.funcMap {
		methodName := k
		router.POST(fmt.Sprintf("%s/%s", basePath, methodName), func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			reqbyte, _ := ioutil.ReadAll(r.Body)
			reqStr := string(reqbyte)
			rspStr, err := methodsMap.Call(methodName, reqStr)
			if err != nil {
				log.Printf("http call error|method:%v|error:%v\n", methodName, err)
				w.WriteHeader(404)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(rspStr))
			log.Printf("http call success|method:%v\n", methodName)

		})
	}
}
