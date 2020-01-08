package main

import (
	"context"
	"github.com/common-go/mongo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
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
		err = t.Create()
	case "update":
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

func main() {

	e := echo.New()
	var r2 RouteList
	_ = ReadRouteFromMongo(&r2)
	AddRoute(r2, e)

	e.Logger.Fatal(e.Start(":1323"))
}
