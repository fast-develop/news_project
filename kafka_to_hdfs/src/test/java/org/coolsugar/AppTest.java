package org.coolsugar;

import static org.junit.Assert.assertTrue;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.URI;
import java.net.URISyntaxException;

import org.apache.commons.compress.utils.IOUtils;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;


import org.junit.Test;


/**
 * Unit test for simple App.
 */
public class AppTest 
{
    /**
     * Rigorous Test :-)
     */
    @Test
    public void shouldAnswerWithTrue () {
        assertTrue( true );

    }

    @Test
    public void get() throws IOException, URISyntaxException {

        Configuration conf = new Configuration();

        FileSystem fs = FileSystem.get(new URI("hdfs://192.168.1.101:9000"), conf);
        InputStream in = fs.open(new Path("/log/2019-05-29.txt"));
        FileOutputStream out = new FileOutputStream("log.txt");
        IOUtils.copy(in, out);

    }


    @Test
    public void put() throws IOException, URISyntaxException, InterruptedException{

        FileSystem fs = FileSystem.get(
                new URI("hdfs://192.168.1.101:9000"), new Configuration(), "root");

        OutputStream out = fs.create(new Path("/c.txt"));
        FileInputStream in = new FileInputStream("a.txt");
        IOUtils.copy(in, out);
    }
}
