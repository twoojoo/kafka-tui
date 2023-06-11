package main

import (
	// "fmt"

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



// func main() {
// 	admin := GetAdminClient()
// 	offset := GetConsumersGroupOffset(admin, "sp-reservation-parser", map[string][]int32{"sp-gpcs-reservations-raw": []int32{0}})
// 	var map1 *sarama.OffsetFetchResponseBlock = offset.Blocks["sp-gpcs-reservations-raw"][0]
// 	fmt.Println(map1.Offset)
// }