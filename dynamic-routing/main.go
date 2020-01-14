package main

import (
	"context"
	"github.com/common-go/mongo"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
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
type ModelJSONList []ModelJSON
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

func DateTimeValidate(fieldName string, valueString string, v *time.Time, errList *[]ErrorMessage) error {
	layout := "2006-01-02T15:04:05.000Z"
	time2, err := time.Parse(layout, valueString)
	if err != nil {
		*errList = append(*errList, ErrorMessage{
			Field:   fieldName,
			Code:    "",
			Message: "Field " + fieldName + ": Invalid Date Time Format",
		})
		return err
	}
	*v = time2
	return nil
}

func IntValidate(k string, valueString string, v *int64, errList *[]ErrorMessage) error {
	x, err := strconv.ParseInt(valueString, 0, 64)
	if err != nil {

		*errList = append(*errList, ErrorMessage{
			Field:   k,
			Code:    "",
			Message: "Field " + k + ": Invalid Integer Format",
		})
		return err
	} else {
		*v = x
	}
	return nil
}
func Validator(source string, mo ModelJSONList, input map[string]interface{}) []ErrorMessage {
	var errList []ErrorMessage
	t := -1
	for i := range mo {
		if mo[i].Name == source {
			t = i
			break
		}
	}
	if t == -1 {
		return []ErrorMessage{{
			Field:   "",
			Code:    "",
			Message: "Table Not Found",
		}}
	}

	for k, v := range input {
		valueString, ok := v.(string)
		contain := false
		if ok {
			for i := range mo[t].Fields {
				if k == mo[t].Fields[i].Source {
					contain = true
					switch mo[t].Fields[i].Type {
					case "time.Time":
						var tmp time.Time
						err := DateTimeValidate(k, valueString, &tmp, &errList)
						if err == nil {
							(input)[k] = tmp
						}
					case "int":
						var tmp int64
						err := IntValidate(k, valueString, &tmp, &errList)
						if err == nil {
							(input)[k] = tmp
						}
					}
					break

				}
			}

			if !contain {
				errList = append(errList, ErrorMessage{
					Field:   "",
					Code:    "",
					Message: "Field " + k + " Not Existed",
				})
			}
		}
	}
	return errList
}

func (info RouteInfo) PathHandler(c echo.Context) error {

	var output string
	input := map[string]interface{}{}
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
		_ = Validator(info.Source, m, input)
		err = t.Create()
	case "update":
		_ = Validator(info.Source, m, input)
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
	collection := db.Collection("sqlServerRoute")
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
func ReadSchemaFromMongo(r *ModelJSONList) error {
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
	var r2 ModelJSON
	for result.Next(ctx) {
		err = result.Decode(&r2)
		if err != nil {
			return err
		}
		*r = append(*r, r2)
	}

	return err
}

var m ModelJSONList

func main() {

	e := echo.New()
	var r2 RouteList

	_ = ReadSchemaFromMongo(&m)
	_ = ReadRouteFromMongo(&r2)
	AddRoute(r2, e)

	e.Logger.Fatal(e.Start(":1323"))
}
