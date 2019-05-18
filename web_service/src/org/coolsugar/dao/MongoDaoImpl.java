package org.coolsugar.dao;

import com.mongodb.*;
import com.mongodb.client.*;
import com.mongodb.client.model.Filters;
import com.mongodb.client.result.DeleteResult;
import org.bson.Document;
import org.bson.conversions.Bson;
import org.bson.types.ObjectId;
import org.coolsugar.util.MongoUtil;

import java.util.ArrayList;
import java.util.List;

import static com.mongodb.client.model.Filters.eq;

public class MongoDaoImpl implements MongoDao {
    private MongoClient mongoClient = null;
    
    public MongoDaoImpl() {
        mongoClient = MongoUtil.getInstance().getMongoClient();
    }

    @Override
    public MongoDatabase getDB(String dbName) {
        if (dbName != null && !"".equals(dbName)) {
            MongoDatabase database = mongoClient.getDatabase(dbName);
            return database;
        }

        return null;
    }


    @Override
    public MongoCollection<Document> getCollection(String dbName, String collName) {
        if (null == collName || "".equals(collName)) {
            return null;
        }
        if (null == dbName || "".equals(dbName)) {
            return null;
        }
        MongoCollection<Document> collection = mongoClient.getDatabase(dbName).getCollection(collName);

        return collection;
    }


    @Override
    public AggregateIterable<Document> gerCollectionRandom(String dbName, String collName, int size) {

        MongoCollection<Document> collection = mongoClient.getDatabase(dbName).getCollection(collName);

        Document sub_match = new Document();
        sub_match.put("size", size);
        Document match = new Document("$sample", sub_match);
        List<Document> aggregateList = new ArrayList<>();
        aggregateList.add(match);

        return collection.aggregate(aggregateList);
    }


    @Override
    public List<String> getAllCollections(String dbName) {
        MongoIterable<String> colls = getDB(dbName).listCollectionNames();
        List<String> _list = new ArrayList<String>();
        for (String s : colls) {
            _list.add(s);
        }

        return _list;
    }


    @Override
    public MongoIterable<String> getAllDBNames() {
        return mongoClient.listDatabaseNames();
    }


    @Override
    public void dropDB(String dbName) {
        getDB(dbName).drop();
    }


    @Override
    public Document findById(MongoCollection<Document> coll, String id) {

        ObjectId _idobj = null;
        try {
            _idobj = new ObjectId(id);
        } catch (Exception e) {
            return null;
        }

        Document myDoc = coll.find(Filters.eq("_id", _idobj)).first();
        return myDoc;
    }


    @Override
    public int getCount(MongoCollection<Document> coll) {
        int count = (int) coll.count();

        return count;
    }


    @Override
    public MongoCursor<Document> find(MongoCollection<Document> coll, Bson filter) {
        return coll.find(filter).iterator();
    }


    @Override
    public MongoCursor<Document> findByPage(MongoCollection<Document> coll, Bson filter, int pageNo, int pageSize) {
        Bson orderBy = new BasicDBObject("_id", 1);

        return coll.find(filter).sort(orderBy).skip((pageNo - 1) * pageSize).limit(pageSize).iterator();
    }


    @Override
    public int deleteById(MongoCollection<Document> coll, String id) {
        int count = 0;
        ObjectId _id = null;
        try {
            _id = new ObjectId(id);
        } catch (Exception e) {
            return 0;
        }

        Bson filter = eq("_id", _id);
        DeleteResult deleteResult = coll.deleteOne(filter);
        count = (int) deleteResult.getDeletedCount();

        return count;
    }


    @Override
    public Document updateById(MongoCollection<Document> coll, String id, Document newdoc) {
        ObjectId _idobj = null;
        try {
            _idobj = new ObjectId(id);
        } catch (Exception e) {
            return null;
        }

        Bson filter = eq("_id", _idobj);
        // coll.replaceOne(filter, newdoc); // 完全替代
        coll.updateOne(filter, new Document("$set", newdoc));

        return newdoc;
    }


    @Override
    public void dropCollection(String dbName, String collName) {
        getDB(dbName).getCollection(collName).drop();
    }


    @Override
    public void close() {
        if (mongoClient != null) {
            mongoClient.close();
            mongoClient = null;
        }
    }

}
