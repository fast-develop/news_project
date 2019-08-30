package org.coolsugar.log.util;

import org.apache.hadoop.conf.Configuration;


public class HdfsConf {

    private static HdfsConf ourInstance = new HdfsConf();
    private Configuration conf = null;

    private HdfsConf () {
        conf = new Configuration();
        try {
            conf.set("fs.defaultFS", ConfigUtil.getInstance().getProperty("hdfsUri"));
            conf.set("dfs.support.append", "true");
            conf.set("dfs.client.block.write.replace-datanode-on-failure.policy", "NEVER");
            conf.set("dfs.client.block.write.replace-datanode-on-failure.enable", "true");

            // 通过java api连接Hadoop集群时，如果集群支持HA方式，那么可以通过如下方式设置来自动切换到活动的master节点上。
            // 其中，ClusterName 是可以任意指定的，跟集群配置无关
//            conf.set("dfs.nameservices", ClusterName);
//            conf.set("dfs.ha.namenodes."+ClusterName, "nn1,nn2");
//            conf.set("dfs.namenode.rpc-address."+ClusterName+".nn1", "172.16.50.24:8020");
//            conf.set("dfs.namenode.rpc-address."+ClusterName+".nn2", "172.16.50.21:8020");
//            //conf.setBoolean(name, value);
//            conf.set("dfs.client.failover.proxy.provider."+ClusterName,
//                    "org.apache.hadoop.hdfs.server.namenode.ha.ConfiguredFailoverProxyProvider");

        } catch (Exception e) {
            System.out.println(e);
        }
    }

    public static HdfsConf getInstance() {
        return ourInstance;
    }

    public Configuration getConf() {
        return conf;
    }
}
