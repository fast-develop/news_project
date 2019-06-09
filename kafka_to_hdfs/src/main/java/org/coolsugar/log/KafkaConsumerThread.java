package org.coolsugar.log;


import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IOUtils;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.coolsugar.log.util.ConfigUtil;
import org.coolsugar.log.util.HdfsConf;
import org.coolsugar.log.util.KafkaProperties;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.text.SimpleDateFormat;
import java.util.*;


public class KafkaConsumerThread extends Thread {

    private static FileSystem hadoopFS = null;

    private String topic = "";


    public void setTopic (String _topic) {
        topic = _topic;
    }

    public void run() {

        /* 定义consumer */
        KafkaConsumer<String, String> consumer =
                new KafkaConsumer<String, String>(KafkaProperties.getInstance().getProperties());

        /* 消费者订阅的topic, 可同时订阅多个 */
//        consumer.subscribe(Arrays.asList("enbook", "cnbook"));

        consumer.subscribe(Collections.singletonList(topic));

        try {
            hadoopFS = FileSystem.get(HdfsConf.getInstance().getConf());

            String topicPathName = "/" + topic;
            Path topicPath = new Path(topicPathName);
            if (!hadoopFS.exists(topicPath)) {
                hadoopFS.mkdirs(topicPath);
            }

            String fileName = topicPathName + "/" +
                    (new SimpleDateFormat("yyyy-MM-dd").format(Calendar.getInstance().getTime())) + ".txt";
            Path logFilePath = new Path(fileName);
            if (!hadoopFS.exists(logFilePath)) {
                FSDataOutputStream output = hadoopFS.create(logFilePath);
                output.close();
            }

            /* 读取数据，读取超时时间为100ms */
            while (true) {
                ConsumerRecords<String, String> records = consumer.poll(1000);
                for (ConsumerRecord<String, String> record : records) {
                    System.out.printf("---------- topic:" + topic + " offset = %d, key = %s, value = %s\n", record.offset(), record.key(), record.value());

                    String lineData = record.value() + "\n";
                    InputStream in = new ByteArrayInputStream(lineData.getBytes("UTF-8"));
                    OutputStream out = hadoopFS.append(logFilePath);

                    IOUtils.copyBytes(in, out, 4096, true);

                }
            }

        } catch (IOException e) {

        } finally {
            try {
                hadoopFS.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }


        consumer.close();
    }

}
