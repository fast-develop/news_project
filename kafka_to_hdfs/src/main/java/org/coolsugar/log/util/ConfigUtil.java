package org.coolsugar.log.util;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

public class ConfigUtil {
    private static ConfigUtil ourInstance = new ConfigUtil();
    private Properties properties = null;

    public static ConfigUtil getInstance() {
        return ourInstance;
    }

    private ConfigUtil() {
        properties = new Properties();
        ClassLoader loader = this.getClass().getClassLoader();
//        String path = loader.getResource("config.properties").getPath();
        //打包成jar包后，获取项目下文件，报找不到文件的异常，原来打成jar包后，需要以流的形式获取jar包里面的文件，改成以下代码即可
        InputStream inputStream = loader.getResourceAsStream("config.properties");

        try {
//            properties.load(new FileInputStream(path));
            properties.load(inputStream);
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public String getProperty(String key) {
        return properties.getProperty(key);
    }


}
