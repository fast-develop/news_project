package org.coolsugar.web.util;

public class BaseFactory {

    private static BaseFactory factory = new BaseFactory();

    public static BaseFactory getFactory() {
        return factory;
    }

    private BaseFactory() {

    }

    public <T> T getInstance(Class<T> cz) {
        String value = ConfigUtil.getInstance().getProperty(cz.getSimpleName());
        try {
            Class c = Class.forName(value);

            return (T) c.newInstance();
        } catch (Exception e) {
            e.printStackTrace();
        }

        return null;
    }
}
