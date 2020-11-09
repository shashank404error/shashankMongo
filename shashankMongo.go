package shashankMongo

import (
	//"fmt"
	"context"
	//"os"
	//"strconv"
	log "github.com/sirupsen/logrus"
	//"github.com/bybrisk/structs"
    //"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)



/*var resultID string
var profileConfig *structs.ProfileConfig
var businessAccount *structs.BusinessAccount
var zones []structs.ZoneInfo
var zoneSingle *structs.ZoneInfo*/

var c *mongo.Client
var errors error
var CtxForDB context.Context
var DatabaseName *mongo.Database

func init(){
	c,errors= mongo.NewClient(options.Client().ApplyURI("mongodb://shashank404error:Y9ivXgMQ5ZrjL4N@parkpoint-shard-00-00.0bxqn.mongodb.net:27017,parkpoint-shard-00-01.0bxqn.mongodb.net:27017,parkpoint-shard-00-02.0bxqn.mongodb.net:27017/parkpoint?ssl=true&replicaSet=atlas-21pobg-shard-0&authSource=admin&retryWrites=true&w=majority"))
	if errors != nil {
		log.Error("Error Connecting Database Client")
		log.Error(errors)
	}
	CtxForDB = context.Background()
	errors = c.Connect(CtxForDB)
	if errors != nil {
		log.Error("error in setting up context")
		log.Error(errors)
	}

	DatabaseName = c.Database("parkpoint")
}

/*func InsertOne(connectionInfo *structs.ConnectToDataBase,collectionString string,customInsertStruct map[string]interface{}) string {
	
	collectionName := databaseName.Collection(collectionString)

	result, insertErr := collectionName.InsertOne(ctx, customInsertStruct)
	if insertErr != nil {
		log.Error("InsertOne ERROR:")
		log.Error(insertErr)
	} else {
	fmt.Println("InsertOne() API result:", result)

	newID := result.InsertedID
	fmt.Println("InsertOne() newID:", newID)
	resultID = newID.(primitive.ObjectID).Hex()
	}
	return resultID

}

func UpdateOneByID(connectionInfo *structs.ConnectToDataBase, collectionString string,docID string,insertKey string, insertValue string) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{insertKey: insertValue}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Error("Update One ERROR:")
			log.Error(err)
		}

	return res.ModifiedCount
}

func UpdateTwoByID(connectionInfo *structs.ConnectToDataBase, collectionString string,docID string,insertKey1 string, insertValue1 string,insertKey2 string, insertValue2 string) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{insertKey1: insertValue1,insertKey2: insertValue2}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Error("UpdateTwoByID ERROR:")
			log.Error(err)
		}

	return res.ModifiedCount
}

func FetchProfileConfiguration(connectionInfo *structs.ConnectToDataBase, collectionString string, filterValue string) *structs.ProfileConfig{

	collectionName := databaseName.Collection(collectionString)
	
	filter := bson.M{"plan": filterValue}
    err:= collectionName.FindOne(ctx, filter).Decode(&profileConfig)
	if err != nil {
		log.Error("FetchProfileConfiguration ERROR:")
		log.Error(err)
	}
    return profileConfig
}

func UpdateProfileConfiguration(connectionInfo *structs.ConnectToDataBase, collectionString string, docID string,config *structs.ProfileConfig) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{"profileConfig": config}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Error("UpdateProfileConfiguration ERROR:")
			log.Error(err)
		}

	fmt.Println("profile created")
	return res.ModifiedCount

}

func FetchProfile(connectionInfo *structs.ConnectToDataBase, collectionString string, docID string) *structs.BusinessAccount{

	collectionName := databaseName.Collection(collectionString)
	
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Error("FetchProfile ERROR:")
		log.Error(err)
	}
	businessAccount.UserID=docID
    return businessAccount
}

//FetchLogin is exported
func FetchLogin(connectionInfo *structs.ConnectToDataBase, collectionString string, username string, password string) (*structs.BusinessAccount, error){
	
	collectionName := databaseName.Collection(collectionString)
	
	filter := bson.M{"username": username,"password": password}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Error("FetchLogin ERROR:")
		log.Error(err)
		return nil, err
	}
	resultID = businessAccount.ID.Hex()
	businessAccount.UserID = resultID
	return businessAccount,err
}

//GetZone is exported
func GetZone(connectionInfo *structs.ConnectToDataBase, collectionString string, docID string) *structs.BusinessAccount{

	collectionName := databaseName.Collection(collectionString)

	cursor, err := collectionName.Find(ctx, bson.M{"businessUid":docID})
	if err != nil {
		log.Error("GetZone ERROR1:")
		log.Error(err)
	}
	if err = cursor.All(ctx, &zones); err != nil {
		log.Error("GetZone ERROR2:")
		log.Error(err)
	}

	for i,v:= range zones{
		zones[i].UserID = v.ID.Hex()
	}
	//fetch other account details
	account:=FetchProfile(connectionInfo,"businessAccounts",docID)
	account.ZoneDetailInfo=zones
    return account
}

func UpdateDeliveryInfo(connectionInfo *structs.ConnectToDataBase, collectionString string, docID string,deliveryStruct []structs.DeliveryDetail) int64 {

	collectionName := databaseName.Collection(collectionString)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$push": bson.M{"deliveryDetail": bson.M{"$each": deliveryStruct }}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Error("UpdateDeliveryInfo ERROR:")
			log.Error(err)
		}

	fmt.Println("Delivery Info assigned to "+docID)
	return res.ModifiedCount

}

func GetFieldByID (connectionInfo *structs.ConnectToDataBase, collectionString string, docID string) primitive.M {

	collectionName := databaseName.Collection(collectionString)

	var document bson.M
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}
	err:= collectionName.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		log.Error("GetFieldByID ERROR:")
		log.Error(err)
	}
	return document
}

func FetchZoneInfo (connectionInfo *structs.ConnectToDataBase, collectionString string , docID string , zoneID string) (*structs.ZoneInfo , string, error) {
	
	collectionName := databaseName.Collection(collectionString)

	filter := bson.M{"name": zoneID,"businessUid": docID}
    err:= collectionName.FindOne(ctx, filter).Decode(&zoneSingle)
	if err != nil {
		log.Error("FetchZoneInfo ERROR:")
		log.Error(err)
		return zoneSingle,"0",err
	}
	var index int
	for index, _ = range zoneSingle.DeliveryDetail {
	   index=index+1
	}
	indexString:=strconv.Itoa(index)
	return zoneSingle,indexString,nil	
}

func UpdateFieldInArray(connectionInfo *structs.ConnectToDataBase,collectionString string,fieldIdentifier string, filter1 string,filter2 string) int64 {
	
	collectionName := databaseName.Collection(collectionString)

	change := bson.M{"$pull": bson.M{"deliveryDetail": bson.M{ "customermob": fieldIdentifier}}}
	filter := bson.M{ "businessUid":filter1,"name":  filter2}
	res,err := collectionName.UpdateOne(ctx,filter, change)
	if err!=nil{
		log.Error("UpdateFieldInArray ERROR:")
		log.Error(err)
		return 0
	}
	fmt.Println("One order delivered to "+fieldIdentifier)
	return res.ModifiedCount
}	

func UpdateOneByFilters(connectionInfo *structs.ConnectToDataBase, collectionString string,filter1 string,filter2 string,insertKey string, insertValue string) int64 {

	collectionName := databaseName.Collection(collectionString)

	filter := bson.M{ "businessUid":filter1,"name":  filter2}

	update := bson.M{"$set": bson.M{insertKey: insertValue}}
	res,err := collectionName.UpdateOne(ctx,filter, update)
	if err!=nil{
		log.Error("UpdateOneByFilters ERROR:")
		log.Error(err)
	}

	return res.ModifiedCount
}

func FetchAndUpdateProfileDataByID(connectionInfo *structs.ConnectToDataBase, collectionString string,docID string ) int64 {

	businessAccount:=FetchProfile(connectionInfo, collectionString, docID)
	deliveryPendingInt, _ := strconv.ParseInt(businessAccount.DeliveryPending, 10, 64)
	newDeliveryPendingString:=	strconv.FormatInt((deliveryPendingInt-1), 10)
	deliveryDeliveredInt, _ := strconv.ParseInt(businessAccount.DeliveryDelivered, 10, 64)
	newDeliverydeliveredString:=	strconv.FormatInt((deliveryDeliveredInt+1), 10)

	res:=UpdateTwoByID(connectionInfo, collectionString,docID,"deliveryPending", newDeliveryPendingString,"deliveryDelivered", newDeliverydeliveredString)
	return res
}

func GetFieldByFilter (connectionInfo *structs.ConnectToDataBase, collectionString string, filterKey string, filterValue string) []primitive.M {

	collectionName := databaseName.Collection(collectionString)

	var documents []bson.M
	filter := bson.M{filterKey: filterValue}
	cursor,err:= collectionName.Find(ctx, filter)
	if err != nil {
		log.Error("GetFieldByFilter ERROR1:")
		log.Error(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
    	var document bson.M
    	if err = cursor.Decode(&document); err != nil {
			log.Error("GetFieldByFilter ERROR2:")
			log.Error(err)
		}
		documents = append(documents,document)
	}
	
	return documents
}
*/