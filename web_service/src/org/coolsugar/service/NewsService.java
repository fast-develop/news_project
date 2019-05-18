package org.coolsugar.service;

import java.awt.image.BufferedImage;
import java.util.List;

public interface NewsService {

    List<String> getNewsList();

    String getNewsDetail(String id);

    BufferedImage getThumb(String thumburl);
}
