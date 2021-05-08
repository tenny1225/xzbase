package xzbase

import (
	gj "github.com/segmentio/objconv/json"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"


	"reflect"
	"strconv"
	"strings"
)

type response struct {
	Code  int64       `json:"code"`
	Msg   string      `json:"msg"`
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

type RouteFunc struct {
	id         string
	params     []string
	method     reflect.Value
	methodName string
	controller controller
	returnType int
}

func (r *RouteFunc) Handle() func(*gin.Context) {

	return func(c *gin.Context) {
		defer c.Abort()

		params := make([]reflect.Value, 0)
		isError := false
		for k, t := range r.params {

			v := c.Param(fmt.Sprintf("id%d", k))

			if t == "int" {
				if n, e := strconv.Atoi(v); e == nil {
					params = append(params, reflect.ValueOf(n))
				} else {
					isError = true
				}
			} else if t == "int64" {
				if n, e := strconv.ParseInt(v, 0, 64); e == nil {
					params = append(params, reflect.ValueOf(n))
				} else {
					isError = true
				}
			} else if t == "float64" {
				if n, e := strconv.ParseFloat(v, 64); e == nil {
					params = append(params, reflect.ValueOf(n))
				} else {
					isError = true
				}
			} else if t == "float32" {
				if n, e := strconv.ParseFloat(v, 32); e == nil {
					params = append(params, reflect.ValueOf(float32(n)))
				} else {
					isError = true
				}
			} else if t == "string" {
				params = append(params, reflect.ValueOf(v))
			}
		}
		if !isError {
			//r.controller.SetContext(c)
			//values := r.method.Call(params)

			valueController := reflect.ValueOf(r.controller)

			ptrType := reflect.TypeOf(r.controller) //获取call的指针的reflect.Type

			trueType := ptrType.Elem() //获取type的真实类型

			ptrValue := reflect.New(trueType) //返回对象的指针对应的reflect.Value

			control := ptrValue.Interface().(controller)
			control.setContext(c)
			if _,ok:=ptrType.Elem().FieldByName("Service");ok{
				ptrValue.Elem().FieldByName("Service").Set(valueController.Elem().FieldByName("Service"))
			}


			values := ptrValue.MethodByName(r.methodName).Call(params)

			if values != nil && len(values) > 0 {
				c.Writer.Header().Set("Content-type", "application/json")
				c.Writer.WriteHeader(200)
				//obj:=values[0].Interface();
				if r.returnType == 1 {
					//c.JSON(200, response{Code: 200, Msg: "success", Data: values[0].Interface()})
					writeJson(c.Writer, response{Code: 200, Msg: "success", Data: values[0].Interface()})
				} else if r.returnType == 2 {
					switch values[0].Interface().(type) {
					case ZError:
						v := values[0].Interface().(ZError)
						//c.JSON(200, response{Code: v.Code, Msg: v.Msg})
						writeJson(c.Writer, response{Code: v.Code, Msg: v.Msg})
						return
					case error:
						v := values[0].Interface().(error)
						//c.JSON(200, response{Code: 1000, Msg: v.Error()})
						writeJson(c.Writer, response{Code: 1000, Msg: v.Error()})
						return
					}

					//c.JSON(200,response{Code:1000,Msg:values[0].Interface().(error).Error()})
					//c.JSON(200, response{Code: 200, Msg: "success", Data: values[0].Interface()})
					writeJson(c.Writer, response{Code: 200, Msg: "success", Data: values[0].Interface()})
				} else if r.returnType == 3 {
					switch values[1].Interface().(type) {
					case ZError:
						v := values[1].Interface().(ZError)
						//c.JSON(200, response{Code: v.Code, Msg: v.Msg})
						writeJson(c.Writer, response{Code: v.Code, Msg: v.Msg})
						return
					case error:
						v := values[1].Interface().(error)
						//c.JSON(200, response{Code: 1000, Msg: v.Error()})
						writeJson(c.Writer, response{Code: 1000, Msg: v.Error()})
						return

					}

					//c.JSON(200, response{Code: 200, Msg: "success", Data: values[0].Interface()})
					writeJson(c.Writer, response{Code: 200, Msg: "success", Data: values[0].Interface()})
				} else if r.returnType == 4 {
					switch values[2].Interface().(type) {
					case ZError:
						v := values[2].Interface().(ZError)
						//c.JSON(200, response{Code: v.Code, Msg: v.Msg})
						writeJson(c.Writer, response{Code: v.Code, Msg: v.Msg})
						return
					case error:
						v := values[2].Interface().(error)
						//c.JSON(200, response{Code: 1000, Msg: v.Error()})
						writeJson(c.Writer, response{Code: 1000, Msg: v.Error()})
						return
					}

					//c.JSON(200, response{Code: 200, Msg: "success", Data: values[0].Interface(), Total: values[1].Interface().(int64)})
					writeJson(c.Writer, response{Code: 200, Msg: "success", Data: values[0].Interface(), Total: values[1].Interface().(int64)})
				}
			}
		}

	}
}
func writeJson(writer gin.ResponseWriter, obj response) {

	b, e := json.Marshal(obj)
	if e != nil {
		b,e=gj.Marshal(obj)
		if e!=nil{
			b, _ = json.Marshal(response{Code: 1000, Msg: e.Error()})
		}


	}
	writer.Write(b)
}


func AddRoute(r *gin.Engine, i controller,s Service) {
	o := reflect.ValueOf(i)
	value := reflect.TypeOf(i)
	elem:=value.Elem()
	if elem!=nil{
		if _, ok := elem.FieldByName("Service"); ok &&s!=nil{
			values := reflect.ValueOf(s).Elem().Call(nil)
			if values != nil && len(values) > 0 {
				o.Elem().FieldByName("Service").Set(values[0])
			}
		}
	}

	name := strings.ToLower(strings.ReplaceAll(strings.Split(value.String(), ".")[1], "Controller", ""))

	rg := r.Group(name)
	num := value.NumMethod()
	for n := 0; n < num; n++ {
		m := value.Method(n)

		methodName := strings.ToLower(m.Name)

		method := o.MethodByName(m.Name)

		action := ""
		if strings.HasPrefix(methodName, "get") {
			action = "GET"
			methodName = strings.TrimPrefix(methodName, "get")
		} else if strings.HasPrefix(methodName, "post") {
			action = "POST"
			methodName = strings.TrimPrefix(methodName, "post")
		} else if strings.HasPrefix(methodName, "put") {
			action = "PUT"
			methodName = strings.TrimPrefix(methodName, "put")
		} else if strings.HasPrefix(methodName, "delete") {
			action = "DELETE"
			methodName = strings.TrimPrefix(methodName, "delete")
		}
		if action == "" {
			continue
		}

		returnType := 0
		if m.Func.Type().NumOut() == 1 {

			if m.Func.Type().Out(0).Name() == "error" {
				returnType = 2
			} else {
				returnType = 1
			}

		} else if m.Func.Type().NumOut() == 2 {
			if m.Func.Type().Out(1).Name() == "error" {
				returnType = 3
			}
		} else if m.Func.Type().NumOut() == 3 {
			if m.Func.Type().Out(2).Name() == "error" {
				returnType = 4
			}
		}
		fun := &RouteFunc{
			id:         name,
			params:     make([]string, 0),
			method:     method,
			methodName: m.Name,
			controller: i,
			returnType: returnType,
		}
		if strings.Contains(methodName, "by") {
			if methodName == "by" {
				methodName = "/:id0"
				typ := m.Func.Type()
				num := typ.NumIn()
				if num > 1 {
					fun.params = append(fun.params, typ.In(1).Name())
				}

			} else {
				list := strings.Split(methodName, "by")
				methodName = ""
				idIndex := 0
				for _, i := range list {
					if i == "" {
						methodName += fmt.Sprintf("/:id%d", idIndex)
						idIndex++
					} else {
						methodName += fmt.Sprintf("/%s", i)
					}

				}
				typ := m.Func.Type()
				num := typ.NumIn()
				for i := 1; i < num; i++ {
					t := typ.In(i).Name()
					fun.params = append(fun.params, t)
				}
			}

		}

		//fmt.Println(name,methodName)
		rg.Handle(action, methodName, fun.Handle())
	}
}
