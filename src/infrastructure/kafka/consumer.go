package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"

	"github.com/vostelmakh/mixturka/src/application/processor/recipe"
)

const TopicBabushkaRecipeV1 = "babushka.recipes.v1"

type Consumer struct {
	consumer  sarama.Consumer
	processor *recipe.Processor
	topic     string
}

func NewConsumer(brokers []string, topic string, processor *recipe.Processor) (*Consumer, error) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  consumer,
		processor: processor,
		topic:     topic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				if err := c.processor.ProcessRecipe(ctx, msg.Value); err != nil {
					log.Printf("Error processing message: %v", err)
				}
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming message: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
