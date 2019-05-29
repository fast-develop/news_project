package org.coolsugar.util;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
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
        String path = loader.getResource("config.properties").getPath();

        try {
            properties.load(new FileInputStream(path));
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
