package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/common-go/mongo"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
)

type RouteInfo struct {
	Id             string `bson:"_id"`
	Source         string `bson:"source"`
	DriverName     string `bson:"driverName"`
	Path           string `bson:"path"`
	Method         string `bson:"method"`
	Detail         string `bson:"detail"`
	DbName         string `bson:"dbName"`
	DataSourceName string `bson:"dataSourceName"`
}
type RouteList []RouteInfo
type DbMethod interface {
	GetAll() error
	GetById() error
	Create() error
	Update() error
	Delete() error
}
type DbType struct {
	Info         RouteInfo
	OutputString *string
	IdParam      string
	InputMap     map[string]interface{}
}
type MongoType DbType
type SqlType DbType

type ErrorMessage struct {
	Field   string `json:"field,omitempty" bson:"field,omitempty" gorm:"column:field"`
	Code    string `json:"code,omitempty" bson:"code,omitempty" gorm:"column:code"`
	Message string `json:"message,omitempty" bson:"message,omitempty" gorm:"column:message"`
}

type ModelJSON struct {
	Env        string          `json:"env"`
	Name       string          `json:"name"`
	Source     string          `json:"source"`
	ConstValue []Const         `json:"const"`
	TypeAlias  []TypeAlias     `json:"type_alias"`
	Fields     []FieldElements `json:"fields"`
}
type Const struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TypeAlias struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FieldElements struct {
	Name       string `json:"name"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	PrimaryKey bool   `json:"primaryKey"`
}

var input = map[string]interface{}{}

func (info RouteInfo) PathHandler(c echo.Context) error {

	var output string

	err := c.Bind(&input)
	if err != nil {
		return err
	}
	delete(input, "id")
	id := c.Param("id")
	tmp := DbType{
		Info:         info,
		OutputString: &output,
		IdParam:      id,
		InputMap:     input,
	}
	var t DbMethod
	switch info.DriverName {
	case "mongo":
		t = MongoType(tmp)
	default:
		t = SqlType(tmp)

	}

	switch info.Detail {
	case "getById":
		err = t.GetById()
	case "getAll":
		err = t.GetAll()
	case "create":
		var tmp []ErrorMessage
		tmp, input = Validate(input, table[info.Source])
		log.Print(tmp)
		if len(tmp) != 0 {
			return c.JSON(http.StatusBadRequest, tmp)
		}
		err = t.Create()
	case "update":
		var tmp []ErrorMessage
		tmp, input = Validate(input, table[info.Source])
		log.Print(tmp)
		if len(tmp) != 0 {
			return c.JSON(http.StatusBadRequest, tmp)
		}
		err = t.Update()
	case "delete":
		err = t.Delete()
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, output)
}

func AddRoute(r RouteList, e *echo.Echo) {
	for i := range r {
		e.Match([]string{r[i].Method}, r[i].Path, r[i].PathHandler)
	}
}

func ReadRouteFromMongo(r *RouteList) error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, "mongodb://localhost:27017", "evaluation")
	if err != nil {
		return err
	}
	collection := db.Collection("route")
	result, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	var r2 RouteInfo
	for result.Next(ctx) {
		err = result.Decode(&r2)
		if err != nil {
			return err
		}
		*r = append(*r, r2)
	}

	return err
}
func ReadSchemaFromMongo() error {

	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, "mongodb://localhost:27017", "evaluation")
	if err != nil {
		return err
	}
	collection := db.Collection("schema")
	result, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	r := []ModelJSON{}
	for result.Next(ctx) {
		ri := ModelJSON{}
		err = result.Decode(&ri)
		if err != nil {
			return err
		}
		r = append(r, ri)
	}
	for i := range r {
		table[r[i].Name] = InnerMap{}
		for j := range r[i].Fields {
			table[r[i].Name][r[i].Fields[j].Source] = r[i].Fields[j].Type
		}
	}

	return err
}

func IsBolean(x interface{}) (bool, bool) {
	switch x {
	case true, "true", "True":
		return true, true
	case false, "false", "False":
		return true, false
	default:
		return false, false
	}
}
func Validate(input map[string]interface{}, instruct InnerMap) ([]ErrorMessage, map[string]interface{}) {
	errList := []ErrorMessage{}
	for k, v := range input {
		ok := true

		switch instruct[k] {
		case "int":
			_, ok = v.(float64)
			if !ok {
				if _, t := v.(string); t {

					if x, err := strconv.ParseFloat(v.(string), 64); err == nil {
						input[k] = x
						ok = true

					}
				}
			}
		case "string":
			continue
		case "time.Time":
			ok = false
			layout := "2006-01-02T15:04:05.000Z"
			if _, t := v.(string); t {
				if x, err := time.Parse(layout, v.(string)); err == nil {
					input[k] = x
					ok = true
				}
			}

		case "boolean":
			x := true
			ok, x = IsBolean(v)
			if ok {
				input[k] = x
			}
		case "email":
			ok = false
			if _, t := v.(string); t {
				ok = valid.IsEmail(v.(string))
			}
		case "url":
			ok = false
			if _, t := v.(string); t {
				ok = valid.IsURL(v.(string))
			}
		case "":
			ok = false

		default:
			if instruct[k][:2] == "[]" {
				tmpMap2 := []map[string]interface{}{}
				x, t := v.([]interface{})
				var errList1 []ErrorMessage
				if t {
					for i := range x {
						x1, t1 := x[i].(map[string]interface{})
						if !t1 {
							t = false
							break
						}
						returnErr, returnMap := Validate(x1, table[instruct[k][2:]])
						tmpMap2 = append(tmpMap2, returnMap)
						errList1 = append(errList1, returnErr...)
					}
					if t {
						input[k] = tmpMap2
						errList = append(errList, errList1...)
					}
				}
				ok = t

			} else {
				x, t := v.(map[string]interface{})
				if t {
					returnErr, returnMap := Validate(x, table[instruct[k]])
					input[k] = returnMap
					errList = append(errList, returnErr...)
				}
				ok = t
			}
		}
		if !ok {
			mess := "Wrong format: " + k + " must be " + instruct[k]
			if instruct[k] == "" {
				mess = "Non-existed field: " + k
			}
			errList = append(errList, ErrorMessage{
				Field:   k,
				Code:    "",
				Message: mess,
			})
		}
	}

	return errList, input
}

type OuterMap map[string]map[string]string
type InnerMap map[string]string

var table = OuterMap{}

func main() {

	e := echo.New()
	var r2 RouteList

	_ = ReadSchemaFromMongo()
	_ = ReadRouteFromMongo(&r2)
	AddRoute(r2, e)

	e.Logger.Fatal(e.Start(":1323"))
}
