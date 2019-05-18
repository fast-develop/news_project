package org.coolsugar.service;

import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import com.mongodb.client.AggregateIterable;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoCursor;
import net.sf.json.JSONObject;
import org.bson.Document;
import org.bson.types.Binary;
import org.coolsugar.dao.MongoDao;
import org.coolsugar.util.BaseFactory;

import javax.imageio.ImageIO;

import static com.mongodb.client.model.Filters.eq;

public class NewsServiceImpl implements NewsService {

    @Override
    public List<String> getNewsList() {

        MongoDao md = BaseFactory.getFactory().getInstance(MongoDao.class);

//        MongoCollection<Document> collection = md.getCollection("meta", "toutiao");
//        MongoCursor<Document> cursor = collection.find().iterator();

        AggregateIterable<Document> collection = md.gerCollectionRandom("meta", "toutiao", 10);
        MongoCursor<Document> cursor = collection.iterator();

        List<String> ret = new ArrayList<>();
        try {
            int index = 0;
            while(cursor.hasNext()) {
                Document doc = cursor.next();

                JSONObject item = new JSONObject();
                item.put("title", doc.getString("title"));
                item.put("_id",  System.currentTimeMillis() + index + doc.getString("docid"));
                item.put("brief", "brief");
                item.put("category", "category");
                item.put("link", doc.getString("url"));
                item.put("thumb", doc.getString("image_url"));
                item.put("publisher", doc.getString("author"));
                item.put("pubData", doc.getString("time"));
                ret.add(item.toString());
                index++;

            }
        } finally {
            cursor.close();
        }

        return ret;
    }


    @Override
    public String getNewsDetail(String id) {

        MongoDao md = BaseFactory.getFactory().getInstance(MongoDao.class);
        MongoCollection<Document> collection = md.getCollection("pages", "toutiao");
        Document doc = md.find(collection, eq("docid", id)).next();

        JSONObject item = new JSONObject();
        item.put("id", id);
        item.put("text", doc.getString("content"));


        return item.toString();
    }

    @Override
    public BufferedImage getThumb(String thumburl) {
        MongoDao md = BaseFactory.getFactory().getInstance(MongoDao.class);
        MongoCollection<Document> collection = md.getCollection("meta", "thumb");
        Document doc = md.find(collection, eq("thumb_url", thumburl + '\n')).next();
        Binary bin = doc.get("data", org.bson.types.Binary.class);

        ByteArrayInputStream bais = new ByteArrayInputStream(bin.getData());
        try {
            return ImageIO.read(bais);
        } catch (IOException e) {
            e.printStackTrace();
        }


        return null;
    }


}
