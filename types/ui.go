package types

import (
	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/twoojoo/ktui/components"
)

type UI struct {
	//graphic stuff
	App            *tview.Application
	Theme          *Theme
	MainContainer  *tview.Flex
	MainView       *tview.Flex
	SidePane       *tview.List
	CentralView    *tview.Flex
	BrokersTable   *components.SearchableTable
	ConsumersTable *components.SearchableTable
	TopicsTable    *components.SearchableTable
	TopicDetail    *tview.Flex
	IsModalOpen    bool
	Modal          *tview.Flex

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
