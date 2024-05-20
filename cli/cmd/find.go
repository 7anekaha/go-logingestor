package cmd

import (
	"context"
	"fmt"
	// "log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// findCmd represents the find command
var (

	//flags
	filters      Filters
	regexPattern string
	findCmd      = &cobra.Command{
		Use:   "find",
		Short: "find logs in the database",
		Long:  `Find logs in the database based on the filters provided or the regex pattern.`,
		Run:   FindLogs,
	}
)

func init() {
	filters = make(Filters)

	findCmd.Flags().VarP(&filters, "filter", "f", "Filters to apply to the logs (e.g. level=info, resourceId=1234, timestamp=<2024-01-01 00:00:00)")
	findCmd.Flags().StringVarP(&regexPattern, "regex", "r", "", "Regex pattern to search in the entire document (e.g. '.*error.*')")

	rootCmd.AddCommand(findCmd)

}

const (
	user       = "nico"
	pass       = "secret"
	hostName   = "localhost"
	port       = "27017"
	DB         = "logs-db"
	COLLECTION = "logs-collection"
)

type Log struct {
	Level      string    `json:"level" bson:"level"`
	Message    string    `json:"message" bson:"message"`
	ResourceId string    `json:"resourceId" bson:"resourceId"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
	TraceId    string    `json:"traceId" bson:"traceId"`
	SpanId     string    `json:"spanId" bson:"spanId"`
	Commit     string    `json:"commit" bson:"commit"`
	Metadata   Metadata  `json:"metadata" bson:"metadata"`
}

type Metadata struct {
	ParentResourceId string `json:"parentResourceId" bson:"parentResourceId"`
}

type Filters map[string]string

func (f *Filters) String() string {
	var result []string
	for key, value := range *f {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(result, ", ")
}

func (f *Filters) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid filter format, expected key value")
	}
	(*f)[parts[0]] = parts[1]
	return nil
}

func (f *Filters) Type() string {
	return "Filters"
}

func FindLogs(cmd *cobra.Command, args []string) {

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, hostName, port))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("error connecting to mongodb")
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic("error disconnecting from mongodb")
		}
	}()

	// Get all logs
	collection := client.Database(DB).Collection(COLLECTION)
	filter := bson.M{}

	// Apply regular expression filter if provided
	if regexPattern != "" {
		regexQuery := bson.M{"$regex": regexPattern, "$options": "i"} // "i" for case-insensitive matching
		filter["$or"] = []bson.M{
			{"level": regexQuery},
			{"message": regexQuery},
			{"resourceId": regexQuery},
			{"traceId": regexQuery},
			{"spanId": regexQuery},
			{"commit": regexQuery},
			{"metadata.parentResourceId": regexQuery},
		}
	}

	for key, value := range filters {

		if key == "timestamp" {
			signal := value[0]
			timeValue, err := time.Parse("2006-01-02 15:04:05", value[1:])
			if err != nil {
				panic("error parsing timestamp value. Please use the format: <,>,yyyy-mm-dd hh:mm:ss")
			}
			if signal == '>' {
				filter["timestamp"] = bson.M{"$gt": timeValue}
			} else if signal == '<' {
				filter["timestamp"] = bson.M{"$lt": timeValue}
			} else {
				panic("error parsing timestamp value. Please use the format: <,>,yyyy-mm-dd hh:mm:ss")
			}
		} else {
			filter[key] = value
		}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic("error finding logs")
	}
	defer cursor.Close(context.Background())

	var logs []Log
	if err := cursor.All(context.Background(), &logs); err != nil {
		panic("error decoding logs")
	}

	// Set table formatting
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Time", "Level", "Message", "Resource ID", "Trace ID", "Span ID", "Commit", "Parent Resource ID")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, currLog := range logs {
		tbl.AddRow(currLog.Timestamp, currLog.Level, currLog.Message, currLog.ResourceId, currLog.TraceId, currLog.SpanId, currLog.Commit, currLog.Metadata.ParentResourceId)
	}

	tbl.Print()
}
