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


import org.coolsugar.log.util.ConfigUtil;
import org.junit.Before;
import org.junit.Test;


/**
 * Unit test for simple App.
 */
public class TestHdfs {

    @Before
    public void before () {
        System.setProperty("hadoop.home.dir", "/usr/local/Cellar/hadoop/3.1.1");
    }

    /**
     * Rigorous Test :-)
     */
    @Test
    public void shouldAnswerWithTrue () {
        assertTrue( true );

    }

    @Test
    public void test1 () {

        try {
            HdfsTools.Ls();
//            HdfsTools.Get();
//            HdfsTools.Delete();
//            HdfsTools.Mkdir("/aa");
        } catch (Exception e) {
            System.out.println(e);
        }
    }
}
