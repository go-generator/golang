package impl

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"../../model"
	m "github.com/common-go/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingServiceImpl struct {
	*m.DefaultViewService
	*m.DefaultSearchService
}

func NewBookingServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *BookingServiceImpl {
	var model model.Booking
	typeOfModel := reflect.TypeOf(model)
	viewService, searchService := m.NewMongoViewSearchService(typeOfModel, db, "booking", searchBuilder, false)
	return &BookingServiceImpl{viewService, searchService}
}

func (b *BookingServiceImpl) SaveDraft(cxt context.Context, booking model.Booking) (*model.Booking, error) {
	if len(booking.BookingId) == 0 {
		id := uuid.New()
		booking.BookingId = strings.Replace(id.String(), "-", "", -1)
	}
	booking.Status = "N"
	var query = bson.M{
		"_id": booking.BookingId,
	}
	_, err := m.UpsertOne(cxt, b.DefaultViewService.Collection, query, booking)
	if err == nil {
		return &booking, nil
	}
	return nil, err
}

func (b *BookingServiceImpl) Cancel(cxt context.Context, bookingId string) (bool, error) {
	result, err := b.GetById(cxt, bookingId)
	if err != nil {
		return false, err
	}
	if booking, ok := result.(*model.Booking); ok {
		if booking.Status == model.News {
			var query = bson.M{
				"_id": bookingId,
			}
			_, err1 := m.DeleteOne(cxt, b.DefaultViewService.Collection, query)
			if err1 != nil {
				return false, err1
			}
			return true, nil
		}
	}
	return false, nil
}

func (b *BookingServiceImpl) GetFreeLocationByBookable(ctx context.Context, bookableId string, date string) ([]bool, error) {
	result := make([]bool, 48)
	location, _ := time.LoadLocation("UTC")
	fromDate, err := time.ParseInLocation(layoutDate, date, location)
	if err != nil {
		return nil, err
	}
	dateEnd := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 23, 59, 59, 0, location)
	var query = bson.M{
		"bookableId": bookableId,
		"status":     "A",
		"startBookingTime": bson.M{
			"$gte": fromDate,
			"$lte": dateEnd,
		},
	}
	var model1 []model.Booking
	//typeOfModel := reflect.TypeOf(model1)
	_, err = m.FindAndDecode(ctx, b.Collection, query, &model1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(model1); i++ {
		end := model1[i].EndBookingTime
		start := model1[i].StartBookingTime
		minutes := end.Sub(*start).Minutes()
		hour := model1[i].StartBookingTime.Hour()
		minuteStart := model1[i].StartBookingTime.Minute()
		indexStartCheck := hour * 2
		if minuteStart == 30 {
			indexStartCheck++
		}
		minutesSpilit := minutes / 30
		for j := indexStartCheck; j < indexStartCheck+int(minutesSpilit); j++ {
			result[j] = true
		}
	}
	return result, nil
}

func (b *BookingServiceImpl) Submit(ctx context.Context, booking model.Booking) (*model.Booking, error) {
	if len(booking.BookingId) == 0 {
		id := uuid.New()
		booking.BookingId = strings.Replace(id.String(), "-", "", -1)
	}
	var query = bson.M{
		"_id": booking.BookingId,
	}
	objDate := booking.StartBookingTime.Format("2006-01-02")
	result1, _ := b.GetFreeLocationByBookable(ctx, booking.BookableId, objDate)
	end := booking.EndBookingTime
	start := booking.StartBookingTime
	minutes := end.Sub(*start).Minutes()
	hour := booking.StartBookingTime.Hour()
	minuteStart := booking.StartBookingTime.Minute()
	indexStartCheck := hour * 2
	if minuteStart == 30 {
		indexStartCheck++
	}
	minutesSpilit := minutes / 30
	if len(result1) <= int(minutesSpilit) {
		minutesSpilit = float64(len(result1))
	}
	smallArray := result1[indexStartCheck : indexStartCheck+int(minutesSpilit+1)]
	indexTrue := false
	for _, v := range smallArray {
		if v == true {
			indexTrue = true
			break
		}
	}
	if indexTrue {
		booking.Status = "C"
	} else {
		booking.Status = "A"
	}
	_, err := m.UpsertOne(ctx, b.DefaultViewService.Collection, query, booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (b *BookingServiceImpl) GetFreeLocationByBookableList(ctx context.Context, bookableIdList []interface{}, date string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, v := range bookableIdList {
		res, _ := b.GetFreeLocationByBookable(ctx, v.(string), date)
		result = append(result, res)
	}
	return result, nil
}

func (b *BookingServiceImpl) GetLocationFreeInTime(ctx context.Context, startDate string, endDate string) ([]interface{}, error) {
	location, _ := time.LoadLocation("UTC")
	fromDate, _ := time.ParseInLocation(layoutDate, startDate, location)
	toDate, _ := time.ParseInLocation(layoutDate, endDate, location)
	pipeline := []bson.M{
		{"$lookup": bson.M{"from": "bookable", "localField": "bookableId", "foreignField": "_id", "as": "bookables"}},
	}
	a, _ := b.Collection.Aggregate(ctx, pipeline)
	temporaryBytes := []byte{}
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	results := make([]bson.M, 0)
	for _, v := range resp {
		if v["status"] == "A" {
			results = append(results, v)
		}
	}
	listnobookable := make([]interface{}, 0)
	listbookable := results[0]["bookables"]
	listLocationFree := make([]interface{}, 0)
	listbookDocuments := make([]map[string]interface{}, 0)
	temporaryBytes, _ = json.Marshal(listbookable)
	_ = json.Unmarshal(temporaryBytes, &listbookDocuments)

	for i := 0; i < len(listbookDocuments); i++ {
		for j := 0; j < len(results); j++ {
			if results[j]["bookableId"] == listbookDocuments[i]["_id"] {
				//bookalbeid := listbookDocuments[i]["_id"]
				if fromDate.Sub(results[j]["startBookingTime"].(primitive.DateTime).Time()) < 0 && toDate.Sub(results[j]["startBookingTime"].(primitive.DateTime).Time()) <= 0 {
					// Check it is exist in listnobookable?
					index := indexOf(listbookDocuments[i], listnobookable)
					if index < 0 {
						index1 := indexOf(listbookDocuments[i], listLocationFree)
						if index1 < 0 {
							listLocationFree = append(listLocationFree, listbookDocuments[i])
						}
					}
				} else if fromDate.Sub(results[j]["startBookingTime"].(primitive.DateTime).Time()) > 0 && fromDate.Sub(results[j]["endBookingTime"].(primitive.DateTime).Time()) >= 0 {
					index := indexOf(listbookDocuments[i], listnobookable)
					if index < 0 {
						index1 := indexOf(listbookDocuments[i], listLocationFree)
						if index1 < 0 {
							listLocationFree = append(listLocationFree, listbookDocuments[i])
						}
					}
				} else { // Delete exist bookable
					index1 := indexOf(listbookDocuments[i], listnobookable)
					if index1 < 0 {
						listnobookable = append(listnobookable, listbookDocuments[i])
					}
					index := indexOf(listbookDocuments[i], listLocationFree)
					if index > -1 {
						listLocationFree = remove(listLocationFree, index)
					}
				}
			}
		}
	}
	return listLocationFree, nil
}

func indexOf(element interface{}, data []interface{}) int {
	for k, v := range data {
		if reflect.DeepEqual(element, v) {
			return k
		}
	}
	return -1
}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
