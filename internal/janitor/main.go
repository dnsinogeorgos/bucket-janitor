package janitor

import (
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dnsinogeorgos/bucket-janitor/internal/aws"
	"github.com/olekukonko/tablewriter"
)

type token struct{}

type Janitor struct {
	wg          *sync.WaitGroup
	client      *s3.Client
	downloader  *manager.Downloader
	buckets     []string
	objectChan  chan Object
	objectSem   chan token
	headerChan  chan Header
	headerSem   chan token
	typeSetChan chan TypeSet
	counter     *sync.Map
}

func New(configPath string) *Janitor {
	config := aws.CreateLoad(configPath)
	var wg sync.WaitGroup

	j := &Janitor{
		&wg,
		config.Client,
		config.Downloader,
		config.S3Buckets,
		make(chan Object),
		make(chan token, config.Concurrency),
		make(chan Header),
		make(chan token, config.Concurrency),
		make(chan TypeSet),
		&sync.Map{},
	}

	return j
}

func (j *Janitor) Run() {
	j.wg.Add(4)
	go j.listObjects()
	go j.retrieveHeaders()
	go j.detectTypes()
	go j.countTypes()

	j.wg.Wait()

	results := prepareResultSet(j.counter)
	printTable(results)
}

func printTable(rs resultSet) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Mime", "Magic", "Count"})

	for _, r := range rs {
		table.Append([]string{r.Mime, r.Magic, strconv.FormatUint(*r.Occurences, 10)})
	}

	table.Render()
}
