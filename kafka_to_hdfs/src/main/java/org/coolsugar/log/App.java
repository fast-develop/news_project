package org.coolsugar.log;

import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.coolsugar.log.util.ConfigUtil;
import org.coolsugar.log.util.HdfsConf;

import java.io.IOException;

/**
 * Hello world!
 *
 */
public class App {

    private static Boolean isDebug = false;

    private static void debug(String str) {
        if (isDebug) {
            System.out.println(str);
        }
    }

    private static void useage() {
        System.out.println("* kafka写入到hdfs的Java工具使用说明 ");
        System.out.println("# java -cp kafkatohdfs.jar KafkaToHdfs KAFKA_HOST KAFKA_GROUP KAFKA_TOPIC HDFS_URI HDFS_DIRECTORY IS_DEBUG");
        System.out.println("*  参数说明:");
        System.out.println("*   KAFKA_HOST      : 代表kafka的主机名或IP:port，例如：namenode:2181,datanode1:2181,datanode2:2181");
        System.out.println("*   KAFKA_GROUP     : 代表kafka的组，例如：test-consumer-group");
        System.out.println("*   KAFKA_TOPIC     : 代表kafka的topic名称 ，例如：usertags");
        System.out.println("*   HDFS_URI        : 代表hdfs链接uri ，例如：hdfs://namenode:9000");
        System.out.println("*   HDFS_DIRECTORY  : 代表hdfs目录名称 ，例如：usertags");
        System.out.println("*  可选参数:");
        System.out.println("*   IS_DEBUG        : 代表是否开启调试模式，true是，false否，默认为false");
    }

    public static void main( String[] args ) {

        System.out.println("开始启动服务...");


        //创建好相应的目录
        try {
            FileSystem hadoopFS = FileSystem.get(HdfsConf.getInstance().getConf());
            //如果hdfs的对应的目录不存在，则进行创建
            String hdfsDir = ConfigUtil.getInstance().getProperty("hdfsDir");
            if (!hadoopFS.exists(new Path("/" + hdfsDir))) {
                hadoopFS.mkdirs(new Path("/" + hdfsDir));
            }
            hadoopFS.close();
        } catch (IOException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }

        KafkaConsumer selfObj = new KafkaConsumer();
        selfObj.start();

        System.out.println("服务启动完毕，监听执行中");
    }
}
