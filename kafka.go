package main

import (
	// "fmt"

	"fmt"

	"github.com/Shopify/sarama"
)

func GetAdminClient() *sarama.ClusterAdmin {
	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	return &admin
}

func GetTopics(admin *sarama.ClusterAdmin) map[string]sarama.TopicDetail {
	topics, err := (*admin).ListTopics()
	if err != nil {
		panic(err)
	}

	return topics
}

func GetTopicsMetadata(admin *sarama.ClusterAdmin, topics map[string]sarama.TopicDetail) []*sarama.TopicMetadata {
	topicStrings := make([]string, 0, len(topics))
	for k := range topics {
		topicStrings = append(topicStrings, k)
	}

	metadata, err := (*admin).DescribeTopics(topicStrings)
	if err != nil {
		panic(err)
	}

	return metadata
}

func GetConsumersGroups(admin *sarama.ClusterAdmin) map[string]string {
	consumerGroups, err := (*admin).ListConsumerGroups()
	if err != nil {
		panic(err)
	}

	return consumerGroups
}

func GetConsumersGroupOffset(admin *sarama.ClusterAdmin, group string, topicsPartitions map[string][]int32) *sarama.OffsetFetchResponse {
	offset, err := (*admin).ListConsumerGroupOffsets(group, topicsPartitions)
	if err != nil {
		panic(err)
	}

	return offset
}

func GetConsumersGroupsDescription(admin *sarama.ClusterAdmin, groups map[string]string) []*sarama.GroupDescription {
	names := make([]string, 0, len(groups))
	for k := range groups {
		names = append(names, k)
	}

	desc, err := (*admin).DescribeConsumerGroups(names)
	if err != nil {
		panic(err)
	}

	return desc
}

func GetTopicsSize(admin *sarama.ClusterAdmin, topics []string) map[string]int {
	var broker int32 = 0

	desc, err := (*admin).DescribeLogDirs([]int32{broker})
	if err != nil {
		panic(err)
	}


	out := map[string]int{}
	for _, t := range desc[0][0].Topics {
		if includes[string](&topics, t.Topic) {
			size := 0

			for _,p := range t.Partitions {
				size += int(p.Size)
			}

			out[t.Topic] = size
		}
	}

	return out
}

// func main() {
// 	admin := GetAdminClient()
// 	desc := GetLogDirsDescriptions(admin)
// 	// var map1 *sarama.OffsetFetchResponseBlock = offset.Blocks["sp-gpcs-reservations-raw"][0]
// 	topic := desc[0][0].Topics[0].Partitions
// 	fmt.Println()
// }

func includes[T comparable](slice *[]T, value T) bool {
	for _, item := range *slice {
		if item == value {
			return true
		}
	}
	return false
}
