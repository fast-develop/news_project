package org.coolsugar;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.IOUtils;
import org.coolsugar.log.util.ConfigUtil;

import java.io.ByteArrayInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.URI;

public class HdfsTools {

    private static final String HDFSUri = ConfigUtil.getInstance().getProperty("hdfsUri");


    public static void Get() throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs = FileSystem.get(new URI(ConfigUtil.getInstance().getProperty("hdfsUri")), conf);
        InputStream in = fs.open(new Path("/news_detail_log/2019-06-08.txt.txt"));
        FileOutputStream out = new FileOutputStream("log.txt");
        IOUtils.copyBytes(in, out, conf);
    }


    public static void Put() throws Exception{
        Configuration conf=new Configuration();
//        conf.set("dfs.replication","1");

        FileSystem fs = FileSystem.get(new URI(HDFSUri),conf,"root");
        OutputStream out=fs.create(new Path("/b.txt"));
        ByteArrayInputStream in=new ByteArrayInputStream("hello hdfs".getBytes());

        IOUtils.copyBytes(in, out, conf);
    }


    public static void Delete()throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs=FileSystem.get(new URI(HDFSUri),conf,"root");
//        fs.delete(new Path("/c.txt"),true);
        fs.delete(new Path("/park02"),false);
        fs.close();
    }


    public static void Mkdir(String path)throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs=FileSystem.get(new URI(HDFSUri),conf,"root");

//        fs.mkdirs(new Path("/park02"));
    }



    public static void Ls()throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs=FileSystem.get(new URI(HDFSUri),conf,"root");
        RemoteIterator<LocatedFileStatus> rt=fs.listFiles(new Path("/"), true);
        while(rt.hasNext()){
            System.out.println("file: " + rt.next());
        }
    }


    public static void Rename() throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs=FileSystem.get(new URI(HDFSUri), conf, "root");
        fs.rename(new Path("/park"), new Path("/park01"));
    }


    public static void CopyFromLoaclFileSystem() throws Exception{
        Configuration conf=new Configuration();
        FileSystem fs=FileSystem.get(new URI(HDFSUri),conf,"root");
        BlockLocation[] data=fs.getFileBlockLocations(new Path("/park01/1.txt"), 0, Integer.MAX_VALUE);
        for(BlockLocation bl:data){
            System.out.println(bl);
        }
    }


}
