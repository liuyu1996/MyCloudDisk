package mq

import (
	"MyCloudDisk/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)


var done chan bool

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
}


//创建RabbitMQ结构体实例
func NewRabbitMQ() *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: config.Configs.TransOSSQueueName,
		Exchange:config.Configs.TransExchangeName,
		Key:config.Configs.TransOSSRoutingKey,
		Mqurl:config.Configs.RabbitURL}
	var err error
	//创建RabbitMQ连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "获取channel失败")

	return rabbitmq
}


func (r *RabbitMQ) StartConsume(cName string, callback func(msg[]byte) bool)  {

	//1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	//保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
		//是否自动删除
		false,
		//是否具有排他
		false,
		//是否阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//1.通过channel.consume获得消息通道
	msgs, err := r.channel.Consume(r.QueueName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//2.循环获取消息的队列
	//3.调用callback方法处理新的消息
	done = make(chan bool)
	go func() {
		for msg:= range msgs{
			ok := callback(msg.Body)
			if !ok {
				//TODO:将任务写到另一个队列， 用于异常情况重试
			}
		}
	}()
	<-done
}



func (r *RabbitMQ) Publish(msg[]byte) bool {
	//1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	r.Lock()
	defer r.Unlock()
	//保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
		//是否自动删除
		false,
		//是否具有排他
		false,
		//是否阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		return false
	}
	//2.执行消息发布动作
	err = r.channel.Publish(r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body:msg,
		})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

//断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s",message, err)
		panic(fmt.Sprint("%s:%s",message, err))
	}
}
