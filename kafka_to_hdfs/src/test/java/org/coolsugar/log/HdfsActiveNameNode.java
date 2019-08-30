package org.coolsugar.log;

import java.io.IOException;
import java.net.URI;
import java.util.ArrayList;
import java.util.List;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileStatus;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.hdfs.server.namenode.ha.proto.HAZKInfoProtos;
import org.apache.hadoop.hdfs.server.namenode.ha.proto.HAZKInfoProtos.ActiveNodeInfo;
import org.apache.log4j.Logger;
import org.apache.zookeeper.Watcher;
import org.apache.zookeeper.ZooKeeper;
import org.apache.zookeeper.data.Stat;

import com.google.protobuf.InvalidProtocolBufferException;

/**
 * hdfs公共服务类
 * @author author
 *
 */
public class HdfsActiveNameNode {
    private static Logger log = Logger.getLogger(HdfsActiveNameNode.class.getName());

    private static final String ZOOKEEPER_IP = "192.168.1.100;192.168.1.101;192.168.1.102";
    private static final int ZOOKEEPER_PORT = 2181;
    private static final int ZOOKEEPER_TIMEOUT = 30000;
    private static final String DATA_DIR = "/hadoop";

    /**
     * 判断文件是否存在
     * @param dirs hdfs路径
     * @return 存在ture 否则false
     */
    public static boolean isExist(List<String> dirs,String ddate) {
        if(dirs==null||dirs.size()<1){
            log.info("Please Check your sign directory configure is correct!");
            return false;
        }
        log.info("Check File or Directory, uri is:"+dirs.toString());
        Configuration conf = new Configuration();
        List<String> flags = new ArrayList<String>();
        FileStatus status=null;
        String hostname = getHostname(ZOOKEEPER_IP,ZOOKEEPER_PORT,ZOOKEEPER_TIMEOUT,DATA_DIR);
        log.info("According to Zookeeper get the active namenode domain:"+hostname);
        try {
            for(String dir:dirs){
                String url = dir.replaceAll("url", hostname);
                FileSystem fs = FileSystem.get(URI.create(url+"dt="+ddate), conf);
                status = fs.getFileStatus(new Path(url+"dt="+ddate));
                if(status==null){
                    flags.add("false");
                }else{
                    flags.add("true");
                }
                fs.close();
            }
            if(flags.contains("false")){
                return false;
            }else{
                return true;
            }
        } catch (IllegalArgumentException e) {
            // TODO Auto-generated catch block
            log.error(e);
            return false;
        } catch (IOException e) {
            // TODO Auto-generated catch block
            log.error(e);
            return false;
        }
    }

    /**
     * ͨ通过zookeeper获取active namenode地址
     * @param ZOOKEEPER_IP ip地址
     * @param ZOOKEEPER_PORT 端口
     * @param ZOOKEEPER_TIMEOUT 超时时间
     * @return 地址
     */
    public static String getHostname(String ZOOKEEPER_IP, int ZOOKEEPER_PORT,
                                     int ZOOKEEPER_TIMEOUT,String DATA_DIR) {
        String hostname = null;
        Watcher watcher = new Watcher() {
            public void process(org.apache.zookeeper.WatchedEvent event) {
                log.info("event:"+event.toString());
            }
        };
        ZooKeeper zk = null;
        byte[] data1 = null;
        String[] iparr = ZOOKEEPER_IP.split(";");
        for (String ip : iparr) {
            try {
                zk = new ZooKeeper(ip + ":" + ZOOKEEPER_PORT,
                        ZOOKEEPER_TIMEOUT, watcher);
                data1 = zk.getData(DATA_DIR, true, new Stat());
            } catch (Exception e) {
                // TODO Auto-generated catch block
                log.info("This ip is not active..."+ip);
                continue;
            }
            if (data1 != null) {
                log.info("This ip is normal..."+ip);
                ActiveNodeInfo activeNodeInfo=null;
                try {
                    activeNodeInfo = HAZKInfoProtos.ActiveNodeInfo.parseFrom(data1);
                } catch (InvalidProtocolBufferException e) {
                    // TODO Auto-generated catch block
                    log.error(e);
                }
                hostname = activeNodeInfo.getHostname();
                return hostname;
            }
        }
        return hostname;
    }
}



