package org.coolsugar.log.util;

import org.apache.kafka.common.serialization.StringDeserializer;

import java.util.Properties;

public class KafkaProperties {

    private static KafkaProperties ourInstance = new KafkaProperties();
    private Properties properties = null;

    private KafkaProperties() {
        try {
            properties = new Properties();
            properties.put("zookeeper.connect", ConfigUtil.getInstance().getProperty("zookeeperHost"));
            /* 新版 kafka 要配此地址，定义kakfa 服务的地址，不需要将所有broker指定上 */
            properties.put("bootstrap.servers", "152.136.128.83:9092");
            /* 制定consumer group */
            properties.put("group.id", ConfigUtil.getInstance().getProperty("kafkaGroup"));
            /* 在当前没有 producer 时，新的 consumer group 可以通过打开此开关来从 offset 0 获取数据 */
            properties.put("auto.offset.reset", "earliest");
            /* 是否自动确认offset */
            properties.put("enable.auto.commit", "true");
            /* 自动确认offset的时间间隔,默认值为5000ms */
            properties.put("auto.commit.interval.ms", "1000");
            properties.put("session.timeout.ms", "30000");
            /* key的序列化类 */
            properties.put("key.deserializer", StringDeserializer.class.getName());
            /* value的序列化类 */
            properties.put("value.deserializer", StringDeserializer.class.getName());

        } catch (Exception e) {
            System.out.println(e);
        }
    }

    public static KafkaProperties getInstance() {
        return ourInstance;
    }

    public Properties getProperties () {
        return properties;
    }
}
