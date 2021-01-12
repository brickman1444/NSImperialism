package dynamodbwrapper

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const CELL_TABLE_NAME string = "nsimperialism-cell"

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

func (cell DatabaseCell) ToString() string {
	return fmt.Sprintf("ID: %v Resident: %v", cell.ID, cell.Resident)
}

func initializeDatabase() {
	databaseContext := context.TODO()

	awsConfig, err := config.LoadDefaultConfig(databaseContext)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig)

	getItemsResponse, err := dynamodbClient.Scan(databaseContext, &dynamodb.ScanInput{
		TableName: aws.String(CELL_TABLE_NAME),
	})

	if err != nil {
		log.Fatalf("failed to get items, %v", err)
	}

	records := []DatabaseCell{}
	err = attributevalue.UnmarshalListOfMaps(getItemsResponse.Items, &records)
	if err != nil {
		log.Println("failed to unmarshal Items, %w", err)
	}
	for _, record := range records {
		log.Println(record.ToString())
	}

	itemToPut := DatabaseCell{
		ID:       "A",
		Resident: "testlandia",
	}
	itemToPutMap, err := attributevalue.MarshalMap(itemToPut)
	if err != nil {
		log.Fatalf("failed to marshal item, %v", err)
	}

	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
		Item:      itemToPutMap,
	})
	if err != nil {
		log.Fatalf("failed to put item, %v", err)
	}

	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: itemToPut.ID,
			},
		},
	})

	updatedItem := DatabaseCell{}
	err = attributevalue.UnmarshalMap(getItemOutput.Item, &updatedItem)
	if err != nil {
		log.Fatalf("failed to unmarshal item, %v", err)
	}

	log.Println(updatedItem.ToString())
}

func PutCell(cell DatabaseCell) error {

	itemToPutMap, err := attributevalue.MarshalMap(cell)
	if err != nil {
		return err
	}

	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
		Item:      itemToPutMap,
	})
	return err
}

func GetCell(territoryName string) (DatabaseCell, error) {
	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
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
