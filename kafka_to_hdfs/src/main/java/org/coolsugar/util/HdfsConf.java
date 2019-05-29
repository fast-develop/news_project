package org.coolsugar.util;

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
