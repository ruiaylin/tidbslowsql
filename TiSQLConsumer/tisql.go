package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	cp "tisql/components"
	"tisql/message"
	"tisql/model"
	tidbTools "tisql/utils"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

var (
	brokerList = flag.String("brokers", os.Getenv("KAFKA_PEERS_NEW1"), "The comma separated list of brokers in the Kafka cluster")
	topic      = flag.String("topic", "tidbslowquery", "REQUIRED: the topic to consume")
	partitions = flag.String("partitions", "all", "The partitions to consume, can be 'all' or comma-separated numbers")
	offset     = flag.String("offset", "newest", "The offset to start with. Can be `oldest`, `newest`")
	verbose    = flag.Bool("verbose", true, "Whether to turn on sarama logging")
	bufferSize = flag.Int("buffer-size", 256, "The buffer size of the message channel.")
	logger     = log.New(os.Stderr, "", log.LstdFlags)
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

// consumer
//  |
// \/
// chan1
// chan2
// chan3
// go run do check

func main() {

	cp.DBWriter.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.SlowSQL{})
	cp.DBWriter.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.SlowQuery{})
	cp.DBWriter.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.SlowQueryInfo{})

	flag.Parse()
	if *brokerList == "" {
		printUsageErrorAndExit("You have to provide -brokers as a comma-separated list, or set the KAFKA_PEERS environment variable.")
	}
	if *topic == "" {
		printUsageErrorAndExit("-topic is required")
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *verbose {
		sarama.Logger = logger
	}
	var initialOffset int64
	switch *offset {
	case "oldest":
		initialOffset = sarama.OffsetOldest
	case "newest":
		initialOffset = sarama.OffsetNewest
	default:
		printUsageErrorAndExit("-offset should be `oldest` or `newest`")
	}

	logger.Println("ruiaylin : Create consumer for slow query log")
	c, err := sarama.NewConsumer(strings.Split(*brokerList, ","), nil)
	if err != nil {
		printErrorAndExit(69, "Failed to start consumer: %s", err)
	}
	logger.Println("ruiaylin : Get partitionList through getPartitions() ")
	partitionList, err := getPartitions(c)
	if err != nil {
		printErrorAndExit(69, "Failed to get the list of partitions: %s", err)
	}
	var (
		messages = make(chan *sarama.ConsumerMessage, *bufferSize)
		closing  = make(chan struct{})
		wg       sync.WaitGroup
	)
	// what is the function
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		logger.Println("Initiating shutdown of consumer...")
		fmt.Printf("ruiaylin : Initiating shutdown of consumer...")
		close(closing)
	}()
	// create one group & one PartitionConsumer for one partition
	for _, partition := range partitionList {
		logger.Println("ruiaylin : Create Consumer for every Partition ")
		pc, err := c.ConsumePartition(*topic, partition, initialOffset)
		if err != nil {
			printErrorAndExit(69, "Failed to start consumer for partition %d: %s", partition, err)
		}
		go func(pc sarama.PartitionConsumer) {
			<-closing
			pc.AsyncClose()
		}(pc)
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for message := range pc.Messages() {
				messages <- message
			}
		}(pc)
	}
	// closure for processing  message in the messages channel
	// for chanel
	dbh := &tidbTools.DBHelper{}
	dbh.NewDBH()
	for i := 0; i < 4; i++ {
		go func() {
			for tMsg := range messages {
				// unmarshal data
				msg := &message.Message{}
				proto.Unmarshal(tMsg.Value, msg)
				fmt.Printf("Partition:\t%d\n", tMsg.Partition)
				fmt.Printf("Offset:\t%d\n", tMsg.Offset)
				fmt.Printf("Key:\t%s\n", string(tMsg.Key))
				fmt.Println()
				mQuery := tidbTools.InitTiDBQuery(msg)
				mQuery.SetFields(msg.Fields)
				mQuery.GenID()
				mQuery.FormatSQL()
				// mQuery.PrintInfo()
				dbh.HandleSQL(*mQuery)
				dbh.HandleQuery(*mQuery)
			}
		}()
	}
	wg.Wait()
	logger.Println("Done consuming topic", *topic)
	close(messages)

	if err := c.Close(); err != nil {
		logger.Println("Failed to close consumer: ", err)
	}
}

// get the partitions
func getPartitions(c sarama.Consumer) ([]int32, error) {
	if *partitions == "all" {
		return c.Partitions(*topic)
	}

	tmp := strings.Split(*partitions, ",")
	var pList []int32
	for i := range tmp {
		val, err := strconv.ParseInt(tmp[i], 10, 32)
		if err != nil {
			return nil, err
		}
		pList = append(pList, int32(val))
	}
	return pList, nil
}

func printErrorAndExit(code int, format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

func printUsageErrorAndExit(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Available command line options:")
	flag.PrintDefaults()
	os.Exit(64)
}
