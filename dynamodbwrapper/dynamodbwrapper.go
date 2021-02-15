package dynamodbwrapper

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brickman1444/NSImperialism/databasemap"
)

var CellDoesntExistError = errors.New("Cell doesn't exist")
var MapDoesntExistError = errors.New("Map doesn't exist")

var dynamodbClient *dynamodb.Client = nil
var databaseContext = context.TODO()

func Initialize() {
	awsConfig, err := config.LoadDefaultConfig(databaseContext)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	dynamodbClient = dynamodb.NewFromConfig(awsConfig)
}

func getTableName(enviromentVariable string, productionTableName string) string {
	environmentVariableValue, doesEnvironmentVariableExist := os.LookupEnv(enviromentVariable)
	if doesEnvironmentVariableExist {
		return environmentVariableValue
	}

	return productionTableName
}

func cellTableName() string {
	return getTableName("CELL_TABLE_NAME", "nsimperialism-cell")
}

func PutCell(cell databasemap.DatabaseCell) error {

	itemToPutMap, err := attributevalue.MarshalMap(cell)
	if err != nil {
		return err
	}

	log.Println("DynamoDB: Put on cell table:", cell.ID)
	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(cellTableName()),
		Item:      itemToPutMap,
	})
	return err
}

func GetCell(territoryName string) (databasemap.DatabaseCell, error) {

	log.Println("DynamoDB: Get on cell table:", territoryName)
	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(cellTableName()),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: territoryName,
			},
		},
	})
	if err != nil {
		return databasemap.DatabaseCell{}, err
	}

	if len(getItemOutput.Item) == 0 {
		return databasemap.DatabaseCell{}, CellDoesntExistError
	}

	gotItem := databasemap.DatabaseCell{}
	err = attributevalue.UnmarshalMap(getItemOutput.Item, &gotItem)
	if err != nil {
		return databasemap.DatabaseCell{}, err
	}

	return gotItem, nil
}

type DatabaseWar struct {
	Attacker      string
	Defender      string
	Score         int
	ID            string
	TerritoryName string
	IsOngoing     bool
}

func warTableName() string {
	return getTableName("WAR_TABLE_NAME", "nsimperialism-war")
}

func PutWars(wars []DatabaseWar) error {

	for _, warToAdd := range wars {
		itemToPutMap, err := attributevalue.MarshalMap(warToAdd)
		if err != nil {
			return err
		}

		log.Println("DynamoDB: Put on war table:", warToAdd.ID)
		_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
			TableName: aws.String(warTableName()),
			Item:      itemToPutMap,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetWars() ([]DatabaseWar, error) {

	log.Println("DynamoDB: Scan on all of war table")
	scanOutput, err := dynamodbClient.Scan(databaseContext, &dynamodb.ScanInput{
		TableName:      aws.String(warTableName()),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	warsToReturn := []DatabaseWar{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &warsToReturn)
	if err != nil {
		return nil, err
	}

	return warsToReturn, nil
}

func mapTableName() string {
	return getTableName("MAP_TABLE_NAME", "nsimperialism-map")
}

func GetMap(ID string) (databasemap.DatabaseMap, error) {
	log.Println("DynamoDB: Get on map table")
	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(mapTableName()),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: ID,
			},
		},
	})

	if err != nil {
		return databasemap.NewBlankDatabaseMap(), err
	}

	if len(getItemOutput.Item) == 0 {
		return databasemap.NewBlankDatabaseMap(), MapDoesntExistError
	}

	gotItem := databasemap.NewBlankDatabaseMap()
	err = attributevalue.UnmarshalMap(getItemOutput.Item, &gotItem)
	if err != nil {
		return databasemap.NewBlankDatabaseMap(), err
	}

	return gotItem, nil
}

func PutMap(item databasemap.DatabaseMap) error {

	itemToPutMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	log.Println("DynamoDB: Put on map table")
	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(mapTableName()),
		Item:      itemToPutMap,
	})
	return err
}

func GetAllMapIDs() ([]string, error) {

	log.Println("DynamoDB: Scan on all of map table")
	scanOutput, err := dynamodbClient.Scan(databaseContext, &dynamodb.ScanInput{
		TableName:      aws.String(mapTableName()),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	maps := []databasemap.DatabaseMap{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &maps)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	for _, databaseMap := range maps {
		ids = append(ids, databaseMap.ID)
	}

	return ids, nil
}
