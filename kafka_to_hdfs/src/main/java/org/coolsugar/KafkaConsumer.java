package org.coolsugar;


import kafka.consumer.Consumer;
import kafka.consumer.ConsumerConfig;
import kafka.consumer.ConsumerIterator;
import kafka.consumer.KafkaStream;
import kafka.javaapi.consumer.ConsumerConnector;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IOUtils;
import org.coolsugar.util.ConfigUtil;
import org.coolsugar.util.HdfsConf;
import org.coolsugar.util.KafkaConsumerConf;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.HashMap;
import java.util.List;
import java.util.Map;


public class KafkaConsumer extends Thread {

    private ConsumerConnector consumer = null;

    private static FileSystem hadoopFS = null;


    public void run() {

        consumer = Consumer.createJavaConsumerConnector(KafkaConsumerConf.getInstance().getConf());

        Map<String, Integer> topicCountMap = new HashMap<String, Integer>();
        String kafkaTopic = ConfigUtil.getInstance().getProperty("kafkaTopic");
        topicCountMap.put(kafkaTopic, new Integer(1));
        Map<String, List<KafkaStream<byte[], byte[]>>> consumerMap = consumer.createMessageStreams(topicCountMap);
        KafkaStream<byte[], byte[]> stream = consumerMap.get(kafkaTopic).get(0);
        ConsumerIterator<byte[], byte[]> it = stream.iterator();

        while (it.hasNext()) {
            String tmp = new String(it.next().message());
            String fileContent = null;

            if (!tmp.endsWith("\n"))
                fileContent = new String(tmp + "\n");
            else
                fileContent = tmp;

            try {
                hadoopFS = FileSystem.get(HdfsConf.getInstance().getConf());

                String fileName = "/" + ConfigUtil.getInstance().getProperty("hdfsDir") + "/" +
                        (new SimpleDateFormat("yyyy-MM-dd").format(Calendar.getInstance().getTime())) + ".txt";

                Path dst = new Path(fileName);

                if (!hadoopFS.exists(dst)) {
                    FSDataOutputStream output = hadoopFS.create(dst);
                    output.close();
                }

                InputStream in = new ByteArrayInputStream(fileContent.getBytes("UTF-8"));
                OutputStream out = hadoopFS.append(dst);
                IOUtils.copyBytes(in, out, 4096, true);

            } catch (IOException e) {
                e.printStackTrace();
            } finally {
                try {
                    hadoopFS.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
        consumer.shutdown();
    }


}
