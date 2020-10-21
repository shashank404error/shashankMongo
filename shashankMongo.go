package shashankMongo

import (
	"fmt"
	"context"
	"os"
	"strconv"
	"log"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectToDataBase struct {
	CustomApplyURI string 
	DatabaseName string 
}

type ProfileConfig struct{
	Zone int64 `bson: "zone" json: "zone"`
	MessagePlan int64 `bson: "messageplan" json: "messageplan"`
	Tracking bool `bson: "tracking" json: "tracking"`
	ZoneID []string `bson: "zoneid" json: "zoneid"`
}

type BusinessAccount struct{
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PicURL string `bson: "picurl" json: "picurl"`
	UserName string `bson: "username" json: "username"`
	BusinessName string `bson: "businessname" json: "businessname"`
	Password string `bson: "password" json: "password"`
	City string `bson: "city" json: "city"`
	BusinessPlan string `bson: "businessplan" json: "businessplan"`
	ProfileConfig ProfileConfig `bson: "profileConfig" json: "profileConfig"`
	DeliveryPending string `bson: "deliveryPending" json: "deliveryPending"`
	DeliveryDelivered string `bson: "deliveryDelivered" json: "deliveryDelivered"`
	UserID string
	ZoneDetailInfo []ZoneInfo
}

type ZoneInfo struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PicURL string `bson: "picurl" json: "picurl"`
	Name string `bson: "name" json: "name"`
	BusinessUID string `bson: "businessUid" json: "businessUid"`
	DeliveryInZone string `bson: "deliveryInZone" json: "deliveryInZone"`
	UserID string
	DeliveryDetail []DeliveryDetail `bson: "deliveryDetail" json: "deliveryDetail"`
	Longitude string `bson: "longitude" json: "longitude"`
	Latitude string `bson: "latitude" json: "latitude"`
	Error string
}

type DeliveryDetail struct {
	CustomerName string `bson: "cusotmername" json: "cusotmername"`
	DistanceFromYou float64 
	ETA float64
	CustomerMob string `bson: "cusotmermob" json: "cusotmermob"`
	PicURL string `bson: "picurl" json: "picurl"`
	Address string `bson: "address" json: "address"`
	Latitude float64 `bson: "latitude" json: "latitude"`
	Longitude float64 `bson: "longitude" json: "longitude"`
	LongLat string `bson: "longlat" json: "longlat"`
	DeliveryInZone int64
	BusinessId string
} 

type MapBoxResp struct {
	Code   string `json:"code"`
	Routes []RoutesResp
}

type RoutesResp struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type DeliveredDetail struct {
	Longitude string `bson:"longitude" json:"longitude"`
	Latitude string `bson:"latitude" json:"latitude"`
	CustomerName string `bson:"customerName" json:"customerName"`
	CustomerMob string `bson:"customerMob" json:"customerMob"`
	ZoneID string `bson:"zoneId" json:"zoneId"`
	DateOfDelivery string `bson:"dateOfDelivery" json:"dateOfDelivery"`
	Address string `bson:"address" json:"address"`
	BusinessUid string `bson:"businessUid" json:"businessUid"`
}

type DeliveredAndAccount struct {
	DeliveredDetails []DeliveredDetail
	BusinessAccount *BusinessAccount
}

var resultID string
var profileConfig *ProfileConfig
var businessAccount *BusinessAccount
var zones []ZoneInfo
var zoneSingle *ZoneInfo

var c *mongo.Client
var errors error
var ctx context.Context
var databaseName *mongo.Database

func init(){
	c,errors= mongo.NewClient(options.Client().ApplyURI("mongodb://shashank404error:Y9ivXgMQ5ZrjL4N@parkpoint-shard-00-00.0bxqn.mongodb.net:27017,parkpoint-shard-00-01.0bxqn.mongodb.net:27017,parkpoint-shard-00-02.0bxqn.mongodb.net:27017/parkpoint?ssl=true&replicaSet=atlas-21pobg-shard-0&authSource=admin&retryWrites=true&w=majority"))
	if errors != nil {
		fmt.Println("error in client")
		log.Fatal(errors)
	}
	ctx = context.Background()
	errors = c.Connect(ctx)
	if errors != nil {
		fmt.Println("error in context")
		log.Fatal(errors)
	}

	databaseName = c.Database("parkpoint")
}

func InsertOne(connectionInfo *ConnectToDataBase,collectionString string,customInsertStruct map[string]interface{}) string {
	
	collectionName := databaseName.Collection(collectionString)

	result, insertErr := collectionName.InsertOne(ctx, customInsertStruct)
	if insertErr != nil {
	fmt.Println("InsertOne ERROR:", insertErr)
	os.Exit(1) // safely exit script on error
	} else {
	fmt.Println("InsertOne() API result:", result)

	newID := result.InsertedID
	fmt.Println("InsertOne() newID:", newID)
	resultID = newID.(primitive.ObjectID).Hex()
	}
	return resultID

}

func UpdateOneByID(connectionInfo *ConnectToDataBase, collectionString string,docID string,insertKey string, insertValue string) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{insertKey: insertValue}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	return res.ModifiedCount
}

func UpdateTwoByID(connectionInfo *ConnectToDataBase, collectionString string,docID string,insertKey1 string, insertValue1 string,insertKey2 string, insertValue2 string) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{insertKey1: insertValue1,insertKey2: insertValue2}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	return res.ModifiedCount
}

func FetchProfileConfiguration(connectionInfo *ConnectToDataBase, collectionString string, filterValue string) *ProfileConfig{

	collectionName := databaseName.Collection(collectionString)
	
	filter := bson.M{"plan": filterValue}
    err:= collectionName.FindOne(ctx, filter).Decode(&profileConfig)
	if err != nil {
		log.Fatal(err)
	}
    return profileConfig
}

func UpdateProfileConfiguration(connectionInfo *ConnectToDataBase, collectionString string, docID string,config *ProfileConfig) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{"profileConfig": config}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	fmt.Println("profile created")
	return res.ModifiedCount

}

func FetchProfile(connectionInfo *ConnectToDataBase, collectionString string, docID string) *BusinessAccount{

	collectionName := databaseName.Collection(collectionString)
	
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Fatal(err)
	}
	businessAccount.UserID=docID
    return businessAccount
}

//FetchLogin is exported
func FetchLogin(connectionInfo *ConnectToDataBase, collectionString string, username string, password string) (*BusinessAccount, error){
	
	collectionName := databaseName.Collection(collectionString)
	
	filter := bson.M{"username": username,"password": password}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resultID = businessAccount.ID.Hex()
	businessAccount.UserID = resultID
	return businessAccount,err
}

//GetZone is exported
func GetZone(connectionInfo *ConnectToDataBase, collectionString string, docID string) *BusinessAccount{

	collectionName := databaseName.Collection(collectionString)

	cursor, err := collectionName.Find(ctx, bson.M{"businessUid":docID})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &zones); err != nil {
		log.Fatal(err)
	}

	for i,v:= range zones{
		zones[i].UserID = v.ID.Hex()
	}
	//fetch other account details
	account:=FetchProfile(connectionInfo,"businessAccounts",docID)
	account.ZoneDetailInfo=zones
    return account
}

func UpdateDeliveryInfo(connectionInfo *ConnectToDataBase, collectionString string, docID string,deliveryStruct []DeliveryDetail) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$push": bson.M{"deliveryDetail": bson.M{"$each": deliveryStruct }}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	fmt.Println("Delivery Info assigned to "+docID)
	return res.ModifiedCount

}

func GetFieldByID (connectionInfo *ConnectToDataBase, collectionString string, docID string) primitive.M {

	collectionName := databaseName.Collection(collectionString)

	var document bson.M
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}
	err:= collectionName.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		log.Fatal(err)
	}
	return document
}

func FetchZoneInfo (connectionInfo *ConnectToDataBase, collectionString string , docID string , zoneID string) (*ZoneInfo , string, error) {
	
	collectionName := databaseName.Collection(collectionString)

	filter := bson.M{"name": zoneID,"businessUid": docID}
    err:= collectionName.FindOne(ctx, filter).Decode(&zoneSingle)
	if err != nil {
		log.Println(err)
		return zoneSingle,"0",err
	}
	var index int
	for index, _ = range zoneSingle.DeliveryDetail {
	   index=index+1
	}
	indexString:=strconv.Itoa(index)
	return zoneSingle,indexString,nil	
}

func UpdateFieldInArray(connectionInfo *ConnectToDataBase,collectionString string,fieldIdentifier string, filter1 string,filter2 string) int64 {
	
	collectionName := databaseName.Collection(collectionString)

	change := bson.M{"$pull": bson.M{"deliveryDetail": bson.M{ "customermob": fieldIdentifier}}}
	filter := bson.M{ "businessUid":filter1,"name":  filter2}
	res,err := collectionName.UpdateOne(ctx,filter, change)
	if err!=nil{
		fmt.Println(err)
		return 0
	}
	fmt.Println("One order delivered to "+fieldIdentifier)
	return res.ModifiedCount
}	

func UpdateOneByFilters(connectionInfo *ConnectToDataBase, collectionString string,filter1 string,filter2 string,insertKey string, insertValue string) int64 {

	collectionName := databaseName.Collection(collectionString)

	filter := bson.M{ "businessUid":filter1,"name":  filter2}

	update := bson.M{"$set": bson.M{insertKey: insertValue}}
	res,err := collectionName.UpdateOne(ctx,filter, update)
	if err!=nil{
		log.Fatal(err)
	}

	return res.ModifiedCount
}

func FetchAndUpdateProfileDataByID(connectionInfo *ConnectToDataBase, collectionString string,docID string ) int64 {

	businessAccount:=FetchProfile(connectionInfo, collectionString, docID)
	deliveryPendingInt, _ := strconv.ParseInt(businessAccount.DeliveryPending, 10, 64)
	newDeliveryPendingString:=	strconv.FormatInt((deliveryPendingInt-1), 10)
	deliveryDeliveredInt, _ := strconv.ParseInt(businessAccount.DeliveryDelivered, 10, 64)
	newDeliverydeliveredString:=	strconv.FormatInt((deliveryDeliveredInt+1), 10)

	res:=UpdateTwoByID(connectionInfo, collectionString,docID,"deliveryPending", newDeliveryPendingString,"deliveryDelivered", newDeliverydeliveredString)
	return res
}

func GetFieldByFilter (connectionInfo *ConnectToDataBase, collectionString string, filterKey string, filterValue string) []primitive.M {

	collectionName := databaseName.Collection(collectionString)

	var documents []bson.M
	filter := bson.M{filterKey: filterValue}
	cursor,err:= collectionName.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
    	var document bson.M
    	if err = cursor.Decode(&document); err != nil {
       	 log.Fatal(err)
		}
		documents = append(documents,document)
	}
	
	return documents
}