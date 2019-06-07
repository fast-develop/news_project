package org.coolsugar;

import java.util.*;
import java.util.concurrent.ExecutionException;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.security.JaasUtils;
import org.coolsugar.log.util.ConfigUtil;
import org.coolsugar.log.util.KafkaProperties;
import org.junit.Test;

import kafka.admin.AdminUtils;
import kafka.admin.RackAwareMode;
import kafka.utils.ZkUtils;

public class TestKafka {

    @Test
    public void create_topic(){
        ZkUtils zkUtils = ZkUtils.apply("152.136.128.83:2181", 30000, 30000, JaasUtils.isZkSecurityEnabled());
        AdminUtils.createTopic(zkUtils, "enbook", 1, 1, new Properties(), RackAwareMode.Enforced$.MODULE$);
        zkUtils.close();
    }


    @Test
    public void delete_topic(){
        ZkUtils zkUtils = ZkUtils.apply("152.136.128.83:2181", 30000, 30000, JaasUtils.isZkSecurityEnabled());
        AdminUtils.deleteTopic(zkUtils, "cnbook");
        zkUtils.close();
    }

    @Test
    public void producer() throws InterruptedException, ExecutionException{
        Properties props=new Properties();
        props.put("key.serializer", "org.apache.kafka.common.serialization.IntegerSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("bootstrap.servers","152.136.128.83:9092");

        Producer<Integer, String> kafkaProducer = new KafkaProducer<Integer, String>(props);
        for(int i=100;i<200;i++){
//            ProducerRecord<Integer, String> message = new ProducerRecord<Integer, String>("enbook","" + i);
            ProducerRecord<Integer, String> message = new ProducerRecord<Integer, String>("dnbook","" + i);
            kafkaProducer.send(message);
        }

        kafkaProducer.close();
    }

    @Test
    public void consumer () {
        /* 定义consumer */
        KafkaConsumer<String, String> consumer =
                new KafkaConsumer<>(KafkaProperties.getInstance().getProperties());

        /* 消费者订阅的topic, 可同时订阅多个 */
//        consumer.subscribe(Arrays.asList("enbook", "cnbook"));
        consumer.subscribe(Arrays.asList(ConfigUtil.getInstance().getProperty("kafkaTopic")));

        /* 读取数据，读取超时时间为100ms */
        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(100);
            for (ConsumerRecord<String, String> record : records)
                System.out.printf("########### offset = %d, key = %s, value = %s\n", record.offset(), record.key(), record.value());
        }

    }

}
