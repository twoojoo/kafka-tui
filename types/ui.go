package types

import (
	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/utils"
)

type UI struct {
	//graphic stuff
	Theme          *Theme
	App            *tview.Application
	SidePane       *tview.List
	View           *tview.Flex
	BrokersTable   *utils.SearchableTable
	ConsumersTable *utils.SearchableTable
	TopicsTable    *utils.SearchableTable
	TopicDetail    *tview.Flex

	UpdateFunc func(*UI)

	//kafka stuff
	AdminClient                *sarama.ClusterAdmin
	Brokers                    []*sarama.Broker
	ControllerId               int32
	Topics                     map[string]sarama.TopicDetail
	TopicsMetadata             []*sarama.TopicMetadata
	TopicsSize                 map[string]int
	ConsumerGroups             map[string]string
	ConsumerGroupsOffsets      map[string]*sarama.OffsetFetchResponse
	ConsumerGroupsDescriptions []*sarama.GroupDescription
}

type Theme struct {
	Background      tcell.Color
	Foreground      tcell.Color
	PrimaryColor    tcell.Color
	InEvidenceColor tcell.Color
	InfoColor       tcell.Color
	ErrorColor      tcell.Color
	WarningColor    tcell.Color
}