package org.coolsugar.log.util;

import kafka.consumer.ConsumerConfig;

import java.util.Properties;

public class KafkaConf {

    private static KafkaConf ourInstance = new KafkaConf();
    private ConsumerConfig conf = null;

    private KafkaConf() {
        try {
            Properties props = new Properties();
            props.put("zookeeper.connect", ConfigUtil.getInstance().getProperty("zookeeperHost"));
            props.put("group.id", ConfigUtil.getInstance().getProperty("kafkaGroup"));

            props.put("zookeeper.session.timeout.ms", "10000");
            props.put("zookeeper.sync.time.ms", "200");
            props.put("auto.commit.interval.ms", "1000");
            props.put("auto.offset.reset", "smallest");
            props.put("format", "binary");
            props.put("auto.commit.enable", "true");
            props.put("serializer.class", "kafka.serializer.StringEncoder");

            conf = new ConsumerConfig(props);
        } catch (Exception e) {
            System.out.println(e);
        }
    }

    public static KafkaConf getInstance() {
        return ourInstance;
    }

    public ConsumerConfig getConf () {
        return conf;
    }
}