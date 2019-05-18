package org.coolsugar.dao;

import com.mongodb.DBObject;
import com.mongodb.client.*;
import org.bson.Document;
import org.bson.conversions.Bson;

import java.util.ArrayList;
import java.util.List;

public interface MongoDao {

    /**
     * 获取DB实例 - 指定DB
     *
     * @param dbName
     * @return
     */
    MongoDatabase getDB(String dbName);

    /**
     * 获取collection对象 - 指定Collection
     *
     * @param dbName
     * @param collName
     * @return
     */
    MongoCollection<Document> getCollection(String dbName, String collName);

    /**
     * 查询DB下的所有表名
     *
     * @param dbName
     * @return
     */
    List<String> getAllCollections(String dbName);

    /**
     * 获取所有数据库名称列表
     *
     * @return
     */
    MongoIterable<String> getAllDBNames();

    /**
     * 删除一个数据库
     *
     * @param dbName
     */
    void dropDB(String dbName);

    /**
     * 查找对象 - 根据主键_id
     *
     * @param coll
     * @param id
     * @return
     */
    Document findById(MongoCollection<Document> coll, String id);

    /**
     * 统计数
     *
     * @param coll
     * @return
     */
    int getCount(MongoCollection<Document> coll);

    /**
     * 条件查询
     *
     * @param coll
     * @param filter
     * @return
     */
    MongoCursor<Document> find(MongoCollection<Document> coll, Bson filter);


    /**
     * 分页查询
     *
     * @param coll
     * @param filter
     * @param pageNo
     * @param pageSize
     * @return
     */
    MongoCursor<Document> findByPage(MongoCollection<Document> coll, Bson filter, int pageNo, int pageSize);

    /**
     * 通过ID删除
     *
     * @param coll
     * @param id
     * @return
     */
    int deleteById(MongoCollection<Document> coll, String id);

    /**
     * 通过ID更改
     *
     * @param coll
     * @param id
     * @param newdoc
     * @return
     */
    Document updateById(MongoCollection<Document> coll, String id, Document newdoc);

    /**
     * 删除 collection
     *
     * @param dbName
     * @param collName
     */
    public void dropCollection(String dbName, String collName);

    /**
     * 关闭Mongodb
     */
    void close();

    /**
     * 随机获取
     *
     * @param dbName
     * @param collName
     * @param size
     * @return
     */
    AggregateIterable<Document> gerCollectionRandom(String dbName, String collName, int size);


}
