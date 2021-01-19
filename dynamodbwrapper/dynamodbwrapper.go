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
)

var CellDoesntExistError = errors.New("Cell doesn't exist")

var dynamodbClient *dynamodb.Client = nil
var databaseContext = context.TODO()

func Initialize() {
	awsConfig, err := config.LoadDefaultConfig(databaseContext)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	dynamodbClient = dynamodb.NewFromConfig(awsConfig)
}

type DatabaseCell struct {
	ID       string
	Resident string
}

func cellTableName() string {
	environmentVariableValue, doesEnvironmentVariableExist := os.LookupEnv("CELL_TABLE_NAME")
	if doesEnvironmentVariableExist {
		return environmentVariableValue
	}

	return "nsimperialism-cell"
}

func PutCell(cell DatabaseCell) error {

	itemToPutMap, err := attributevalue.MarshalMap(cell)
	if err != nil {
		return err
	}

	log.Println("DynamoDB: Put on cell table")
	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(cellTableName()),
		Item:      itemToPutMap,
	})
	return err
}

func GetCell(territoryName string) (DatabaseCell, error) {

	log.Println("DynamoDB: Get on cell table")
	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(cellTableName()),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: territoryName,
			},
		},
	})
	if err != nil {
		return DatabaseCell{}, err
	}

	if len(getItemOutput.Item) == 0 {
		return DatabaseCell{}, CellDoesntExistError
	}

	gotItem := DatabaseCell{}
	err = attributevalue.UnmarshalMap(getItemOutput.Item, &gotItem)
	if err != nil {
		return DatabaseCell{}, err
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
	environmentVariableValue, doesEnvironmentVariableExist := os.LookupEnv("WAR_TABLE_NAME")
	if doesEnvironmentVariableExist {
		return environmentVariableValue
	}

	return "nsimperialism-war"
}

func PutWars(wars []DatabaseWar) error {

	for _, warToAdd := range wars {
		itemToPutMap, err := attributevalue.MarshalMap(warToAdd)
		if err != nil {
			return err
		}

		log.Println("DynamoDB: Put on war table")
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

	log.Println("DynamoDB: Scan on war table")
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
