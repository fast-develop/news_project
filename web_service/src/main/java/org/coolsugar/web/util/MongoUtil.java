package org.coolsugar.web.util;

import com.mongodb.MongoClient;

public class MongoUtil {
    private static MongoUtil ourInstance = new MongoUtil();
    private MongoClient mongoClient = null;

    public static MongoUtil getInstance() {
        return ourInstance;
    }

    private MongoUtil() {
        String host = ConfigUtil.getInstance().getProperty("mongodb_ip");
        String port = ConfigUtil.getInstance().getProperty("mongodb_port");

        // 连接到 mongodb 服务
        mongoClient = new MongoClient(host, Integer.parseInt(port));
    }


    public MongoClient getMongoClient() {
        return mongoClient;
    }
}
